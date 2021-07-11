package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/iancoleman/strcase"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/unionj-cloud/go-doudou/astutils"
	v3 "github.com/unionj-cloud/go-doudou/openapi/v3"
	"github.com/unionj-cloud/go-doudou/sliceutils"
	"github.com/unionj-cloud/go-doudou/stringutils"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"text/template"
)

var votmpl = `package {{.Pkg}}

{{- range $k, $v := .Schemas }}
{{ $v.Description | toComment }}
type {{$k | toCamel}} struct {
{{- range $pk, $pv := $v.Properties }}
	{{ $pv.Description | toComment }}
	{{- if stringContains $v.Required $pk }}
	// required
	{{- end }}
	{{ $pk | toCamel}} {{$pv | toGoType }} ` + "`" + `json:"{{$pk}}{{if $.Omit}},omitempty{{end}}" url:"{{$pk}}"` + "`" + `
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
	_querystring "github.com/google/go-querystring/query"
	"github.com/unionj-cloud/go-doudou/fileutils"
	"github.com/unionj-cloud/go-doudou/stringutils"
	ddhttp "github.com/unionj-cloud/go-doudou/svc/http"
	"io"
	"mime/multipart"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

type {{.Meta.Name}}Client struct {
	provider ddhttp.IServiceProvider
	client   *resty.Client
}

func (receiver *{{.Meta.Name}}Client) SetProvider(provider ddhttp.IServiceProvider) {
	receiver.provider = provider
}

func (receiver *{{.Meta.Name}}Client) SetClient(client *resty.Client) {
	receiver.client = client
}

{{- range $m := .Meta.Methods }}
	{{- range $c := $m.Comments }}
	// {{$c}}
	{{- end }}
	func (receiver *{{$.Meta.Name}}Client) {{$m.Name}}(ctx context.Context, {{ range $i, $p := $m.Params}}
    {{- if $i}},{{end}}
	{{- range $c := $p.Comments }}
	// {{$c}}
	{{- end }}
    {{ $p.Name}} {{$p.Type}}
    {{- end }}) ({{(index $m.Results 0).Name}} {{(index $m.Results 0).Type}}, err error) {
		var (
			_server string
			_err error
		)
		if _server, _err = receiver.provider.SelectServer(); _err != nil {
			err = errors.Wrap(_err, "")
			return
		}

		_req := receiver.client.R()
		_req.SetContext(ctx)

		{{- if $m.QueryParams }}
			_queryParams, _ := _querystring.Values({{$m.QueryParams.Name}})
			_req.SetQueryParamsFromValues(_queryParams)
		{{- end }}
		{{- if $m.PathVars }}
			{{- range $p := $m.PathVars }}
				_req.SetPathParam("{{$p.Name}}", fmt.Sprintf("%v", {{$p.Name}}))
			{{- end }}
		{{- end }}
		{{- if $m.HeaderVars }}
			{{- range $p := $m.HeaderVars }}
				_req.SetHeader("{{$p.Name}}", fmt.Sprintf("%v", {{$p.Name}}))
			{{- end }}
		{{- end }}
		{{- if $m.BodyParams }}
			_bodyParams, _ := _querystring.Values({{$m.BodyParams.Name}})
			_req.SetFormDataFromValues(_bodyParams)
		{{- end }}
		{{- if $m.BodyJson }}
			_req.SetBody({{$m.BodyJson.Name}})
		{{- end }}
		{{- if $m.Files }}
			{{- range $p := $m.Files }}
				{{- if contains $p.Type "["}}
				for _, _fh := range {{$p.Name}} {
					_f, _err := _fh.Open()
					if _err != nil {
						err = errors.Wrap(_err, "")
						return
					}
					_req.SetFileReader("{{$p.Name}}", _fh.Filename, _f)
				}
				{{- else}}
				_f, _err := {{$p.Name}}.Open()
				if _err != nil {
					err = errors.Wrap(_err, "")
					return
				}
				_req.SetFileReader("{{$p.Name}}", {{$p.Name}}.Filename, _f)
				{{- end}}
			{{- end }}
		{{- end }}

		{{- range $r := $m.Results }}
			{{- if eq $r.Type "*os.File" }}
				_req.SetDoNotParseResponse(true)
			{{- end }}
		{{- end }}

		_resp, _err := _req.{{$m.Name | restyMethod}}(_server + "{{$m.Path}}")
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
				_output := config.GddOutput.Load()
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

	return svcClient
}
`

