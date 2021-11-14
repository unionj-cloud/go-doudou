package ddl

import (
	"context"
	"fmt"
	mapset "github.com/deckarep/golang-set"
	"github.com/pkg/errors"
	"github.com/unionj-cloud/go-doudou/ddl/columnenum"
	"github.com/unionj-cloud/go-doudou/ddl/ddlast"
	"github.com/unionj-cloud/go-doudou/ddl/extraenum"
	"github.com/unionj-cloud/go-doudou/ddl/sortenum"
	"github.com/unionj-cloud/go-doudou/sliceutils"
	"github.com/unionj-cloud/go-doudou/stringutils"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"reflect"
	"strings"
	"time"

	// here must import mysql
	_ "github.com/go-sql-driver/mysql"
	"github.com/iancoleman/strcase"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"github.com/unionj-cloud/go-doudou/astutils"
	"github.com/unionj-cloud/go-doudou/ddl/codegen"
	"github.com/unionj-cloud/go-doudou/ddl/config"
	"github.com/unionj-cloud/go-doudou/ddl/table"
	"path/filepath"
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
	if !d.Reverse {
		tables = struct2Table(timeoutCtx, d, existTables, db)
	} else {
		tables = table2struct(timeoutCtx, d, existTables, db)
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

func table2struct(ctx context.Context, d Ddl, existTables []string, db *sqlx.DB) (tables []table.Table) {
	var err error
	if err = os.MkdirAll(d.Dir, os.ModePerm); err != nil {
		panic(fmt.Sprintf("%+v", err))
	}
	for _, t := range existTables {
		if stringutils.IsNotEmpty(d.Pre) && !strings.HasPrefix(t, d.Pre) {
			continue
		}
		var dbIndice []table.DbIndex
		if err = db.SelectContext(ctx, &dbIndice, fmt.Sprintf("SHOW INDEXES FROM %s", t)); err != nil {
			panic(fmt.Sprintf("%+v", err))
		}

		idxMap := make(map[string][]table.DbIndex)

		for _, idx := range dbIndice {
			if val, exists := idxMap[idx.KeyName]; exists {
				val = append(val, idx)
				idxMap[idx.KeyName] = val
			} else {
				idxMap[idx.KeyName] = []table.DbIndex{
					idx,
				}
			}
		}

		indexes, colIdxMap := idxListAndMap(idxMap)

		var columns []table.DbColumn
		if err = db.SelectContext(ctx, &columns, fmt.Sprintf("SHOW FULL COLUMNS FROM %s", t)); err != nil {
			panic(fmt.Sprintf("%+v", err))
		}

		fks := foreignKeys(ctx, db, d.Conf.Schema, t)
		fkMap := make(map[string]table.ForeignKey)
		for _, item := range fks {
			fkMap[item.Fk] = item
		}

		var cols []table.Column
		var fields []astutils.FieldMeta
		for _, item := range columns {
			col := dbColumn2Column(item, colIdxMap, t, fkMap[item.Field])
			fields = append(fields, col.Meta)
			cols = append(cols, col)
		}

		domain := astutils.StructMeta{
			Name:   strcase.ToCamel(strings.TrimPrefix(t, d.Pre)),
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
			Fks:     fks,
		})

		dfile := filepath.Join(d.Dir, strings.ToLower(domain.Name)+".go")
		if _, err = os.Stat(dfile); os.IsNotExist(err) {
			if err = codegen.GenDomainGo(d.Dir, domain); err != nil {
				panic(fmt.Sprintf("%+v", err))
			}
		} else {
			logrus.Warnf("file %s already exists", dfile)
		}
	}
	return
}

func foreignKeys(ctx context.Context, db *sqlx.DB, schema, t string) (fks []table.ForeignKey) {
	var (
		dbForeignKeys []table.DbForeignKey
		err           error
	)
	rawSql := `
		SELECT TABLE_NAME,COLUMN_NAME,CONSTRAINT_NAME, REFERENCED_TABLE_NAME,REFERENCED_COLUMN_NAME
		FROM INFORMATION_SCHEMA.KEY_COLUMN_USAGE
		WHERE TABLE_SCHEMA = ? AND REFERENCED_TABLE_SCHEMA = ? AND TABLE_NAME = ?
	`
	if err = db.SelectContext(ctx, &dbForeignKeys, db.Rebind(rawSql), schema, schema, t); err != nil {
		panic(fmt.Sprintf("%+v", err))
	}
	for _, item := range dbForeignKeys {
		var (
			dbActions []table.DbAction
			dbAction  table.DbAction
		)
		rawSql = `
			select CONSTRAINT_NAME, UPDATE_RULE, DELETE_RULE, TABLE_NAME, REFERENCED_TABLE_NAME 
			from information_schema.REFERENTIAL_CONSTRAINTS 
			where CONSTRAINT_SCHEMA=? and TABLE_NAME=? and CONSTRAINT_NAME=?
		`
		if err = db.SelectContext(ctx, &dbActions, db.Rebind(rawSql), schema, t, item.ConstraintName); err != nil {
			panic(fmt.Sprintf("%+v", err))
		}
		if len(dbActions) > 0 {
			dbAction = dbActions[0]
		}
		fks = append(fks, table.ForeignKey{
			Table:           t,
			Constraint:      item.ConstraintName,
			Fk:              item.ColumnName,
			ReferencedTable: item.ReferencedTableName,
			ReferencedCol:   item.ReferencedColumnName,
			UpdateRule:      dbAction.UpdateRule,
			DeleteRule:      dbAction.DeleteRule,
		})
	}
	return
}

func idxListAndMap(idxMap map[string][]table.DbIndex) ([]table.Index, map[string][]table.IndexItem) {
	var indexes []table.Index
	colIdxMap := make(map[string][]table.IndexItem)
	for k, v := range idxMap {
		if len(v) == 0 {
			continue
		}
		items := make([]table.IndexItem, len(v))
		for i, idx := range v {
			var sor sortenum.Sort
			if idx.Collation == "B" {
				sor = sortenum.Desc
			} else {
				sor = sortenum.Asc
			}
			items[i] = table.IndexItem{
				Unique: !v[0].NonUnique,
				Name:   k,
				Column: idx.ColumnName,
				Order:  idx.SeqInIndex,
				Sort:   sor,
			}
			if val, exists := colIdxMap[idx.ColumnName]; exists {
				val = append(val, items[i])
				colIdxMap[idx.ColumnName] = val
			} else {
				colIdxMap[idx.ColumnName] = []table.IndexItem{
					items[i],
				}
			}
		}
		indexes = append(indexes, table.Index{
			Unique: !v[0].NonUnique,
			Name:   k,
			Items:  items,
		})
	}
	return indexes, colIdxMap
}

func dbColumn2Column(item table.DbColumn, colIdxMap map[string][]table.IndexItem, t string, fk table.ForeignKey) table.Column {
	extra := item.Extra
	if strings.Contains(extra, "auto_increment") {
		extra = ""
	}
	extra = strings.TrimSpace(strings.TrimPrefix(extra, "DEFAULT_GENERATED"))
	if stringutils.IsNotEmpty(item.Comment) {
		extra += fmt.Sprintf(" comment '%s'", item.Comment)
	}
	extra = strings.TrimSpace(extra)
	var defaultVal string
	if item.Default != nil {
		defaultVal = *item.Default
	}
	col := table.Column{
		Table:         t,
		Name:          item.Field,
		Type:          columnenum.ColumnType(item.Type),
		Default:       defaultVal,
		Pk:            table.CheckPk(item.Key),
		Nullable:      table.CheckNull(item.Null),
		Unsigned:      table.CheckUnsigned(item.Type),
		Autoincrement: table.CheckAutoincrement(item.Extra),
		Extra:         extraenum.Extra(extra),
		AutoSet:       table.CheckAutoSet(defaultVal),
		Indexes:       colIdxMap[item.Field],
		Fk:            fk,
	}
	col.Meta = table.NewFieldFromColumn(col)
	return col
}

func struct2Table(ctx context.Context, d Ddl, existTables []string, db *sqlx.DB) (tables []table.Table) {
	var (
		files []string
		err   error
		tx    *sqlx.Tx
		root  *ast.File
	)
	if err = filepath.Walk(d.Dir, astutils.Visit(&files)); err != nil {
		panic(fmt.Sprintf("%+v", err))
	}
	sc := astutils.NewStructCollector(astutils.ExprString)
	for _, file := range files {
		fset := token.NewFileSet()
		if root, err = parser.ParseFile(fset, file, nil, parser.ParseComments); err != nil {
			panic(fmt.Sprintf("%+v", err))
		}
		ast.Walk(sc, root)
	}

	flattened := ddlast.FlatEmbed(sc.Structs)
	for _, sm := range flattened {
		tables = append(tables, table.NewTableFromStruct(sm, d.Pre))
	}

	if tx, err = db.BeginTxx(ctx, nil); err != nil {
		panic(fmt.Sprintf("%+v", err))
	}
	defer func() {
		if r := recover(); r != nil {
			if _err := tx.Rollback(); _err != nil {
				err = errors.Wrap(_err, "")
			}
			panic(fmt.Sprintf("%+v", err))
		}
	}()

	for _, t := range tables {
		if sliceutils.StringContains(existTables, t.Name) {
			var columns []table.DbColumn
			if err = tx.SelectContext(ctx, &columns, fmt.Sprintf("desc %s", t.Name)); err != nil {
				panic(fmt.Sprintf("%+v", err))
			}
			var existColumnNames []interface{}
			for _, dbCol := range columns {
				existColumnNames = append(existColumnNames, dbCol.Field)
			}
			existColSet := mapset.NewSetFromSlice(existColumnNames)

			for _, col := range t.Columns {
				if existColSet.Contains(col.Name) {
					if err = table.ChangeColumn(ctx, tx, col); err != nil {
						panic(fmt.Sprintf("%+v", err))
					}
				} else {
					if err = table.AddColumn(ctx, tx, col); err != nil {
						panic(fmt.Sprintf("%+v", err))
					}
				}
			}
			updateIndexFromStruct(ctx, tx, t)
		} else {
			if err = table.CreateTable(ctx, tx, t); err != nil {
				panic(fmt.Sprintf("%+v", err))
			}
		}
	}
	if err = tx.Commit(); err != nil {
		panic(fmt.Sprintf("%+v", err))
	}
	return
}

func updateIndexFromStruct(ctx context.Context, tx *sqlx.Tx, t table.Table) {
	var dbIndexes []table.DbIndex
	if err := tx.SelectContext(ctx, &dbIndexes, fmt.Sprintf("SHOW INDEXES FROM %s", t.Name)); err != nil {
		panic(fmt.Sprintf("%+v", err))
	}

	keyIndexMap := make(map[string][]table.DbIndex)
	for _, index := range dbIndexes {
		if index.KeyName == "PRIMARY" {
			continue
		}
		if val, exists := keyIndexMap[index.KeyName]; exists {
			val = append(val, index)
			keyIndexMap[index.KeyName] = val
		} else {
			keyIndexMap[index.KeyName] = []table.DbIndex{index}
		}
	}

	for _, index := range t.Indexes {
		if current, exists := keyIndexMap[index.Name]; exists {
			copied := table.NewIndexFromDbIndexes(current)
			if reflect.DeepEqual(index, copied) {
				continue
			}
			index.Table = t.Name
			if err := table.DropAddIndex(ctx, tx, index); err != nil {
				panic(fmt.Sprintf("%+v", err))
			}
		} else {
			index.Table = t.Name
			if err := table.AddIndex(ctx, tx, index); err != nil {
				panic(fmt.Sprintf("%+v", err))
			}
		}
	}

	var idxKeys []string
	for _, index := range t.Indexes {
		idxKeys = append(idxKeys, index.Name)
	}
	for k, v := range keyIndexMap {
		if !sliceutils.StringContains(idxKeys, k) {
			index := table.NewIndexFromDbIndexes(v)
			index.Table = t.Name
			if err := table.DropIndex(ctx, tx, index); err != nil {
				panic(fmt.Sprintf("%+v", err))
			}
		}
	}
}
