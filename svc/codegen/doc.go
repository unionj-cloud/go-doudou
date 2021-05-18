package codegen

import (
	"encoding/json"
	"fmt"
	"github.com/iancoleman/strcase"
	"github.com/sirupsen/logrus"
	"github.com/unionj-cloud/go-doudou/astutils"
	v3 "github.com/unionj-cloud/go-doudou/openapi/v3"
	"github.com/unionj-cloud/go-doudou/sliceutils"
	"github.com/unionj-cloud/go-doudou/stringutils"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strings"
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
	switch field.Type {
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
	case "*multipart.FileHeader":
		return v3.File
	default:
		if strings.HasPrefix(field.Type, "map[") {
			elem := field.Type[strings.Index(field.Type, "]")+1:]
			elem = strings.TrimPrefix(elem, "*")
			return &v3.Schema{
				Type: v3.ObjectT,
				AdditionalProperties: schemaOf(astutils.FieldMeta{
					Type: elem,
				}),
			}
		}
		if strings.HasPrefix(field.Type, "[") {
			elem := field.Type[strings.Index(field.Type, "]")+1:]
			elem = strings.TrimPrefix(elem, "*")
			return &v3.Schema{
				Type: v3.ArrayT,
				Items: schemaOf(astutils.FieldMeta{
					Type: elem,
				}),
			}
		}
		if !strings.Contains(field.Type, ".") {
			return &v3.Schema{
				Ref: "#/components/schemas/" + strings.TrimPrefix(field.Type, "*"),
			}
		}
		if strings.Contains(field.Type, "*vo.") || strings.Contains(field.Type, "vo.") {
			title := strings.TrimPrefix(strings.TrimPrefix(field.Type, "*"), "vo.")
			return &v3.Schema{
				Ref: "#/components/schemas/" + title,
			}
		}
		return v3.Any
	}
}

func possibleSchemas(vofile string) []v3.Schema {
	fset := token.NewFileSet()
	root, err := parser.ParseFile(fset, vofile, nil, 0)
	if err != nil {
		panic(err)
	}
	var sc astutils.StructCollector
	ast.Walk(&sc, root)
	var schemas []v3.Schema
	for _, item := range sc.Structs {
		properties := make(map[string]*v3.Schema)
		for _, field := range item.Fields {
			properties[field.Name] = schemaOf(field)
		}
		schemas = append(schemas, v3.Schema{
			Title:      item.Name,
			Type:       v3.ObjectT,
			Properties: properties,
		})
	}
	return schemas
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
		if stringutils.IsEmpty(pschema.Ref) && pschema.Type != v3.ObjectT && pschema.Type != v3.ArrayT {
			params = append(params, v3.Parameter{
				Name:   strcase.ToLowerCamel(item.Name),
				In:     v3.InQuery,
				Schema: *pschema,
			})
		} else {
			var content v3.Content
			mt := v3.MediaType{
				Schema: *pschema,
			}
			if strings.HasSuffix(method.Name, "Form") {
				if hasFile(item) {
					reflect.ValueOf(content).FieldByName("FormData").Set(reflect.ValueOf(mt))
				} else {
					reflect.ValueOf(content).FieldByName("FormUrl").Set(reflect.ValueOf(mt))
				}
			} else {
				reflect.ValueOf(content).FieldByName("Json").Set(reflect.ValueOf(mt))
			}
			ret.RequestBody = v3.RequestBody{
				Content:  content,
				Required: true,
			}
		}
	}
	ret.Parameters = params
	// responseBody TODO
	return ret
}

func pathOf(method astutils.MethodMeta, svcname string) v3.Path {
	ret := v3.Path{
		Endpoint: fmt.Sprintf("/%s/%s", strings.ToLower(svcname), pattern(method.Name)),
	}
	op := operationOf(method)
	hm := httpMethod(method.Name)
	reflect.ValueOf(ret).FieldByName(strings.Title(strings.ToLower(hm))).Set(reflect.ValueOf(op))
	return ret
}

func pathsOf(ic astutils.InterfaceCollector) map[string]v3.Path {
	if len(ic.Interfaces) <= 0 {
		return nil
	}
	pathmap := make(map[string]v3.Path)
	inter := ic.Interfaces[0]
	for _, method := range inter.Methods {
		v3path := pathOf(method, inter.Name)
		pathmap[v3path.Endpoint] = v3path
	}
	return pathmap
}

func GenDoc(dir string, ic astutils.InterfaceCollector) {
	var (
		err       error
		svcname   string
		docfile   string
		vofile    string
		fi        os.FileInfo
		api       v3.Api
		data      []byte
		possibles []v3.Schema
		vos       []string
		paths     map[string]v3.Path
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
	possibles = possibleSchemas(vofile)
	vos = vosOf(ic)
	for _, item := range possibles {
		if sliceutils.StringContains(vos, item.Title) {
			schemas[item.Title] = item
		}
	}
	paths = pathsOf(ic)
	api = v3.Api{
		Openapi: "3.0.2",
		Paths:   paths,
		Components: v3.Components{
			Schemas: schemas,
		},
	}
	data, err = json.Marshal(api)
	err = ioutil.WriteFile(docfile, data, os.ModePerm)
	if err != nil {
		panic(err)
	}
}
