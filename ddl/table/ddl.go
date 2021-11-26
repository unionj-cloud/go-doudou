package table

import (
	"context"
	"fmt"
	mapset "github.com/deckarep/golang-set"
	"github.com/go-git/go-git/v5/utils/merkletrie/index"
	"github.com/iancoleman/strcase"
	"github.com/jmoiron/sqlx"
	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/unionj-cloud/go-doudou/astutils"
	"github.com/unionj-cloud/go-doudou/ddl/columnenum"
	"github.com/unionj-cloud/go-doudou/ddl/config"
	"github.com/unionj-cloud/go-doudou/ddl/ddlast"
	"github.com/unionj-cloud/go-doudou/ddl/extraenum"
	"github.com/unionj-cloud/go-doudou/ddl/sortenum"
	"github.com/unionj-cloud/go-doudou/ddl/wrapper"
	"github.com/unionj-cloud/go-doudou/pathutils"
	"github.com/unionj-cloud/go-doudou/sliceutils"
	"github.com/unionj-cloud/go-doudou/stringutils"
	"github.com/unionj-cloud/go-doudou/test"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"reflect"
	"strings"
)

// CreateTable create table from Table
func CreateTable(ctx context.Context, db wrapper.Querier, t Table) error {
	var (
		statement string
		err       error
	)
	if statement, err = t.CreateSql(); err != nil {
		return err
	}
	fmt.Println(statement)
	if _, err = db.ExecContext(ctx, statement); err != nil {
		return err
	}
	return err
}

// ChangeColumn change a column definition by Column
func ChangeColumn(ctx context.Context, db wrapper.Querier, col Column) error {
	var (
		statement string
		err       error
	)
	if statement, err = col.ChangeColumnSql(); err != nil {
		return err
	}
	fmt.Println(statement)
	if _, err = db.ExecContext(ctx, statement); err != nil {
		return err
	}
	return err
}

// AddColumn add a column by Column
func AddColumn(ctx context.Context, db wrapper.Querier, col Column) error {
	var (
		statement string
		err       error
	)
	if statement, err = col.AddColumnSql(); err != nil {
		return err
	}
	fmt.Println(statement)
	if _, err = db.ExecContext(ctx, statement); err != nil {
		return err
	}
	return err
}

// dropAddIndex drop and then add an existing index with the same key_name
func dropAddIndex(ctx context.Context, db wrapper.Querier, idx Index) error {
	var err error
	if err = dropIndex(ctx, db, idx); err != nil {
		return errors.Wrap(err, "")
	}
	if err = addIndex(ctx, db, idx); err != nil {
		return errors.Wrap(err, "")
	}
	return nil
}

// addIndex add a new index
func addIndex(ctx context.Context, db wrapper.Querier, idx Index) error {
	var (
		statement string
		err       error
	)
	if statement, err = idx.AddIndexSql(); err != nil {
		return errors.Wrap(err, "")
	}
	fmt.Println(statement)
	if _, err = db.ExecContext(ctx, statement); err != nil {
		return errors.Wrap(err, "")
	}
	return nil
}

// dropIndex drop an existing index
func dropIndex(ctx context.Context, db wrapper.Querier, idx Index) error {
	var (
		statement string
		err       error
	)
	if statement, err = idx.DropIndexSql(); err != nil {
		return errors.Wrap(err, "")
	}
	fmt.Println(statement)
	if _, err = db.ExecContext(ctx, statement); err != nil {
		return errors.Wrap(err, "")
	}
	return nil
}

// dropAddFk drop and then add an existing foreign key with the same constraint
func dropAddFk(ctx context.Context, db wrapper.Querier, idx Index) error {
	var err error
	if err = dropIndex(ctx, db, idx); err != nil {
		return errors.Wrap(err, "")
	}
	if err = addIndex(ctx, db, idx); err != nil {
		return errors.Wrap(err, "")
	}
	return nil
}

// addFk add a new foreign key
func addFk(ctx context.Context, db wrapper.Querier, fk ForeignKey) error {
	var (
		statement string
		err       error
	)
	if statement, err = fk.AddFkSql(); err != nil {
		return errors.Wrap(err, "")
	}
	fmt.Println(statement)
	if _, err = db.ExecContext(ctx, statement); err != nil {
		return errors.Wrap(err, "")
	}
	return nil
}

// dropFk drop an existing foreign key
func dropFk(ctx context.Context, db wrapper.Querier, idx Index) error {
	var (
		statement string
		err       error
	)
	if statement, err = idx.DropIndexSql(); err != nil {
		return errors.Wrap(err, "")
	}
	fmt.Println(statement)
	if _, err = db.ExecContext(ctx, statement); err != nil {
		return errors.Wrap(err, "")
	}
	return nil
}

