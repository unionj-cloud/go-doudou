package codegen

import (
	"github.com/pkg/errors"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/unionj-cloud/go-doudou/cmd/internal/astutils"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestGenMain(t *testing.T) {
	MkdirAll = os.MkdirAll
	Open = os.Open
	Create = os.Create
	Stat = os.Stat
	dir := testDir + "main"
	InitSvc(dir)
	defer os.RemoveAll(dir)
	ic := astutils.BuildInterfaceCollector(filepath.Join(dir, "svc.go"), astutils.ExprString)
	GenMain(dir, ic)
	expect := `package main

import (
	ddhttp "github.com/unionj-cloud/go-doudou/framework/http"
	service "testdatamain"
    "testdatamain/config"
	"testdatamain/transport/httpsrv"
)

func main() {
	conf := config.LoadFromEnv()
    svc := service.NewTestdatamain(conf)
	handler := httpsrv.NewTestdatamainHandler(svc)
	srv := ddhttp.NewDefaultHttpSrv()
	srv.AddRoute(httpsrv.Routes(handler)...)
	srv.Run()
}
`
	file := filepath.Join(dir, "cmd", "main.go")
	f, err := os.Open(file)
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

func TestGenMainPanic_Stat(t *testing.T) {
	Convey("Test GenMain panic from Stat", t, func() {
		dir := testDir + "main"
		InitSvc(dir)
		defer os.RemoveAll(dir)
		MkdirAll = os.MkdirAll
		Open = os.Open
		Create = os.Create
		Stat = os.Stat
		Stat = func(name string) (os.FileInfo, error) {
			return nil, errors.New("mock Stat error")
		}

		svcfile := filepath.Join(dir, "svc.go")
		ic := astutils.BuildInterfaceCollector(svcfile, astutils.ExprString)

		So(func() {
			GenMain(dir, ic)
		}, ShouldNotPanic)
	})
}

func TestGenMainPanic_Create(t *testing.T) {
	Convey("Test GenMain panic from Create", t, func() {
		dir := testDir + "main"
		InitSvc(dir)
		defer os.RemoveAll(dir)
		MkdirAll = os.MkdirAll
		Open = os.Open
		Create = os.Create
		Stat = os.Stat
		Create = func(name string) (*os.File, error) {
			return nil, errors.New("mock Create error")
		}
		svcfile := filepath.Join(dir, "svc.go")
		ic := astutils.BuildInterfaceCollector(svcfile, astutils.ExprString)

		So(func() {
			GenMain(dir, ic)
		}, ShouldPanic)
	})
}

func TestGenMainPanic_Open(t *testing.T) {
	Convey("Test GenMain panic from Open", t, func() {
		dir := testDir + "main"
		InitSvc(dir)
		defer os.RemoveAll(dir)
		MkdirAll = os.MkdirAll
		Open = os.Open
		Create = os.Create
		Stat = os.Stat
		Open = func(name string) (*os.File, error) {
			return nil, errors.New("mock Open error")
		}
		svcfile := filepath.Join(dir, "svc.go")
		ic := astutils.BuildInterfaceCollector(svcfile, astutils.ExprString)

		So(func() {
			GenMain(dir, ic)
		}, ShouldPanic)
	})
}

func TestGenMainPanic_MkdirAll(t *testing.T) {
	Convey("Test GenMain panic from MkdirAll", t, func() {
		dir := testDir + "main"
		InitSvc(dir)
		defer os.RemoveAll(dir)
		MkdirAll = os.MkdirAll
		Open = os.Open
		Create = os.Create
		Stat = os.Stat
		MkdirAll = func(path string, perm os.FileMode) error {
			return errors.New("mock MkdirAll error")
		}
		svcfile := filepath.Join(dir, "svc.go")
		ic := astutils.BuildInterfaceCollector(svcfile, astutils.ExprString)

		So(func() {
			GenMain(dir, ic)
		}, ShouldPanic)
	})
}
