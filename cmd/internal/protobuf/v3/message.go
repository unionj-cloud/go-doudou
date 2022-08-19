package v3

import (
	"encoding/json"
	"fmt"
	"github.com/unionj-cloud/go-doudou/cmd/internal/astutils"
	"github.com/unionj-cloud/go-doudou/toolkit/stringutils"
	"regexp"
	"strings"
	"unicode"
)

var MessageStore = make(map[string]*Message)

var EnumStore = make(map[string]*Message)

var importStore = make(map[string]struct{})

// Message represents protobuf message definition
type Message struct {
	Name     string
	Fields   []*Field
	Comments []string
}

func (m *Message) String() string {
	return m.Name
}

// NewMessage returns message instance from astutils.StructMeta
func NewMessage(structmeta astutils.StructMeta) *Message {
	var fields []*Field
	for i, field := range structmeta.Fields {
		fields = append(fields, NewField(field, i+1))
	}
	return &Message{
		Name:     structmeta.Name,
		Fields:   fields,
		Comments: structmeta.Comments,
	}
}

// Field represents protobuf message field definition
type Field struct {
	Name     string
	Type     *Message
	Number   int
	Comments []string
	JsonName string
}

func NewField(field astutils.FieldMeta, index int) *Field {
	return &Field{
		Name:     field.Name,
		Type:     MessageOf(field.Type),
		Number:   index,
		Comments: field.Comments,
		JsonName: field.DocName,
	}
}

var (
	Double = &Message{
		Name: "double",
	}
	Float = &Message{
		Name: "float",
	}
	Int32 = &Message{
		Name: "int32",
	}
	Int64 = &Message{
		Name: "int64",
	}
	Uint32 = &Message{
		Name: "uint32",
	}
	Uint64 = &Message{
		Name: "uint64",
	}
	Bool = &Message{
		Name: "bool",
	}
	String = &Message{
		Name: "string",
	}
	Bytes = &Message{
		Name: "bytes",
	}
	Any = &Message{
		Name: "google.protobuf.Any",
	}
)

func MessageOf(ft string) *Message {
	if astutils.IsVarargs(ft) {
		ft = astutils.ToSlice(ft)
	}
	ft = strings.TrimLeft(ft, "*")
	switch ft {
	case "int", "int8", "int16", "int32", "byte", "rune", "complex64", "complex128":
		return Int32
	case "uint", "uint8", "uint16", "uint32":
		return Uint32
	case "int64":
		return Int64
	case "uint64", "uintptr":
		return Uint64
	case "bool":
		return Bool
	case "string", "error", "[]rune":
		return String
	case "[]byte", "v3.FileModel", "os.File":
		return Bytes
	case "float32":
		return Float
	case "float64":
		return Double
	default:
		return handleDefaultCase(ft)
	}
}

func handleDefaultCase(ft string) *Message {
	if strings.HasPrefix(ft, "map[") {
		elem := ft[strings.Index(ft, "]")+1:]
		key := ft[4:strings.Index(ft, "]")]
		keyMessage := MessageOf(key)
		if keyMessage == Float || keyMessage == Double || keyMessage == Bytes {
			panic("floating point types and bytes cannot be key_type of maps, please refer to https://developers.google.com/protocol-buffers/docs/proto3#maps")
		}
		elemMessage := MessageOf(elem)
		if strings.HasPrefix(elemMessage.Name, "map<") {
			panic("the value_type cannot be another map, please refer to https://developers.google.com/protocol-buffers/docs/proto3#maps")
		}
		return &Message{
			Name: fmt.Sprintf("map<%s, %s>", keyMessage, elemMessage),
		}
	}
	if strings.HasPrefix(ft, "[") {
		elem := ft[strings.Index(ft, "]")+1:]
		elemMessage := MessageOf(elem)
		if strings.HasPrefix(elemMessage.Name, "map<") {
			panic("map fields cannot be repeated, please refer to https://developers.google.com/protocol-buffers/docs/proto3#maps")
		}
		return &Message{
			Name: fmt.Sprintf("repeated %s", elemMessage),
		}
	}
	re := regexp.MustCompile(`anonystruct«(.*)»`)
	if re.MatchString(ft) {
		result := re.FindStringSubmatch(ft)
		var structmeta astutils.StructMeta
		json.Unmarshal([]byte(result[1]), &structmeta)
		message := NewMessage(structmeta)
		return message
	}
	var title string
	if !strings.Contains(ft, ".") {
		title = ft
	}
	if stringutils.IsEmpty(title) {
		title = ft[strings.LastIndex(ft, ".")+1:]
	}
	if stringutils.IsNotEmpty(title) || unicode.IsUpper(rune(title[0])) {
		if m, ok := MessageStore[title]; ok {
			return m
		}
		if m, ok := EnumStore[title]; ok {
			return m
		}
	}
	importStore["google/protobuf/any.proto"] = struct{}{}
	return Any
}
