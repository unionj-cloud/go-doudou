package astutils

import (
	"fmt"
	"github.com/iancoleman/strcase"
	"github.com/sirupsen/logrus"
	"github.com/unionj-cloud/go-doudou/toolkit/stringutils"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
)

// InterfaceCollector collect interfaces by parsing source code
type InterfaceCollector struct {
	Interfaces []InterfaceMeta
	Package    PackageMeta
	exprString func(ast.Expr) string
	cmap       ast.CommentMap
}

// Visit traverse each node from source code
func (ic *InterfaceCollector) Visit(n ast.Node) ast.Visitor {
	return ic.Collect(n)
}

// Collect collects all interfaces from source code
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
			comments := doc2Comments(spec.Doc)
			for _, item := range spec.Specs {
				typeSpec := item.(*ast.TypeSpec)
				typeName := typeSpec.Name.Name
				switch specType := typeSpec.Type.(type) {
				case *ast.InterfaceType:
					ic.Interfaces = append(ic.Interfaces, InterfaceMeta{
						Name:     typeName,
						Methods:  ic.field2Methods(specType.Methods.List),
						Comments: comments,
					})
				}
			}
		}
	}
	return nil
}

func (ic *InterfaceCollector) field2Methods(list []*ast.Field) []MethodMeta {
	var methods []MethodMeta
	for _, method := range list {
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

		ft, _ := method.Type.(*ast.FuncType)
		var params []FieldMeta
		if ft.Params != nil {
			params = ic.field2Params(ft.Params.List)
		}
		var results []FieldMeta
		if ft.Results != nil {
			results = ic.field2Results(ft.Results.List)
		}
		methods = append(methods, MethodMeta{
			Name:     mn,
			Params:   params,
			Results:  results,
			Comments: mComments,
		})
	}
	return methods
}

func (ic *InterfaceCollector) field2Params(list []*ast.Field) []FieldMeta {
	var params []FieldMeta
	pkeymap := make(map[string]int)
	for _, param := range list {
		pt := ic.exprString(param.Type)
		if len(param.Names) > 0 {
			for i, name := range param.Names {
				field := FieldMeta{
					Name: name.Name,
					Type: pt,
				}
				var cnode ast.Node
				if i == 0 {
					cnode = param
				} else {
					cnode = name
				}
				if cmts, exists := ic.cmap[cnode]; exists {
					for _, comment := range cmts {
						field.Comments = append(field.Comments, strings.TrimSpace(strings.TrimPrefix(comment.Text(), "//")))
					}
				}
				params = append(params, field)
			}
			continue
		}
		var pComments []string
		if cmts, exists := ic.cmap[param]; exists {
			for _, comment := range cmts {
				pComments = append(pComments, strings.TrimSpace(strings.TrimPrefix(comment.Text(), "//")))
			}
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
			Comments: pComments,
		})
	}
	return params
}

func (ic *InterfaceCollector) field2Results(list []*ast.Field) []FieldMeta {
	var results []FieldMeta
	rkeymap := make(map[string]int)
	for _, result := range list {
		var rComments []string
		if cmts, exists := ic.cmap[result]; exists {
			for _, comment := range cmts {
				rComments = append(rComments, strings.TrimSpace(strings.TrimPrefix(comment.Text(), "//")))
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
	return results
}

func doc2Comments(doc *ast.CommentGroup) []string {
	var comments []string
	if doc != nil {
		for _, comment := range doc.List {
			comments = append(comments, strings.TrimSpace(strings.TrimPrefix(comment.Text, "//")))
		}
	}
	return comments
}

// NewInterfaceCollector initializes an InterfaceCollector
func NewInterfaceCollector(exprString func(ast.Expr) string) *InterfaceCollector {
	return &InterfaceCollector{
		exprString: exprString,
	}
}

// BuildInterfaceCollector initializes an InterfaceCollector and collects interfaces
func BuildInterfaceCollector(file string, exprString func(ast.Expr) string) InterfaceCollector {
	ic := NewInterfaceCollector(exprString)
	fset := token.NewFileSet()
	root, err := parser.ParseFile(fset, file, nil, parser.ParseComments)
	if err != nil {
		logrus.Panicln(err)
	}
	ic.cmap = ast.NewCommentMap(fset, root, root.Comments)
	ast.Walk(ic, root)
	return *ic
}
