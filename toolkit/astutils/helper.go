package astutils

import (
	"github.com/pkg/errors"
	"github.com/unionj-cloud/go-doudou/v2/toolkit/errorx"
	"go/ast"
	"go/parser"
	"go/token"
	"golang.org/x/tools/go/packages"
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

func GetPkgPath(filePath string) string {
	pkgs, err := packages.Load(&packages.Config{
		Mode: packages.NeedName,
		Dir:  filePath,
	})
	if err != nil {
		errorx.Wrap(err)
		return ""
	}
	if len(pkgs) == 0 {
		errorx.Wrap(errors.New("no package found"))
		return ""
	}
	return pkgs[0].PkgPath
}
