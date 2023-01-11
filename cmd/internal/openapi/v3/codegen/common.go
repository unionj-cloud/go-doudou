package codegen

import (
	"bytes"
	"fmt"
	"github.com/iancoleman/strcase"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/unionj-cloud/go-doudou/v2/toolkit/astutils"
	"github.com/unionj-cloud/go-doudou/v2/toolkit/copier"
	v3 "github.com/unionj-cloud/go-doudou/v2/toolkit/openapi/v3"
	"github.com/unionj-cloud/go-doudou/v2/toolkit/sliceutils"
	"github.com/unionj-cloud/go-doudou/v2/toolkit/stringutils"
	"github.com/unionj-cloud/go-doudou/v2/version"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"text/template"
)

type OpenAPICodeGenerator struct {
	Schemas          map[string]v3.Schema
	RequestBodies    map[string]v3.RequestBody
	Responses        map[string]v3.Response
	Omitempty        bool
	SvcName, ModName string
	Comments         []string
	DtoPkg           string
	ApiInfo          *v3.Info
}

func toComment(comment string, title ...string) string {
	if stringutils.IsEmpty(comment) {
		return ""
	}
	b := new(strings.Builder)
	lines := strings.Split(comment, "\n")
	for i, line := range lines {
		if len(title) > 0 && i == 0 {
			b.WriteString(fmt.Sprintf("// %s %s\n", title[0], line))
		} else {
			b.WriteString(fmt.Sprintf("// %s\n", line))
		}
	}
	return strings.TrimSuffix(b.String(), "\n")
}

func clean(str string) string {
	return strings.TrimSpace(strings.ReplaceAll(strings.ReplaceAll(str, "«", ""), "»", ""))
}

func toCamel(str string) string {
	return strcase.ToCamel(clean(str))
}

func (receiver OpenAPICodeGenerator) APIComments() []string {
	info := receiver.ApiInfo
	comments := []string{info.Title}
	comments = append(comments, strings.Split(info.Description, "\n")...)
	comments = append(comments, info.TermsOfService, info.Version)
	if info.Contact != nil {
		comments = append(comments, info.Contact.Email)
	}
	if info.License != nil {
		comments = append(comments, info.License.Name, info.License.URL)
	}
	return comments
}

func (receiver *OpenAPICodeGenerator) object2Struct(schema *v3.Schema) string {
	if schema.AdditionalProperties != nil {
		result := receiver.additionalProperties2Map(schema.AdditionalProperties)
		if stringutils.IsNotEmpty(result) {
			return result
		}
	}
	if len(schema.Properties) == 0 {
		return "interface{}"
	}
	b := new(strings.Builder)
	b.WriteString("struct {\n")
	for k, v := range schema.Properties {
		if stringutils.IsNotEmpty(v.Description) {
			descs := strings.Split(v.Description, "\n")
			for _, desc := range descs {
				b.WriteString(fmt.Sprintf("  // %s\n", desc))
			}
		}
		if sliceutils.StringContains(schema.Required, k) {
			b.WriteString("  // required\n")
		}
		jsontag := k
		if receiver.Omitempty {
			jsontag += ",omitempty"
		}
		if sliceutils.StringContains(schema.Required, k) {
			b.WriteString(fmt.Sprintf("  %s %s `json:\"%s\" url:\"%s\"`\n", strcase.ToCamel(k), receiver.toGoType(v), jsontag, k))
		} else {
			b.WriteString(fmt.Sprintf("  %s %s `json:\"%s\" url:\"%s\"`\n", strcase.ToCamel(k), "*"+receiver.toGoType(v), jsontag, k))
		}
	}
	b.WriteString("}")
	return b.String()
}

func string2Go(schema *v3.Schema) string {
	switch schema.Format {
	case v3.DateTimeF:
		return "time.Time"
	case v3.BinaryF:
		return "v3.FileModel"
	default:
		return "string"
	}
}

// integer2Go converts integer schema to golang basic type
//	Int32F    Format = "int32"
//	Int64F    Format = "int64"
//	FloatF    Format = "float"
//	DoubleF   Format = "double"
//	DateTimeF Format = "date-time"
//	BinaryF   Format = "binary"
func integer2Go(schema *v3.Schema) string {
	switch schema.Format {
	case v3.Int32F:
		return "int"
	case v3.Int64F:
		return "int64"
	default:
		return "int"
	}
}

