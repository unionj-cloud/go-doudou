package v3

import (
	"encoding/json"
	"github.com/unionj-cloud/go-doudou/cmd/internal/astutils"
	"github.com/unionj-cloud/go-doudou/toolkit/copier"
	. "github.com/unionj-cloud/go-doudou/toolkit/openapi/v3"
	"github.com/unionj-cloud/go-doudou/toolkit/sliceutils"
	"github.com/unionj-cloud/go-doudou/toolkit/stringutils"
	"regexp"
	"strings"
	"unicode"
)

// Schemas from components of OpenAPI3.0 json document
var Schemas = make(map[string]Schema)

var Enums = make(map[string]astutils.EnumMeta)

// SchemaNames schema names from components of OpenAPI3.0 json document
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
	ft := field.Type
	if IsVarargs(ft) {
		ft = ToSlice(ft)
	}
	ft = strings.TrimLeft(ft, "*")
	switch ft {
	case "int", "int8", "int16", "int32", "uint", "uint8", "uint16", "uint32", "byte", "rune", "complex64", "complex128":
		return Int
	case "int64", "uint64", "uintptr":
		return Int64
	case "bool":
		return Bool
	case "string", "error", "[]rune", "[]byte":
		return String
	case "float32":
		return Float32
	case "float64":
		return Float64
	case "multipart.FileHeader", "v3.FileModel":
		return File
	default:
		return handleDefaultCase(ft)
	}
}

func handleDefaultCase(ft string) *Schema {
	if strings.HasPrefix(ft, "map[") {
		elem := ft[strings.Index(ft, "]")+1:]
		return &Schema{
			Type: ObjectT,
			AdditionalProperties: SchemaOf(astutils.FieldMeta{
				Type: elem,
			}),
		}
	}
	if strings.HasPrefix(ft, "[") {
		elem := ft[strings.Index(ft, "]")+1:]
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
	}
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
		if enumMeta, ok := Enums[title]; ok {
			enumSchema := &Schema{
				Type: StringT,
				Enum: enumMeta.Values,
			}
			if len(enumMeta.Values) > 0 {
				enumSchema.Default = enumMeta.Values[0]
			}
			return enumSchema
		}
	}
	return Any
}

var castFuncMap = map[string]string{
	"bool":          "ToBool",
	"float64":       "ToFloat64",
	"float32":       "ToFloat32",
	"int64":         "ToInt64",
	"int32":         "ToInt32",
	"int16":         "ToInt16",
	"int8":          "ToInt8",
	"int":           "ToInt",
	"uint":          "ToUint",
	"uint8":         "ToUint8",
	"uint16":        "ToUint16",
	"uint32":        "ToUint32",
	"uint64":        "ToUint64",
	"error":         "ToError",
	"[]byte":        "ToByteSlice",
	"[]rune":        "ToRuneSlice",
	"[]interface{}": "ToInterfaceSlice",
	"[]bool":        "ToBoolSlice",
	"[]int":         "ToIntSlice",
	"[]float64":     "ToFloat64Slice",
	"[]float32":     "ToFloat32Slice",
	"[]int64":       "ToInt64Slice",
	"[]int32":       "ToInt32Slice",
	"[]int16":       "ToInt16Slice",
	"[]int8":        "ToInt8Slice",
	"[]uint":        "ToUintSlice",
	"[]uint8":       "ToUint8Slice",
	"[]uint16":      "ToUint16Slice",
	"[]uint32":      "ToUint32Slice",
	"[]uint64":      "ToUint64Slice",
	"[]error":       "ToErrorSlice",
	"[][]byte":      "ToByteSliceSlice",
	"[][]rune":      "ToRuneSliceSlice",
}

func IsSupport(t string) bool {
	if IsVarargs(t) {
		t = ToSlice(t)
	}
	_, exists := castFuncMap[strings.TrimLeft(t, "*")]
	return exists
}

func IsOptional(t string) bool {
	return strings.HasPrefix(t, "*") || strings.HasPrefix(t, "...")
}

func IsSlice(t string) bool {
	return strings.Contains(t, "[") || strings.HasPrefix(t, "...")
}

func IsVarargs(t string) bool {
	return strings.HasPrefix(t, "...")
}

func ToSlice(t string) string {
	return "[]" + strings.TrimPrefix(t, "...")
}

func CastFunc(t string) string {
	if IsVarargs(t) {
		t = ToSlice(t)
	}
	return castFuncMap[strings.TrimLeft(t, "*")]
}

// CopySchema as SchemaOf returns pointer, so deepcopy the schema the pointer points
func CopySchema(field astutils.FieldMeta) Schema {
	var schema Schema
	err := copier.DeepCopy(SchemaOf(field), &schema)
	if err != nil {
		panic(err)
	}
	return schema
}

func RefAddDoc(schema *Schema, doc string) {
	if stringutils.IsNotEmpty(schema.Ref) {
		title := strings.TrimPrefix(schema.Ref, "#/components/schemas/")
		temp := Schemas[title]
		temp.Description = strings.Join([]string{doc, temp.Description}, "\n")
		Schemas[title] = temp
	} else {
		schema.Description = doc
	}
}

// NewSchema new schema from astutils.StructMeta
func NewSchema(structmeta astutils.StructMeta) Schema {
	properties := make(map[string]*Schema)
	var required []string
	for _, field := range structmeta.Fields {
		fschema := CopySchema(field)
		RefAddDoc(&fschema, strings.Join(field.Comments, "\n"))
		properties[field.DocName] = &fschema
		if !strings.HasPrefix(field.Type, "*") {
			required = append(required, field.DocName)
		}
	}
	return Schema{
		Title:       structmeta.Name,
		Type:        ObjectT,
		Properties:  properties,
		Description: strings.Join(structmeta.Comments, "\n"),
		Required:    required,
	}
}

// IsBuiltin check whether field is built-in type https://pkg.go.dev/builtin or not
func IsBuiltin(field astutils.FieldMeta) bool {
	simples := []interface{}{Int, Int64, Bool, String, Float32, Float64}
	types := []interface{}{IntegerT, StringT, BooleanT, NumberT}
	pschema := SchemaOf(field)
	if sliceutils.Contains(simples, pschema) || (sliceutils.Contains(types, pschema.Type) && pschema.Format != BinaryF) {
		return true
	}
	if pschema.Type == ArrayT && (sliceutils.Contains(simples, pschema.Items) || (sliceutils.Contains(types, pschema.Items.Type) && pschema.Items.Format != BinaryF)) {
		return true
	}
	return false
}

// IsEnum check whether field is enum
func IsEnum(field astutils.FieldMeta) bool {
	pschema := SchemaOf(field)
	return len(pschema.Enum) > 0 || (pschema.Type == ArrayT && len(pschema.Items.Enum) > 0)
}

// ElementType get element type string from slice
func ElementType(t string) string {
	if IsVarargs(t) {
		return strings.TrimPrefix(t, "...")
	}
	return t[strings.Index(t, "]")+1:]
}
