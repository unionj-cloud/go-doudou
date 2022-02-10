package astutils

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/unionj-cloud/go-doudou/toolkit/pathutils"
	"go/ast"
	"go/parser"
	"go/token"
	"testing"
)

func TestInterfaceCollector(t *testing.T) {
	file := pathutils.Abs("testdata/svc.go")
	fset := token.NewFileSet()
	root, err := parser.ParseFile(fset, file, nil, parser.ParseComments)
	if err != nil {
		panic(err)
	}
	//spew.Dump(root)
	sc := NewInterfaceCollector(ExprString)
	sc.cmap = ast.NewCommentMap(fset, root, root.Comments)
	ast.Walk(sc, root)
	fmt.Println(sc)
}

func TestBuildInterfaceCollector(t *testing.T) {
	file := pathutils.Abs("testdata/svc.go")
	ic := BuildInterfaceCollector(file, ExprString)
	assert.NotNil(t, ic)
}
