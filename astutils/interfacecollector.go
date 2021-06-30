package astutils

import (
	"fmt"
	"github.com/iancoleman/strcase"
	"github.com/sirupsen/logrus"
	"github.com/unionj-cloud/go-doudou/stringutils"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
)

type InterfaceCollector struct {
	Interfaces []InterfaceMeta
	Package    PackageMeta
	exprString func(ast.Expr) string
}

func (ic *InterfaceCollector) Visit(n ast.Node) ast.Visitor {
	return ic.Collect(n)
}

func (ic *InterfaceCollector) Collect(n ast.Node) ast.Visitor {
	switch spec := n.(type) {
	case *ast.Package:
		return ic
	case *ast.File: // actually it is package name
		ic.Package = PackageMeta{
			Name: spec.Name.Name,
		}
		return ic
	case *ast.GenDecl:
		if spec.Tok == token.TYPE {
			var comments []string
			if spec.Doc != nil {
				for _, comment := range spec.Doc.List {
					comments = append(comments, strings.TrimSpace(strings.TrimPrefix(comment.Text, "//")))
				}
			}
			for _, item := range spec.Specs {
				typeSpec := item.(*ast.TypeSpec)
				typeName := typeSpec.Name.Name
				switch specType := typeSpec.Type.(type) {
				case *ast.InterfaceType:
					var methods []MethodMeta
					for _, method := range specType.Methods.List {
						if len(method.Names) == 0 {
							panic("no method name")
						}
						mn := method.Names[0].Name

						var mComments []string
						if method.Doc != nil {
							for _, comment := range method.Doc.List {
								mComments = append(mComments, strings.TrimSpace(strings.TrimPrefix(comment.Text, "//")))
							}
						}

						var ft *ast.FuncType
						var ok bool
						if ft, ok = method.Type.(*ast.FuncType); !ok {
							panic("not funcType")
						}
						var params, results []FieldMeta
						pkeymap := make(map[string]int)
						for _, param := range ft.Params.List {
							var pComments []string
							if param.Doc != nil {
								for _, comment := range param.Doc.List {
									pComments = append(pComments, strings.TrimSpace(strings.TrimPrefix(comment.Text, "//")))
								}
							}
							pt := ic.exprString(param.Type)
							if len(param.Names) > 0 {
								for _, name := range param.Names {
									params = append(params, FieldMeta{
										Name:     name.Name,
										Type:     pt,
										Tag:      "",
										Comments: pComments,
									})
								}
								continue
							}
							var pn string
							elemt := strings.TrimPrefix(pt, "*")
							if stringutils.IsNotEmpty(elemt) {
								if strings.Contains(elemt, "[") {
									elemt = elemt[strings.Index(elemt, "]")+1:]
									elemt = strings.TrimPrefix(elemt, "*")
								}
								splits := strings.Split(elemt, ".")
								_key := "p" + strcase.ToLowerCamel(splits[len(splits)-1][0:1])
								if _, exists := pkeymap[_key]; exists {
									pkeymap[_key]++
									pn = _key + fmt.Sprintf("%d", pkeymap[_key])
								} else {
									pkeymap[_key]++
									pn = _key
								}
							}
							params = append(params, FieldMeta{
								Name:     pn,
								Type:     pt,
								Tag:      "",
								Comments: pComments,
							})
						}
						if ft.Results != nil {
							rkeymap := make(map[string]int)
							for _, result := range ft.Results.List {
								var rComments []string
								if result.Doc != nil {
									for _, comment := range result.Doc.List {
										rComments = append(rComments, strings.TrimSpace(strings.TrimPrefix(comment.Text, "//")))
									}
								}
								rt := ic.exprString(result.Type)
								if len(result.Names) > 0 {
									for _, name := range result.Names {
										results = append(results, FieldMeta{
											Name:     name.Name,
											Type:     rt,
											Tag:      "",
											Comments: rComments,
										})
									}
									continue
								}
								var rn string
								elemt := strings.TrimPrefix(rt, "*")
								if stringutils.IsNotEmpty(elemt) {
									if strings.Contains(elemt, "[") {
										elemt = elemt[strings.Index(elemt, "]")+1:]
										elemt = strings.TrimPrefix(elemt, "*")
									}
									splits := strings.Split(elemt, ".")
									_key := "r" + strcase.ToLowerCamel(splits[len(splits)-1][0:1])
									if _, exists := rkeymap[_key]; exists {
										rkeymap[_key]++
										rn = _key + fmt.Sprintf("%d", rkeymap[_key])
									} else {
										rkeymap[_key]++
										rn = _key
									}
								}
								results = append(results, FieldMeta{
									Name:     rn,
									Type:     rt,
									Tag:      "",
									Comments: rComments,
								})
							}
						}
						methods = append(methods, MethodMeta{
							Name:     mn,
							Params:   params,
							Results:  results,
							Comments: mComments,
						})
					}

					ic.Interfaces = append(ic.Interfaces, InterfaceMeta{
						Name:     typeName,
						Methods:  methods,
						Comments: comments,
					})
				}
			}
		}
	}
	return nil
}

func NewInterfaceCollector(exprString func(ast.Expr) string) *InterfaceCollector {
	return &InterfaceCollector{
		exprString: exprString,
	}
}

func BuildInterfaceCollector(file string, exprString func(ast.Expr) string) InterfaceCollector {
	ic := NewInterfaceCollector(exprString)
	fset := token.NewFileSet()
	root, err := parser.ParseFile(fset, file, nil, parser.ParseComments)
	if err != nil {
		logrus.Panicln(err)
	}
	ast.Walk(ic, root)
	return *ic
}
