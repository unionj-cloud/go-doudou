package astutils

import (
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
	"strings"
)

func IsSlice(t string) bool {
	return strings.Contains(t, "[") || strings.HasPrefix(t, "...")
}

func IsVarargs(t string) bool {
	return strings.HasPrefix(t, "...")
}

func ToSlice(t string) string {
	return "[]" + strings.TrimPrefix(t, "...")
}

// ElementType get element type string from slice
func ElementType(t string) string {
	if IsVarargs(t) {
		return strings.TrimPrefix(t, "...")
	}
	return t[strings.Index(t, "]")+1:]
}

func CollectStructsInFolder(dir string, sc *StructCollector) {
	dir, _ = filepath.Abs(dir)
	var files []string
	err := filepath.Walk(dir, Visit(&files))
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		if filepath.Ext(file) != ".go" {
			continue
		}
		root, err := parser.ParseFile(token.NewFileSet(), file, nil, parser.ParseComments)
		if err != nil {
			panic(err)
		}
		ast.Walk(sc, root)
	}
}
