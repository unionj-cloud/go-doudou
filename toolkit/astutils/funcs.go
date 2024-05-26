package astutils

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
	"github.com/unionj-cloud/go-doudou/v2/toolkit/constants"
	"github.com/unionj-cloud/go-doudou/v2/toolkit/stringutils"
	"go/ast"
	"go/format"
	"golang.org/x/tools/imports"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"text/template"
	"unicode"
)

func GetImportStatements(input []byte) []byte {
	reg := regexp.MustCompile("(?s)import \\((.*?)\\)")
	if !reg.Match(input) {
		return nil
	}
	matches := reg.FindSubmatch(input)
	return matches[1]
}

func AppendImportStatements(src []byte, appendImports []byte) []byte {
	reg := regexp.MustCompile("(?s)import \\((.*?)\\)")
	if !reg.Match(src) {
		return src
	}
	matches := reg.FindSubmatch(src)
	old := matches[1]
	re := regexp.MustCompile(`[\r\n]+`)
	splits := re.Split(string(old), -1)
	oldmap := make(map[string]struct{})
	for _, item := range splits {
		oldmap[strings.TrimSpace(item)] = struct{}{}
	}
	splits = re.Split(string(appendImports), -1)
	var newimps []string
	for _, item := range splits {
		key := strings.TrimSpace(item)
		if _, ok := oldmap[key]; !ok {
			newimps = append(newimps, "\t"+key)
		}
	}
	if len(newimps) == 0 {
		return src
	}
	appendImports = []byte(constants.LineBreak + strings.Join(newimps, constants.LineBreak) + constants.LineBreak)
	return reg.ReplaceAllFunc(src, func(i []byte) []byte {
		old = append([]byte("import ("), old...)
		old = append(old, appendImports...)
		old = append(old, []byte(")")...)
		return old
	})
}

func GrpcRelatedModify(src []byte, metaName string, grpcSvcName string) []byte {
	expr := fmt.Sprintf(`type %sImpl struct {`, metaName)
	reg := regexp.MustCompile(expr)
	unimpl := fmt.Sprintf("pb.Unimplemented%sServer", grpcSvcName)
	if !strings.Contains(string(src), unimpl) {
		appendUnimpl := []byte(constants.LineBreak + unimpl + constants.LineBreak)
		src = reg.ReplaceAllFunc(src, func(i []byte) []byte {
			return append([]byte(expr), appendUnimpl...)
		})
	}
	var_pb := fmt.Sprintf("var _ pb.%sServer = (*%sImpl)(nil)", grpcSvcName, metaName)
	if !strings.Contains(string(src), var_pb) {
		appendVarPb := []byte(constants.LineBreak + var_pb + constants.LineBreak)
		src = reg.ReplaceAllFunc(src, func(i []byte) []byte {
			return append(appendVarPb, []byte(expr)...)
		})
	}
	return src
}

func RestRelatedModify(src []byte, metaName string) []byte {
	expr := fmt.Sprintf(`type %sImpl struct {`, metaName)
	reg := regexp.MustCompile(expr)
	var_ := fmt.Sprintf("var _ %s = (*%sImpl)(nil)", metaName, metaName)
	if !strings.Contains(string(src), var_) {
		appendVarPb := []byte(constants.LineBreak + var_ + constants.LineBreak)
		src = reg.ReplaceAllFunc(src, func(i []byte) []byte {
			return append(appendVarPb, []byte(expr)...)
		})
	}
	return src
}

// FixImport format source code and add missing import syntax automatically
func FixImport(src []byte, file string) {
	var (
		res []byte
		err error
	)
	if res, err = imports.Process(file, src, &imports.Options{
		TabWidth:  8,
		TabIndent: true,
		Comments:  true,
		Fragment:  true,
	}); err != nil {
		lines := strings.Split(string(src), "\n")
		errLine, _ := strconv.Atoi(strings.Split(err.Error(), ":")[1])
		startLine, endLine := errLine-5, errLine+5
		fmt.Println("Format fail:", errLine, err)
		if startLine < 0 {
			startLine = 0
		}
		if endLine > len(lines)-1 {
			endLine = len(lines) - 1
		}
		for i := startLine; i <= endLine; i++ {
			fmt.Println(i, lines[i])
		}
		errors.WithStack(fmt.Errorf("cannot format file: %w", err))
	} else {
		_ = ioutil.WriteFile(file, res, os.ModePerm)
		return
	}
	_ = ioutil.WriteFile(file, src, os.ModePerm)
}

