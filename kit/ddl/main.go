package main

import (
	"cloud/unionj/papilio/kit/astutils"
	"cloud/unionj/papilio/kit/stringutils"
	"flag"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
)

type DbConfig struct {
	host    string
	port    string
	user    string
	passwd  string
	schema  string
	charset string
}

var file = flag.String("file", "/Users/wubin1989/workspace/cloud/papilio/kit/ddl/example/models/user.go", "name of file")

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
		//log.Fatalln(err)
	}
	fmt.Println(db)

	flag.Parse()
	log.Println(*file)
	if stringutils.IsEmpty(*file) {
		log.Fatal("file flag should not be empty")
	}

	fset := token.NewFileSet()
	root, err := parser.ParseFile(fset, *file, nil, 0)
	if err != nil {
		panic(err)
	}
	var sc astutils.StructCollector
	ast.Walk(&sc, root)
	fmt.Println(sc.Structs)

}
