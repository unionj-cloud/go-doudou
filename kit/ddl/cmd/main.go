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
	"github.com/sirupsen/logrus"
	"github.com/unionj-cloud/go-doudou/kit/astutils"
	"github.com/unionj-cloud/go-doudou/kit/ddl/codegen"
	"github.com/unionj-cloud/go-doudou/kit/ddl/columnenum"
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

var dir = flag.String("domain", "/Users/wubin1989/workspace/cloud/go-doudou/kit/ddl/cmd/domain", "path of domain folder")
var reverse = flag.Bool("reverse", false, "If true, generate domain code from database. If false, update or create database tables from domain code."+
	"Default is false")
var dao = flag.Bool("dao", false, "If true, generate dao code. Default is false.")

func init() {
	customFormatter := new(logrus.TextFormatter)
	customFormatter.TimestampFormat = "2006-01-02 15:04:05"
	customFormatter.FullTimestamp = true
	logrus.SetFormatter(customFormatter)
}

func main() {
	var db *sqlx.DB
	err := godotenv.Load(pathutils.Abs(".env"))
	if err != nil {
		logrus.Panicln("Error loading .env file", err)
	}
	var dbConfig DbConfig
	err = envconfig.Process("db", &dbConfig)
	if err != nil {
		logrus.Panicln("Error processing env", err)
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
		logrus.Panicln(err)
	}
	defer db.Close()
	db.MapperFunc(strcase.ToSnake)
	db = db.Unsafe()

	flag.Parse()
	logrus.Println(*dir)
	logrus.Println(*reverse)
	logrus.Println(*dao)

	var existTables []string
	if err = db.Select(&existTables, "show tables"); err != nil {
		logrus.Panicln(err)
	}

	var tables []table.Table
	if !*reverse {
		if stringutils.IsEmpty(*dir) {
			if wd, err := os.Getwd(); err != nil {
				logrus.Panicln(err)
			} else {
				*dir = filepath.Join(wd, "domain")
			}
		}
		if !filepath.IsAbs(*dir) {
			if wd, err := os.Getwd(); err != nil {
				logrus.Panicln(err)
			} else {
				*dir = filepath.Join(wd, *dir)
			}
		}

		var files []string
		err = filepath.Walk(*dir, astutils.Visit(&files))
		if err != nil {
			logrus.Panicln(err)
		}
		var sc astutils.StructCollector
		for _, file := range files {
			fset := token.NewFileSet()
			root, err := parser.ParseFile(fset, file, nil, parser.ParseComments)
			if err != nil {
				logrus.Panicln(err)
			}
			ast.Walk(&sc, root)
		}

		flattened := sc.FlatEmbed()
		for _, sm := range flattened {
			tables = append(tables, table.NewTableFromStruct(sm))
		}
		for _, t := range tables {
			if sliceutils.StringContains(existTables, t.Name) {
				var columns []table.DbColumn
				if err = db.Select(&columns, fmt.Sprintf("desc %s", t.Name)); err != nil {
					logrus.Panicln(err)
				}
				var existColumnNames []interface{}
				for _, dbCol := range columns {
					existColumnNames = append(existColumnNames, dbCol.Field)
				}
				existColSet := mapset.NewSetFromSlice(existColumnNames)

				for _, col := range t.Columns {
					if existColSet.Contains(col.Name) {
						if err = table.ChangeColumn(db, col); err != nil {
							logrus.Infof("FATAL: %+v\n", err)
						}
					} else {
						if err = table.AddColumn(db, col); err != nil {
							logrus.Infof("FATAL: %+v\n", err)
						}
					}
				}
			} else {
				if err = table.CreateTable(db, t); err != nil {
					logrus.Errorf("FATAL: %+v\n", err)
				}
			}
		}
	} else {
		dfolder := pathutils.Abs("domain")
		if err = os.MkdirAll(dfolder, os.ModePerm); err != nil {
			logrus.Panicln(err)
		}
		for _, t := range existTables {
			var dbIndice []table.DbIndex
			if err = db.Select(&dbIndice, fmt.Sprintf("SHOW INDEXES FROM %s", t)); err != nil {
				logrus.Panicln(err)
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
			colIdxMap := make(map[string][]table.IndexItem)
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
						Unique: !v[0].Non_unique,
						Name:   k,
						Column: idx.Column_name,
						Order:  idx.Seq_in_index,
						Sort:   sort,
					}
					if val, exists := colIdxMap[idx.Column_name]; exists {
						val = append(val, items[i])
						colIdxMap[idx.Column_name] = val
					} else {
						colIdxMap[idx.Column_name] = []table.IndexItem{
							items[i],
						}
					}
				}
				indexes = append(indexes, table.Index{
					Unique: !v[0].Non_unique,
					Name:   k,
					Items:  items,
				})
			}

			var columns []table.DbColumn
			if err = db.Select(&columns, fmt.Sprintf("SHOW FULL COLUMNS FROM %s", t)); err != nil {
				logrus.Panicln(err)
			}

			var cols []table.Column
			var fields []astutils.FieldMeta
			for _, item := range columns {
				extra := item.Extra
				if strings.Contains(extra, "auto_increment") {
					extra = ""
				}
				extra = strings.TrimSpace(strings.TrimPrefix(extra, "DEFAULT_GENERATED"))
				if stringutils.IsNotEmpty(item.Comment) {
					extra += fmt.Sprintf(" comment '%s'", item.Comment)
				}
				extra = strings.TrimSpace(extra)

				col := table.Column{
					Table:         t,
					Name:          item.Field,
					Type:          columnenum.ColumnType(item.Type),
					Default:       item.Default,
					Pk:            table.CheckPk(item.Key),
					Nullable:      table.CheckNull(item.Null),
					Unsigned:      table.CheckUnsigned(item.Type),
					Autoincrement: table.CheckAutoincrement(item.Extra),
					Extra:         extraenum.Extra(extra),
					AutoSet:       table.CheckAutoSet(item.Default),
					Indexes:       colIdxMap[item.Field],
				}
				col.Meta = table.NewFieldFromColumn(col)
				fields = append(fields, col.Meta)
				cols = append(cols, col)
			}

			domain := astutils.StructMeta{
				Name:   strcase.ToCamel(t),
				Fields: fields,
			}

			var pkColumn table.Column
			for _, column := range cols {
				if column.Pk {
					pkColumn = column
					break
				}
			}

			tables = append(tables, table.Table{
				Name:    t,
				Columns: cols,
				Pk:      pkColumn.Name,
				Indexes: indexes,
				Meta:    domain,
			})

			dfile := filepath.Join(dfolder, strings.ToLower(domain.Name)+".go")
			if _, err = os.Stat(dfile); os.IsNotExist(err) {
				if err = codegen.GenDomainGo(dfolder, domain); err != nil {
					logrus.Errorf("FATAL: %+v\n", err)
				}
			} else {
				logrus.Warnf("file %s already exists", dfile)
			}
		}
	}

	if *dao {
		for _, t := range tables {
			if err = codegen.GenDaoGo(*dir, t); err != nil {
				logrus.Errorf("FATAL: %+v\n", err)
				break
			}
			if err = codegen.GenDaoImplGo(*dir, t); err != nil {
				logrus.Errorf("FATAL: %+v\n", err)
				break
			}
			if err = codegen.GenDaoSql(*dir, t); err != nil {
				logrus.Errorf("FATAL: %+v\n", err)
				break
			}
		}
	}

}
