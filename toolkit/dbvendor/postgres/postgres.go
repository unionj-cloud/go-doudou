package postgres

import (
	"context"
	"database/sql"
	"github.com/unionj-cloud/go-doudou/v2/toolkit/dbvendor"
	"fmt"
	"github.com/pkg/errors"
	"github.com/unionj-cloud/go-doudou/v2/framework/database"
	"gorm.io/gorm"
)

func init() {
	dbvendor.Registry.Register(database.DriverPostgres, &Vendor{})
}

var _ dbvendor.IVendor = (*Vendor)(nil)

type Vendor struct {
}

func (v *Vendor) ToColumnType(goType string, autoincrementing bool) string {
	switch goType {
	case "int", "int16", "int32":
		if autoincrementing {
			return SerialType
		}
		return IntType
	case "int64":
		if autoincrementing {
			return BigSerialType
		}
		return BigintType
	case "float32":
		return FloatType
	case "float64":
		return DoubleType
	case "string":
		return VarcharType
	case "text":
		return LongtextType
	case "bool", "int8":
		return TinyintType
	case "time.Time":
		return DatetimeType
	case "decimal.Decimal":
		return "decimal(6,2)"
	case "types.JSONText":
		return JSONType
	}
	panic(fmt.Sprintf("no available type %s", goType))
}

func (v *Vendor) CreateTable(ctx context.Context, db *gorm.DB, t dbvendor.Table) error {
	var (
		statement string
		err       error
	)
	if statement, err = dbvendor.String(createTable, createTable, t, nil); err != nil {
		return err
	}
	return db.WithContext(ctx).Exec(statement).Error
}

func (v *Vendor) DropTable(ctx context.Context, db *gorm.DB, t dbvendor.Table) error {
	var (
		statement string
		err       error
	)
	if statement, err = dbvendor.String(dropTable, dropTable, t, nil); err != nil {
		return err
	}
	return db.WithContext(ctx).Exec(statement).Error
}

func (v *Vendor) ChangeColumn(ctx context.Context, db *gorm.DB, col dbvendor.Column) error {
	var (
		statement string
		err       error
	)
	if statement, err = dbvendor.StringBlock(alterTable, alterTable, "change", col, nil); err != nil {
		return err
	}
	return db.WithContext(ctx).Exec(statement).Error
}

func (v *Vendor) AddColumn(ctx context.Context, db *gorm.DB, col dbvendor.Column) error {
	var (
		statement string
		err       error
	)
	if statement, err = dbvendor.StringBlock(alterTable, alterTable, "add", col, nil); err != nil {
		return err
	}
	return db.WithContext(ctx).Exec(statement).Error
}

func (v *Vendor) DropColumn(ctx context.Context, db *gorm.DB, col dbvendor.Column) error {
	var (
		statement string
		err       error
	)
	if statement, err = dbvendor.StringBlock(alterTable, alterTable, "drop", col, nil); err != nil {
		return err
	}
	return db.WithContext(ctx).Exec(statement).Error
}

func (v *Vendor) Insert(ctx context.Context, db *gorm.DB, dml dbvendor.DMLSchema, args ...interface{}) (int64, error) {
	var (
		statement string
		err       error
	)
	if statement, err = v.GetInsertStatement(dml); err != nil {
		return 0, errors.WithStack(err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		return 0, errors.WithStack(err)
	}
	var id int64
	if err = sqlDB.QueryRowContext(ctx, statement, args...).Scan(&id); err != nil {
		return 0, errors.WithStack(err)
	}
	return id, nil
}

func (v *Vendor) GetInsertStatement(dml dbvendor.DMLSchema) (statement string, err error) {
	if statement, err = dbvendor.String(insertInto, insertInto, dml, dbvendor.Dollar); err != nil {
		return "", errors.WithStack(err)
	}
	return statement, nil
}

func (v *Vendor) GetUpdateStatement(dml dbvendor.DMLSchema) (statement string, err error) {
	if statement, err = dbvendor.String(updateTable, updateTable, dml, dbvendor.Dollar); err != nil {
		return "", errors.WithStack(err)
	}
	return statement, nil
}

func (v *Vendor) Update(ctx context.Context, db *gorm.DB, dml dbvendor.DMLSchema, args ...interface{}) error {
	var (
		statement string
		err       error
	)
	if statement, err = v.GetUpdateStatement(dml); err != nil {
		return errors.WithStack(err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		return errors.WithStack(err)
	}
	_, err = sqlDB.ExecContext(ctx, statement, args...)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (v *Vendor) Delete(ctx context.Context, db *gorm.DB, dml dbvendor.DMLSchema, args ...interface{}) error {
	var (
		statement string
		err       error
	)
	if statement, err = dbvendor.String(deleteFrom, deleteFrom, dml, dbvendor.Dollar); err != nil {
		return errors.WithStack(err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		return errors.WithStack(err)
	}
	_, err = sqlDB.ExecContext(ctx, statement, args...)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (v *Vendor) SelectById(ctx context.Context, db *gorm.DB, dml dbvendor.DMLSchema, args ...interface{}) (map[string]interface{}, error) {
	var (
		statement string
		err       error
		rows      *sql.Rows
	)
	if statement, err = dbvendor.String(selectFromById, selectFromById, dml, dbvendor.Dollar); err != nil {
		return nil, errors.WithStack(err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if rows, err = sqlDB.QueryContext(ctx, statement, args...); err != nil {
		return nil, errors.WithStack(err)
	}
	result := make([]map[string]interface{}, 0)
	dbvendor.Scan(rows, &result)
	if len(result) == 0 {
		return nil, errors.WithStack(sql.ErrNoRows)
	}
	return result[0], nil
}
