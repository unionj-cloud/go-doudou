package table

import (
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

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
