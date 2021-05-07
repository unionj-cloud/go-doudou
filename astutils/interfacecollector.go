package astutils

import (
	"go/ast"
	"go/token"
)

type MethodMeta struct {
	Name     string
	Params   []FieldMeta
	Results  []FieldMeta
	Comments []string
}

type InterfaceMeta struct {
	Name     string
	Methods  []MethodMeta
	Comments []string
}

type InterfaceCollector struct {
	Interfaces []InterfaceMeta
	Package    PackageMeta
}

func (ic *InterfaceCollector) Visit(n ast.Node) ast.Visitor {
	return ic.Collect(n)
}

func (sc *InterfaceCollector) Collect(n ast.Node) ast.Visitor {
	switch spec := n.(type) {
	case *ast.Package:
		return sc
	case *ast.File: // actually it is package name
		sc.Package = PackageMeta{
			Name: spec.Name.Name,
		}
		return sc
	case *ast.GenDecl:
		if spec.Tok == token.TYPE {
			var comments []string
			if spec.Doc != nil {
				for _, comment := range spec.Doc.List {
					comments = append(comments, comment.Text)
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
						if method.Comment != nil {
							for _, comment := range method.Comment.List {
								mComments = append(mComments, comment.Text)
							}
						}

						var ft *ast.FuncType
						var ok bool
						if ft, ok = method.Type.(*ast.FuncType); !ok {
							panic("not funcType")
						}
						var params, results []FieldMeta
						for _, param := range ft.Params.List {
							var pn string
							if len(param.Names) > 0 {
								pn = param.Names[0].Name
							}
							pt := exprString(param.Type)
							var pComments []string
							if param.Comment != nil {
								for _, comment := range param.Comment.List {
									pComments = append(pComments, comment.Text)
								}
							}
							params = append(params, FieldMeta{
								Name:     pn,
								Type:     pt,
								Tag:      "",
								Comments: pComments,
							})
						}
						for _, result := range ft.Results.List {
							var rn string
							if len(result.Names) > 0 {
								rn = result.Names[0].Name
							}
							rt := exprString(result.Type)
							var rComments []string
							if result.Comment != nil {
								for _, comment := range result.Comment.List {
									rComments = append(rComments, comment.Text)
								}
							}
							results = append(results, FieldMeta{
								Name:     rn,
								Type:     rt,
								Tag:      "",
								Comments: rComments,
							})
						}
						methods = append(methods, MethodMeta{
							Name:     mn,
							Params:   params,
							Results:  results,
							Comments: mComments,
						})
					}

					sc.Interfaces = append(sc.Interfaces, InterfaceMeta{
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
