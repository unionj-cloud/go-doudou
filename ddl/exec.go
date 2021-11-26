package ddl

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"strings"
	"time"

	// here must import mysql
	_ "github.com/go-sql-driver/mysql"
	"github.com/iancoleman/strcase"
	"github.com/jmoiron/sqlx"
	"github.com/unionj-cloud/go-doudou/ddl/codegen"
	"github.com/unionj-cloud/go-doudou/ddl/config"
	"github.com/unionj-cloud/go-doudou/ddl/table"
)

// Ddl is for ddl command
type Ddl struct {
	Dir     string
	Reverse bool
	Dao     bool
	Pre     string
	Df      string
	Conf    config.DbConfig
}

// Exec executes the logic for ddl command
// if Reverse is true, it will generate code from database tables,
// otherwise it will update database tables from structs defined in domain pkg
func (d Ddl) Exec() {
	var db *sqlx.DB
	var err error
	conf := d.Conf
	conn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s",
		conf.User,
		conf.Passwd,
		conf.Host,
		conf.Port,
		conf.Schema,
		conf.Charset)
	conn += `&loc=Asia%2FShanghai&parseTime=True`
	db, err = sqlx.Connect("mysql", conn)
	if err != nil {
		panic(fmt.Sprintf("%+v", err))
	}
	defer db.Close()
	db.MapperFunc(strcase.ToSnake)
	db = db.Unsafe()

	var existTables []string
	if err = db.Select(&existTables, "show tables"); err != nil {
		panic(fmt.Sprintf("%+v", err))
	}

	var tables []table.Table
	timeoutCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	_ = os.MkdirAll(d.Dir, os.ModePerm)
	if !d.Reverse {
		tables = table.Struct2Table(timeoutCtx, d.Dir, d.Pre, existTables, db, d.Conf.Schema)
	} else {
		tables = table.Table2struct(timeoutCtx, d.Pre, d.Conf.Schema, existTables, db)
		for _, item := range tables {
			dfile := filepath.Join(d.Dir, strings.ToLower(item.Meta.Name)+".go")
			if _, err = os.Stat(dfile); os.IsNotExist(err) {
				if err = codegen.GenDomainGo(d.Dir, item.Meta); err != nil {
					panic(fmt.Sprintf("%+v", err))
				}
			} else {
				logrus.Warnf("file %s already exists", dfile)
			}
		}
	}

	if d.Dao {
		genDao(d, tables)
	}
}

func genDao(d Ddl, tables []table.Table) {
	var err error
	if err = codegen.GenBaseGo(d.Dir, d.Df); err != nil {
		panic(fmt.Sprintf("%+v", err))
	}
	for _, t := range tables {
		if err = codegen.GenDaoGo(d.Dir, t, d.Df); err != nil {
			panic(fmt.Sprintf("%+v", err))
		}
		if err = codegen.GenDaoImplGo(d.Dir, t, d.Df); err != nil {
			panic(fmt.Sprintf("%+v", err))
		}
		if err = codegen.GenDaoSQL(d.Dir, t, d.Df); err != nil {
			panic(fmt.Sprintf("%+v", err))
		}
	}
}