func number2Go(schema *v3.Schema) string {
	switch schema.Format {
	case v3.FloatF:
		return "float32"
	case v3.DoubleF:
		return "float64"
	default:
		return "float64"
	}
}

// toGoType converts schema to golang type
//	IntegerT Type = "integer"
//	StringT  Type = "string"
//	BooleanT Type = "boolean"
//	NumberT  Type = "number"
//	ObjectT  Type = "object"
//	ArrayT   Type = "array"
func (receiver *OpenAPICodeGenerator) toGoType(schema *v3.Schema) string {
	if stringutils.IsNotEmpty(schema.Ref) {
		refName := strings.TrimPrefix(schema.Ref, "#/components/schemas/")
		if realSchema, exists := receiver.Schemas[refName]; exists {
			if realSchema.Type == v3.ObjectT && realSchema.AdditionalProperties != nil {
				result := receiver.additionalProperties2Map(realSchema.AdditionalProperties)
				if stringutils.IsNotEmpty(result) {
					return result
				}
			}
		}
		dtoName := toCamel(clean(refName))
		if stringutils.IsNotEmpty(receiver.DtoPkg) {
			dtoName = receiver.DtoPkg + "." + dtoName
		}
		return dtoName
	}
	switch schema.Type {
	case v3.IntegerT:
		return integer2Go(schema)
	case v3.StringT:
		return string2Go(schema)
	case v3.BooleanT:
		return "bool"
	case v3.NumberT:
		return number2Go(schema)
	case v3.ObjectT:
		return receiver.object2Struct(schema)
	case v3.ArrayT:
		return "[]" + receiver.toGoType(schema.Items)
	default:
		return "interface{}"
	}
}

func (receiver *OpenAPICodeGenerator) toOptionalGoType(schema *v3.Schema) string {
	if stringutils.IsNotEmpty(schema.Ref) {
		refName := strings.TrimPrefix(schema.Ref, "#/components/schemas/")
		if realSchema, exists := receiver.Schemas[refName]; exists {
			if realSchema.Type == v3.ObjectT && realSchema.AdditionalProperties != nil {
				result := receiver.additionalProperties2Map(realSchema.AdditionalProperties)
				if stringutils.IsNotEmpty(result) {
					return result
				}
			}
		}
		dtoName := toCamel(clean(refName))
		if stringutils.IsNotEmpty(receiver.DtoPkg) {
			dtoName = receiver.DtoPkg + "." + dtoName
		}
		return "*" + dtoName
	}
	switch schema.Type {
	case v3.IntegerT:
		return "*" + integer2Go(schema)
	case v3.StringT:
		return "*" + string2Go(schema)
	case v3.BooleanT:
		return "*bool"
	case v3.NumberT:
		return "*" + number2Go(schema)
	case v3.ObjectT:
		result := receiver.object2Struct(schema)
		if strings.HasPrefix(result, "struct {") {
			return "*" + result
		}
		return result
	case v3.ArrayT:
		return "[]" + receiver.toGoType(schema.Items)
	default:
		return "interface{}"
	}
}

func (receiver *OpenAPICodeGenerator) additionalProperties2Map(additionalProperties interface{}) string {
	if additionalProperties == nil {
		return ""
	}
	if value, ok := additionalProperties.(map[string]interface{}); ok {
		var additionalSchema v3.Schema
		copier.DeepCopy(value, &additionalSchema)
		if stringutils.IsNotEmpty(additionalSchema.XMapType) {
			return additionalSchema.XMapType
		}
		return "map[string]" + receiver.toGoType(&additionalSchema)
	}
	return ""
}

