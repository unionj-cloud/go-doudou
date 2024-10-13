package table

import (
	"context"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"time"

	mapset "github.com/deckarep/golang-set"
	"github.com/iancoleman/strcase"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"github.com/unionj-cloud/go-doudou/v2/cmd/internal/ddl/columnenum"
	"github.com/unionj-cloud/go-doudou/v2/cmd/internal/ddl/config"
	"github.com/unionj-cloud/go-doudou/v2/cmd/internal/ddl/ddlast"
	"github.com/unionj-cloud/go-doudou/v2/cmd/internal/ddl/extraenum"
	"github.com/unionj-cloud/go-doudou/v2/cmd/internal/ddl/sortenum"
	"github.com/unionj-cloud/toolkit/astutils"
	"github.com/unionj-cloud/toolkit/caller"
	"github.com/unionj-cloud/toolkit/envconfig"
	"github.com/unionj-cloud/toolkit/pathutils"
	"github.com/unionj-cloud/toolkit/sliceutils"
	"github.com/unionj-cloud/toolkit/sqlext/wrapper"
	"github.com/unionj-cloud/toolkit/stringutils"
	"github.com/unionj-cloud/toolkit/zlogger"
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
		return errors.Wrap(err, caller.NewCaller().String())
	}
	if err = addIndex(ctx, db, idx); err != nil {
		return errors.Wrap(err, caller.NewCaller().String())
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
		return errors.Wrap(err, caller.NewCaller().String())
	}
	fmt.Println(statement)
	if _, err = db.ExecContext(ctx, statement); err != nil {
		return errors.Wrap(err, caller.NewCaller().String())
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
		return errors.Wrap(err, caller.NewCaller().String())
	}
	fmt.Println(statement)
	if _, err = db.ExecContext(ctx, statement); err != nil {
		return errors.Wrap(err, caller.NewCaller().String())
	}
	return nil
}

// dropAddFk drop and then add an existing foreign key with the same constraint
func dropAddFk(ctx context.Context, db wrapper.Querier, fk ForeignKey) error {
	var err error
	if err = dropFk(ctx, db, fk); err != nil {
		return errors.Wrap(err, caller.NewCaller().String())
	}
	if err = addFk(ctx, db, fk); err != nil {
		return errors.Wrap(err, caller.NewCaller().String())
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
		return errors.Wrap(err, caller.NewCaller().String())
	}
	fmt.Println(statement)
	if _, err = db.ExecContext(ctx, statement); err != nil {
		return errors.Wrap(err, caller.NewCaller().String())
	}
	return nil
}

// dropFk drop an existing foreign key
func dropFk(ctx context.Context, db wrapper.Querier, fk ForeignKey) error {
	var (
		statement string
		err       error
	)
	if statement, err = fk.DropFkSql(); err != nil {
		return errors.Wrap(err, caller.NewCaller().String())
	}
	fmt.Println(statement)
	if _, err = db.ExecContext(ctx, statement); err != nil {
		return errors.Wrap(err, caller.NewCaller().String())
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
			panic(errors.Wrap(err, caller.NewCaller().String()))
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
			panic(errors.Wrap(err, caller.NewCaller().String()))
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

		entity := astutils.StructMeta{
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
			Meta:    entity,
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
		panic(errors.Wrap(err, caller.NewCaller().String()))
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
			panic(errors.Wrap(err, caller.NewCaller().String()))
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
		panic(errors.Wrap(err, caller.NewCaller().String()))
	}
	sc := astutils.NewStructCollector(astutils.ExprString)
	for _, file := range files {
		fset := token.NewFileSet()
		if root, err = parser.ParseFile(fset, file, nil, parser.ParseComments); err != nil {
			panic(errors.Wrap(err, caller.NewCaller().String()))
		}
		ast.Walk(sc, root)
	}

	flattened := ddlast.FlatEmbed(sc.Structs)
	for _, sm := range flattened {
		tables = append(tables, NewTableFromStruct(sm, pre))
	}

	if tx, err = db.BeginTxx(ctx, nil); err != nil {
		panic(errors.Wrap(err, caller.NewCaller().String()))
	}
	defer func() {
		if r := recover(); r != nil {
			if _err := tx.Rollback(); _err != nil {
				err = errors.Wrap(_err, "")
			}
			panic(errors.Wrap(err, caller.NewCaller().String()))
		}
	}()

	if _, err = tx.ExecContext(ctx, `SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0;`); err != nil {
		panic(errors.Wrap(err, caller.NewCaller().String()))
	}

	for _, t := range tables {
		if sliceutils.StringContains(existTables, t.Name) {
			var columns []DbColumn
			if err = tx.SelectContext(ctx, &columns, fmt.Sprintf("desc %s", t.Name)); err != nil {
				panic(errors.Wrap(err, caller.NewCaller().String()))
			}
			var existColumnNames []interface{}
			for _, dbCol := range columns {
				existColumnNames = append(existColumnNames, dbCol.Field)
			}
			existColSet := mapset.NewSetFromSlice(existColumnNames)

			for _, col := range t.Columns {
				if existColSet.Contains(col.Name) {
					if err = ChangeColumn(ctx, tx, col); err != nil {
						panic(errors.Wrap(err, caller.NewCaller().String()))
					}
				} else {
					if err = AddColumn(ctx, tx, col); err != nil {
						panic(errors.Wrap(err, caller.NewCaller().String()))
					}
				}
			}
			fks := foreignKeys(ctx, tx, schema, t.Name)
			updateIndexFromStruct(ctx, tx, t, fks)
			updateFkFromStruct(ctx, tx, t, fks)
		} else {
			if err = CreateTable(ctx, tx, t); err != nil {
				panic(errors.Wrap(err, caller.NewCaller().String()))
			}
		}
	}

	if _, err = tx.ExecContext(ctx, `SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS;`); err != nil {
		panic(errors.Wrap(err, caller.NewCaller().String()))
	}
	_ = tx.Commit()
	return
}

func updateFkFromStruct(ctx context.Context, tx *sqlx.Tx, t Table, fks []ForeignKey) {
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
			if err := dropAddFk(ctx, tx, fk); err != nil {
				panic(errors.Wrap(err, caller.NewCaller().String()))
			}
		} else {
			if err := addFk(ctx, tx, fk); err != nil {
				panic(errors.Wrap(err, caller.NewCaller().String()))
			}
		}
	}

	var constraints []string
	for _, fk := range t.Fks {
		constraints = append(constraints, fk.Constraint)
	}
	for k, v := range fkMap {
		if !sliceutils.StringContains(constraints, k) {
			if err := dropFk(ctx, tx, v); err != nil {
				panic(errors.Wrap(err, caller.NewCaller().String()))
			}
		}
	}
}

