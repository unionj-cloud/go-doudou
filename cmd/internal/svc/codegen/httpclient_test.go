package codegen

import (
	"github.com/iancoleman/strcase"
	"github.com/pkg/errors"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/unionj-cloud/go-doudou/cmd/internal/astutils"
	"os"
	"path/filepath"
	"testing"
)

func TestGenGoClient(t *testing.T) {
	dir := testDir + "client1"
	InitSvc(dir)
	defer os.RemoveAll(dir)
	svcfile := filepath.Join(dir, "svc.go")
	ic := astutils.BuildInterfaceCollector(svcfile, astutils.ExprString)

	type args struct {
		dir string
		ic  astutils.InterfaceCollector
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "",
			args: args{
				dir: dir,
				ic:  ic,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			GenGoClient(tt.args.dir, tt.args.ic, "", 1, strcase.ToLowerCamel)
		})
	}
}

func TestGenGoClient2(t *testing.T) {
	svcfile := filepath.Join(testDir, "svc.go")
	ic := astutils.BuildInterfaceCollector(svcfile, astutils.ExprString)

	type args struct {
		dir string
		ic  astutils.InterfaceCollector
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "",
			args: args{
				dir: testDir,
				ic:  ic,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			GenGoClient(tt.args.dir, tt.args.ic, "", 1, strcase.ToLowerCamel)
		})
	}
}

func TestGenGoClientPanic_Stat(t *testing.T) {
	Convey("Test GenGoClient panic from Stat", t, func() {
		MkdirAll = os.MkdirAll
		Open = os.Open
		Create = os.Create
		Stat = os.Stat
		Stat = func(name string) (os.FileInfo, error) {
			return nil, errors.New("mock Stat error")
		}

		svcfile := filepath.Join(testDir, "svc.go")
		ic := astutils.BuildInterfaceCollector(svcfile, astutils.ExprString)

		So(func() {
			GenGoClient(testDir, ic, "", 1, strcase.ToLowerCamel)
		}, ShouldPanic)
	})
}

func TestGenGoClientPanic_Create(t *testing.T) {
	Convey("Test GenGoClient panic from Create", t, func() {
		MkdirAll = os.MkdirAll
		Open = os.Open
		Create = os.Create
		Stat = os.Stat
		Create = func(name string) (*os.File, error) {
			return nil, errors.New("mock Create error")
		}
		svcfile := filepath.Join(testDir, "svc.go")
		ic := astutils.BuildInterfaceCollector(svcfile, astutils.ExprString)

		So(func() {
			GenGoClient(testDir, ic, "", 1, strcase.ToLowerCamel)
		}, ShouldPanic)
	})
}

func TestGenGoClientPanic_Open(t *testing.T) {
	Convey("Test GenGoClient panic from Open", t, func() {
		MkdirAll = os.MkdirAll
		Open = os.Open
		Create = os.Create
		Stat = os.Stat
		Open = func(name string) (*os.File, error) {
			return nil, errors.New("mock Open error")
		}
		svcfile := filepath.Join(testDir, "svc.go")
		ic := astutils.BuildInterfaceCollector(svcfile, astutils.ExprString)

		So(func() {
			GenGoClient(testDir, ic, "", 1, strcase.ToLowerCamel)
		}, ShouldPanic)
	})
}

func TestGenGoClientPanic_MkdirAll(t *testing.T) {
	Convey("Test GenGoClient panic from MkdirAll", t, func() {
		MkdirAll = os.MkdirAll
		Open = os.Open
		Create = os.Create
		Stat = os.Stat
		MkdirAll = func(path string, perm os.FileMode) error {
			return errors.New("mock MkdirAll error")
		}
		svcfile := filepath.Join(testDir, "svc.go")
		ic := astutils.BuildInterfaceCollector(svcfile, astutils.ExprString)

		So(func() {
			GenGoClient(testDir, ic, "", 1, strcase.ToLowerCamel)
		}, ShouldPanic)
	})
}
