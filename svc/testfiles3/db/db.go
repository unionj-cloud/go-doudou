package db

import (
	"testfiles3/config"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/iancoleman/strcase"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

func NewDb(conf config.DbConfig) (*sqlx.DB, error) {
	conn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s",
		conf.User,
		conf.Passwd,
		conf.Host,
		conf.Port,
		conf.Schema,
		conf.Charset)
	conn += "&loc=Asia%2FShanghai&parseTime=True"

	db, err := sqlx.Connect(conf.Driver, conn)
	if err != nil {
		return nil, errors.Wrap(err, "database connection failed")
	}
	db.MapperFunc(strcase.ToSnake)
	return db, nil
}