func updateIndexFromStruct(ctx context.Context, tx *sqlx.Tx, t Table, fks []ForeignKey) {
	var dbIndexes []DbIndex
	if err := tx.SelectContext(ctx, &dbIndexes, fmt.Sprintf("SHOW INDEXES FROM %s", t.Name)); err != nil {
		panic(errors.Wrap(err, caller.NewCaller().String()))
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

	for _, idx := range t.Indexes {
		if current, exists := keyIndexMap[idx.Name]; exists {
			copied := NewIndexFromDbIndexes(current)
			if reflect.DeepEqual(idx, copied) {
				continue
			}
			idx.Table = t.Name
			if err := dropAddIndex(ctx, tx, idx); err != nil {
				panic(errors.Wrap(err, caller.NewCaller().String()))
			}
		} else {
			idx.Table = t.Name
			if err := addIndex(ctx, tx, idx); err != nil {
				panic(errors.Wrap(err, caller.NewCaller().String()))
			}
		}
	}

	var idxKeys []string
	for _, idx := range t.Indexes {
		idxKeys = append(idxKeys, idx.Name)
	}
	for k, v := range keyIndexMap {
		if !sliceutils.StringContains(idxKeys, k) {
			shouldDrop := true
			if len(v) == 1 {
				idx := v[0]
				for _, fk := range fks {
					if fk.Table == idx.Table && fk.Fk == idx.ColumnName {
						shouldDrop = false
						break
					}
				}
			}
			if shouldDrop {
				idx := NewIndexFromDbIndexes(v)
				idx.Table = t.Name
				if err := dropIndex(ctx, tx, idx); err != nil {
					panic(errors.Wrap(err, caller.NewCaller().String()))
				}
			}
		}
	}
}

func Setup() (func(), *sqlx.DB, error) {
	var terminateContainer func() // variable to store function to terminate container
	var host string
	var port int
	var err error
	terminateContainer, host, port, err = setupMySQLContainer(zlogger.Logger, pathutils.Abs("../testdata/sql"), "")
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
		return nil, nil, errors.Wrap(err, "[go-doudou] Error processing env")
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
		return nil, nil, errors.Wrap(err, caller.NewCaller().String())
	}
	db.MapperFunc(strcase.ToSnake)
	db = db.Unsafe()
	return terminateContainer, db, nil
}

func setupMySQLContainer(logger zerolog.Logger, initdb string, dbname string) (func(), string, int, error) {
	logger.Info().Msg("setup MySQL Container")
	ctx := context.Background()
	if stringutils.IsEmpty(dbname) {
		dbname = "test"
	}
	req := testcontainers.ContainerRequest{
		Image:        "mysql:latest",
		ExposedPorts: []string{"3306/tcp", "33060/tcp"},
		Env: map[string]string{
			"MYSQL_ROOT_PASSWORD": "1234",
			"MYSQL_DATABASE":      dbname,
		},
		Mounts:     make(testcontainers.ContainerMounts, 0),
		WaitingFor: wait.ForLog("port: 3306  MySQL Community Server - GPL").WithStartupTimeout(60 * time.Second),
	}

	// TODO
	//req.Mounts = append(req.Mounts, testcontainers.ContainerMount{
	//	Source:   nil,
	//	Target:   "/docker-entrypoint-initdb.d",
	//	ReadOnly: false,
	//})

	mysqlC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})

	if err != nil {
		logger.Error().Msgf("error starting mysql container: %s", err)
		panic(fmt.Sprintf("%v", err))
	}

	closeContainer := func() {
		logger.Info().Msg("terminating container")
		err := mysqlC.Terminate(ctx)
		if err != nil {
			logger.Error().Msgf("error terminating mysql container: %s", err)
			panic(fmt.Sprintf("%v", err))
		}
	}

	host, _ := mysqlC.Host(ctx)
	p, _ := mysqlC.MappedPort(ctx, "3306/tcp")
	port := p.Int()

	return closeContainer, host, port, nil
}
