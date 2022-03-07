package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/iancoleman/strcase"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/unionj-cloud/go-doudou/cmd/internal/astutils"
	v3 "github.com/unionj-cloud/go-doudou/toolkit/openapi/v3"
	"github.com/unionj-cloud/go-doudou/toolkit/sliceutils"
	"github.com/unionj-cloud/go-doudou/toolkit/stringutils"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"
)

var votmpl = `package {{.Pkg}}

{{- range $k, $v := .Schemas }}
{{ toComment $v.Description ($k | toCamel)}}
type {{$k | toCamel}} struct {
{{- range $pk, $pv := $v.Properties }}
	{{ $pv.Description | toComment }}
	{{- if stringContains $v.Required $pk }}
	// required
	{{ $pk | toCamel}} {{$pv | toGoType }} ` + "`" + `json:"{{$pk}}{{if $.Omit}},omitempty{{end}}" url:"{{$pk}}"` + "`" + `
	{{- else }}
	{{ $pk | toCamel}} {{$pv | toGoType | toOptional }} ` + "`" + `json:"{{$pk}}{{if $.Omit}},omitempty{{end}}" url:"{{$pk}}"` + "`" + `
	{{- end }}
{{- end }}
}
{{- end }}
`

var httptmpl = `package {{.Pkg}}

import (
	"context"
	"encoding/json"
	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
	"github.com/opentracing-contrib/go-stdlib/nethttp"
	"github.com/opentracing/opentracing-go"
	"github.com/unionj-cloud/go-doudou/framework/registry"
	_querystring "github.com/google/go-querystring/query"
	"github.com/unionj-cloud/go-doudou/toolkit/fileutils"
	"github.com/unionj-cloud/go-doudou/toolkit/stringutils"
	ddhttp "github.com/unionj-cloud/go-doudou/framework/http"
	v3 "github.com/unionj-cloud/go-doudou/toolkit/openapi/v3"
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

func New{{.Meta.Name}}(opts ...ddhttp.DdClientOption) *{{.Meta.Name}}Client {
	{{- if .Env }}
	defaultProvider := ddhttp.NewServiceProvider("{{.Env}}")
	{{- else }}
	defaultProvider := ddhttp.NewServiceProvider("{{.Meta.Name | toUpper}}")
	{{- end }}
	defaultClient := ddhttp.NewClient()

	svcClient := &{{.Meta.Name}}Client{
		provider: defaultProvider,
		client:   defaultClient,
	}

	for _, opt := range opts {
		opt(svcClient)
	}

	svcClient.client.OnBeforeRequest(func(_ *resty.Client, request *resty.Request) error {
		request.URL = svcClient.provider.SelectServer() + request.URL
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

func restyMethod(method string) string {
	return strings.Title(strings.ToLower(httpMethod(method)))
}

func isOptional(t string) bool {
	return strings.HasPrefix(t, "*")
}

func genGoHTTP(paths map[string]v3.Path, svcname, dir, env, pkg string) {
	_ = os.MkdirAll(dir, os.ModePerm)
	output := filepath.Join(dir, svcname+"client.go")
	fi, err := os.Stat(output)
	if err != nil && !os.IsNotExist(err) {
		panic(err)
	}
	if fi != nil {
		logrus.Warningln("file " + svcname + "client.go will be overwritten")
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
	funcMap["restyMethod"] = restyMethod
	funcMap["toUpper"] = strings.ToUpper
	funcMap["isOptional"] = isOptional
	tpl, _ := template.New("http.go.tmpl").Funcs(funcMap).Parse(httptmpl)
	var sqlBuf bytes.Buffer
	_ = tpl.Execute(&sqlBuf, struct {
		Meta astutils.InterfaceMeta
		Env  string
		Pkg  string
	}{
		Meta: api2Interface(paths, svcname),
		Env:  env,
		Pkg:  pkg,
	})
	source := strings.TrimSpace(sqlBuf.String())
	astutils.FixImport([]byte(source), output)
}

func api2Interface(paths map[string]v3.Path, svcname string) astutils.InterfaceMeta {
	var meta astutils.InterfaceMeta
	meta.Name = strcase.ToCamel(svcname)
	for endpoint, path := range paths {
		if path.Get != nil {
			if method, err := operation2Method(endpoint, "Get", path.Get, path.Parameters); err == nil {
				meta.Methods = append(meta.Methods, method)
			} else {
				logrus.Errorln(err)
			}
		}
		if path.Post != nil {
			if method, err := operation2Method(endpoint, "Post", path.Post, path.Parameters); err == nil {
				meta.Methods = append(meta.Methods, method)
			} else {
				logrus.Errorln(err)
			}
		}
		if path.Put != nil {
			if method, err := operation2Method(endpoint, "Put", path.Put, path.Parameters); err == nil {
				meta.Methods = append(meta.Methods, method)
			} else {
				logrus.Errorln(err)
			}
		}
		if path.Delete != nil {
			if method, err := operation2Method(endpoint, "Delete", path.Delete, path.Parameters); err == nil {
				meta.Methods = append(meta.Methods, method)
			} else {
				logrus.Errorln(err)
			}
		}
	}
	return meta
}

func operation2Method(endpoint, httpMethod string, operation *v3.Operation, gparams []v3.Parameter) (astutils.MethodMeta, error) {
	var files, params []astutils.FieldMeta
	var bodyJSON, bodyParams, qparams *astutils.FieldMeta
	comments := commentLines(operation)
	qSchema, pathvars, headervars := globalParams(gparams)
	operationParams(operation.Parameters, &qSchema, &pathvars, &headervars)

	if len(qSchema.Properties) > 0 {
		qparams = schema2Field(&qSchema, "queryParams")
		if qSchema.Type == v3.ObjectT && len(qSchema.Required) == 0 {
			qparams.Type = toOptional(qparams.Type)
		}
	}

	if httpMethod != "Get" && operation.RequestBody != nil {
		bodyJSON, bodyParams, files = requestBody(operation)
	}

	if operation.Responses == nil {
		return astutils.MethodMeta{}, errors.Errorf("response definition not found in api %s %s", httpMethod, endpoint)
	}

	if operation.Responses.Resp200 == nil {
		return astutils.MethodMeta{}, errors.Errorf("200 response definition not found in api %s %s", httpMethod, endpoint)
	}

	results, err := responseBody(endpoint, httpMethod, operation)
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

	return astutils.MethodMeta{
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
	}, nil
}

func operationParams(parameters []v3.Parameter, qSchema *v3.Schema, pathvars, headervars *[]astutils.FieldMeta) {
	for _, item := range parameters {
		switch item.In {
		case v3.InQuery:
			qSchema.Properties[item.Name] = item.Schema
			if item.Required {
				qSchema.Required = append(qSchema.Required, item.Name)
			}
		case v3.InPath:
			*pathvars = append(*pathvars, parameter2Field(item))
		case v3.InHeader:
			*headervars = append(*headervars, parameter2Field(item))
		default:
			panic(fmt.Errorf("not support %s parameter yet", item.In))
		}
	}
}

func responseBody(endpoint, httpMethod string, operation *v3.Operation) (results []astutils.FieldMeta, err error) {
	if stringutils.IsNotEmpty(operation.Responses.Resp200.Ref) {
		key := strings.TrimPrefix(operation.Responses.Resp200.Ref, "#/components/responses/")
		if response, exists := responses[key]; exists {
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
		results = append(results, *schema2Field(content.JSON.Schema, "ret"))
	} else if content.Stream != nil {
		results = append(results, astutils.FieldMeta{
			Name: "_downloadFile",
			Type: "*os.File",
		})
	} else if content.TextPlain != nil {
		results = append(results, *schema2Field(content.TextPlain.Schema, "ret"))
	} else if content.Default != nil {
		results = append(results, *schema2Field(content.Default.Schema, "ret"))
	} else {
		return nil, errors.Errorf("200 response content definition not support yet in api %s %s", httpMethod, endpoint)
	}
	return
}

func requestBody(operation *v3.Operation) (bodyJSON, bodyParams *astutils.FieldMeta, files []astutils.FieldMeta) {
	resolveSchemaFromRef(operation)

	content := operation.RequestBody.Content
	if content.JSON != nil {
		bodyJSON = schema2Field(content.JSON.Schema, "bodyJSON")
		if !operation.RequestBody.Required {
			bodyJSON.Type = toOptional(bodyJSON.Type)
		}
	} else if content.FormURL != nil {
		bodyParams = schema2Field(content.FormURL.Schema, "bodyParams")
		if !operation.RequestBody.Required {
			bodyParams.Type = toOptional(bodyParams.Type)
		}
	} else if content.FormData != nil {
		bodyParams, files = parseFormData(content.FormData)
		if !operation.RequestBody.Required {
			bodyParams.Type = toOptional(bodyParams.Type)
		}
	} else if content.Stream != nil {
		f := astutils.FieldMeta{
			Name: "file",
			Type: "v3.FileModel",
		}
		if !operation.RequestBody.Required {
			f.Type = toOptional(f.Type)
		}
		files = append(files, f)
	} else if content.TextPlain != nil {
		bodyJSON = schema2Field(content.TextPlain.Schema, "bodyJSON")
		if !operation.RequestBody.Required {
			bodyJSON.Type = toOptional(bodyJSON.Type)
		}
	} else if content.Default != nil {
		bodyJSON = schema2Field(content.Default.Schema, "bodyJSON")
		if !operation.RequestBody.Required {
			bodyJSON.Type = toOptional(bodyJSON.Type)
		}
	}
	return
}

func parseFormData(formData *v3.MediaType) (bodyParams *astutils.FieldMeta, files []astutils.FieldMeta) {
	schema := *formData.Schema
	if stringutils.IsNotEmpty(schema.Ref) {
		schema = schemas[strings.TrimPrefix(formData.Schema.Ref, "#/components/schemas/")]
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
			gotype = toOptional(gotype)
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
		bodyParams = schema2Field(&aSchema, "bodyParams")
	}
	return
}

// resolveSchemaFromRef resolves schema from ref
func resolveSchemaFromRef(operation *v3.Operation) {
	if stringutils.IsNotEmpty(operation.RequestBody.Ref) {
		// #/components/requestBodies/Raw3
		key := strings.TrimPrefix(operation.RequestBody.Ref, "#/components/requestBodies/")
		if requestBody, exists := requestBodies[key]; exists {
			operation.RequestBody = &requestBody
		} else {
			panic(fmt.Errorf("requestBody %s not exists", operation.RequestBody.Ref))
		}
	}
}

func globalParams(gparams []v3.Parameter) (v3.Schema, []astutils.FieldMeta, []astutils.FieldMeta) {
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
			pathvars = append(pathvars, parameter2Field(item))
		case v3.InHeader:
			headervars = append(headervars, parameter2Field(item))
		default:
			panic(fmt.Errorf("not support %s parameter yet", item.In))
		}
	}
	return qSchema, pathvars, headervars
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

func schema2Field(schema *v3.Schema, name string) *astutils.FieldMeta {
	var comments []string
	if stringutils.IsNotEmpty(schema.Description) {
		comments = append(comments, strings.Split(schema.Description, "\n")...)
	}
	return &astutils.FieldMeta{
		Name:     name,
		Type:     toGoType(schema),
		Comments: comments,
	}
}

func parameter2Field(param v3.Parameter) astutils.FieldMeta {
	var comments []string
	if stringutils.IsNotEmpty(param.Description) {
		comments = append(comments, strings.Split(param.Description, "\n")...)
	}
	t := toGoType(param.Schema)
	if param.Required {
		comments = append(comments, "required")
	} else {
		t = toOptional(t)
	}
	return astutils.FieldMeta{
		Name:     param.Name,
		Type:     t,
		Comments: comments,
	}
}

// toGoType converts schema to golang type
//	IntegerT Type = "integer"
//	StringT  Type = "string"
//	BooleanT Type = "boolean"
//	NumberT  Type = "number"
//	ObjectT  Type = "object"
//	ArrayT   Type = "array"
func toGoType(schema *v3.Schema) string {
	if stringutils.IsNotEmpty(schema.Ref) {
		return clean(strings.TrimPrefix(schema.Ref, "#/components/schemas/"))
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
		return object2Struct(schema)
	case v3.ArrayT:
		return "[]" + toGoType(schema.Items)
	default:
		return "interface{}"
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

func object2Struct(schema *v3.Schema) string {
	if stringutils.IsNotEmpty(schema.Title) {
		if _, exists := schemas[schema.Title]; exists {
			return schema.Title
		}
	}
	if schema.AdditionalProperties != nil {
		if ap, ok := schema.AdditionalProperties.(*v3.Schema); ok {
			return "map[string]" + toGoType(ap)
		}
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
		if omitempty {
			jsontag += ",omitempty"
		}
		if sliceutils.StringContains(schema.Required, k) {
			b.WriteString(fmt.Sprintf("  %s %s `json:\"%s\" url:\"%s\"`\n", strcase.ToCamel(k), toGoType(v), jsontag, k))
		} else {
			b.WriteString(fmt.Sprintf("  %s %s `json:\"%s\" url:\"%s\"`\n", strcase.ToCamel(k), "*"+toGoType(v), jsontag, k))
		}
	}
	b.WriteString("}")
	return b.String()
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

