package codegen

import (
	"bufio"
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

var dbTmpl = `package db

import (
	"{{.ConfigPackage}}"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/iancoleman/strcase"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

func NewDb(conf config.DbConfig) (*sqlx.DB, error) {
	conn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s",
		conf.User,
		conf.Passwd,
		conf.Host,
		conf.Port,
		conf.Schema,
		conf.Charset)
	conn += "&loc=Asia%2FShanghai&parseTime=True"

	db, err := sqlx.Connect(conf.Driver, conn)
	if err != nil {
		return nil, errors.Wrap(err, "database connection failed")
	}
	db.MapperFunc(strcase.ToSnake)
	return db, nil
}
`

func GenDb(dir string) {
	var (
		err       error
		dbfile    string
		f         *os.File
		tpl       *template.Template
		dbDir     string
		modfile   string
		modName   string
		firstLine string
	)
	dbDir = filepath.Join(dir, "db")
	if err = os.MkdirAll(dbDir, 0644); err != nil {
		panic(err)
	}

	dbfile = filepath.Join(dbDir, "db.go")
	if _, err = os.Stat(dbfile); os.IsNotExist(err) {
		modfile = filepath.Join(dir, "go.mod")
		if f, err = os.Open(modfile); err != nil {
			panic(err)
		}
		reader := bufio.NewReader(f)
		if firstLine, err = reader.ReadString('\n'); err != nil {
			panic(err)
		}
		modName = strings.TrimSpace(strings.TrimPrefix(firstLine, "module"))

		if f, err = os.Create(dbfile); err != nil {
			panic(err)
		}
		defer f.Close()

		if tpl, err = template.New("db.go.tmpl").Parse(dbTmpl); err != nil {
			panic(err)
		}
		if err = tpl.Execute(f, struct {
			ConfigPackage string
		}{
			ConfigPackage: modName + "/config",
		}); err != nil {
			panic(err)
		}
	} else {
		logrus.Warnf("file %s already exists", dbfile)
	}
}
