package astutils

import (
	"fmt"
	"github.com/unionj-cloud/go-doudou/pathutils"
	"go/ast"
	"go/parser"
	"go/token"
	"testing"
)

func TestInterfaceCollector(t *testing.T) {
	file := pathutils.Abs("testfiles/svc.go")
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
