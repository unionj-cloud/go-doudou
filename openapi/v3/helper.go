package v3

import (
	"encoding/json"
	"github.com/unionj-cloud/go-doudou/astutils"
	"github.com/unionj-cloud/go-doudou/copier"
	"github.com/unionj-cloud/go-doudou/sliceutils"
	"github.com/unionj-cloud/go-doudou/stringutils"
	"regexp"
	"strings"
	"unicode"
)

var Schemas map[string]Schema
var SchemaNames []string

// SchemaOf reference https://golang.org/pkg/builtin/
// type bool
// type byte
// type complex128
// type complex64
// type error
// type float32
// type float64
// type int
// type int16
// type int32
// type int64
// type int8
// type rune
// type string
// type uint
// type uint16
// type uint32
// type uint64
// type uint8
// type uintptr
func SchemaOf(field astutils.FieldMeta) *Schema {
	ft := strings.TrimPrefix(field.Type, "*")
	switch ft {
	case "int", "int8", "int16", "int32", "uint", "uint8", "uint16", "uint32", "byte", "rune", "complex64", "complex128":
		return Int
	case "int64", "uint64", "uintptr":
		return Int64
	case "bool":
		return Bool
	case "string", "error":
		return String
	case "float32":
		return Float32
	case "float64":
		return Float64
	case "multipart.FileHeader":
		return File
	default:
		if strings.HasPrefix(ft, "map[") {
			elem := ft[strings.Index(ft, "]")+1:]
			elem = strings.TrimPrefix(elem, "*")
			return &Schema{
				Type: ObjectT,
				AdditionalProperties: SchemaOf(astutils.FieldMeta{
					Type: elem,
				}),
			}
		}
		if strings.HasPrefix(ft, "[") {
			elem := ft[strings.Index(ft, "]")+1:]
			elem = strings.TrimPrefix(elem, "*")
			return &Schema{
				Type: ArrayT,
				Items: SchemaOf(astutils.FieldMeta{
					Type: elem,
				}),
			}
		}
		re := regexp.MustCompile(`anonystruct«(.*)»`)
		if re.MatchString(ft) {
			result := re.FindStringSubmatch(ft)
			var structmeta astutils.StructMeta
			json.Unmarshal([]byte(result[1]), &structmeta)
			schema := NewSchema(structmeta)
			return &schema
		} else {
			var title string
			if !strings.Contains(ft, ".") {
				title = ft
			}
			if stringutils.IsEmpty(title) {
				title = ft[strings.LastIndex(ft, ".")+1:]
			}
			if stringutils.IsNotEmpty(title) {
				if unicode.IsUpper(rune(title[0])) {
					if sliceutils.StringContains(SchemaNames, title) {
						return &Schema{
							Ref: "#/components/schemas/" + title,
						}
					}
				}
			}
		}
		return Any
	}
}

func CopySchema(field astutils.FieldMeta) Schema {
	var schema Schema
	err := copier.DeepCopy(SchemaOf(field), &schema)
	if err != nil {
		panic(err)
	}
	return schema
}

func NewSchema(structmeta astutils.StructMeta) Schema {
	properties := make(map[string]*Schema)
	for _, field := range structmeta.Fields {
		fschema := CopySchema(field)
		fschema.Description = strings.Join(field.Comments, "\n")
		properties[field.DocName] = &fschema
	}
	return Schema{
		Title:       structmeta.Name,
		Type:        ObjectT,
		Properties:  properties,
		Description: strings.Join(structmeta.Comments, "\n"),
	}
}

func IsBuiltin(field astutils.FieldMeta) bool {
	simples := []interface{}{Int, Int64, Bool, String, Float32, Float64}
	pschema := SchemaOf(field)
	if pschema == nil {
		return false
	}
	return sliceutils.Contains(simples, pschema) || (pschema.Type == ArrayT && sliceutils.Contains(simples, pschema.Items))
}
