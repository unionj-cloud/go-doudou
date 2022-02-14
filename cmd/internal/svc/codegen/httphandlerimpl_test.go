package codegen

import (
	"github.com/iancoleman/strcase"
	"github.com/unionj-cloud/go-doudou/cmd/internal/astutils"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestGenHttpHandlerImplWithImpl(t *testing.T) {
	dir := "testdata"
	defer os.RemoveAll(filepath.Join(dir, "transport"))
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
			GenHttpHandlerImplWithImpl(tt.args.dir, tt.args.ic, true, strcase.ToLowerCamel)
		})
	}
}

func TestGenHttpHandlerImpl(t *testing.T) {
	dir := testDir + "handlerImpl12"
	InitSvc(dir)
	defer os.RemoveAll(dir)
	ic := astutils.BuildInterfaceCollector(filepath.Join(dir, "svc.go"), astutils.ExprString)
	GenHttpHandlerImpl(dir, ic)
	expect := `package httpsrv

import (
	"net/http"
	service "testdatahandlerImpl12"
)

type TestdatahandlerImpl12HandlerImpl struct {
	testdatahandlerImpl12 service.TestdatahandlerImpl12
}

func (receiver *TestdatahandlerImpl12HandlerImpl) PageUsers(w http.ResponseWriter, r *http.Request) {
	panic("implement me")
}

func NewTestdatahandlerImpl12Handler(testdatahandlerImpl12 service.TestdatahandlerImpl12) TestdatahandlerImpl12Handler {
	return &TestdatahandlerImpl12HandlerImpl{
		testdatahandlerImpl12,
	}
}
`
	file := filepath.Join(dir, "transport", "httpsrv", "handlerimpl.go")
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

func TestGenHttpHandlerImplWithImpl2(t *testing.T) {
	svcfile := testDir + "/svc.go"
	ic := astutils.BuildInterfaceCollector(svcfile, astutils.ExprString)
	defer os.RemoveAll(testDir + "/transport")
	type args struct {
		dir string
		ic  astutils.InterfaceCollector
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "1",
			args: args{
				dir: testDir,
				ic:  ic,
			},
		},
		{
			name: "2",
			args: args{
				dir: testDir,
				ic:  ic,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			GenHttpHandlerImplWithImpl(tt.args.dir, tt.args.ic, true, strcase.ToLowerCamel)
		})
	}
}