func (receiver *OpenAPICodeGenerator) GenGoDto(schemas map[string]v3.Schema, output, pkg, tmpl string) {
	if err := os.MkdirAll(filepath.Dir(output), os.ModePerm); err != nil {
		panic(err)
	}
	funcMap := make(map[string]interface{})
	funcMap["toCamel"] = toCamel
	funcMap["toGoType"] = receiver.toGoType
	funcMap["toComment"] = toComment
	funcMap["toOptionalGoType"] = receiver.toOptionalGoType
	funcMap["stringContains"] = sliceutils.StringContains
	filterMap := make(map[string]v3.Schema)
	for k, v := range schemas {
		result := receiver.additionalProperties2Map(v.AdditionalProperties)
		if stringutils.IsEmpty(result) {
			filterMap[k] = v
		}
	}
	tpl, _ := template.New("dto.go.tmpl").Funcs(funcMap).Parse(tmpl)
	var sqlBuf bytes.Buffer
	_ = tpl.Execute(&sqlBuf, struct {
		Schemas map[string]v3.Schema
		Omit    bool
		Pkg     string
		Version string
	}{
		Schemas: filterMap,
		Omit:    receiver.Omitempty,
		Pkg:     pkg,
		Version: version.Release,
	})
	source := strings.TrimSpace(sqlBuf.String())
	astutils.FixImport([]byte(source), output)
}

// TODO example2Schema converts example to *v3.Schema
func (receiver *OpenAPICodeGenerator) example2Schema(example interface{}, exampleType v3.ExampleType) *v3.Schema {
	return v3.Any
}

func (receiver *OpenAPICodeGenerator) schema2Field(schema *v3.Schema, name string, example interface{}, exampleType v3.ExampleType) *astutils.FieldMeta {
	if schema == nil {
		schema = receiver.example2Schema(example, exampleType)
	}
	var comments []string
	if stringutils.IsNotEmpty(schema.Description) {
		comments = append(comments, strings.Split(schema.Description, "\n")...)
	}
	return &astutils.FieldMeta{
		Name:     name,
		Type:     receiver.toGoType(schema),
		Comments: comments,
	}
}

func (receiver *OpenAPICodeGenerator) responseBody(endpoint, httpMethod string, operation *v3.Operation) (results []astutils.FieldMeta, err error) {
	if stringutils.IsNotEmpty(operation.Responses.Resp200.Ref) {
		key := strings.TrimPrefix(operation.Responses.Resp200.Ref, "#/components/responses/")
		if response, exists := receiver.Responses[key]; exists {
			operation.Responses.Resp200 = &response
		} else {
			panic(fmt.Errorf("response %s not exists", operation.Responses.Resp200.Ref))
		}
	}

	content := operation.Responses.Resp200.Content
	if content == nil {
		return nil, errors.Errorf("200 response content definition not found in api %s %s", httpMethod, endpoint)
	}

	if content.JSON != nil {
		results = append(results, *receiver.schema2Field(content.JSON.Schema, "ret", content.JSON.Example, v3.JSON_EXAMPLE))
	} else if content.Stream != nil {
		results = append(results, astutils.FieldMeta{
			Name: "_downloadFile",
			Type: "*os.File",
		})
	} else if content.TextPlain != nil {
		results = append(results, *receiver.schema2Field(content.TextPlain.Schema, "ret", content.TextPlain.Example, v3.TEXT_EXAMPLE))
	} else if content.Default != nil {
		results = append(results, *receiver.schema2Field(content.Default.Schema, "ret", content.TextPlain.Example, v3.TEXT_EXAMPLE))
	} else {
		return nil, errors.Errorf("200 response content definition not support yet in api %s %s", httpMethod, endpoint)
	}
	return
}

func commentLines(operation *v3.Operation) []string {
	var comments []string
	if stringutils.IsNotEmpty(operation.Summary) {
		comments = append(comments, strings.Split(operation.Summary, "\n")...)
	}
	if stringutils.IsNotEmpty(operation.Description) {
		comments = append(comments, strings.Split(operation.Description, "\n")...)
	}
	return comments
}

func (receiver *OpenAPICodeGenerator) parameter2Field(param v3.Parameter) astutils.FieldMeta {
	var comments []string
	if stringutils.IsNotEmpty(param.Description) {
		comments = append(comments, strings.Split(param.Description, "\n")...)
	}
	t := receiver.toGoType(param.Schema)
	if param.Required {
		comments = append(comments, "required")
	} else {
		t = v3.ToOptional(t)
	}
	return astutils.FieldMeta{
		Name:     param.Name,
		Type:     t,
		Comments: comments,
	}
}

