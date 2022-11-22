package v3

import (
	"encoding/json"
	"fmt"
	"github.com/iancoleman/strcase"
	"github.com/unionj-cloud/go-doudou/v2/cmd/internal/astutils"
	"github.com/unionj-cloud/go-doudou/v2/toolkit/sliceutils"
	"github.com/unionj-cloud/go-doudou/v2/toolkit/stringutils"
	"reflect"
	"regexp"
	"strings"
	"unicode"
)

var _ ProtobufType = (*Enum)(nil)
var _ ProtobufType = (*Message)(nil)

type ProtobufType interface {
	GetName() string
	String() string
	Inner() bool
}

var MessageStore = make(map[string]Message)

var EnumStore = make(map[string]Enum)

var ImportStore = make(map[string]struct{})

var MessageNames []string

type EnumField struct {
	Name   string
	Number int
}

func (receiver ProtoGenerator) newEnumField(field string, index int) EnumField {
	return EnumField{
		Name:   strings.ToUpper(strcase.ToSnake(field)),
		Number: index,
	}
}

type Enum struct {
	Name   string
	Fields []EnumField
}

func (e Enum) Inner() bool {
	return false
}

func (e Enum) String() string {
	return e.Name
}

func (e Enum) GetName() string {
	return e.Name
}

func (receiver ProtoGenerator) NewEnum(enumMeta astutils.EnumMeta) Enum {
	var fields []EnumField
	for i, field := range enumMeta.Values {
		fields = append(fields, receiver.newEnumField(field, i))
	}
	return Enum{
		Name:   strcase.ToCamel(enumMeta.Name),
		Fields: fields,
	}
}

// Message represents protobuf message definition
type Message struct {
	Name       string
	Fields     []Field
	Comments   []string
	IsInner    bool
	IsScalar   bool
	IsMap      bool
	IsRepeated bool
	IsTopLevel bool
	// IsImported denotes the message will be imported from third-party, such as from google/protobuf
	IsImported bool
}

func (m Message) Inner() bool {
	return m.IsInner
}

func (m Message) GetName() string {
	return m.Name
}

func (m Message) String() string {
	switch {
	case reflect.DeepEqual(m, Any):
		return "anypb.Any"
	case reflect.DeepEqual(m, Empty):
		return "emptypb.Empty"
	default:
		return m.Name
	}
}

// NewMessage returns message instance from astutils.StructMeta
func (receiver ProtoGenerator) NewMessage(structmeta astutils.StructMeta) Message {
	var fields []Field
	for i, field := range structmeta.Fields {
		fields = append(fields, receiver.newField(field, i+1))
	}
	return Message{
		Name:       strcase.ToCamel(structmeta.Name),
		Fields:     fields,
		Comments:   structmeta.Comments,
		IsTopLevel: true,
	}
}

// Field represents protobuf message field definition
type Field struct {
	Name     string
	Type     ProtobufType
	Number   int
	Comments []string
	JsonName string
}

func (receiver ProtoGenerator) newField(field astutils.FieldMeta, index int) Field {
	t := receiver.MessageOf(field.Type)
	if t.Inner() {
		message := t.(Message)
		message.Name = strcase.ToCamel(field.Name)
		t = message
	}
	fieldName := receiver.fieldNamingFunc(field.Name)
	return Field{
		Name:     fieldName,
		Type:     t,
		Number:   index,
		Comments: field.Comments,
		JsonName: fieldName,
	}
}

var (
	Double = Message{
		Name:     "double",
		IsScalar: true,
	}
	Float = Message{
		Name:     "float",
		IsScalar: true,
	}
	Int32 = Message{
		Name:     "int32",
		IsScalar: true,
	}
	Int64 = Message{
		Name:     "int64",
		IsScalar: true,
	}
	Uint32 = Message{
		Name:     "uint32",
		IsScalar: true,
	}
	Uint64 = Message{
		Name:     "uint64",
		IsScalar: true,
	}
	Bool = Message{
		Name:     "bool",
		IsScalar: true,
	}
	String = Message{
		Name:     "string",
		IsScalar: true,
	}
	Bytes = Message{
		Name:     "bytes",
		IsScalar: true,
	}
	Any = Message{
		Name:       "google.protobuf.Any",
		IsTopLevel: true,
		IsImported: true,
	}
	Empty = Message{
		Name:       "google.protobuf.Empty",
		IsTopLevel: true,
		IsImported: true,
	}
)

func (receiver ProtoGenerator) MessageOf(ft string) ProtobufType {
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
		return receiver.handleDefaultCase(ft)
	}
}

var anonystructre *regexp.Regexp

func init() {
	anonystructre = regexp.MustCompile(`anonystruct«(.*)»`)
}

func (receiver ProtoGenerator) handleDefaultCase(ft string) ProtobufType {
	if strings.HasPrefix(ft, "map[") {
		elem := ft[strings.Index(ft, "]")+1:]
		key := ft[4:strings.Index(ft, "]")]
		keyMessage := receiver.MessageOf(key)
		if reflect.DeepEqual(keyMessage, Float) || reflect.DeepEqual(keyMessage, Double) || reflect.DeepEqual(keyMessage, Bytes) {
			panic("floating point types and bytes cannot be key_type of maps, please refer to https://developers.google.com/protocol-buffers/docs/proto3#maps")
		}
		elemMessage := receiver.MessageOf(elem)
		if strings.HasPrefix(elemMessage.GetName(), "map<") {
			panic("the value_type cannot be another map, please refer to https://developers.google.com/protocol-buffers/docs/proto3#maps")
		}
		return Message{
			Name:  fmt.Sprintf("map<%s, %s>", keyMessage.GetName(), elemMessage.GetName()),
			IsMap: true,
		}
	}
	if strings.HasPrefix(ft, "[") {
		elem := ft[strings.Index(ft, "]")+1:]
		elemMessage := receiver.MessageOf(elem)
		if strings.HasPrefix(elemMessage.GetName(), "map<") {
			panic("map fields cannot be repeated, please refer to https://developers.google.com/protocol-buffers/docs/proto3#maps")
		}
		messageName := elemMessage.GetName()
		if strings.Contains(elemMessage.GetName(), "repeated ") {
			messageName = messageName[strings.LastIndex(messageName, ".")+1:]
			messageName = "Nested" + strcase.ToCamel(messageName)
			fieldName := receiver.fieldNamingFunc(messageName)
			MessageStore[messageName] = Message{
				Name: messageName,
				Fields: []Field{
					{
						Name:     fieldName,
						Type:     elemMessage,
						Number:   1,
						JsonName: fieldName,
					},
				},
				IsInner: true,
			}
		}
		return Message{
			Name:       fmt.Sprintf("repeated %s", messageName),
			IsRepeated: true,
		}
	}
	if anonystructre.MatchString(ft) {
		result := anonystructre.FindStringSubmatch(ft)
		var structmeta astutils.StructMeta
		json.Unmarshal([]byte(result[1]), &structmeta)
		message := receiver.NewMessage(structmeta)
		message.IsInner = true
		message.IsTopLevel = false
		return message
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
			if sliceutils.StringContains(MessageNames, title) {
				return Message{
					Name:       strcase.ToCamel(title),
					IsTopLevel: true,
				}
			}
		}
		if e, ok := EnumStore[title]; ok {
			return e
		}
	}
	ImportStore["google/protobuf/any.proto"] = struct{}{}
	return Any
}
