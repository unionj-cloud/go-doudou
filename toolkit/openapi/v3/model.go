package v3

import (
	"io"
)

// Contact https://spec.openapis.org/oas/v3.0.3#contact-object
type Contact struct {
	Email string `json:"email,omitempty"`
}

// License https://spec.openapis.org/oas/v3.0.3#license-object
type License struct {
	Name string `json:"name,omitempty"`
	URL  string `json:"url,omitempty"`
}

// Info https://spec.openapis.org/oas/v3.0.3#info-object
type Info struct {
	Title          string   `json:"title,omitempty"`
	Description    string   `json:"description,omitempty"`
	TermsOfService string   `json:"termsOfService,omitempty"`
	Contact        *Contact `json:"contact,omitempty"`
	License        *License `json:"license,omitempty"`
	Version        string   `json:"version,omitempty"`
}

// Server https://spec.openapis.org/oas/v3.0.3#server-object
type Server struct {
	URL string `json:"url,omitempty"`
}

// ExternalDocs https://spec.openapis.org/oas/v3.0.3#external-documentation-object
type ExternalDocs struct {
	Description string `json:"description,omitempty"`
	URL         string `json:"url,omitempty"`
}

// Tag https://spec.openapis.org/oas/v3.0.3#tag-object
type Tag struct {
	Name         string        `json:"name,omitempty"`
	Description  string        `json:"description,omitempty"`
	ExternalDocs *ExternalDocs `json:"externalDocs,omitempty"`
}

// In represents parameter position
type In string

const (
	// InQuery query string parameter
	InQuery In = "query"
	// InPath TODO not implemented yet
	InPath In = "path"
	// InHeader TODO not implemented yet
	InHeader In = "header"
	// InCookie TODO not implemented yet
	InCookie In = "cookie"
)

// Example https://spec.openapis.org/oas/v3.0.3#example-object
type Example struct {
	// TODO
}

// Encoding https://spec.openapis.org/oas/v3.0.3#encoding-object
type Encoding struct {
	// TODO
}

// MediaType https://spec.openapis.org/oas/v3.0.3#media-type-object
type MediaType struct {
	Schema   *Schema             `json:"schema,omitempty"`
	Example  interface{}         `json:"example,omitempty"`
	Examples map[string]Example  `json:"examples,omitempty"`
	Encoding map[string]Encoding `json:"encoding,omitempty"`
}

// Content REQUIRED. The content of the request body. The key is a media type or [media type range]appendix-D)
// and the value describes it. For requests that match multiple keys, only the most specific key is applicable.
// e.g. text/plain overrides text/*
type Content struct {
	TextPlain *MediaType `json:"text/plain,omitempty"`
	JSON      *MediaType `json:"application/json,omitempty"`
	FormURL   *MediaType `json:"application/x-www-form-urlencoded,omitempty"`
	Stream    *MediaType `json:"application/octet-stream,omitempty"`
	FormData  *MediaType `json:"multipart/form-data,omitempty"`
	Default   *MediaType `json:"*/*,omitempty"`
}

// Parameter https://spec.openapis.org/oas/v3.0.3#parameter-object
type Parameter struct {
	Name            string      `json:"name,omitempty"`
	In              In          `json:"in,omitempty"`
	Description     string      `json:"description,omitempty"`
	Required        bool        `json:"required,omitempty"`
	Deprecated      bool        `json:"deprecated,omitempty"`
	Example         interface{} `json:"example,omitempty"`
	Schema          *Schema     `json:"schema,omitempty"`
	Style           string      `json:"style,omitempty"`
	Explode         bool        `json:"explode,omitempty"`
	AllowReserved   bool        `json:"allowReserved,omitempty"`
	Content         *Content    `json:"content,omitempty"`
	AllowEmptyValue bool        `json:"allowEmptyValue,omitempty"`
}

// RequestBody https://spec.openapis.org/oas/v3.0.3#request-body-object
type RequestBody struct {
	Description string   `json:"description,omitempty"`
	Content     *Content `json:"content,omitempty"`
	Required    bool     `json:"required,omitempty"`
	Ref         string   `json:"$ref,omitempty"`
}

// Header https://spec.openapis.org/oas/v3.0.3#header-object
type Header struct {
	Ref         string      `json:"$ref,omitempty"`
	Description string      `json:"description,omitempty"`
	Required    bool        `json:"required,omitempty"`
	Deprecated  bool        `json:"deprecated,omitempty"`
	Example     interface{} `json:"example,omitempty"`
	Schema      *Schema     `json:"schema,omitempty"`
}

