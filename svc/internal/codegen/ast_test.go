package codegen

import (
	"github.com/stretchr/testify/assert"
	"github.com/unionj-cloud/go-doudou/astutils"
	"github.com/unionj-cloud/go-doudou/pathutils"
	"go/ast"
	"go/parser"
	"go/token"
	"testing"
)

func TestStructCollector_Alias(t *testing.T) {
	file := pathutils.Abs("testfiles/vo/alias.go")
	fset := token.NewFileSet()
	root, err := parser.ParseFile(fset, file, nil, parser.ParseComments)
	if err != nil {
		panic(err)
	}
	sc := astutils.NewStructCollector(ExprStringP)
	assert.Panics(t, func() {
		ast.Walk(sc, root)
	})
}
