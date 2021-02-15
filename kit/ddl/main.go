package main

import (
	"cloud/unionj/papilio/kit/astutils"
	"cloud/unionj/papilio/kit/sliceutils"
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

var dir = flag.String("models", "/Users/wubin1989/workspace/cloud/papilio/kit/ddl/example/models", "path of models folder")

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

func visit(files *[]string) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Fatal(err)
		}
		if !info.IsDir() {
			*files = append(*files, path)
		}
		return nil
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
	err = filepath.Walk(*dir, visit(&files))
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

	structMap := make(map[string]astutils.StructMeta)
	for _, structMeta := range sc.Structs {
		if _, exists := structMap[structMeta.Name]; !exists {
			structMap[structMeta.Name] = structMeta
		}
	}

	var tables []astutils.StructMeta
	for _, structMeta := range sc.Structs {
		if sliceutils.IsEmpty(structMeta.Comments) {
			continue
		}
		if structMeta.Comments[0] != "//papi:table" {
			continue
		}
		_structMeta := astutils.StructMeta{
			Name:     structMeta.Name,
			Fields:   make([]astutils.FieldMeta, 0),
			Comments: make([]string, len(structMeta.Comments)),
		}
		copy(_structMeta.Comments, structMeta.Comments)

		fieldMap := make(map[string]astutils.FieldMeta)
		embedFieldMap := make(map[string]astutils.FieldMeta)
		for _, fieldMeta := range structMeta.Fields {
			if fieldMeta.Type == "embed" {
				if embeded, exists := structMap[fieldMeta.Name]; exists {
					for _, field := range embeded.Fields {
						if _, _exists := embedFieldMap[field.Name]; !_exists {
							embedFieldMap[field.Name] = field
						}
					}
				}
			} else {
				_structMeta.Fields = append(_structMeta.Fields, fieldMeta)
				fieldMap[fieldMeta.Name] = fieldMeta
			}
		}

		for key, field := range embedFieldMap {
			if _, exists := fieldMap[key]; !exists {
				_structMeta.Fields = append(_structMeta.Fields, field)
			}
		}
		tables = append(tables, _structMeta)
	}

	fmt.Println(tables)

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




}
