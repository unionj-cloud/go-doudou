package astutils

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	"github.com/unionj-cloud/go-doudou/v2/toolkit/pathutils"
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
