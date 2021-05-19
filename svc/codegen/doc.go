package codegen

import (
	"encoding/json"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"unicode"

	"github.com/iancoleman/strcase"
	"github.com/sirupsen/logrus"
	"github.com/unionj-cloud/go-doudou/astutils"
	v3 "github.com/unionj-cloud/go-doudou/openapi/v3"
	"github.com/unionj-cloud/go-doudou/stringutils"
)

var schemas map[string]v3.Schema

/**
bool

string

int  int8  int16  int32  int64
uint uint8 uint16 uint32 uint64 uintptr

byte // alias for uint8

rune // alias for int32
     // represents a Unicode code point

float32 float64

complex64 complex128
*/
func schemaOf(field astutils.FieldMeta) *v3.Schema {
	ft := strings.TrimPrefix(field.Type, "*")
	switch ft {
	case "int", "int8", "int16", "int32", "uint", "uint8", "uint16", "uint32", "byte", "rune":
		return v3.Int
	case "int64", "uint64", "uintptr":
		return v3.Int64
	case "bool":
		return v3.Bool
	case "string":
		return v3.String
	case "float32":
		return v3.Float32
	case "float64":
		return v3.Float64
	case "complex64", "complex128":
		return v3.Any
	case "multipart.FileHeader":
		return v3.File
	case "":
		return v3.Any
	default:
		if strings.Contains(ft, "map[") {
			elem := ft[strings.Index(ft, "]")+1:]
			elem = strings.TrimPrefix(elem, "*")
			return &v3.Schema{
				Type: v3.ObjectT,
				AdditionalProperties: schemaOf(astutils.FieldMeta{
					Type: elem,
				}),
			}
		}
		if strings.Contains(ft, "[") {
			elem := ft[strings.Index(ft, "]")+1:]
			elem = strings.TrimPrefix(elem, "*")
			return &v3.Schema{
				Type: v3.ArrayT,
				Items: schemaOf(astutils.FieldMeta{
					Type: elem,
				}),
			}
		}
		if !strings.Contains(ft, ".") {
			title := ft
			if unicode.IsUpper(rune(title[0])) {
				return &v3.Schema{
					Ref: "#/components/schemas/" + title,
				}
			}
		}
		if strings.Contains(ft, "vo.") {
			title := strings.TrimPrefix(ft, "vo.")
			return &v3.Schema{
				Ref: "#/components/schemas/" + title,
			}
		}
		return v3.Any
	}
}

func schemasOf(vofile string) []v3.Schema {
	fset := token.NewFileSet()
	root, err := parser.ParseFile(fset, vofile, nil, 0)
	if err != nil {
		panic(err)
	}
	var sc astutils.StructCollector
	ast.Walk(&sc, root)
	var ret []v3.Schema
	for _, item := range sc.Structs {
		if unicode.IsLower(rune(item.Name[0])) {
			continue
		}
		properties := make(map[string]*v3.Schema)
		for _, field := range item.Fields {
			properties[strcase.ToLowerCamel(field.Name)] = schemaOf(field)
		}
		ret = append(ret, v3.Schema{
			Title:      item.Name,
			Type:       v3.ObjectT,
			Properties: properties,
		})
	}
	return ret
}

func vosOf(ic astutils.InterfaceCollector) []string {
	if len(ic.Interfaces) <= 0 {
		return nil
	}
	vomap := make(map[string]int)
	var vos []string
	inter := ic.Interfaces[0]
	for _, method := range inter.Methods {
		for _, field := range method.Params {
			if strings.Contains(field.Type, "*vo.") || strings.Contains(field.Type, "vo.") {
				title := strings.TrimPrefix(strings.TrimPrefix(field.Type, "*"), "vo.")
				if _, ok := vomap[title]; !ok {
					vomap[title] = 1
					vos = append(vos, title)
				}
			}
		}
		for _, field := range method.Results {
			if strings.Contains(field.Type, "*vo.") || strings.Contains(field.Type, "vo.") {
				title := strings.TrimPrefix(strings.TrimPrefix(field.Type, "*"), "vo.")
				if _, ok := vomap[title]; !ok {
					vomap[title] = 1
					vos = append(vos, title)
				}
			}
		}
	}
	return vos
}

