package codegen

import (
	v3 "github.com/unionj-cloud/go-doudou/openapi/v3"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/unionj-cloud/go-doudou/astutils"
)

func buildIc(svcfile string) astutils.InterfaceCollector {
	if _, err := os.Stat(svcfile); os.IsNotExist(err) {
		logrus.Panicln(svcfile + " file cannot be found. Execute command go-doudou svc init first!")
	}
	var ic astutils.InterfaceCollector
	fset := token.NewFileSet()
	root, err := parser.ParseFile(fset, svcfile, nil, parser.ParseComments)
	if err != nil {
		logrus.Panicln(err)
	}
	ast.Walk(&ic, root)
	return ic
}

func TestGenDoc(t *testing.T) {
	type args struct {
		dir string
		ic  astutils.InterfaceCollector
	}
	dir := "/Users/wubin1989/workspace/cloud/go-doudou/example/user-svc"
	svcfile := filepath.Join(dir, "svc.go")
	ic := buildIc(svcfile)

	tests := []struct {
		name string
		args args
	}{
		{
			name: "1",
			args: args{
				dir,
				ic,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			GenDoc(tt.args.dir, tt.args.ic)
		})
	}
}

func Test_schemasOf(t *testing.T) {
	type args struct {
		vofile string
	}
	tests := []struct {
		name string
		args args
		want []v3.Schema
	}{
		{
			name: "Test_schemasOf",
			args: args{
				vofile: "/Users/wubin1989/workspace/chengdutreeyee/team3-cloud-analyse/vo/vo.go",
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := schemasOf(tt.args.vofile); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("schemasOf() = %v, want %v", got, tt.want)
			}
		})
	}
}
