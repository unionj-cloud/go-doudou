package codegen

import (
	"github.com/pkg/errors"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/unionj-cloud/go-doudou/v2/cmd/internal/astutils"
	"os"
	"path/filepath"
	"testing"
)

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