func toMethod(endpoint string) string {
	endpoint = strings.ReplaceAll(strings.ReplaceAll(endpoint, "{", ""), "}", "")
	endpoint = strings.ReplaceAll(strings.Trim(endpoint, "/"), "/", "_")
	nosymbolreg := regexp.MustCompile(`[^a-zA-Z0-9_]`)
	endpoint = nosymbolreg.ReplaceAllLiteralString(endpoint, "")
	endpoint = strcase.ToCamel(endpoint)
	numberstartreg := regexp.MustCompile(`^[0-9]+`)
	if numberstartreg.MatchString(endpoint) {
		startNumbers := numberstartreg.FindStringSubmatch(endpoint)
		endpoint = numberstartreg.ReplaceAllLiteralString(endpoint, "")
		endpoint += startNumbers[0]
	}
	return endpoint
}

func httpMethod(method string) string {
	httpMethods := []string{"GET", "POST", "PUT", "DELETE"}
	snake := strcase.ToSnake(method)
	splits := strings.Split(snake, "_")
	head := strings.ToUpper(splits[0])
	for _, m := range httpMethods {
		if head == m {
			return m
		}
	}
	return "POST"
}

func RestyMethod(method string) string {
	return strings.Title(strings.ToLower(httpMethod(method)))
}

func IsOptional(t string) bool {
	return strings.HasPrefix(t, "*")
}

func (receiver *OpenAPICodeGenerator) Api2Interface(paths map[string]v3.Path, svcName string, operationConverter IOperationConverter) astutils.InterfaceMeta {
	var meta astutils.InterfaceMeta
	meta.Name = strcase.ToCamel(svcName)
	meta.Comments = receiver.APIComments()
	for endpoint, path := range paths {
		if path.Get != nil {
			if method, err := operationConverter.ConvertOperation(endpoint, "Get", path.Get, path.Parameters); err == nil {
				meta.Methods = append(meta.Methods, method)
			} else {
				logrus.Errorln(err)
			}
		}
		if path.Post != nil {
			if method, err := operationConverter.ConvertOperation(endpoint, "Post", path.Post, path.Parameters); err == nil {
				meta.Methods = append(meta.Methods, method)
			} else {
				logrus.Errorln(err)
			}
		}
		if path.Put != nil {
			if method, err := operationConverter.ConvertOperation(endpoint, "Put", path.Put, path.Parameters); err == nil {
				meta.Methods = append(meta.Methods, method)
			} else {
				logrus.Errorln(err)
			}
		}
		if path.Delete != nil {
			if method, err := operationConverter.ConvertOperation(endpoint, "Delete", path.Delete, path.Parameters); err == nil {
				meta.Methods = append(meta.Methods, method)
			} else {
				logrus.Errorln(err)
			}
		}
	}
	sort.SliceStable(meta.Methods, func(i, j int) bool {
		return meta.Methods[i].Name < meta.Methods[j].Name
	})
	return meta
}

type ClientOperationConverter struct {
	_         [0]int
	Generator *OpenAPICodeGenerator
}

func (receiver *ClientOperationConverter) form(operation *v3.Operation) (bodyParams *astutils.FieldMeta) {
	receiver.Generator.resolveSchemaFromRef(operation)
	content := operation.RequestBody.Content
	if content.JSON != nil {
	} else if content.FormURL != nil {
		bodyParams = receiver.Generator.schema2Field(content.FormURL.Schema, "bodyParams", content.FormURL.Example, v3.TEXT_EXAMPLE)
		if !operation.RequestBody.Required && bodyParams != nil {
			bodyParams.Type = v3.ToOptional(bodyParams.Type)
		}
	} else if content.FormData != nil {
		bodyParams, _ = receiver.Generator.parseFormData(content.FormData)
		if !operation.RequestBody.Required && bodyParams != nil {
			bodyParams.Type = v3.ToOptional(bodyParams.Type)
		}
	}
	return
}

