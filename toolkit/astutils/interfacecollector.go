package astutils

import (
	"fmt"
	"github.com/iancoleman/strcase"
	"github.com/sirupsen/logrus"
	"github.com/unionj-cloud/go-doudou/v2/toolkit/sliceutils"
	"github.com/unionj-cloud/go-doudou/v2/toolkit/stringutils"
	"go/ast"
	"go/parser"
	"go/token"
	"net/http"
	"regexp"
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

// GetShelves_ShelfBooks_Book
// shelves/:shelf/books/:book
func Pattern(method string) (httpMethod string, endpoint string) {
	httpMethods := []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete}
	re1, err := regexp.Compile("_?[A-Z]")
	if err != nil {
		panic(err)
	}
	method = re1.ReplaceAllStringFunc(method, func(s string) string {
		if strings.HasPrefix(s, "_") {
			return "/:" + strings.ToLower(strings.TrimPrefix(s, "_"))
		} else {
			return "/" + strings.ToLower(s)
		}
	})
	splits := strings.Split(method, "/")[1:]
	httpMethod = httpMethods[1]
	head := strings.ToUpper(splits[0])
	if sliceutils.StringContains(httpMethods, head) {
		httpMethod = head
		splits = splits[1:]
	}
	return httpMethod, strings.Join(splits, "/")
}

func pathVariables(endpoint string) (ret []string) {
	splits := strings.Split(endpoint, "/")
	pvs := sliceutils.StringFilter(splits, func(item string) bool {
		return stringutils.IsNotEmpty(item) && strings.HasPrefix(item, ":")
	})
	for _, v := range pvs {
		ret = append(ret, strings.TrimPrefix(v, ":"))
	}
	return
}

func (ic *InterfaceCollector) field2Methods(list []*ast.Field) []MethodMeta {
	var methods []MethodMeta
	for _, method := range list {
		if len(method.Names) == 0 {
			panic("no method name")
		}
		mn := method.Names[0].Name

		var mComments []string
		var annotations []Annotation
		if method.Doc != nil {
			for _, comment := range method.Doc.List {
				mComments = append(mComments, strings.TrimSpace(strings.TrimPrefix(comment.Text, "//")))
				annotations = append(annotations, GetAnnotations(comment.Text)...)
			}
		}

		ft, _ := method.Type.(*ast.FuncType)
		var params []FieldMeta
		if ft.Params != nil {
			params = ic.field2Params(ft.Params.List)
		}

		queryVars := make([]FieldMeta, 0)
		httpMethod, endpoint := Pattern(mn)
		pvs := pathVariables(endpoint)
		for i := range params {
			if sliceutils.StringContains(pvs, params[i].Name) {
				params[i].IsPathVariable = true
			}
			if httpMethod == http.MethodGet && !params[i].IsPathVariable && params[i].Type != "context.Context" {
				queryVars = append(queryVars, params[i])
			}
		}

		var results []FieldMeta
		if ft.Results != nil {
			results = ic.field2Results(ft.Results.List)
		}
		methods = append(methods, MethodMeta{
			Name:            mn,
			Params:          params,
			Results:         results,
			Comments:        mComments,
			Annotations:     annotations,
			HasPathVariable: len(pvs) > 0,
			HttpMethod:      httpMethod,
			QueryVars:       queryVars,
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
					for _, comment := range cmts[0].List {
						field.Comments = append(field.Comments, strings.TrimSpace(strings.TrimPrefix(comment.Text, "//")))
						field.Annotations = append(field.Annotations, GetAnnotations(comment.Text)...)
					}
				}
				var validateTags []string
				for _, item := range field.Annotations {
					if item.Name == "@validate" {
						validateTags = append(validateTags, item.Params...)
					}
				}
				field.ValidateTag = strings.Join(validateTags, ",")
				params = append(params, field)
			}
			continue
		}
		var pComments []string
		var annotations []Annotation
		if cmts, exists := ic.cmap[param]; exists {
			for _, comment := range cmts[0].List {
				pComments = append(pComments, strings.TrimSpace(strings.TrimPrefix(comment.Text, "//")))
				annotations = append(annotations, GetAnnotations(comment.Text)...)
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
		var validateTags []string
		for _, item := range annotations {
			if item.Name == "@validate" {
				validateTags = append(validateTags, item.Params...)
			}
		}
		params = append(params, FieldMeta{
			Name:        pn,
			Type:        pt,
			Comments:    pComments,
			Annotations: annotations,
			ValidateTag: strings.Join(validateTags, ","),
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
			for _, comment := range cmts[0].List {
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
