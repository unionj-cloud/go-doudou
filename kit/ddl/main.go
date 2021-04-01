package main

import (
	"flag"
	"fmt"
	mapset "github.com/deckarep/golang-set"
	_ "github.com/go-sql-driver/mysql"
	"github.com/iancoleman/strcase"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"github.com/unionj-cloud/go-doudou/kit/astutils"
	"github.com/unionj-cloud/go-doudou/kit/ddl/cmd"
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
	host    string
	port    string
	user    string
	passwd  string
	schema  string
	charset string
}

var dir = flag.String("domain", "/Users/wubin1989/workspace/cloud/go-doudou/kit/ddl/example/domain", "path of domain folder")

var dbConfig DbConfig

func init() {
	err := godotenv.Load(pathutils.Abs(".env"))
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	dbConfig = DbConfig{
		host:    os.Getenv("DB_HOST"),
		port:    os.Getenv("DB_PORT"),
		user:    os.Getenv("DB_USER"),
		passwd:  os.Getenv("DB_PASSWD"),
		schema:  os.Getenv("DB_SCHEMA"),
		charset: os.Getenv("DB_CHARSET"),
	}
}

func main() {
	flag.Parse()
	log.Println(*dir)
	if stringutils.IsEmpty(*dir) {
		log.Fatal("dir flag should not be empty")
	}

	var files []string
	var err error
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
	fmt.Println(sc.Structs)

	flattened := sc.FlatEmbed()
	fmt.Println(flattened)

	conn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s",
		dbConfig.user,
		dbConfig.passwd,
		dbConfig.host,
		dbConfig.port,
		dbConfig.schema,
		dbConfig.charset)
	conn += `&loc=Asia%2FShanghai&parseTime=True`
	db, err := sqlx.Connect("mysql", conn)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(db)
	db.MapperFunc(strcase.ToSnake)

	var existTables []string
	if err = db.Select(&existTables, "show tables"); err != nil {
		panic(err)
	}
	fmt.Println(existTables)

	for _, sm := range flattened {
		t := table.NewTableFromStruct(sm)
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
				log.Infof("FATAL: %+v\n", err)
			}
		}
	}
}
