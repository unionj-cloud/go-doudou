package table

import (
	"context"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/unionj-cloud/go-doudou/ddl/wrapper"
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
	logrus.Infoln(statement)
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
	logrus.Infoln(statement)
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
	logrus.Infoln(statement)
	if _, err = db.ExecContext(ctx, statement); err != nil {
		return err
	}
	return err
}

// DropAddIndex drop and then add an existing index with the same key_name
func DropAddIndex(ctx context.Context, db wrapper.Querier, idx Index) error {
	var err error
	if err = DropIndex(ctx, db, idx); err != nil {
		return errors.Wrap(err, "")
	}
	if err = AddIndex(ctx, db, idx); err != nil {
		return errors.Wrap(err, "")
	}
	return nil
}

// AddIndex add an new index
func AddIndex(ctx context.Context, db wrapper.Querier, idx Index) error {
	var (
		statement string
		err       error
	)
	if statement, err = idx.AddIndexSql(); err != nil {
		return errors.Wrap(err, "")
	}
	logrus.Infoln(statement)
	if _, err = db.ExecContext(ctx, statement); err != nil {
		return errors.Wrap(err, "")
	}
	return nil
}

// DropIndex drop an existing index
func DropIndex(ctx context.Context, db wrapper.Querier, idx Index) error {
	var (
		statement string
		err       error
	)
	if statement, err = idx.DropIndexSql(); err != nil {
		return errors.Wrap(err, "")
	}
	logrus.Infoln(statement)
	if _, err = db.ExecContext(ctx, statement); err != nil {
		return errors.Wrap(err, "")
	}
	return nil
}
