package table

import (
	"cloud/unionj/papilio/kit/astutils"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
)

const testDir = "/Users/wubin1989/workspace/cloud/papilio/kit/ddl/example/models"

func ExampleNewTableFromStruct() {
	var files []string
	var err error
	err = filepath.Walk(testDir, astutils.Visit(&files))
	if err != nil {
		panic(err)
	}
	var sc astutils.StructCollector
	for _, file := range files {
		fset := token.NewFileSet()
		root, err := parser.ParseFile(fset, file, nil, parser.ParseComments)
		if err != nil {
			panic(err)
		}
		ast.Walk(&sc, root)
	}
	flattened := sc.FlatEmbed()

	for _, sm := range flattened {
		table := NewTableFromStruct(sm)
		fmt.Println(table)
	}
	//Output:
}
