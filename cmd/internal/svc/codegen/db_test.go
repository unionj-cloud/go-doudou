package codegen

import (
	"github.com/pkg/errors"
	. "github.com/smartystreets/goconvey/convey"
	"io/ioutil"
	"os"
	"path/filepath"
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
	configfile := filepath.Join(dir, "db", "db.go")
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

func TestGenDbPanic(t *testing.T) {
	Convey("GenDb should panic from MkdirAll", t, func() {
		MkdirAll = os.MkdirAll
		Open = os.Open
		Create = os.Create
		Stat = os.Stat
		MkdirAll = func(path string, perm os.FileMode) error {
			return errors.New("mock MkdirAll error")
		}
		dir := testDir + "db"
		InitSvc(dir)
		defer os.RemoveAll(dir)
		So(func() {
			GenDb(dir)
		}, ShouldPanic)
	})
}

func TestGenDbPanic2(t *testing.T) {
	Convey("GenDb should panic from Open", t, func() {
		MkdirAll = os.MkdirAll
		Open = os.Open
		Create = os.Create
		Stat = os.Stat
		Open = func(name string) (*os.File, error) {
			return nil, errors.New("mock Open error")
		}
		dir := testDir + "db"
		InitSvc(dir)
		defer os.RemoveAll(dir)
		So(func() {
			GenDb(dir)
		}, ShouldPanic)
	})
}

func TestGenDbPanic3(t *testing.T) {
	Convey("GenDb should panic from Create", t, func() {
		MkdirAll = os.MkdirAll
		Open = os.Open
		Create = os.Create
		Stat = os.Stat
		Create = func(name string) (*os.File, error) {
			return nil, errors.New("mock Create error")
		}
		dir := testDir + "db"
		InitSvc(dir)
		defer os.RemoveAll(dir)
		So(func() {
			GenDb(dir)
		}, ShouldPanic)
	})
}

func TestGenDbWarn(t *testing.T) {
	Convey("GenDb should warn", t, func() {
		MkdirAll = os.MkdirAll
		Open = os.Open
		Create = os.Create
		Stat = os.Stat
		Stat = func(name string) (os.FileInfo, error) {
			return nil, errors.New("mock Stat error")
		}
		dir := testDir + "db"
		InitSvc(dir)
		defer os.RemoveAll(dir)
		So(func() {
			GenDb(dir)
		}, ShouldNotPanic)
	})
}

func TestGenDbPanic4(t *testing.T) {
	Convey("GenDb should panic from dbTmpl", t, func() {
		MkdirAll = os.MkdirAll
		Open = os.Open
		Create = os.Create
		Stat = os.Stat
		dbTmpl = `{{test tmpl}`
		dir := testDir + "db"
		InitSvc(dir)
		defer os.RemoveAll(dir)
		So(func() {
			GenDb(dir)
		}, ShouldPanic)
	})
}