func (receiver *ClientOperationConverter) ConvertOperation(endpoint, httpMethod string, operation *v3.Operation, gparams []v3.Parameter) (astutils.MethodMeta, error) {
	var files, params []astutils.FieldMeta
	var bodyJSON, bodyParams, qparams *astutils.FieldMeta
	comments := commentLines(operation)
	qSchema, pathvars, headervars := receiver.globalParams(gparams)
	receiver.operationParams(operation.Parameters, &qSchema, &pathvars, &headervars)

	if len(qSchema.Properties) > 0 {
		qparams = receiver.Generator.schema2Field(&qSchema, "queryParams", nil, v3.UNKNOWN_EXAMPLE)
		if qSchema.Type == v3.ObjectT && len(qSchema.Required) == 0 {
			qparams.Type = v3.ToOptional(qparams.Type)
		}
	}

	if httpMethod != "Get" && operation.RequestBody != nil {
		bodyJSON, files = receiver.Generator.requestBody(operation)
		bodyParams = receiver.form(operation)
	}

	if operation.Responses == nil {
		return astutils.MethodMeta{}, errors.Errorf("response definition not found in api %s %s", httpMethod, endpoint)
	}

	if operation.Responses.Resp200 == nil {
		return astutils.MethodMeta{}, errors.Errorf("200 response definition not found in api %s %s", httpMethod, endpoint)
	}

	results, err := receiver.Generator.responseBody(endpoint, httpMethod, operation)
	if err != nil {
		return astutils.MethodMeta{}, err
	}

	if qparams != nil {
		params = append(params, *qparams)
	}

	params = append(params, pathvars...)
	params = append(params, headervars...)

	if bodyParams != nil {
		params = append(params, *bodyParams)
	}

	if bodyJSON != nil {
		params = append(params, *bodyJSON)
	}

	params = append(params, files...)

	ret := astutils.MethodMeta{
		Name:        httpMethod + toMethod(endpoint),
		Params:      params,
		Results:     results,
		PathVars:    pathvars,
		HeaderVars:  headervars,
		BodyParams:  bodyParams,
		BodyJSON:    bodyJSON,
		Files:       files,
		Comments:    comments,
		Path:        endpoint,
		QueryParams: qparams,
	}
	return ret, nil
}

func (receiver *ClientOperationConverter) operationParams(parameters []v3.Parameter, qSchema *v3.Schema, pathvars, headervars *[]astutils.FieldMeta) {
	for _, item := range parameters {
		switch item.In {
		case v3.InQuery:
			qSchema.Properties[item.Name] = item.Schema
			if item.Required {
				qSchema.Required = append(qSchema.Required, item.Name)
			}
		case v3.InPath:
			*pathvars = append(*pathvars, receiver.Generator.parameter2Field(item))
		case v3.InHeader:
			*headervars = append(*headervars, receiver.Generator.parameter2Field(item))
		default:
			panic(fmt.Errorf("not support %s parameter yet", item.In))
		}
	}
}

func (receiver *OpenAPICodeGenerator) requestBody(operation *v3.Operation) (bodyJSON *astutils.FieldMeta, files []astutils.FieldMeta) {
	receiver.resolveSchemaFromRef(operation)

	content := operation.RequestBody.Content
	if content.JSON != nil {
		bodyJSON = receiver.schema2Field(content.JSON.Schema, "bodyJSON", content.JSON.Example, v3.JSON_EXAMPLE)
		if !operation.RequestBody.Required && bodyJSON != nil {
			bodyJSON.Type = v3.ToOptional(bodyJSON.Type)
		}
	} else if content.FormData != nil {
		_, files = receiver.parseFormData(content.FormData)
	} else if content.Stream != nil {
		f := astutils.FieldMeta{
			Name: "file",
			Type: "v3.FileModel",
		}
		if !operation.RequestBody.Required {
			f.Type = v3.ToOptional(f.Type)
		}
		files = append(files, f)
	} else if content.TextPlain != nil {
		bodyJSON = receiver.schema2Field(content.TextPlain.Schema, "bodyJSON", content.TextPlain.Example, v3.TEXT_EXAMPLE)
		if !operation.RequestBody.Required && bodyJSON != nil {
			bodyJSON.Type = v3.ToOptional(bodyJSON.Type)
		}
	} else if content.Default != nil {
		bodyJSON = receiver.schema2Field(content.Default.Schema, "bodyJSON", content.TextPlain.Example, v3.TEXT_EXAMPLE)
		if !operation.RequestBody.Required && bodyJSON != nil {
			bodyJSON.Type = v3.ToOptional(bodyJSON.Type)
		}
	}
	return
}

