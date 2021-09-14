package table

import (
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

// CreateTable create table from Table
func CreateTable(db *sqlx.DB, t Table) error {
	var (
		statement string
		err       error
	)
	if statement, err = t.CreateSql(); err != nil {
		return err
	}
	logrus.Infoln(statement)
	if _, err = db.Exec(statement); err != nil {
		return err
	}
	return err
}

// ChangeColumn change a column definition by Column
func ChangeColumn(db *sqlx.DB, col Column) error {
	var (
		statement string
		err       error
	)
	if statement, err = col.ChangeColumnSql(); err != nil {
		return err
	}
	logrus.Infoln(statement)
	if _, err = db.Exec(statement); err != nil {
		return err
	}
	return err
}

// AddColumn add a column by Column
func AddColumn(db *sqlx.DB, col Column) error {
	var (
		statement string
		err       error
	)
	if statement, err = col.AddColumnSql(); err != nil {
		return err
	}
	logrus.Infoln(statement)
	if _, err = db.Exec(statement); err != nil {
		return err
	}
	return err
}
