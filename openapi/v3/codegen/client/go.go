package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/iancoleman/strcase"
	"github.com/sirupsen/logrus"
	"github.com/unionj-cloud/go-doudou/astutils"
	v3 "github.com/unionj-cloud/go-doudou/openapi/v3"
	"github.com/unionj-cloud/go-doudou/stringutils"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

var votmpl = `package client

//go:generate go-doudou name --file $GOFILE

{{- range $k, $v := .Schemas }}
{{ $v.Description | toComment }}
type {{$k | toCamel}} struct {
{{- range $pk, $pv := $v.Properties }}
	{{ $pv.Description | toComment }}
	{{ $pk | toCamel}} {{$pv | toGoType }}
{{- end }}
}
{{- end }}
`

var httptmpl = `package client

import (
	"context"
	"encoding/json"
	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
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

{{- range $m := .Meta.Methods }}
	func (receiver *{{$.Meta.Name}}Client) {{$m.Name}}({{- range $i, $p := $m.Params}}
    {{- if $i}},{{end}}
    {{- $p.Name}} {{$p.Type}}
    {{- end }}) ({{- range $i, $r := $m.Results}}
                     {{- if $i}},{{end}}
                     {{- $r.Name}} {{$r.Type}}
                     {{- end }}) {
		var (
			_server string
			_err error
		)
		if _server, _err = receiver.provider.SelectServer(); _err != nil {
			{{- range $r := $m.Results }}
				{{- if eq $r.Type "error" }}
					{{ $r.Name }} = errors.Wrap(_err, "")
				{{- end }}
			{{- end }}
			return
		}
		_urlValues := url.Values{}
		_req := receiver.client.R()
		{{- range $p := $m.Params }}
		{{- if contains $p.Type "*multipart.FileHeader" }}
		{{- if contains $p.Type "["}}
		for _, _fh := range {{$p.Name}} {
			_f, _err := _fh.Open()
			if _err != nil {
				{{- range $r := $m.Results }}
					{{- if eq $r.Type "error" }}
						{{ $r.Name }} = errors.Wrap(_err, "")
					{{- end }}
				{{- end }}
				return
			}
			_req.SetFileReader("{{$p.Name}}", _fh.Filename, _f)
		}
		{{- else}}
		_f, _err := {{$p.Name}}.Open()
		if _err != nil {
			{{- range $r := $m.Results }}
				{{- if eq $r.Type "error" }}
					{{ $r.Name }} = errors.Wrap(_err, "")
				{{- end }}
			{{- end }}
			return
		}
		_req.SetFileReader("{{$p.Name}}", {{$p.Name}}.Filename, _f)
		{{- end}}
		{{- else if eq $p.Type "context.Context" }}
		_req.SetContext({{$p.Name}})
		{{- else if not (isBuiltin $p)}}
		_req.SetBody({{$p.Name}})
		{{- else if contains $p.Type "["}}
		for _, _item := range {{$p.Name}} {
			_urlValues.Add("{{$p.Name}}", fmt.Sprintf("%v", _item))
		}
		{{- else }}
		_urlValues.Set("{{$p.Name}}", fmt.Sprintf("%v", {{$p.Name}}))
		{{- end }}
		{{- end }}

		{{- range $r := $m.Results }}
			{{- if eq $r.Type "*os.File" }}
				_req.SetDoNotParseResponse(true)
			{{- end }}
		{{- end }}

		{{- if eq ($m.Name | httpMethod) "GET" }}
		_resp, _err := _req.SetQueryParamsFromValues(_urlValues).
			Get(_server + "/{{$.Meta.Name | lower}}/{{$m.Name | pattern}}")
		{{- else }}
		if _req.Body != nil {
			_req.SetQueryParamsFromValues(_urlValues)
		} else {
			_req.SetFormDataFromValues(_urlValues)
		}
		_resp, _err := _req.{{$m.Name | restyMethod}}(_server + "/{{$.Meta.Name | lower}}/{{$m.Name | pattern}}")
		{{- end }}
		if _err != nil {
			{{- range $r := $m.Results }}
				{{- if eq $r.Type "error" }}
					{{ $r.Name }} = errors.Wrap(_err, "")
				{{- end }}
			{{- end }}
			return
		}
		if _resp.IsError() {
			{{- range $r := $m.Results }}
				{{- if eq $r.Type "error" }}
					{{ $r.Name }} = errors.New(_resp.String())
				{{- end }}
			{{- end }}
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
					{{- range $r := $m.Results }}
						{{- if eq $r.Type "error" }}
							{{ $r.Name }} = errors.Wrap(_err, "")
						{{- end }}
					{{- end }}
					return
				}
				_outFile, _err := os.Create(_file)
				if _err != nil {
					{{- range $r := $m.Results }}
						{{- if eq $r.Type "error" }}
							{{ $r.Name }} = errors.Wrap(_err, "")
						{{- end }}
					{{- end }}
					return
				}
				defer _outFile.Close()
				defer _resp.RawBody().Close()
				_, _err = io.Copy(_outFile, _resp.RawBody())
				if _err != nil {
					{{- range $r := $m.Results }}
						{{- if eq $r.Type "error" }}
							{{ $r.Name }} = errors.Wrap(_err, "")
						{{- end }}
					{{- end }}
					return
				}
				{{ $r.Name }} = _outFile
				return
				{{- $done = true }}	
			{{- end }}
		{{- end }}
		{{- if not $done }}
			var _result struct {
				{{- range $r := $m.Results }}
				{{- if eq $r.Type "error" }}
				{{ $r.Name | toCamel }} string ` + "`" + `json:"{{ $r.Name | toLowerCamel }}"` + "`" + `
				{{- else }}
				{{ $r.Name | toCamel }} {{ $r.Type }} ` + "`" + `json:"{{ $r.Name | toLowerCamel }}"` + "`" + `
				{{- end }}
				{{- end }}
			}
			if _err = json.Unmarshal(_resp.Body(), &_result); _err != nil {
				{{- range $r := $m.Results }}
					{{- if eq $r.Type "error" }}
						{{ $r.Name }} = errors.Wrap(_err, "")
					{{- end }}
				{{- end }}
				return
			}
			{{- range $r := $m.Results }}
				{{- if eq $r.Type "error" }}
					if stringutils.IsNotEmpty(_result.{{ $r.Name | toCamel }}) {
						{{ $r.Name }} = errors.New(_result.{{ $r.Name | toCamel }})
						return
					}
				{{- end }}
			{{- end }}
			return {{range $i, $r := $m.Results }}{{- if $i}},{{end}}{{ if eq $r.Type "error" }}nil{{else}}_result.{{ $r.Name | toCamel }}{{end}}{{- end }}
		{{- end }}    
	}
{{- end }}

type {{.Meta.Name}}ClientOption func(*{{.Meta.Name}}Client)

func WithProvider(provider ddhttp.IServiceProvider) {{.Meta.Name}}ClientOption {
	return func(a *{{.Meta.Name}}Client) {
		a.provider = provider
	}
}

func WithClient(client *resty.Client) {{.Meta.Name}}ClientOption {
	return func(a *{{.Meta.Name}}Client) {
		a.client = client
	}
}

func New{{.Meta.Name}}(opts ...{{.Meta.Name}}ClientOption) *{{.Meta.Name}}Client {
	defaultProvider := ddhttp.NewServiceProvider("{{.Meta.Name}}")
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

func genGoHttp(api v3.Api, output, svcname string) {
	funcMap := make(map[string]interface{})
	funcMap["toMethod"] = toMethod
	funcMap["toLowerCamel"] = strcase.ToLowerCamel
	funcMap["toCamel"] = strcase.ToCamel
	funcMap["lower"] = strings.ToLower
	funcMap["contains"] = strings.Contains
	tpl, _ := template.New("http.go.tmpl").Funcs(funcMap).Parse(httptmpl)
	var sqlBuf bytes.Buffer
	_ = tpl.Execute(&sqlBuf, struct {
		Meta astutils.InterfaceMeta
	}{
		Meta: api2Interface(api, svcname),
	})
	source := strings.TrimSpace(sqlBuf.String())
	astutils.FixImport([]byte(source), output)
}

func api2Interface(api v3.Api, svcname string) astutils.InterfaceMeta {
	var meta astutils.InterfaceMeta
	meta.Name = svcname
	for endpoint, path := range api.Paths {
		if path.Get != nil {
			meta.Methods = append(meta.Methods, operation2Method("Get"+toMethod(endpoint), path.Get))
		}
		if path.Post != nil {
			meta.Methods = append(meta.Methods, operation2Method("Post"+toMethod(endpoint), path.Post))
		}
		if path.Put != nil {
			meta.Methods = append(meta.Methods, operation2Method("Put"+toMethod(endpoint), path.Put))
		}
		if path.Delete != nil {
			meta.Methods = append(meta.Methods, operation2Method("Delete"+toMethod(endpoint), path.Delete))
		}
	}
	return meta
}

func operation2Method(name string, operation *v3.Operation) astutils.MethodMeta {
	//var params, results, pathvars []astutils.FieldMeta
	var comments []string
	//if stringutils.IsNotEmpty(operation.Summary) {
	//	comments = append(comments, strings.Split(operation.Summary, "\n")...)
	//}
	//if stringutils.IsNotEmpty(operation.Description) {
	//	comments = append(comments, strings.Split(operation.Description, "\n")...)
	//}
	//
	//for _, item :=range operation.Parameters {
	//
	//}

	return astutils.MethodMeta{
		Name:     name,
		Params:   nil,
		Results:  nil,
		PathVars: nil,
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
			return schema.Title
		}
		if schema.AdditionalProperties != nil {
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
			b.WriteString(fmt.Sprintf("  %s %s\n", strcase.ToCamel(k), toGoType(v)))
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

func genGoVo(schemas map[string]v3.Schema, output string) {
	funcMap := make(map[string]interface{})
	funcMap["toCamel"] = strcase.ToCamel
	funcMap["toGoType"] = toGoType
	funcMap["toComment"] = toComment
	tpl, _ := template.New("vo.go.tmpl").Funcs(funcMap).Parse(votmpl)
	var sqlBuf bytes.Buffer
	_ = tpl.Execute(&sqlBuf, struct {
		Schemas map[string]v3.Schema
	}{
		Schemas: schemas,
	})
	source := strings.TrimSpace(sqlBuf.String())
	astutils.FixImport([]byte(source), output)
}

func GenGoClient(dir string, file string, svcname string) {
	var (
		err       error
		httpfile  string
		f         *os.File
		clientDir string
		fi        os.FileInfo
		api       v3.Api
		vofile    string
	)
	clientDir = filepath.Join(dir, "client")
	if err = os.MkdirAll(clientDir, 0644); err != nil {
		panic(err)
	}

	httpfile = filepath.Join(clientDir, "http.go")
	fi, err = os.Stat(httpfile)
	if err != nil && !os.IsNotExist(err) {
		panic(err)
	}
	if fi != nil {
		logrus.Warningln("file http.go will be overwrited")
	}
	if f, err = os.Create(httpfile); err != nil {
		panic(err)
	}
	defer f.Close()

	api = loadApi(file)

	genGoHttp(api, httpfile, svcname)

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
	genGoVo(api.Components.Schemas, vofile)
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