func (receiver *OpenAPICodeGenerator) parseFormData(formData *v3.MediaType) (bodyParams *astutils.FieldMeta, files []astutils.FieldMeta) {
	schema := *formData.Schema
	if stringutils.IsNotEmpty(schema.Ref) {
		schema = receiver.Schemas[strings.TrimPrefix(formData.Schema.Ref, "#/components/schemas/")]
	}
	aSchema := v3.Schema{
		Type:       v3.ObjectT,
		Properties: make(map[string]*v3.Schema),
	}
	for k, v := range schema.Properties {
		var gotype string
		if v.Type == v3.StringT && v.Format == v3.BinaryF {
			gotype = "v3.FileModel"
		} else if v.Type == v3.ArrayT && v.Items.Type == v3.StringT && v.Items.Format == v3.BinaryF {
			gotype = "[]v3.FileModel"
		}
		if stringutils.IsNotEmpty(gotype) && !sliceutils.StringContains(schema.Required, k) {
			gotype = v3.ToOptional(gotype)
		}
		if stringutils.IsNotEmpty(gotype) {
			files = append(files, astutils.FieldMeta{
				Name: k,
				Type: gotype,
			})
			continue
		}
		aSchema.Properties[k] = v
		if sliceutils.StringContains(schema.Required, k) {
			aSchema.Required = append(aSchema.Required, k)
		}
	}
	if len(aSchema.Properties) > 0 {
		bodyParams = receiver.schema2Field(&aSchema, "bodyParams", nil, v3.UNKNOWN_EXAMPLE)
	}
	return
}

// resolveSchemaFromRef resolves schema from ref
func (receiver *OpenAPICodeGenerator) resolveSchemaFromRef(operation *v3.Operation) {
	if stringutils.IsNotEmpty(operation.RequestBody.Ref) {
		// #/components/requestBodies/Raw3
		key := strings.TrimPrefix(operation.RequestBody.Ref, "#/components/requestBodies/")
		if requestBody, exists := receiver.RequestBodies[key]; exists {
			operation.RequestBody = &requestBody
		} else {
			panic(fmt.Errorf("requestBody %s not exists", operation.RequestBody.Ref))
		}
	}
}

func (receiver *ClientOperationConverter) globalParams(gparams []v3.Parameter) (v3.Schema, []astutils.FieldMeta, []astutils.FieldMeta) {
	var pathvars, headervars []astutils.FieldMeta
	qSchema := v3.Schema{
		Type:       v3.ObjectT,
		Properties: make(map[string]*v3.Schema),
	}
	for _, item := range gparams {
		switch item.In {
		case v3.InQuery:
			qSchema.Properties[item.Name] = item.Schema
			if item.Required {
				qSchema.Required = append(qSchema.Required, item.Name)
			}
		case v3.InPath:
			pathvars = append(pathvars, receiver.Generator.parameter2Field(item))
		case v3.InHeader:
			headervars = append(headervars, receiver.Generator.parameter2Field(item))
		default:
			panic(fmt.Errorf("not support %s parameter yet", item.In))
		}
	}
	return qSchema, pathvars, headervars
}

