package tests

import (
	"github.com/unionj-cloud/go-doudou/pathutils"
	"github.com/unionj-cloud/go-doudou/svc"
	. "github.com/unionj-cloud/go-doudou/svc/codegen"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/unionj-cloud/go-doudou/astutils"
)

var testDir string

func init() {
	testDir = pathutils.Abs("testfiles")
}

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
	dir := testDir + "doc1"
	receiver := svc.Svc{
		Dir: dir,
	}
	receiver.Init()
	defer os.RemoveAll(dir)
	type args struct {
		dir string
		ic  astutils.InterfaceCollector
	}
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
