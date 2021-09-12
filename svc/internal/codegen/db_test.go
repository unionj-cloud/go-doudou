package codegen

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestGenDb(t *testing.T) {
	dir := testDir + "db"
	InitSvc(dir)
	defer os.RemoveAll(dir)
	GenDb(dir)
	expect := `package db

import (
	"testdatadb/config"
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
`
	configfile := dir + "/db/db.go"
	f, err := os.Open(configfile)
	if err != nil {
		t.Fatal(err)
	}
	content, err := ioutil.ReadAll(f)
	if err != nil {
		t.Fatal(err)
	}
	if string(content) != expect {
		t.Errorf("want %s, got %s\n", expect, string(content))
	}
}