func hasFile(field astutils.FieldMeta) bool {
	if strings.Contains(field.Type, "*vo.") || strings.Contains(field.Type, "vo.") {
		title := strings.TrimPrefix(strings.TrimPrefix(field.Type, "*"), "vo.")
		if schema, ok := schemas[title]; ok {
			props := schema.Properties
			for _, v := range props {
				if v == v3.File {
					return true
				} else if v.Type == v3.ArrayT && v.Items == v3.File {
					return true
				}
			}
		}
	}
	return false
}

func operationOf(method astutils.MethodMeta) v3.Operation {
	var ret v3.Operation
	var params []v3.Parameter
	for _, item := range method.Params {
		pschema := schemaOf(item)
		if pschema == v3.Any {
			continue
		}
		if stringutils.IsEmpty(pschema.Ref) && pschema.Type != v3.ObjectT && pschema.Type != v3.ArrayT {
			params = append(params, v3.Parameter{
				Name:   strcase.ToLowerCamel(item.Name),
				In:     v3.InQuery,
				Schema: pschema,
			})
		} else {
			var content v3.Content
			mt := &v3.MediaType{
				Schema: pschema,
			}
			if strings.HasSuffix(method.Name, "Form") {
				if hasFile(item) {
					reflect.ValueOf(&content).Elem().FieldByName("FormData").Set(reflect.ValueOf(mt))
				} else {
					reflect.ValueOf(&content).Elem().FieldByName("FormUrl").Set(reflect.ValueOf(mt))
				}
			} else {
				reflect.ValueOf(&content).Elem().FieldByName("Json").Set(reflect.ValueOf(mt))
			}
			ret.RequestBody = &v3.RequestBody{
				Content:  &content,
				Required: true,
			}
		}
	}
	ret.Parameters = params
	var respContent v3.Content
	var hasFile bool
	for _, item := range method.Results {
		if item.Type == "*os.File" {
			hasFile = true
			break
		}
	}
	if hasFile {
		respContent.Stream = &v3.MediaType{
			Schema: v3.File,
		}
	} else if len(method.Results) > 0 {
		respContent.Json = &v3.MediaType{
			Schema: schemaOf(method.Results[0]),
		}
	} else {
		respContent.Json = &v3.MediaType{
			Schema: v3.Any,
		}
	}
	ret.Responses = &v3.Responses{
		Resp200: &v3.Response{
			Content: &respContent,
		},
	}
	return ret
}

func pathOf(method astutils.MethodMeta) v3.Path {
	var ret v3.Path
	op := operationOf(method)
	hm := httpMethod(method.Name)
	reflect.ValueOf(&ret).Elem().FieldByName(strings.Title(strings.ToLower(hm))).Set(reflect.ValueOf(&op))
	return ret
}

func pathsOf(ic astutils.InterfaceCollector) map[string]v3.Path {
	if len(ic.Interfaces) <= 0 {
		return nil
	}
	pathmap := make(map[string]v3.Path)
	inter := ic.Interfaces[0]
	for _, method := range inter.Methods {
		v3path := pathOf(method)
		endpoint := fmt.Sprintf("/%s/%s", strings.ToLower(inter.Name), pattern(method.Name))
		pathmap[endpoint] = v3path
	}
	return pathmap
}

// Currently not suport alias type in vo file. TODO
func GenDoc(dir string, ic astutils.InterfaceCollector) {
	var (
		err     error
		svcname string
		docfile string
		vofile  string
		fi      os.FileInfo
		api     v3.Api
		data    []byte
		vos     []v3.Schema
		paths   map[string]v3.Path
	)
	schemas = make(map[string]v3.Schema)
	svcname = ic.Interfaces[0].Name
	docfile = filepath.Join(dir, strings.ToLower(svcname)+"_openapi3.json")
	fi, err = os.Stat(docfile)
	if err != nil && !os.IsNotExist(err) {
		panic(err)
	}
	if fi != nil {
		logrus.Warningln("file " + docfile + " will be overwrited")
	}
	vofile = filepath.Join(dir, "vo/vo.go")
	vos = schemasOf(vofile)
	for _, item := range vos {
		schemas[item.Title] = item
	}
	paths = pathsOf(ic)
	api = v3.Api{
		Openapi: "3.0.2",
		Paths:   paths,
		Components: &v3.Components{
			Schemas: schemas,
		},
	}
	data, err = json.Marshal(api)
	err = ioutil.WriteFile(docfile, data, os.ModePerm)
	if err != nil {
		panic(err)
	}
}
