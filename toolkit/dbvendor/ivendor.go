package dbvendor

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/pkg/errors"
	"github.com/unionj-cloud/go-doudou/v2/toolkit/errorx"
	"github.com/unionj-cloud/go-doudou/v2/toolkit/templateutils"
	"github.com/unionj-cloud/go-doudou/v2/toolkit/zlogger"
	"gorm.io/gorm"
)

var Registry = &registry{
	vendors: map[string]IVendor{},
}

type registry struct {
	vendors map[string]IVendor
}

func (receiver *registry) Register(driver string, vendor IVendor) {
	receiver.vendors[driver] = vendor
}

func (receiver *registry) GetVendor(driver string) IVendor {
	vendor, ok := receiver.vendors[driver]
	if !ok {
		errorx.Panic(fmt.Sprintf("Unsupported driver %s", driver))
	}
	return vendor
}

type IVendor interface {
	CreateTable(ctx context.Context, db *gorm.DB, t Table) error
	DropTable(ctx context.Context, db *gorm.DB, t Table) error
	ChangeColumn(ctx context.Context, db *gorm.DB, col Column) error
	AddColumn(ctx context.Context, db *gorm.DB, col Column) error
	DropColumn(ctx context.Context, db *gorm.DB, col Column) error
	ToColumnType(goType string, autoincrementing bool) string

	Insert(ctx context.Context, db *gorm.DB, dml DMLSchema, args ...interface{}) (int64, error)
	Update(ctx context.Context, db *gorm.DB, dml DMLSchema, args ...interface{}) error
	Delete(ctx context.Context, db *gorm.DB, dml DMLSchema, args ...interface{}) error
	SelectById(ctx context.Context, db *gorm.DB, dml DMLSchema, args ...interface{}) (map[string]interface{}, error)
	GetInsertStatement(dml DMLSchema) (statement string, err error)
	GetBatchInsertStatement(dml DMLSchema, rows []interface{}) (statement string, err error)
	GetUpdateStatement(dml DMLSchema) (statement string, err error)
}

type DMLSchema struct {
	Schema        string
	TablePrefix   string
	TableName     string
	InsertColumns []Column
	UpdateColumns []Column
	Pk            Column
}

// Column define a column
type Column struct {
	TablePrefix   string
	Table         string
	Name          string
	OldName       string
	Type          string
	Default       *string
	Pk            bool
	Nullable      bool
	Unsigned      bool
	Autoincrement bool
	Extra         string
	Comment       string
	// 关联表名
	Foreign string
}

// Table defines a table
type Table struct {
	TablePrefix string
	Name        string
	Columns     []Column
	BizColumns  []Column
	Pk          string
	Joins       []string
	// 父表
	Inherited string
}

func String(tmplname, tmpl string, data interface{}, pf PlaceholderFormat) (string, error) {
	result, err := templateutils.String(tmplname, tmpl, data)
	if err != nil {
		return "", errors.WithStack(err)
	}
	if pf != nil {
		result, err = pf.ReplacePlaceholders(result)
		if err != nil {
			return "", errors.WithStack(err)
		}
	}
	zlogger.Debug().Msg(result)
	return result, nil
}

func StringBlock(tmplname, tmpl string, block string, data interface{}, pf PlaceholderFormat) (string, error) {
	result, err := templateutils.StringBlock(tmplname, tmpl, block, data)
	if err != nil {
		return "", errors.WithStack(err)
	}
	if pf != nil {
		result, err = pf.ReplacePlaceholders(result)
		if err != nil {
			return "", errors.WithStack(err)
		}
	}
	zlogger.Debug().Msg(result)
	return result, nil
}

const (
	// Update used for update_at column
	Update = "on update CURRENT_TIMESTAMP"
)

func Scan(rows *sql.Rows, result *[]map[string]interface{}) {
	fields, _ := rows.Columns()
	for rows.Next() {
		scans := make([]interface{}, len(fields))
		data := make(map[string]interface{})

		for i := range scans {
			scans[i] = &scans[i]
		}
		rows.Scan(scans...)
		for i, v := range scans {
			data[fields[i]] = v
		}
		*result = append(*result, data)
	}
}
