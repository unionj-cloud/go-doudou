package codegen

import (
	"encoding/json"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/unionj-cloud/go-doudou/astutils"
	v3 "github.com/unionj-cloud/go-doudou/openapi/v3"
)

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
		return &v3.Int
	case "int64", "uint64", "uintptr":
		return &v3.Int64
	case "bool":
		return &v3.Bool
	case "string":
		return &v3.String
	case "float32":
		return &v3.Float32
	case "float64":
		return &v3.Float64
	case "complex64", "complex128":
		return &v3.Any
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
		if strings.HasPrefix(field.Type, "vo.") {
			return &v3.Schema{
				Ref: "#/components/schemas/" + strings.TrimPrefix(field.Type, "vo."),
			}
		}
		return &v3.Any
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

func GenDoc(dir string, ic astutils.InterfaceCollector) {
	var (
		err     error
		svcname string
		docfile string
		fi      os.FileInfo
		api     v3.Api
		data    []byte
	)
	svcname = ic.Interfaces[0].Name
	docfile = filepath.Join(dir, strings.ToLower(svcname)+"_openapi3.json")
	fi, err = os.Stat(docfile)
	if err != nil && !os.IsNotExist(err) {
		panic(err)
	}
	if fi != nil {
		logrus.Warningln("file " + docfile + " will be overwrited")
	}
	api = v3.Api{
		Openapi:      "3.0.2",
		Info:         v3.Info{},
		Servers:      nil,
		Tags:         nil,
		Paths:        nil,
		Components:   v3.Components{},
		ExternalDocs: v3.ExternalDocs{},
	}
	data, err = json.Marshal(api)
	err = ioutil.WriteFile(docfile, data, os.ModePerm)
	if err != nil {
		panic(err)
	}
}
