package astutils

import (
	v3 "github.com/unionj-cloud/go-doudou/toolkit/openapi/v3"
	"go/ast"
	"go/token"
	"strings"
)

type EnumCollector struct {
	Methods    map[string][]MethodMeta
	Package    PackageMeta
	exprString func(ast.Expr) string
	Consts     map[string][]string
	Enums      map[string]EnumMeta
}

func IsEnum(methods []MethodMeta) bool {
	methodMap := make(map[string]struct{})
	for _, item := range methods {
		methodMap[item.String()] = struct{}{}
	}
	for _, item := range v3.IEnumMethods {
		if _, ok := methodMap[item]; !ok {
			return false
		}
	}
	return true
}

// Visit traverse each node from source code
func (sc *EnumCollector) Visit(n ast.Node) ast.Visitor {
	return sc.Collect(n)
}

// Collect collects all structs from source code
func (sc *EnumCollector) Collect(n ast.Node) ast.Visitor {
	switch spec := n.(type) {
	case *ast.Package:
		return sc
	case *ast.File: // actually it is package name
		sc.Package = PackageMeta{
			Name: spec.Name.Name,
		}
		return sc
	case *ast.FuncDecl:
		if spec.Recv != nil {
			typeName := strings.TrimPrefix(sc.exprString(spec.Recv.List[0].Type), "*")
			methods, _ := sc.Methods[typeName]
			methods = append(methods, GetMethodMeta(spec))
			if sc.Methods == nil {
				sc.Methods = make(map[string][]MethodMeta)
			}
			sc.Methods[typeName] = methods
		}
	case *ast.GenDecl:
		if spec.Tok == token.CONST {
			var comments []string
			if spec.Doc != nil {
				for _, comment := range spec.Doc.List {
					comments = append(comments, strings.TrimSpace(strings.TrimPrefix(comment.Text, "//")))
				}
			}
			var typeName string
			for _, item := range spec.Specs {
				valueSpec := item.(*ast.ValueSpec)
				if len(valueSpec.Names) == 0 {
					continue
				}
				switch specType := valueSpec.Type.(type) {
				case *ast.Ident:
					typeName = specType.Name
				case nil:
					if len(valueSpec.Values) > 0 {
						switch valueExpr := valueSpec.Values[0].(type) {
						case *ast.BasicLit:
							switch valueExpr.Kind {
							case token.INT:
								typeName = "int"
							case token.FLOAT:
								typeName = "float64"
							case token.IMAG:
								typeName = "complex128"
							case token.CHAR:
								typeName = "rune"
							case token.STRING:
								typeName = "string"
							default:
								continue
							}
						}
					}
				}
				sc.Consts[typeName] = append(sc.Consts[typeName], valueSpec.Names[0].Name)
			}
		}
	}
	return nil
}

// NewEnumCollector initializes an EnumCollector
func NewEnumCollector(exprString func(ast.Expr) string) *EnumCollector {
	return &EnumCollector{
		Methods:    make(map[string][]MethodMeta),
		Package:    PackageMeta{},
		exprString: exprString,
		Consts:     make(map[string][]string),
		Enums:      make(map[string]EnumMeta),
	}
}