// GetMethodMeta get method name then new MethodMeta struct from *ast.FuncDecl
func GetMethodMeta(spec *ast.FuncDecl) MethodMeta {
	methodName := ExprString(spec.Name)
	mm := NewMethodMeta(spec.Type, ExprString)
	mm.Name = methodName
	return mm
}

// NewMethodMeta new MethodMeta struct from *ast.FuncDecl
func NewMethodMeta(ft *ast.FuncType, exprString func(ast.Expr) string) MethodMeta {
	var params, results []FieldMeta
	for _, param := range ft.Params.List {
		pt := exprString(param.Type)
		if len(param.Names) > 0 {
			for _, name := range param.Names {
				params = append(params, FieldMeta{
					Name: name.Name,
					Type: pt,
					Tag:  "",
				})
			}
			continue
		}
		params = append(params, FieldMeta{
			Name: "",
			Type: pt,
			Tag:  "",
		})
	}
	if ft.Results != nil {
		for _, result := range ft.Results.List {
			rt := exprString(result.Type)
			if len(result.Names) > 0 {
				for _, name := range result.Names {
					results = append(results, FieldMeta{
						Name: name.Name,
						Type: rt,
						Tag:  "",
					})
				}
				continue
			}
			results = append(results, FieldMeta{
				Name: "",
				Type: rt,
				Tag:  "",
			})
		}
	}
	return MethodMeta{
		Params:  params,
		Results: results,
	}
}

// NewStructMeta new StructMeta from *ast.StructType
func NewStructMeta(structType *ast.StructType, exprString func(ast.Expr) string) StructMeta {
	var fields []FieldMeta
	re := regexp.MustCompile(`json:"(.*?)"`)
	for _, field := range structType.Fields.List {
		var fieldComments []string
		if field.Doc != nil {
			for _, comment := range field.Doc.List {
				fieldComments = append(fieldComments, strings.TrimSpace(strings.TrimPrefix(comment.Text, "//")))
			}
		}

		fieldType := exprString(field.Type)

		var tag string
		var docName string
		if field.Tag != nil {
			tag = strings.Trim(field.Tag.Value, "`")
			if re.MatchString(tag) {
				docName = strings.TrimSuffix(re.FindStringSubmatch(tag)[1], ",omitempty")
			}
		}

		if len(field.Names) > 0 {
			for _, name := range field.Names {
				_docName := docName
				if stringutils.IsEmpty(_docName) {
					_docName = name.Name
				}
				fields = append(fields, FieldMeta{
					Name:     name.Name,
					Type:     fieldType,
					Tag:      tag,
					Comments: fieldComments,
					IsExport: unicode.IsUpper(rune(name.Name[0])),
					DocName:  _docName,
				})
			}
		} else {
			splits := strings.Split(fieldType, ".")
			name := splits[len(splits)-1]
			fieldType = "embed:" + fieldType
			_docName := docName
			if stringutils.IsEmpty(_docName) {
				_docName = name
			}
			fields = append(fields, FieldMeta{
				Name:     name,
				Type:     fieldType,
				Tag:      tag,
				Comments: fieldComments,
				IsExport: unicode.IsUpper(rune(name[0])),
				DocName:  _docName,
			})
		}
	}
	return StructMeta{
		Fields: fields,
	}
}

// PackageMeta wraps package info
type PackageMeta struct {
	Name string
}

