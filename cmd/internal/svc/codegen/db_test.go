package codegen

import (
	"os"
	"testing"

	"github.com/pkg/errors"
	. "github.com/smartystreets/goconvey/convey"
)

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
