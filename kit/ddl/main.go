package main

import (
	"flag"
	"fmt"
	mapset "github.com/deckarep/golang-set"
	_ "github.com/go-sql-driver/mysql"
	"github.com/iancoleman/strcase"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	log "github.com/sirupsen/logrus"
	"github.com/unionj-cloud/go-doudou/kit/astutils"
	"github.com/unionj-cloud/go-doudou/kit/ddl/cmd"
	"github.com/unionj-cloud/go-doudou/kit/ddl/codegen/dao"
	"github.com/unionj-cloud/go-doudou/kit/ddl/table"
	"github.com/unionj-cloud/go-doudou/kit/pathutils"
	"github.com/unionj-cloud/go-doudou/kit/sliceutils"
	"github.com/unionj-cloud/go-doudou/kit/stringutils"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
)

type DbConfig struct {
	Host    string
	Port    string
	User    string
	Passwd  string
	Schema  string
	Charset string
}

var dir = flag.String("domain", "", "path of domain folder")

func main() {
	var db *sqlx.DB
	err := godotenv.Load(pathutils.Abs(".env"))
	if err != nil {
		log.Fatal("Error loading .env file", err)
	}
	var dbConfig DbConfig
	err = envconfig.Process("db", &dbConfig)
	if err != nil {
		log.Fatal("Error processing env", err)
	}

	conn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s",
		dbConfig.User,
		dbConfig.Passwd,
		dbConfig.Host,
		dbConfig.Port,
		dbConfig.Schema,
		dbConfig.Charset)
	conn += `&loc=Asia%2FShanghai&parseTime=True`
	db, err = sqlx.Connect("mysql", conn)
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()
	db.MapperFunc(strcase.ToSnake)

	flag.Parse()
	log.Println(*dir)
	if stringutils.IsEmpty(*dir) {
		if wd, err := os.Getwd(); err != nil {
			log.Fatal(err)
		} else {
			*dir = filepath.Join(wd, "domain")
		}
	}
	if !filepath.IsAbs(*dir) {
		if wd, err := os.Getwd(); err != nil {
			log.Fatal(err)
		} else {
			*dir = filepath.Join(wd, *dir)
		}
	}

	var files []string
	err = filepath.Walk(*dir, astutils.Visit(&files))
	if err != nil {
		panic(err)
	}
	var sc astutils.StructCollector
	for _, file := range files {
		fset := token.NewFileSet()
		root, err := parser.ParseFile(fset, file, nil, parser.ParseComments)
		if err != nil {
			panic(err)
		}
		ast.Walk(&sc, root)
	}

	flattened := sc.FlatEmbed()

	var existTables []string
	if err = db.Select(&existTables, "show tables"); err != nil {
		panic(err)
	}

	var tables []table.Table
	for _, sm := range flattened {
		tables = append(tables, table.NewTableFromStruct(sm))
	}
	for _, t := range tables {
		if sliceutils.StringContains(existTables, t.Name) {
			var columns []table.DbColumn
			if err = db.Select(&columns, fmt.Sprintf("desc %s", t.Name)); err != nil {
				panic(err)
			}
			var existColumnNames []interface{}
			for _, dbCol := range columns {
				existColumnNames = append(existColumnNames, dbCol.Field)
			}
			existColSet := mapset.NewSetFromSlice(existColumnNames)

			for _, col := range t.Columns {
				if existColSet.Contains(col.Name) {
					if err = cmd.ChangeColumn(db, col); err != nil {
						log.Infof("FATAL: %+v\n", err)
					}
				} else {
					if err = cmd.AddColumn(db, col); err != nil {
						log.Infof("FATAL: %+v\n", err)
					}
				}
			}
		} else {
			if err = cmd.CreateTable(db, t); err != nil {
				log.Errorf("FATAL: %+v\n", err)
			}
		}
	}

	for _, t := range tables {
		if err = dao.GenDaoGo(*dir, t); err != nil {
			log.Errorf("FATAL: %+v\n", err)
			break
		}
		if err = dao.GenDaoImplGo(*dir, t); err != nil {
			log.Errorf("FATAL: %+v\n", err)
			break
		}
		if err = dao.GenDaoSql(*dir, t); err != nil {
			log.Errorf("FATAL: %+v\n", err)
			break
		}
	}
}