// FieldMeta wraps field info
type FieldMeta struct {
	Name     string
	Type     string
	Tag      string
	Comments []string
	IsExport bool
	// used in OpenAPI 3.0 spec as property name
	DocName string
	// Annotations of the field
	Annotations []Annotation
	// ValidateTag based on https://github.com/go-playground/validator
	// please refer to its documentation https://pkg.go.dev/github.com/go-playground/validator/v10
	ValidateTag    string
	IsPathVariable bool
}

// StructMeta wraps struct info
type StructMeta struct {
	Name     string
	Fields   []FieldMeta
	Comments []string
	Methods  []MethodMeta
	IsExport bool
	// go-doudou version
	Version string
}

// EnumMeta wraps struct info
type EnumMeta struct {
	Name   string
	Values []string
}

// ExprString return string representation from ast.Expr
func ExprString(expr ast.Expr) string {
	switch _expr := expr.(type) {
	case *ast.Ident:
		return _expr.Name
	case *ast.StarExpr:
		return "*" + ExprString(_expr.X)
	case *ast.SelectorExpr:
		return ExprString(_expr.X) + "." + _expr.Sel.Name
	case *ast.InterfaceType:
		return "interface{}"
	case *ast.ArrayType:
		if _expr.Len == nil {
			return "[]" + ExprString(_expr.Elt)
		}
		return "[" + ExprString(_expr.Len) + "]" + ExprString(_expr.Elt)
	case *ast.BasicLit:
		return _expr.Value
	case *ast.MapType:
		return "map[" + ExprString(_expr.Key) + "]" + ExprString(_expr.Value)
	case *ast.StructType:
		structmeta := NewStructMeta(_expr, ExprString)
		b, _ := json.Marshal(structmeta)
		return "anonystruct«" + string(b) + "»"
	case *ast.FuncType:
		return NewMethodMeta(_expr, ExprString).String()
	case *ast.ChanType:
		var result string
		if _expr.Dir == ast.SEND {
			result += "chan<- "
		} else if _expr.Dir == ast.RECV {
			result += "<-chan "
		} else {
			result += "chan "
		}
		return result + ExprString(_expr.Value)
	case *ast.Ellipsis:
		if _expr.Ellipsis.IsValid() {
			return "..." + ExprString(_expr.Elt)
		}
		panic(fmt.Sprintf("invalid ellipsis expression: %+v\n", expr))
	case *ast.IndexExpr:
		return ExprString(_expr.X) + "[" + ExprString(_expr.Index) + "]"
	case *ast.IndexListExpr:
		typeParams := lo.Map[ast.Expr, string](_expr.Indices, func(item ast.Expr, index int) string {
			return ExprString(item)
		})
		return ExprString(_expr.X) + "[" + strings.Join(typeParams, ", ") + "]"
	default:
		logrus.Infof("not support expression: %+v\n", expr)
		logrus.Infof("not support expression: %+v\n", reflect.TypeOf(expr))
		logrus.Infof("not support expression: %#v\n", reflect.TypeOf(expr))
		logrus.Infof("not support expression: %v\n", reflect.TypeOf(expr).String())
		return ""
		//panic(fmt.Sprintf("not support expression: %+v\n", expr))
	}
}

type Annotation struct {
	Name   string
	Params []string
}

var reAnno = regexp.MustCompile(`@(\S+?)\((.*?)\)`)

func GetAnnotations(text string) []Annotation {
	if !reAnno.MatchString(text) {
		return nil
	}
	var annotations []Annotation
	matches := reAnno.FindAllStringSubmatch(text, -1)
	for _, item := range matches {
		name := fmt.Sprintf(`@%s`, item[1])
		var params []string
		if stringutils.IsNotEmpty(item[2]) {
			params = strings.Split(strings.TrimSpace(item[2]), ",")
		}
		annotations = append(annotations, Annotation{
			Name:   name,
			Params: params,
		})
	}
	return annotations
}

