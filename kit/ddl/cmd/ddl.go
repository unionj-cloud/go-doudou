package cmd

import (
	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
	"github.com/unionj-cloud/go-doudou/kit/ddl/table"
)

func CreateTable(db *sqlx.DB, t table.Table) error {
	var (
		statement string
		err       error
	)
	if statement, err = t.CreateSql(); err != nil {
		return err
	}
	log.Infoln(statement)
	if _, err = db.Exec(statement); err != nil {
		return err
	}
	return err
}

func ChangeColumn(db *sqlx.DB, col table.Column) error {
	var (
		statement string
		err       error
	)
	if statement, err = col.ChangeColumnSql(); err != nil {
		return err
	}
	log.Infoln(statement)
	if _, err = db.Exec(statement); err != nil {
		return err
	}
	return err
}

func AddColumn(db *sqlx.DB, col table.Column) error {
	var (
		statement string
		err       error
	)
	if statement, err = col.AddColumnSql(); err != nil {
		return err
	}
	log.Infoln(statement)
	if _, err = db.Exec(statement); err != nil {
		return err
	}
	return err
}