func Table2struct(ctx context.Context, pre, schema string, existTables []string, db *sqlx.DB) (tables []Table) {
	var err error
	for _, t := range existTables {
		if stringutils.IsNotEmpty(pre) && !strings.HasPrefix(t, pre) {
			continue
		}
		var dbIndice []DbIndex
		if err = db.SelectContext(ctx, &dbIndice, fmt.Sprintf("SHOW INDEXES FROM %s", t)); err != nil {
			panic(fmt.Sprintf("%+v", err))
		}

		idxMap := make(map[string][]DbIndex)

		for _, idx := range dbIndice {
			if val, exists := idxMap[idx.KeyName]; exists {
				val = append(val, idx)
				idxMap[idx.KeyName] = val
			} else {
				idxMap[idx.KeyName] = []DbIndex{
					idx,
				}
			}
		}

		indexes, colIdxMap := idxListAndMap(idxMap)

		var columns []DbColumn
		if err = db.SelectContext(ctx, &columns, fmt.Sprintf("SHOW FULL COLUMNS FROM %s", t)); err != nil {
			panic(fmt.Sprintf("%+v", err))
		}

		fks := foreignKeys(ctx, db, schema, t)
		fkMap := make(map[string]ForeignKey)
		for _, item := range fks {
			fkMap[item.Fk] = item
		}

		var cols []Column
		var fields []astutils.FieldMeta
		for _, item := range columns {
			col := dbColumn2Column(item, colIdxMap, t, fkMap[item.Field])
			fields = append(fields, col.Meta)
			cols = append(cols, col)
		}

		domain := astutils.StructMeta{
			Name:   strcase.ToCamel(strings.TrimPrefix(t, pre)),
			Fields: fields,
		}

		var pkColumn Column
		for _, column := range cols {
			if column.Pk {
				pkColumn = column
				break
			}
		}

		tables = append(tables, Table{
			Name:    t,
			Columns: cols,
			Pk:      pkColumn.Name,
			Indexes: indexes,
			Meta:    domain,
			Fks:     fks,
		})
	}
	return
}

func foreignKeys(ctx context.Context, db wrapper.Querier, schema, t string) (fks []ForeignKey) {
	var (
		dbForeignKeys []DbForeignKey
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
			dbActions []DbAction
			dbAction  DbAction
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
		var rules []string
		if stringutils.IsNotEmpty(dbAction.DeleteRule) {
			rules = append(rules, fmt.Sprintf("ON DELETE %s", dbAction.DeleteRule))
		}
		if stringutils.IsNotEmpty(dbAction.UpdateRule) {
			rules = append(rules, fmt.Sprintf("ON UPDATE %s", dbAction.UpdateRule))
		}
		var fullRule string
		if len(rules) > 0 {
			fullRule = strings.Join(rules, " ")
		}
		fks = append(fks, ForeignKey{
			Table:           t,
			Constraint:      item.ConstraintName,
			Fk:              item.ColumnName,
			ReferencedTable: item.ReferencedTableName,
			ReferencedCol:   item.ReferencedColumnName,
			UpdateRule:      dbAction.UpdateRule,
			DeleteRule:      dbAction.DeleteRule,
			FullRule:        fullRule,
		})
	}
	return
}

func idxListAndMap(idxMap map[string][]DbIndex) ([]Index, map[string][]IndexItem) {
	var indexes []Index
	colIdxMap := make(map[string][]IndexItem)
	for k, v := range idxMap {
		if len(v) == 0 {
			continue
		}
		items := make([]IndexItem, len(v))
		for i, idx := range v {
			var sor sortenum.Sort
			if idx.Collation == "B" {
				sor = sortenum.Desc
			} else {
				sor = sortenum.Asc
			}
			items[i] = IndexItem{
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
				colIdxMap[idx.ColumnName] = []IndexItem{
					items[i],
				}
			}
		}
		indexes = append(indexes, Index{
			Unique: !v[0].NonUnique,
			Name:   k,
			Items:  items,
		})
	}
	return indexes, colIdxMap
}

func dbColumn2Column(item DbColumn, colIdxMap map[string][]IndexItem, t string, fk ForeignKey) Column {
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
	col := Column{
		Table:         t,
		Name:          item.Field,
		Type:          columnenum.ColumnType(item.Type),
		Default:       defaultVal,
		Pk:            CheckPk(item.Key),
		Nullable:      CheckNull(item.Null),
		Unsigned:      CheckUnsigned(item.Type),
		Autoincrement: CheckAutoincrement(item.Extra),
		Extra:         extraenum.Extra(extra),
		AutoSet:       CheckAutoSet(defaultVal),
		Indexes:       colIdxMap[item.Field],
		Fk:            fk,
	}
	col.Meta = NewFieldFromColumn(col)
	return col
}