// Link https://spec.openapis.org/oas/v3.0.3#link-object
type Link struct {
	// TODO
}

// Response https://spec.openapis.org/oas/v3.0.3#response-object
type Response struct {
	Description string   `json:"description"`
	Content     *Content `json:"content,omitempty"`
	// TODO
	Headers map[string]Header `json:"headers,omitempty"`
	Links   map[string]Link   `json:"links,omitempty"`
	Ref     string            `json:"$ref,omitempty"`
}

// Responses https://spec.openapis.org/oas/v3.0.3#responses-object
type Responses struct {
	Resp200 *Response `json:"200,omitempty"`
	Resp400 *Response `json:"400,omitempty"`
	Resp401 *Response `json:"401,omitempty"`
	Resp403 *Response `json:"403,omitempty"`
	Resp404 *Response `json:"404,omitempty"`
	Resp405 *Response `json:"405,omitempty"`
	Default *Response `json:"default,omitempty"`
}

// Callback https://spec.openapis.org/oas/v3.0.3#callback-object
type Callback struct {
	// TODO
}

// Security not implemented yet
type Security struct {
	// TODO
}

// Operation https://spec.openapis.org/oas/v3.0.3#operation-object
type Operation struct {
	Tags         []string            `json:"tags,omitempty"`
	Summary      string              `json:"summary,omitempty"`
	Description  string              `json:"description,omitempty"`
	OperationID  string              `json:"operationId,omitempty"`
	Parameters   []Parameter         `json:"parameters,omitempty"`
	RequestBody  *RequestBody        `json:"requestBody,omitempty"`
	Responses    *Responses          `json:"responses,omitempty"`
	Deprecated   bool                `json:"deprecated,omitempty"`
	ExternalDocs *ExternalDocs       `json:"externalDocs,omitempty"`
	Callbacks    map[string]Callback `json:"callbacks,omitempty"`
	Security     []Security          `json:"security,omitempty"`
	Servers      []Server            `json:"servers,omitempty"`
}

// Path https://spec.openapis.org/oas/v3.0.3#path-item-object
type Path struct {
	Get    *Operation `json:"get,omitempty"`
	Post   *Operation `json:"post,omitempty"`
	Put    *Operation `json:"put,omitempty"`
	Delete *Operation `json:"delete,omitempty"`
	// TODO
	Parameters []Parameter `json:"parameters,omitempty"`
}

// SecurityScheme https://spec.openapis.org/oas/v3.0.3#security-scheme-object
type SecurityScheme struct {
	// TODO
}

// Discriminator https://spec.openapis.org/oas/v3.0.3#discriminator-object
type Discriminator struct {
	PropertyName string            `json:"propertyName,omitempty"`
	Mapping      map[string]string `json:"mapping,omitempty"`
}

// Schema https://spec.openapis.org/oas/v3.0.3#schema-object
type Schema struct {
	Ref              string             `json:"$ref,omitempty"`
	Title            string             `json:"title,omitempty"`
	Type             Type               `json:"type,omitempty"`
	Properties       map[string]*Schema `json:"properties,omitempty"`
	Format           Format             `json:"format,omitempty"`
	Items            *Schema            `json:"items,omitempty"`
	Description      string             `json:"description,omitempty"`
	Default          interface{}        `json:"default,omitempty"`
	Example          interface{}        `json:"example,omitempty"`
	Deprecated       bool               `json:"deprecated,omitempty"`
	Discriminator    *Discriminator     `json:"discriminator,omitempty"`
	Nullable         bool               `json:"nullable,omitempty"`
	Maximum          interface{}        `json:"maximum,omitempty"`
	Minimum          interface{}        `json:"minimum,omitempty"`
	ExclusiveMaximum interface{}        `json:"exclusiveMaximum,omitempty"`
	ExclusiveMinimum interface{}        `json:"exclusiveMinimum,omitempty"`
	MaxLength        int                `json:"maxLength,omitempty"`
	MinLength        int                `json:"minLength,omitempty"`
	Required         []string           `json:"required,omitempty"`
	Enum             []interface{}      `json:"enum,omitempty"`
	AllOf            []*Schema          `json:"allOf,omitempty"`
	OneOf            []*Schema          `json:"oneOf,omitempty"`
	AnyOf            []*Schema          `json:"anyOf,omitempty"`
	Not              []*Schema          `json:"not,omitempty"`
	// AdditionalProperties *Schema or bool
	AdditionalProperties interface{} `json:"additionalProperties,omitempty"`
	Pattern              interface{} `json:"pattern,omitempty"`
	XMapType             string      `json:"x-map-type,omitempty"`
}

