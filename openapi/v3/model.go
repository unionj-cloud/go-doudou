package v3

type Contact struct {
	Email string `json:"email"`
}

type License struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type Info struct {
	Title          string  `json:"title"`
	Description    string  `json:"description"`
	TermsOfService string  `json:"termsOfService"`
	Contact        Contact `json:"contact"`
	License        License `json:"license"`
	Version        string  `json:"version"`
}

type Server struct {
	Url string `json:"url"`
}

type ExternalDocs struct {
	Description string `json:"description"`
	Url         string `json:"url"`
}

type Tag struct {
	Name         string       `json:"name"`
	Description  string       `json:"description"`
	ExternalDocs ExternalDocs `json:"externalDocs"`
}

type In string

const (
	InQuery  In = "query"
	InHeader In = "header"
	InPath   In = "path"
)

type Example struct {
	// TODO
}

type Encoding struct {
	// TODO
}

type MediaType struct {
	Schema   Schema              `json:"schema"`
	Example  interface{}         `json:"example"`
	Examples map[string]Example  `json:"examples"`
	Encoding map[string]Encoding `json:"encoding"`
}

type Content struct {
	TextPlain MediaType `json:"text/plain"`
	Json      MediaType `json:"application/json"`
	FormUrl   MediaType `json:"application/x-www-form-urlencoded"`
	Stream    MediaType `json:"application/octet-stream"`
	FormData  MediaType `json:"multipart/form-data"`
}

type Parameter struct {
	Ref           string      `json:"$ref"`
	Name          string      `json:"name"`
	In            In          `json:"in"`
	Description   string      `json:"description"`
	Required      bool        `json:"required"`
	Deprecated    bool        `json:"deprecated"`
	Example       interface{} `json:"example"`
	Schema        Schema      `json:"schema"`
	Style         string      `json:"style"`
	Explode       bool        `json:"explode"`
	AllowReserved bool        `json:"allowReserved"`
	Content       Content     `json:"content"`
}

type RequestBody struct {
	Description string  `json:"description"`
	Content     Content `json:"content"`
	Required    bool    `json:"required"`
}

type Header struct {
	Ref         string      `json:"$ref"`
	Description string      `json:"description"`
	Required    bool        `json:"required"`
	Deprecated  bool        `json:"deprecated"`
	Example     interface{} `json:"example"`
	Schema      Schema      `json:"schema"`
}

type Link struct {
	// TODO
}

type Response struct {
	Description string            `json:"description"`
	Content     Content           `json:"content"`
	Headers     map[string]Header `json:"headers"`
	Links       map[string]Link   `json:"links"`
}

type Responses struct {
	Resp200 Response `json:"200"`
	Resp400 Response `json:"400"`
	Resp401 Response `json:"401"`
	Resp403 Response `json:"403"`
	Resp404 Response `json:"404"`
	Resp405 Response `json:"405"`
	Default Response `json:"default"`
}

type Callback struct {
	// TODO
}

type Security struct {
	// TODO
}

type Operation struct {
	Tags         []string            `json:"tags"`
	Summary      string              `json:"summary"`
	Description  string              `json:"description"`
	OperationId  string              `json:"operationId"`
	Parameters   []Parameter         `json:"parameters"`
	RequestBody  RequestBody         `json:"requestBody"`
	Responses    Responses           `json:"responses"`
	Deprecated   bool                `json:"deprecated"`
	ExternalDocs ExternalDocs        `json:"externalDocs"`
	Callbacks    map[string]Callback `json:"callbacks"`
	Security     []Security          `json:"security"`
	Servers      []Server            `json:"servers"`
}

type Path struct {
	Endpoint   string      `json:"endpoint"`
	Get        Operation   `json:"get"`
	Post       Operation   `json:"post"`
	Put        Operation   `json:"put"`
	Delete     Operation   `json:"delete"`
	Parameters []Parameter `json:"parameters"`
}

type SecurityScheme struct {
	// TODO
}

type Discriminator struct {
	PropertyName string            `json:"propertyName"`
	Mapping      map[string]string `json:"mapping"`
}

type Schema struct {
	Ref                  string             `json:"$ref"`
	Title                string             `json:"title"`
	Type                 Type               `json:"type"`
	Properties           map[string]*Schema `json:"properties"`
	Format               Format             `json:"format"`
	Items                *Schema            `json:"items"`
	Description          string             `json:"description"`
	Default              interface{}        `json:"default"`
	Example              interface{}        `json:"example"`
	Deprecated           bool               `json:"deprecated"`
	Discriminator        Discriminator      `json:"discriminator"`
	Nullable             bool               `json:"nullable"`
	Maximum              interface{}        `json:"maximum"`
	Minimum              interface{}        `json:"minimum"`
	ExclusiveMaximum     interface{}        `json:"exclusiveMaximum"`
	ExclusiveMinimum     interface{}        `json:"exclusiveMinimum"`
	MaxLength            int                `json:"maxLength"`
	MinLength            int                `json:"minLength"`
	Required             []string           `json:"required"`
	Enum                 []string           `json:"enum"`
	AllOf                []*Schema          `json:"allOf"`
	OneOf                []*Schema          `json:"oneOf"`
	AnyOf                []*Schema          `json:"anyOf"`
	Not                  []*Schema          `json:"not"`
	AdditionalProperties *Schema            `json:"additionalProperties"`
	Pattern              interface{}        `json:"pattern"`
}

type Components struct {
	Schemas         map[string]Schema         `json:"schemas"`
	Responses       map[string]Response       `json:"responses"`
	Parameters      map[string]Parameter      `json:"parameters"`
	Examples        map[string]Example        `json:"examples"`
	RequestBodies   map[string]RequestBody    `json:"requestBodies"`
	Headers         map[string]Header         `json:"headers"`
	SecuritySchemes map[string]SecurityScheme `json:"securitySchemes"`
	Links           map[string]Link           `json:"links"`
	Callbacks       map[string]Callback       `json:"callbacks"`
}

type Api struct {
	Openapi      string          `json:"openapi"`
	Info         Info            `json:"info"`
	Servers      []Server        `json:"servers"`
	Tags         []Tag           `json:"tags"`
	Paths        map[string]Path `json:"paths"`
	Components   Components      `json:"components"`
	ExternalDocs ExternalDocs    `json:"externalDocs"`
}

type Type string

const (
	IntegerT Type = "integer"
	StringT  Type = "string"
	BooleanT Type = "boolean"
	NumberT  Type = "number"
	ObjectT  Type = "object"
	ArrayT   Type = "array"
)

type Format string

const (
	Int32F    Format = "int32"
	Int64F    Format = "int64"
	FloatF    Format = "float"
	DoubleF   Format = "double"
	DateTimeF Format = "date-time"
	BinaryF   Format = "binary"
)

var (
	Any = &Schema{
		Type: ObjectT,
	}
	Int = &Schema{
		Type:   IntegerT,
		Format: Int32F,
	}
	Int64 = &Schema{
		Type:   IntegerT,
		Format: Int64F,
	}
	String = &Schema{
		Type: StringT,
	}
	Time = &Schema{
		Type:   StringT,
		Format: DateTimeF,
	}
	Bool = &Schema{
		Type: BooleanT,
	}
	Float32 = &Schema{
		Type:   NumberT,
		Format: FloatF,
	}
	Float64 = &Schema{
		Type:   NumberT,
		Format: DoubleF,
	}
	File = &Schema{
		Type:   StringT,
		Format: BinaryF,
	}
)
