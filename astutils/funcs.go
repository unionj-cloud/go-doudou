package astutils

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/unionj-cloud/go-doudou/stringutils"
	"go/ast"
	"go/format"
	"golang.org/x/tools/imports"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"
	"unicode"
)

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
		panic(err)
	}
	err = ioutil.WriteFile(file, res, 0644)
	if err != nil {
		panic(err)
	}
}

func GetMethodMeta(spec *ast.FuncDecl) MethodMeta {
	methodName := ExprString(spec.Name)
	mm := NewMethodMeta(spec.Type, ExprString)
	mm.Name = methodName
	return mm
}

func NewMethodMeta(ft *ast.FuncType, exprString func(ast.Expr) string) MethodMeta {
	var params, results []FieldMeta
	for _, param := range ft.Params.List {
		var pComments []string
		if param.Comment != nil {
			for _, comment := range param.Comment.List {
				pComments = append(pComments, comment.Text)
			}
		}
		pt := exprString(param.Type)
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
		params = append(params, FieldMeta{
			Name:     "",
			Type:     pt,
			Tag:      "",
			Comments: pComments,
		})
	}
	if ft.Results != nil {
		for _, result := range ft.Results.List {
			var rComments []string
			if result.Comment != nil {
				for _, comment := range result.Comment.List {
					rComments = append(rComments, comment.Text)
				}
			}
			rt := exprString(result.Type)
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
			results = append(results, FieldMeta{
				Name:     "",
				Type:     rt,
				Tag:      "",
				Comments: rComments,
			})
		}
	}
	return MethodMeta{
		Params:  params,
		Results: results,
	}
}

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

type PackageMeta struct {
	Name string
}

type FieldMeta struct {
	Name     string
	Type     string
	Tag      string
	Comments []string
	IsExport bool
	DocName  string
}

type StructMeta struct {
	Name     string
	Fields   []FieldMeta
	Comments []string
	Methods  []MethodMeta
	IsExport bool
}

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
		} else {
			return "[" + ExprString(_expr.Len) + "]" + ExprString(_expr.Elt)
		}
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
	default:
		panic(fmt.Sprintf("not support expression: %+v\n", expr))
	}
}

type MethodMeta struct {
	Recv     string
	Name     string
	Params   []FieldMeta
	Results  []FieldMeta
	PathVars []FieldMeta
	Comments []string
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
	t, err := template.New("method.tmpl").Parse(methodTmpl)
	if err != nil {
		panic(err)
	}
	var buf bytes.Buffer
	err = t.Execute(&buf, mm)
	if err != nil {
		panic(err)
	}
	var res []byte
	res, err = format.Source(buf.Bytes())
	if err != nil {
		panic(err)
	}
	result := string(res)
	if isAnony {
		return strings.Replace(result, "func placeholder(", "func(", 1)
	}
	return result
}

type InterfaceMeta struct {
	Name     string
	Methods  []MethodMeta
	Comments []string
}

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
	if firstLine, err = reader.ReadString('\n'); err != nil {
		panic(err)
	}
	return strings.TrimSpace(strings.TrimPrefix(firstLine, "module"))
}

func GetImportPath(file string) string {
	dir, _ := os.Getwd()
	return GetMod() + strings.TrimPrefix(file, dir)
}