// Components https://spec.openapis.org/oas/v3.0.3#components-object
type Components struct {
	Schemas       map[string]Schema      `json:"schemas,omitempty"`
	RequestBodies map[string]RequestBody `json:"requestBodies,omitempty"`
	Responses     map[string]Response    `json:"responses,omitempty"`
	// TODO
	Parameters map[string]Parameter `json:"parameters,omitempty"`
	// TODO
	Examples map[string]Example `json:"examples,omitempty"`
	// TODO
	Headers map[string]Header `json:"headers,omitempty"`
	// TODO
	SecuritySchemes map[string]SecurityScheme `json:"securitySchemes,omitempty"`
	// TODO
	Links map[string]Link `json:"links,omitempty"`
	// TODO
	Callbacks map[string]Callback `json:"callbacks,omitempty"`
}

// API https://spec.openapis.org/oas/v3.0.3#openapi-object
type API struct {
	Openapi      string          `json:"openapi,omitempty"`
	Info         *Info           `json:"info,omitempty"`
	Servers      []Server        `json:"servers,omitempty"`
	Tags         []Tag           `json:"tags,omitempty"`
	Paths        map[string]Path `json:"paths,omitempty"`
	Components   *Components     `json:"components,omitempty"`
	ExternalDocs *ExternalDocs   `json:"externalDocs,omitempty"`
}

// Type represents types in OpenAPI3.0 spec
type Type string

const (
	// IntegerT integer
	IntegerT Type = "integer"
	// StringT string
	StringT Type = "string"
	// BooleanT boolean
	BooleanT Type = "boolean"
	// NumberT number
	NumberT Type = "number"
	// ObjectT object
	ObjectT Type = "object"
	// ArrayT array
	ArrayT Type = "array"
)

// Format represents format in OpenAPI3.0 spec
type Format string

const (
	// Int32F int32
	Int32F Format = "int32"
	// Int64F int64
	Int64F Format = "int64"
	// FloatF float
	FloatF Format = "float"
	// DoubleF double
	DoubleF Format = "double"
	// DateTimeF date-time
	DateTimeF Format = "date-time"
	// BinaryF binary
	BinaryF  Format = "binary"
	DecimalF Format = "decimal"
)

var (
	// Any constant schema for object
	Any = &Schema{
		Type: ObjectT,
	}
	// Int constant schema for int
	Int = &Schema{
		Type:   IntegerT,
		Format: Int32F,
	}
	// Int64 constant schema for int64
	Int64 = &Schema{
		Type:   IntegerT,
		Format: Int64F,
	}
	// String constant schema for string
	String = &Schema{
		Type: StringT,
	}
	// Time constant schema for time
	Time = &Schema{
		Type:   StringT,
		Format: DateTimeF,
	}
	// Bool constant schema for bool
	Bool = &Schema{
		Type: BooleanT,
	}
	// Float32 constant schema for float32
	Float32 = &Schema{
		Type:   NumberT,
		Format: FloatF,
	}
	// Float64 constant schema for float64
	Float64 = &Schema{
		Type:   NumberT,
		Format: DoubleF,
	}
	// File constant schema for file
	File = &Schema{
		Type:   StringT,
		Format: BinaryF,
	}
	// FileArray constant schema for file slice
	FileArray = &Schema{
		Type:  ArrayT,
		Items: File,
	}
	Decimal = &Schema{
		Type:   StringT,
		Format: DecimalF,
	}
)

type FileModel struct {
	Filename string
	Reader   io.ReadCloser
}

func (f *FileModel) Close() error {
	return f.Reader.Close()
}

type IEnum interface {
	StringSetter(value string)
	StringGetter() string
	UnmarshalJSON(bytes []byte) error
	MarshalJSON() ([]byte, error)
}

var IEnumMethods = []string{
	"func StringSetter(value string)",
	"func StringGetter() string",
	"func UnmarshalJSON(bytes []byte) error",
	"func MarshalJSON() ([]byte, error)",
}

type ExampleType int

const (
	UNKNOWN_EXAMPLE ExampleType = iota
	JSON_EXAMPLE
	TEXT_EXAMPLE
)