func toMethod(endpoint string) string {
	ret := strings.ReplaceAll(strings.ReplaceAll(endpoint, "{", ""), "}", "")
	ret = strings.ReplaceAll(strings.Trim(ret, "/"), "/", "_")
	return strcase.ToCamel(ret)
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

func genGoHttp(paths map[string]v3.Path, svcname, dir, env, pkg string) {
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		panic(err)
	}
	output := filepath.Join(dir, svcname+"client.go")
	fi, err := os.Stat(output)
	if err != nil && !os.IsNotExist(err) {
		panic(err)
	}
	if fi != nil {
		logrus.Warningln("file " + svcname + "client.go will be overwrited")
	}
	var f *os.File
	if f, err = os.Create(output); err != nil {
		panic(err)
	}
	defer f.Close()

	funcMap := make(map[string]interface{})
	funcMap["toCamel"] = strcase.ToCamel
	funcMap["contains"] = strings.Contains
	funcMap["restyMethod"] = restyMethod
	funcMap["toUpper"] = strings.ToUpper
	tpl, err := template.New("http.go.tmpl").Funcs(funcMap).Parse(httptmpl)
	if err != nil {
		panic(err)
	}
	var sqlBuf bytes.Buffer
	err = tpl.Execute(&sqlBuf, struct {
		Meta astutils.InterfaceMeta
		Env  string
		Pkg  string
	}{
		Meta: api2Interface(paths, svcname),
		Env:  env,
		Pkg:  pkg,
	})
	if err != nil {
		panic(err)
	}
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
	var results, pathvars, headervars, files, params []astutils.FieldMeta
	var bodyJson, bodyParams, qparams *astutils.FieldMeta
	var comments []string
	if stringutils.IsNotEmpty(operation.Summary) {
		comments = append(comments, strings.Split(operation.Summary, "\n")...)
	}
	if stringutils.IsNotEmpty(operation.Description) {
		comments = append(comments, strings.Split(operation.Description, "\n")...)
	}

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

	for _, item := range operation.Parameters {
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

	if len(qSchema.Properties) > 0 {
		qparams = schema2Field(&qSchema, "queryParams")
	}

	if httpMethod != "Get" && operation.RequestBody != nil {
		if stringutils.IsNotEmpty(operation.RequestBody.Ref) {
			// #/components/requestBodies/Raw3
			key := strings.TrimPrefix(operation.RequestBody.Ref, "#/components/requestBodies/")
			if requestBody, exists := requestBodies[key]; exists {
				operation.RequestBody = &requestBody
			} else {
				panic(fmt.Errorf("requestBody %s not exists", operation.RequestBody.Ref))
			}
		}

		content := operation.RequestBody.Content
		if content.Json != nil {
			bodyJson = schema2Field(content.Json.Schema, "bodyJson")
		} else if content.FormUrl != nil {
			bodyParams = schema2Field(content.FormUrl.Schema, "bodyParams")
		} else if content.FormData != nil {
			schema := *content.FormData.Schema
			if stringutils.IsNotEmpty(schema.Ref) {
				schema = schemas[strings.TrimPrefix(content.FormData.Schema.Ref, "#/components/schemas/")]
			}
			aSchema := v3.Schema{
				Type:       v3.ObjectT,
				Properties: make(map[string]*v3.Schema),
			}
			for k, v := range schema.Properties {
				var gotype string
				if v.Type == v3.StringT && v.Format == v3.BinaryF {
					gotype = "*multipart.FileHeader"
				} else if v.Type == v3.ArrayT && v.Items.Type == v3.StringT && v.Items.Format == v3.BinaryF {
					gotype = "[]*multipart.FileHeader"
				} else {
					gotype = toGoType(v)
				}
				if strings.TrimPrefix(gotype, "[]") == "*multipart.FileHeader" {
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
		} else if content.Stream != nil {
			files = append(files, astutils.FieldMeta{
				Name: "_uploadFile",
				Type: "*multipart.FileHeader",
			})
		}
	}

	if operation.Responses == nil {
		return astutils.MethodMeta{}, errors.Errorf("response definition not found in api %s %s", httpMethod, endpoint)
	}

	if operation.Responses.Resp200 == nil {
		return astutils.MethodMeta{}, errors.Errorf("200 response definition not found in api %s %s", httpMethod, endpoint)
	}

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
		return astutils.MethodMeta{}, errors.Errorf("200 response content definition not found in api %s %s", httpMethod, endpoint)
	}

	if content.Json != nil {
		results = append(results, *schema2Field(content.Json.Schema, "ret"))
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
		return astutils.MethodMeta{}, errors.Errorf("200 response content definition not support yet in api %s %s", httpMethod, endpoint)
	}

	if qparams != nil {
		params = append(params, *qparams)
	}

	params = append(params, pathvars...)
	params = append(params, headervars...)

	if bodyParams != nil {
		params = append(params, *bodyParams)
	}

	if bodyJson != nil {
		params = append(params, *bodyJson)
	}

	params = append(params, files...)

	return astutils.MethodMeta{
		Name:        httpMethod + toMethod(endpoint),
		Params:      params,
		Results:     results,
		PathVars:    pathvars,
		HeaderVars:  headervars,
		BodyParams:  bodyParams,
		BodyJson:    bodyJson,
		Files:       files,
		Comments:    comments,
		Path:        endpoint,
		QueryParams: qparams,
	}, nil
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
	if param.Required {
		comments = append(comments, "required")
	}
	return astutils.FieldMeta{
		Name:     param.Name,
		Type:     toGoType(param.Schema),
		Comments: comments,
	}
}

func toGoType(schema *v3.Schema) string {
	if stringutils.IsNotEmpty(schema.Ref) {
		return strings.TrimPrefix(schema.Ref, "#/components/schemas/")
	}
	// IntegerT Type = "integer"
	//	StringT  Type = "string"
	//	BooleanT Type = "boolean"
	//	NumberT  Type = "number"
	//	ObjectT  Type = "object"
	//	ArrayT   Type = "array"
	switch schema.Type {
	case v3.IntegerT:
		// Int32F    Format = "int32"
		//	Int64F    Format = "int64"
		//	FloatF    Format = "float"
		//	DoubleF   Format = "double"
		//	DateTimeF Format = "date-time"
		//	BinaryF   Format = "binary"
		switch schema.Format {
		case v3.Int32F:
			return "int"
		case v3.Int64F:
			return "int64"
		default:
			return "int"
		}
	case v3.StringT:
		switch schema.Format {
		case v3.DateTimeF:
			return "*time.Time"
		case v3.BinaryF:
			return "*os.File"
		default:
			return "string"
		}
	case v3.BooleanT:
		return "bool"
	case v3.NumberT:
		switch schema.Format {
		case v3.FloatF:
			return "float32"
		case v3.DoubleF:
			return "float64"
		default:
			return "float64"
		}
	case v3.ObjectT:
		if stringutils.IsNotEmpty(schema.Title) {
			if _, exists := schemas[schema.Title]; exists {
				return schema.Title
			}
		}
		if schema.AdditionalProperties != nil && !reflect.ValueOf(*schema.AdditionalProperties).IsZero() {
			return "map[string]" + toGoType(schema.AdditionalProperties)
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
			b.WriteString(fmt.Sprintf("  %s %s `json:\"%s\" url:\"%s\"`\n", strcase.ToCamel(k), toGoType(v), jsontag, k))
		}
		b.WriteString("}")
		return b.String()
	case v3.ArrayT:
		return "[]" + toGoType(schema.Items)
	default:
		return "interface{}"
	}
}

func toComment(comment string) string {
	if stringutils.IsEmpty(comment) {
		return ""
	}
	b := new(strings.Builder)
	lines := strings.Split(comment, "\n")
	for _, line := range lines {
		b.WriteString(fmt.Sprintf("// %s\n", line))
	}
	return strings.TrimSuffix(b.String(), "\n")
}

func genGoVo(schemas map[string]v3.Schema, output, pkg string) {
	if err := os.MkdirAll(filepath.Dir(output), os.ModePerm); err != nil {
		panic(err)
	}
	funcMap := make(map[string]interface{})
	funcMap["toCamel"] = strcase.ToCamel
	funcMap["toGoType"] = toGoType
	funcMap["toComment"] = toComment
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

func GenGoClient(dir string, file string, omit bool, env, pkg string) {
	var (
		err       error
		f         *os.File
		clientDir string
		fi        os.FileInfo
		api       v3.Api
		vofile    string
	)
	clientDir = filepath.Join(dir, pkg)
	if err = os.MkdirAll(clientDir, os.ModePerm); err != nil {
		panic(err)
	}

	api = loadApi(file)
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
		genGoHttp(paths, svcname, clientDir, env, pkg)
	}

	vofile = filepath.Join(clientDir, "vo.go")
	fi, err = os.Stat(vofile)
	if err != nil && !os.IsNotExist(err) {
		panic(err)
	}
	if fi != nil {
		logrus.Warningln("file vo.go will be overwrited")
	}
	if f, err = os.Create(vofile); err != nil {
		panic(err)
	}
	defer f.Close()
	genGoVo(api.Components.Schemas, vofile, pkg)
}

func loadApi(file string) v3.Api {
	var (
		docfile *os.File
		err     error
		docraw  []byte
		api     v3.Api
	)
	if docfile, err = os.Open(file); err != nil {
		panic(err)
	}
	defer docfile.Close()
	if docraw, err = ioutil.ReadAll(docfile); err != nil {
		panic(err)
	}
	json.Unmarshal(docraw, &api)
	return api
}
