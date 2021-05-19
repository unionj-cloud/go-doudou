package codegen

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
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
