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
	"github.com/unionj-cloud/go-doudou/kit/ddl/extraenum"
	"github.com/unionj-cloud/go-doudou/kit/ddl/sortenum"
	"github.com/unionj-cloud/go-doudou/kit/ddl/table"
	"github.com/unionj-cloud/go-doudou/kit/pathutils"
	"github.com/unionj-cloud/go-doudou/kit/sliceutils"
	"github.com/unionj-cloud/go-doudou/kit/stringutils"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

type DbConfig struct {
	Host    string
	Port    string
	User    string
	Passwd  string
	Schema  string
	Charset string
}

var dir = flag.String("domain", "/Users/wubin1989/workspace/cloud/go-doudou/kit/ddl/example/domain", "path of domain folder")
var reverse = flag.Bool("reverse", false, "If true, generate domain and dao code from database. If false, update or create database tables from domain code."+
	"Default is false")

func main() {
	var db *sqlx.DB
	err := godotenv.Load(pathutils.Abs(".env"))
	if err != nil {
		log.Panicln("Error loading .env file", err)
	}
	var dbConfig DbConfig
	err = envconfig.Process("db", &dbConfig)
	if err != nil {
		log.Panicln("Error processing env", err)
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
		log.Panicln(err)
	}
	defer db.Close()
	db.MapperFunc(strcase.ToSnake)

	flag.Parse()
	log.Println(*dir)
	log.Println(*reverse)

	var existTables []string
	if err = db.Select(&existTables, "show tables"); err != nil {
		log.Panicln(err)
	}

	if !*reverse {
		if stringutils.IsEmpty(*dir) {
			if wd, err := os.Getwd(); err != nil {
				log.Panicln(err)
			} else {
				*dir = filepath.Join(wd, "domain")
			}
		}
		if !filepath.IsAbs(*dir) {
			if wd, err := os.Getwd(); err != nil {
				log.Panicln(err)
			} else {
				*dir = filepath.Join(wd, *dir)
			}
		}

		var files []string
		err = filepath.Walk(*dir, astutils.Visit(&files))
		if err != nil {
			log.Panicln(err)
		}
		var sc astutils.StructCollector
		for _, file := range files {
			fset := token.NewFileSet()
			root, err := parser.ParseFile(fset, file, nil, parser.ParseComments)
			if err != nil {
				log.Panicln(err)
			}
			ast.Walk(&sc, root)
		}

		flattened := sc.FlatEmbed()

		var tables []table.Table
		for _, sm := range flattened {
			tables = append(tables, table.NewTableFromStruct(sm))
		}
		for _, t := range tables {
			if sliceutils.StringContains(existTables, t.Name) {
				var columns []table.DbColumn
				if err = db.Select(&columns, fmt.Sprintf("desc %s", t.Name)); err != nil {
					log.Panicln(err)
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
	} else {
		dfolder := pathutils.Abs("domain")
		if err = os.MkdirAll(dfolder, os.ModePerm); err != nil {
			log.Panicln(err)
		}
		for _, t := range existTables {
			domain := filepath.Join(dfolder, strings.ToLower(strcase.ToCamel(strings.TrimSuffix(t, "s")))+".go")
			if _, err = os.Stat(domain); os.IsNotExist(err) {
				var f *os.File
				if f, err = os.Create(domain); err != nil {
					log.Panicln(err)
				}
				defer f.Close()

				var dbIndice []table.DbIndex
				if err = db.Select(&dbIndice, fmt.Sprintf("SHOW INDEXES FROM %s", t)); err != nil {
					log.Panicln(err)
				}

				idxMap := make(map[string][]table.DbIndex)

				for _, idx := range dbIndice {
					if val, exists := idxMap[idx.Key_name]; exists {
						val = append(val, idx)
						idxMap[idx.Key_name] = val
					} else {
						idxMap[idx.Key_name] = []table.DbIndex{
							idx,
						}
					}
				}

				var indexes []table.Index
				for k, v := range idxMap {
					if len(v) == 0 {
						continue
					}
					items := make([]table.IndexItem, len(v))
					for i, idx := range v {
						var sort sortenum.Sort
						if idx.Collation == "B" {
							sort = sortenum.Desc
						} else {
							sort = sortenum.Asc
						}
						items[i] = table.IndexItem{
							Column: idx.Column_name,
							Order:  idx.Seq_in_index,
							Sort:   sort,
						}
					}
					indexes = append(indexes, table.Index{
						Unique: !v[0].Non_unique,
						Name:   k,
						Items:  items,
					})
				}

				var tab table.Table
				tab.Name = t
				tab.Indexes = indexes

				var columns []table.DbColumn
				if err = db.Select(&columns, fmt.Sprintf("SHOW FULL COLUMNS FROM %s", t)); err != nil {
					log.Panicln(err)
				}

				var cols []table.Column
				for _, item := range columns {
					if strings.Contains(item.Extra, "auto_increment") {
						item.Extra = ""
					}
					col := table.Column{
						Table:         t,
						Name:          item.Field,
						Type:          table.DbColType2ColumnType(item.Type),
						Default:       item.Default,
						Pk:            table.CheckPk(item.Key),
						Nullable:      table.CheckNull(item.Null),
						Unsigned:      table.CheckUnsigned(item.Type),
						Autoincrement: table.CheckAutoincrement(item.Extra),
						Extra:         extraenum.Extra(item.Extra),
						Meta: astutils.FieldMeta{
							Name:     "",
							Type:     "",
							Tag:      "",
							Comments: nil,
						}, // TODO
						AutoSet: table.CheckAutoSet(item.Default),
					}
					cols = append(cols, col)
				}

				tab.Columns = cols
				// TODO
			} else {
				log.Warnf("file %s already exists", domain)
			}
		}

	}

}
