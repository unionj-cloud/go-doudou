package astutils

import (
	"go/ast"
	"go/parser"
	"go/token"

	"github.com/sirupsen/logrus"
)

// StaticMethodCollector collect methods by parsing source code
type StaticMethodCollector struct {
	Methods    []MethodMeta
	Package    PackageMeta
	exprString func(ast.Expr) string
}

// Visit traverse each node from source code
func (sc *StaticMethodCollector) Visit(n ast.Node) ast.Visitor {
	return sc.Collect(n)
}

// Collect collects all static methods from source code
func (sc *StaticMethodCollector) Collect(n ast.Node) ast.Visitor {
	switch spec := n.(type) {
	case *ast.Package:
		return sc
	case *ast.File: // actually it is package name
		sc.Package = PackageMeta{
			Name: spec.Name.Name,
		}
		return sc
	case *ast.FuncDecl:
		if spec.Recv == nil {
			sc.Methods = append(sc.Methods, GetMethodMeta(spec))
		}
	case *ast.GenDecl:
	}
	return nil
}

type StaticMethodCollectorOption func(collector *StaticMethodCollector)

// NewStaticMethodCollector initializes an StaticMethodCollector
func NewStaticMethodCollector(exprString func(ast.Expr) string, opts ...StaticMethodCollectorOption) *StaticMethodCollector {
	sc := &StaticMethodCollector{
		Methods:    make([]MethodMeta, 0),
		Package:    PackageMeta{},
		exprString: exprString,
	}
	for _, opt := range opts {
		opt(sc)
	}
	return sc
}

// BuildStaticMethodCollector initializes an StaticMethodCollector and collects static methods
func BuildStaticMethodCollector(file string, exprString func(ast.Expr) string) StaticMethodCollector {
	sc := NewStaticMethodCollector(exprString)
	fset := token.NewFileSet()
	root, err := parser.ParseFile(fset, file, nil, parser.ParseComments)
	if err != nil {
		logrus.Panicln(err)
	}
	ast.Walk(sc, root)
	return *sc
}