var httpTmpl = `/**
* Generated by go-doudou {{.Version}}.
* You can edit it as your need.
*/
package {{.Pkg}}

import (
	"context"
	"encoding/json"
	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
	"github.com/opentracing-contrib/go-stdlib/nethttp"
	"github.com/opentracing/opentracing-go"
	"github.com/unionj-cloud/go-doudou/v2/framework/registry"
	_querystring "github.com/google/go-querystring/query"
	"github.com/unionj-cloud/go-doudou/v2/toolkit/fileutils"
	"github.com/unionj-cloud/go-doudou/v2/toolkit/stringutils"
	"github.com/unionj-cloud/go-doudou/v2/framework/restclient"
	v3 "github.com/unionj-cloud/go-doudou/v2/toolkit/openapi/v3"
	"io"
	"mime/multipart"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

type {{.Meta.Name}}Client struct {
	provider registry.IServiceProvider
	client   *resty.Client
	rootPath string
}

func (receiver *{{.Meta.Name}}Client) SetRootPath(rootPath string) {
	receiver.rootPath = rootPath
}

func (receiver *{{.Meta.Name}}Client) SetProvider(provider registry.IServiceProvider) {
	receiver.provider = provider
}

func (receiver *{{.Meta.Name}}Client) SetClient(client *resty.Client) {
	receiver.client = client
}

{{- range $m := .Meta.Methods }}
	{{- range $i, $c := $m.Comments }}
	{{- if eq $i 0}}
	// {{$m.Name}} {{$c}}
	{{- else}}
	// {{$c}}
	{{- end}}
	{{- end }}
	func (receiver *{{$.Meta.Name}}Client) {{$m.Name}}(ctx context.Context, _headers map[string]string, {{ range $i, $p := $m.Params}}
    {{- if $i}},{{end}}
	{{- range $c := $p.Comments }}
	// {{$c}}
	{{- end }}
    {{ $p.Name}} {{$p.Type}}
    {{- end }}) ({{(index $m.Results 0).Name}} {{(index $m.Results 0).Type}}, _resp *resty.Response, err error) {
		var _err error

		_req := receiver.client.R()
		_req.SetContext(ctx)
		if len(_headers) > 0 {
			_req.SetHeaders(_headers)
		}
		{{- if $m.QueryParams }}
			_queryParams, _ := _querystring.Values({{$m.QueryParams.Name}})
			_req.SetQueryParamsFromValues(_queryParams)
		{{- end }}
		{{- if $m.PathVars }}
			{{- range $p := $m.PathVars }}
			{{- if isOptional $p.Type }}
			if {{$p.Name}} != nil { 
				_req.SetPathParam("{{$p.Name}}", fmt.Sprintf("%v", *{{$p.Name}}))
			}
			{{- else }}
			_req.SetPathParam("{{$p.Name}}", fmt.Sprintf("%v", {{$p.Name}}))
			{{- end }}
			{{- end }}
		{{- end }}
		{{- if $m.HeaderVars }}
			{{- range $p := $m.HeaderVars }}
			{{- if isOptional $p.Type }}
			if {{$p.Name}} != nil { 
				_req.SetHeader("{{$p.Name}}", fmt.Sprintf("%v", *{{$p.Name}}))
			}
			{{- else }}
			_req.SetHeader("{{$p.Name}}", fmt.Sprintf("%v", {{$p.Name}}))
			{{- end }}
			{{- end }}
		{{- end }}
		{{- if $m.BodyParams }}
			_bodyParams, _ := _querystring.Values({{$m.BodyParams.Name}})
			_req.SetFormDataFromValues(_bodyParams)
		{{- end }}
		{{- if $m.BodyJSON }}
			_req.SetBody({{$m.BodyJSON.Name}})
		{{- end }}
		{{- if $m.Files }}
			{{- range $p := $m.Files }}
				{{- if contains $p.Type "["}}
				{{- if isOptional $p.Type }}
				if {{$p.Name}} != nil {
					for _, _f := range *{{$p.Name}} {
						_req.SetFileReader("{{$p.Name}}", _f.Filename, _f.Reader)
					}
				}
				{{- else }}
				if len({{$p.Name}}) == 0 {
					err = errors.New("at least one file should be uploaded for parameter {{$p.Name}}")
					return
				}
				for _, _f := range {{$p.Name}} {
					_req.SetFileReader("{{$p.Name}}", _f.Filename, _f.Reader)
				}
				{{- end }}
				{{- else}}
				{{- if isOptional $p.Type }}
				if {{$p.Name}} != nil { 
					_req.SetFileReader("{{$p.Name}}", {{$p.Name}}.Filename, {{$p.Name}}.Reader)
				}
				{{- else }}
				_req.SetFileReader("{{$p.Name}}", {{$p.Name}}.Filename, {{$p.Name}}.Reader)
				{{- end }}
				{{- end }}
			{{- end }}
		{{- end }}

		{{- range $r := $m.Results }}
			{{- if eq $r.Type "*os.File" }}
				_req.SetDoNotParseResponse(true)
			{{- end }}
		{{- end }}

		_resp, _err = _req.{{$m.Name | restyMethod}}("{{$m.Path}}")
		if _err != nil {
			err = errors.Wrap(_err, "")
			return
		}
		if _resp.IsError() {
			err = errors.New(_resp.String())
			return
		}
		{{- $done := false }}
		{{- range $r := $m.Results }}
			{{- if eq $r.Type "*os.File" }}
				_disp := _resp.Header().Get("Content-Disposition")
				_file := strings.TrimPrefix(_disp, "attachment; filename=")
				_output := os.TempDir()
				if stringutils.IsNotEmpty(_output) {
					_file = _output + string(filepath.Separator) + _file
				}
				_file = filepath.Clean(_file)
				if _err = fileutils.CreateDirectory(filepath.Dir(_file)); _err != nil {
					err = errors.Wrap(_err, "")
					return
				}
				_outFile, _err := os.Create(_file)
				if _err != nil {
					err = errors.Wrap(_err, "")
					return
				}
				defer _outFile.Close()
				defer _resp.RawBody().Close()
				_, _err = io.Copy(_outFile, _resp.RawBody())
				if _err != nil {
					err = errors.Wrap(_err, "")
					return
				}
				{{ $r.Name }} = _outFile
				return
				{{- $done = true }}	
			{{- end }}
		{{- end }}
		{{- if not $done }}
			{{- if eq (index $m.Results 0).Type "string" }}
			{{(index $m.Results 0).Name}} = _resp.String()
			{{- else }}
			if _err = json.Unmarshal(_resp.Body(), &{{(index $m.Results 0).Name}}); _err != nil {
				err = errors.Wrap(_err, "")
				return
			}
			{{- end }}
			return
		{{- end }}  
	}
{{- end }}

func New{{.Meta.Name}}(opts ...restclient.RestClientOption) *{{.Meta.Name}}Client {
	{{- if .Env }}
	defaultProvider := restclient.NewServiceProvider("{{.Env}}")
	{{- else }}
	defaultProvider := restclient.NewServiceProvider("{{.Meta.Name | toUpper}}")
	{{- end }}
	defaultClient := restclient.NewClient()

	svcClient := &{{.Meta.Name}}Client{
		provider: defaultProvider,
		client:   defaultClient,
	}

	for _, opt := range opts {
		opt(svcClient)
	}

	svcClient.client.OnBeforeRequest(func(_ *resty.Client, request *resty.Request) error {
		request.URL = svcClient.provider.SelectServer() + svcClient.rootPath + request.URL
		return nil
	})

	svcClient.client.SetPreRequestHook(func(_ *resty.Client, request *http.Request) error {
		traceReq, _ := nethttp.TraceRequest(opentracing.GlobalTracer(), request,
			nethttp.OperationName(fmt.Sprintf("HTTP %s: %s", request.Method, request.URL.Path)))
		*request = *traceReq
		return nil
	})

	svcClient.client.OnAfterResponse(func(_ *resty.Client, response *resty.Response) error {
		nethttp.TracerFromRequest(response.Request.RawRequest).Finish()
		return nil
	})

	return svcClient
}
`

