package main

import (
	"flag"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/unionj-cloud/go-doudou/kit/astutils"
	"github.com/unionj-cloud/go-doudou/kit/ddl/table"
	"github.com/unionj-cloud/go-doudou/kit/sliceutils"
	"github.com/unionj-cloud/go-doudou/kit/stringutils"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
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

var dir = flag.String("models", "/Users/wubin1989/workspace/cloud/go-doudou/kit/ddl/example/models", "path of models folder")

var dbConfig DbConfig

func init() {
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

	var existTables []string
	err = db.Select(&existTables, "show tables")
	if err != nil {
		panic(err)
	}
	fmt.Println(existTables)

	for _, sm := range flattened {
		table := table.NewTableFromStruct(sm)
		fmt.Println(table)
		if sliceutils.StringContains(existTables, table.Name) {
			continue
		}
		var statement string
		if statement, err = table.CreateSql(); err != nil {
			fmt.Printf("FATAL: %+v\n", err)
			return
		}
		if _, err = db.Exec(statement); err != nil {
			fmt.Printf("FATAL: %+v\n", err)
			return
		}
	}
}
