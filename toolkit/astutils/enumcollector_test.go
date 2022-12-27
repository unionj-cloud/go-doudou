package astutils

import (
	"github.com/unionj-cloud/go-doudou/v2/toolkit/pathutils"
	"go/ast"
	"go/parser"
	"go/token"
	"testing"
)

func TestEnum(t *testing.T) {
	file := pathutils.Abs("testdata/enum.go")
	fset := token.NewFileSet()
	root, err := parser.ParseFile(fset, file, nil, parser.ParseComments)
	if err != nil {
		panic(err)
	}
	sc := NewEnumCollector(ExprString)
	ast.Walk(sc, root)
}
