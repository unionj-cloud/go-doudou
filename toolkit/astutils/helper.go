package astutils

import (
	"bufio"
	"github.com/pkg/errors"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
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

//func GetPkgPath1(filePath string) string {
//	pkgs, err := packages.Load(&packages.Config{
//		Mode: packages.NeedName,
//		Dir:  filePath,
//	})
//	if err != nil {
//		panic(err)
//	}
//	if len(pkgs) == 0 {
//		panic(errors.New("no package found"))
//	}
//	if len(pkgs[0].Errors) > 0 {
//		for _, err = range pkgs[0].Errors {
//			panic(err)
//		}
//	}
//	return pkgs[0].PkgPath
//}

func GetPkgPath(filePath string) string {
	modf, err := FindGoMod(filePath)
	if err != nil {
		panic(err)
	}
	defer modf.Close()
	reader := bufio.NewReader(modf)
	firstLine, _ := reader.ReadString('\n')
	modName := strings.TrimSpace(strings.TrimPrefix(firstLine, "module"))
	filePath, _ = filepath.Abs(filePath)
	filePath = filepath.ToSlash(filePath)
	result := filePath[strings.Index(filePath, modName):]
	return result
}

func FindGoMod(filePath string) (*os.File, error) {
	path, _ := filepath.Abs(filePath)
	if os.IsPathSeparator(path[0]) {
		return nil, errors.New("Can not find go.mod")
	}
	var err error
	mf := filepath.Join(path, "go.mod")
	var f *os.File
	if f, err = os.Open(mf); err != nil {
		return FindGoMod(filepath.Dir(path))
	}
	return f, nil
}