func Struct2Table(ctx context.Context, dir, pre string, existTables []string, db *sqlx.DB, schema string) (tables []Table) {
	var (
		files []string
		err   error
		tx    *sqlx.Tx
		root  *ast.File
	)
	if err = filepath.Walk(dir, astutils.Visit(&files)); err != nil {
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
		tables = append(tables, NewTableFromStruct(sm, pre))
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
			var columns []DbColumn
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
					if err = ChangeColumn(ctx, tx, col); err != nil {
						panic(fmt.Sprintf("%+v", err))
					}
				} else {
					if err = AddColumn(ctx, tx, col); err != nil {
						panic(fmt.Sprintf("%+v", err))
					}
				}
			}
			updateIndexFromStruct(ctx, tx, t)
			updateFkFromStruct(ctx, tx, t, schema)
		} else {
			if err = CreateTable(ctx, tx, t); err != nil {
				panic(fmt.Sprintf("%+v", err))
			}
		}
	}
	_ = tx.Commit()
	return
}

func updateFkFromStruct(ctx context.Context, tx *sqlx.Tx, t Table, schema string) {
	fks := foreignKeys(ctx, tx, schema, t.Name)
	fkMap := make(map[string]ForeignKey)
	for _, fk := range fks {
		fkMap[fk.Constraint] = fk
	}
	for _, fk := range t.Fks {
		if current, exists := fkMap[fk.Constraint]; exists {
			current.DeleteRule = ""
			current.UpdateRule = ""
			fk.DeleteRule = ""
			fk.UpdateRule = ""
			if reflect.DeepEqual(fk, current) {
				continue
			}
			index.Table = t.Name
			if err := dropAddFk(ctx, tx, index); err != nil {
				panic(fmt.Sprintf("%+v", err))
			}
		} else {
			if err := addFk(ctx, tx, fk); err != nil {
				panic(fmt.Sprintf("%+v", err))
			}
		}
	}

}

func updateIndexFromStruct(ctx context.Context, tx *sqlx.Tx, t Table) {
	var dbIndexes []DbIndex
	if err := tx.SelectContext(ctx, &dbIndexes, fmt.Sprintf("SHOW INDEXES FROM %s", t.Name)); err != nil {
		panic(fmt.Sprintf("%+v", err))
	}

	keyIndexMap := make(map[string][]DbIndex)
	for _, index := range dbIndexes {
		if index.KeyName == "PRIMARY" {
			continue
		}
		if val, exists := keyIndexMap[index.KeyName]; exists {
			val = append(val, index)
			keyIndexMap[index.KeyName] = val
		} else {
			keyIndexMap[index.KeyName] = []DbIndex{index}
		}
	}

	for _, index := range t.Indexes {
		if current, exists := keyIndexMap[index.Name]; exists {
			copied := NewIndexFromDbIndexes(current)
			if reflect.DeepEqual(index, copied) {
				continue
			}
			index.Table = t.Name
			if err := dropAddIndex(ctx, tx, index); err != nil {
				panic(fmt.Sprintf("%+v", err))
			}
		} else {
			index.Table = t.Name
			if err := addIndex(ctx, tx, index); err != nil {
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
			index := NewIndexFromDbIndexes(v)
			index.Table = t.Name
			if err := dropIndex(ctx, tx, index); err != nil {
				panic(fmt.Sprintf("%+v", err))
			}
		}
	}
}

func Setup() (func(), *sqlx.DB, error) {
	logger := logrus.New()
	var terminateContainer func() // variable to store function to terminate container
	var host string
	var port int
	var err error
	terminateContainer, host, port, err = test.SetupMySQLContainer(logger, pathutils.Abs("../../test/sql"), "")
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to setup MySQL container")
	}
	os.Setenv("DB_HOST", host)
	os.Setenv("DB_PORT", fmt.Sprint(port))
	os.Setenv("DB_USER", "root")
	os.Setenv("DB_PASSWD", "1234")
	os.Setenv("DB_SCHEMA", "test")
	os.Setenv("DB_CHARSET", "utf8mb4")
	var conf config.DbConfig
	err = envconfig.Process("db", &conf)
	if err != nil {
		return nil, nil, errors.Wrap(err, "Error processing env")
	}
	conn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s",
		conf.User,
		conf.Passwd,
		conf.Host,
		conf.Port,
		conf.Schema,
		conf.Charset)
	conn += `&loc=Asia%2FShanghai&parseTime=True`
	var db *sqlx.DB
	db, err = sqlx.Connect("mysql", conn)
	if err != nil {
		return nil, nil, errors.Wrap(err, "")
	}
	db.MapperFunc(strcase.ToSnake)
	db = db.Unsafe()
	return terminateContainer, db, nil
}