func (receiver *OpenAPICodeGenerator) GenGoHTTP(paths map[string]v3.Path, svcName, dir, env, pkg string, operationConverter IOperationConverter) {
	_ = os.MkdirAll(dir, os.ModePerm)
	output := filepath.Join(dir, svcName+"client.go")
	fi, err := os.Stat(output)
	if err != nil && !os.IsNotExist(err) {
		panic(err)
	}
	if fi != nil {
		logrus.Warningln("file " + svcName + "client.go will be overwritten")
	}
	var f *os.File
	if f, err = os.Create(output); err != nil {
		panic(err)
	}
	defer func(f *os.File) {
		_ = f.Close()
	}(f)

	funcMap := make(map[string]interface{})
	funcMap["toCamel"] = strcase.ToCamel
	funcMap["contains"] = strings.Contains
	funcMap["restyMethod"] = RestyMethod
	funcMap["toUpper"] = strings.ToUpper
	funcMap["isOptional"] = IsOptional
	tpl, _ := template.New("http.go.tmpl").Funcs(funcMap).Parse(httpTmpl)
	var sqlBuf bytes.Buffer
	_ = tpl.Execute(&sqlBuf, struct {
		Meta    astutils.InterfaceMeta
		Env     string
		Pkg     string
		Version string
	}{
		Meta:    receiver.Api2Interface(paths, svcName, operationConverter),
		Env:     env,
		Pkg:     pkg,
		Version: version.Release,
	})
	source := strings.TrimSpace(sqlBuf.String())
	astutils.FixImport([]byte(source), output)
}