func genGoVo(schemas map[string]v3.Schema, output, pkg string) {
	if err := os.MkdirAll(filepath.Dir(output), os.ModePerm); err != nil {
		panic(err)
	}
	funcMap := make(map[string]interface{})
	funcMap["toCamel"] = toCamel
	funcMap["toGoType"] = toGoType
	funcMap["toComment"] = toComment
	funcMap["toOptional"] = toOptional
	funcMap["stringContains"] = sliceutils.StringContains
	tpl, _ := template.New("vo.go.tmpl").Funcs(funcMap).Parse(votmpl)
	var sqlBuf bytes.Buffer
	_ = tpl.Execute(&sqlBuf, struct {
		Schemas map[string]v3.Schema
		Omit    bool
		Pkg     string
	}{
		Schemas: schemas,
		Omit:    omitempty,
		Pkg:     pkg,
	})
	source := strings.TrimSpace(sqlBuf.String())
	astutils.FixImport([]byte(source), output)
}

var schemas map[string]v3.Schema
var requestBodies map[string]v3.RequestBody
var responses map[string]v3.Response
var omitempty bool

// GenGoClient generate go http client code from OpenAPI3.0 json document
func GenGoClient(dir string, file string, omit bool, env, pkg string) {
	var (
		err       error
		f         *os.File
		clientDir string
		fi        os.FileInfo
		api       v3.API
		vofile    string
	)
	clientDir = filepath.Join(dir, pkg)
	if err = os.MkdirAll(clientDir, os.ModePerm); err != nil {
		panic(err)
	}
	api = loadAPI(file)
	schemas = api.Components.Schemas
	requestBodies = api.Components.RequestBodies
	responses = api.Components.Responses
	omitempty = omit
	svcmap := make(map[string]map[string]v3.Path)
	for endpoint, path := range api.Paths {
		svcname := strings.Split(strings.Trim(endpoint, "/"), "/")[0]
		if value, exists := svcmap[svcname]; exists {
			value[endpoint] = path
		} else {
			svcmap[svcname] = make(map[string]v3.Path)
			svcmap[svcname][endpoint] = path
		}
	}

	for svcname, paths := range svcmap {
		genGoHTTP(paths, svcname, clientDir, env, pkg)
	}

	vofile = filepath.Join(clientDir, "vo.go")
	fi, err = os.Stat(vofile)
	if err != nil && !os.IsNotExist(err) {
		panic(err)
	}
	if fi != nil {
		logrus.Warningln("file vo.go will be overwritten")
	}
	if f, err = os.Create(vofile); err != nil {
		panic(err)
	}
	defer f.Close()
	genGoVo(api.Components.Schemas, vofile, pkg)
}

func loadAPI(file string) v3.API {
	var (
		docfile *os.File
		err     error
		docraw  []byte
		api     v3.API
	)
	if strings.HasPrefix(file, "http") {
		link := file
		client := resty.New()
		client.SetRedirectPolicy(resty.FlexibleRedirectPolicy(15))
		root, _ := os.Getwd()
		client.SetOutputDirectory(root)
		filename := ".openapi3"
		_, err := client.R().
			SetOutput(filename).
			Get(link)
		if err != nil {
			panic(err)
		}
		file = filepath.Join(root, filename)
		defer os.Remove(file)
	}
	if docfile, err = os.Open(file); err != nil {
		panic(err)
	}
	defer func(docfile *os.File) {
		_ = docfile.Close()
	}(docfile)
	if docraw, err = ioutil.ReadAll(docfile); err != nil {
		panic(err)
	}
	if err = json.Unmarshal(docraw, &api); err != nil {
		panic(err)
	}
	return api
}

func toOptional(t string) string {
	if !strings.HasPrefix(t, "*") {
		return "*" + t
	}
	return t
}