// MethodMeta represents an api
type MethodMeta struct {
	// Recv method receiver
	Recv string
	// Name method name
	Name string
	// Params when generate client code from openapi3 spec json file, Params holds all method input parameters.
	// when generate client code from service interface in svc.go file, if there is struct type param, this struct type param will put into request body,
	// then others will be put into url as query string. if there is no struct type param and the api is a get request, all will be put into url as query string.
	// if there is no struct type param and the api is Not a get request, all will be put into request body as application/x-www-form-urlencoded data.
	// specially, if there is one or more v3.FileModel or []v3.FileModel params,
	// all will be put into request body as multipart/form-data data.
	Params []FieldMeta
	// Results response
	Results   []FieldMeta
	QueryVars []FieldMeta
	// PathVars not support when generate client code from service interface in svc.go file
	// when generate client code from openapi3 spec json file, PathVars is parameters in url as path variable.
	PathVars []FieldMeta
	// HeaderVars not support when generate client code from service interface in svc.go file
	// when generate client code from openapi3 spec json file, HeaderVars is parameters in header.
	HeaderVars []FieldMeta
	// BodyParams not support when generate client code from service interface in svc.go file
	// when generate client code from openapi3 spec json file, BodyParams is parameters in request body as query string.
	BodyParams *FieldMeta
	// BodyJSON not support when generate client code from service interface in svc.go file
	// when generate client code from openapi3 spec json file, BodyJSON is parameters in request body as json.
	BodyJSON *FieldMeta
	// Files not support when generate client code from service interface in svc.go file
	// when generate client code from openapi3 spec json file, Files is parameters in request body as multipart file.
	Files []FieldMeta
	// Comments of the method
	Comments []string
	// Path api path
	// not support when generate client code from service interface in svc.go file
	Path string
	// QueryParams not support when generate client code from service interface in svc.go file
	// when generate client code from openapi3 spec json file, QueryParams is parameters in url as query string.
	QueryParams *FieldMeta
	// Annotations of the method
	Annotations     []Annotation
	HasPathVariable bool
	// HttpMethod only accepts GET, PUT, POST, DELETE
	HttpMethod string
}

const methodTmpl = `func {{ if .Recv }}(receiver {{.Recv}}){{ end }} {{.Name}}({{- range $i, $p := .Params}}
    {{- if $i}},{{end}}
    {{- $p.Name}} {{$p.Type}}
    {{- end }}) ({{- range $i, $r := .Results}}
                     {{- if $i}},{{end}}
                     {{- $r.Name}} {{$r.Type}}
                     {{- end }})`

func (mm MethodMeta) String() string {
	if stringutils.IsNotEmpty(mm.Recv) && stringutils.IsEmpty(mm.Name) {
		panic("not valid code")
	}
	var isAnony bool
	if stringutils.IsEmpty(mm.Name) {
		isAnony = true
		mm.Name = "placeholder"
	}
	t, _ := template.New("method.tmpl").Parse(methodTmpl)
	var buf bytes.Buffer
	_ = t.Execute(&buf, mm)
	var res []byte
	res, _ = format.Source(buf.Bytes())
	result := string(res)
	if isAnony {
		return strings.Replace(result, "func placeholder(", "func(", 1)
	}
	return result
}

// InterfaceMeta wraps interface info
type InterfaceMeta struct {
	Name     string
	Methods  []MethodMeta
	Comments []string
}

// Visit visit each files
func Visit(files *[]string) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			logrus.Panicln(err)
		}
		if !info.IsDir() {
			*files = append(*files, path)
		}
		return nil
	}
}

// GetMod get module name from go.mod file
func GetMod() string {
	var (
		f         *os.File
		err       error
		firstLine string
	)
	dir, _ := os.Getwd()
	mod := filepath.Join(dir, "go.mod")
	if f, err = os.Open(mod); err != nil {
		panic(err)
	}
	reader := bufio.NewReader(f)
	firstLine, _ = reader.ReadString('\n')
	return strings.TrimSpace(strings.TrimPrefix(firstLine, "module"))
}

// GetImportPath get import path of pkg from dir
func GetImportPath(dir string) string {
	wd, _ := os.Getwd()
	return GetMod() + strings.ReplaceAll(strings.TrimPrefix(dir, wd), `\`, `/`)
}
