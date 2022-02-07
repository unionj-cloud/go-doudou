package codegen

import (
	"bufio"
	"bytes"
	"github.com/iancoleman/strcase"
	"github.com/sirupsen/logrus"
	"github.com/unionj-cloud/go-doudou/astutils"
	"github.com/unionj-cloud/go-doudou/copier"
	v3 "github.com/unionj-cloud/go-doudou/openapi/v3"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

var clientTmpl = `package client

import (
	"context"
	"encoding/json"
	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
	"github.com/unionj-cloud/go-doudou/fileutils"
	"github.com/unionj-cloud/go-doudou/stringutils"
	"github.com/unionj-cloud/go-doudou/svc/registry"
	ddhttp "github.com/unionj-cloud/go-doudou/svc/http"
	v3 "github.com/unionj-cloud/go-doudou/openapi/v3"
	"github.com/opentracing-contrib/go-stdlib/nethttp"
	"github.com/opentracing/opentracing-go"
	"io"
	"net/http"
	"mime/multipart"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"{{.VoPackage}}"
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
	func (receiver *{{$.Meta.Name}}Client) {{$m.Name}}({{- range $i, $p := $m.Params}}
    {{- if $i}},{{end}}
    {{- $p.Name}} {{$p.Type}}
    {{- end }}) (_resp *resty.Response, {{- range $i, $r := $m.Results}}
                     {{- if $i}},{{end}}
                     {{- $r.Name}} {{$r.Type}}
                     {{- end }}) {
		var _err error
		_urlValues := url.Values{}
		_req := receiver.client.R()
		{{- range $p := $m.Params }}
		{{- if or (eq $p.Type "*multipart.FileHeader") (eq $p.Type "[]*multipart.FileHeader") }}
		{{- if contains $p.Type "["}}
		for _, _fh := range {{$p.Name}} {
			_f, _err := _fh.Open()
			if _err != nil {
				{{- range $r := $m.Results }}
					{{- if eq $r.Type "error" }}
				{{ $r.Name }} = errors.Wrap(_err, "error")
					{{- end }}
				{{- end }}
				return
			}
			_req.SetFileReader("{{$p.Name}}", _fh.Filename, _f)
		}
		{{- else}}
		if _f, _err := {{$p.Name}}.Open(); _err != nil {
			{{- range $r := $m.Results }}
				{{- if eq $r.Type "error" }}
			{{ $r.Name }} = errors.Wrap(_err, "error")
				{{- end }}
			{{- end }}
			return
		} else {
			_req.SetFileReader("{{$p.Name}}", {{$p.Name}}.Filename, _f)
		}
		{{- end}}
		{{- else if or (eq $p.Type "v3.FileModel") (eq $p.Type "*v3.FileModel") (eq $p.Type "[]v3.FileModel") (eq $p.Type "*[]v3.FileModel") (eq $p.Type "...v3.FileModel") }}
		{{- if isSlice $p.Type }}
		{{- if isOptional $p.Type }}
		if {{$p.Name}} != nil {
			{{- if isVarargs $p.Type }}
			for _, _f := range {{$p.Name}} {
			{{- else }}
			for _, _f := range *{{$p.Name}} {
			{{- end }}
				_req.SetFileReader("{{$p.Name}}", _f.Filename, _f.Reader)
			}
		}
		{{- else }}
		if len({{$p.Name}}) == 0 {
			{{- range $r := $m.Results }}
				{{- if eq $r.Type "error" }}
			{{ $r.Name }} = errors.New("at least one file should be uploaded for parameter {{$p.Name}}")
				{{- end }}
			{{- end }}
			return
		}
		for _, _f := range {{$p.Name}} {
			_req.SetFileReader("{{$p.Name}}", _f.Filename, _f.Reader)
		}
		{{- end }}
		{{- else }}
		{{- if isOptional $p.Type }}
		if {{$p.Name}} != nil { 
			_req.SetFileReader("{{$p.Name}}", {{$p.Name}}.Filename, {{$p.Name}}.Reader)
		}
		{{- else }}
		_req.SetFileReader("{{$p.Name}}", {{$p.Name}}.Filename, {{$p.Name}}.Reader)
		{{- end }}
		{{- end }}
		{{- else if eq $p.Type "context.Context" }}
		_req.SetContext({{$p.Name}})
		{{- else if not (isBuiltin $p)}}
		_req.SetBody({{$p.Name}})
		{{- else if isSlice $p.Type }}
		{{- if isOptional $p.Type }}
		if {{$p.Name}} != nil { 
			{{- if isVarargs $p.Type }}
			for _, _item := range {{$p.Name}} {
			{{- else }}
			for _, _item := range *{{$p.Name}} {
			{{- end }}
				_urlValues.Add("{{$p.Name}}", fmt.Sprintf("%v", _item))
			}
		}
		{{- else }}
		if len({{$p.Name}}) == 0 {
			{{- range $r := $m.Results }}
				{{- if eq $r.Type "error" }}
			{{ $r.Name }} = errors.New("size of parameter {{$p.Name}} should be greater than zero")
				{{- end }}
			{{- end }}
			return
		}
		for _, _item := range {{$p.Name}} {
			_urlValues.Add("{{$p.Name}}", fmt.Sprintf("%v", _item))
		}
		{{- end}}
		{{- else }}
		{{- if isOptional $p.Type }}
		if {{$p.Name}} != nil { 
			_urlValues.Set("{{$p.Name}}", fmt.Sprintf("%v", *{{$p.Name}}))
		}
		{{- else }}
		_urlValues.Set("{{$p.Name}}", fmt.Sprintf("%v", {{$p.Name}}))
		{{- end }}
		{{- end }}
		{{- end }}

		{{- range $r := $m.Results }}
			{{- if eq $r.Type "*os.File" }}
				_req.SetDoNotParseResponse(true)
			{{- end }}
		{{- end }}

		{{- if eq $.RoutePatternStrategy 1}}
		_path := "/{{$.Meta.Name | lower}}/{{$m.Name | noSplitPattern}}"
		{{- else }}
		_path := "/{{$m.Name | pattern}}"
		{{- end }}

		{{- if eq ($m.Name | httpMethod) "GET" }}
		_resp, _err = _req.SetQueryParamsFromValues(_urlValues).
			Get(_path)
		{{- else }}
		if _req.Body != nil {
			_req.SetQueryParamsFromValues(_urlValues)
		} else {
			_req.SetFormDataFromValues(_urlValues)
		}
		_resp, _err = _req.{{$m.Name | restyMethod}}(_path)
		{{- end }}
		if _err != nil {
			{{- range $r := $m.Results }}
				{{- if eq $r.Type "error" }}
			{{ $r.Name }} = errors.Wrap(_err, "error")
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
					{{ $r.Name }} = errors.Wrap(_err, "error")
						{{- end }}
					{{- end }}
					return
				}
				_outFile, _err := os.Create(_file)
				if _err != nil {
					{{- range $r := $m.Results }}
						{{- if eq $r.Type "error" }}
					{{ $r.Name }} = errors.Wrap(_err, "error")
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
					{{ $r.Name }} = errors.Wrap(_err, "error")
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
				{{- if ne $r.Type "error" }}
				{{ $r.Name | toCamel }} {{ $r.Type }} ` + "`" + `json:"{{ $r.Name | convertCase }}"` + "`" + `
				{{- end }}
				{{- end }}
			}
			if _err = json.Unmarshal(_resp.Body(), &_result); _err != nil {
				{{- range $r := $m.Results }}
					{{- if eq $r.Type "error" }}
				{{ $r.Name }} = errors.Wrap(_err, "error")
					{{- end }}
				{{- end }}
				return
			}
			return _resp, {{range $i, $r := $m.Results }}{{- if $i}},{{end}}{{ if eq $r.Type "error" }}nil{{else}}_result.{{ $r.Name | toCamel }}{{end}}{{- end }}
		{{- end }}    
	}
{{- end }}

func New{{.Meta.Name}}Client(opts ...ddhttp.DdClientOption) *{{.Meta.Name}}Client {
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
			nethttp.OperationName(fmt.Sprintf("HTTP %s: %s", request.Method, request.RequestURI)))
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

func restyMethod(method string) string {
	return strings.Title(strings.ToLower(httpMethod(method)))
}

// GenGoClient generates golang http client code from result of parsing svc.go file in project root path
func GenGoClient(dir string, ic astutils.InterfaceCollector, env string, routePatternStrategy int, caseconvertor func(string) string) {
	var (
		err        error
		clientfile string
		f          *os.File
		tpl        *template.Template
		sqlBuf     bytes.Buffer
		clientDir  string
		fi         os.FileInfo
		source     string
		modfile    string
		modName    string
		firstLine  string
		modf       *os.File
		meta       astutils.InterfaceMeta
	)
	clientDir = filepath.Join(dir, "client")
	if err = os.MkdirAll(clientDir, os.ModePerm); err != nil {
		panic(err)
	}

	clientfile = filepath.Join(clientDir, "client.go")
	fi, err = os.Stat(clientfile)
	if err != nil && !os.IsNotExist(err) {
		panic(err)
	}
	if fi != nil {
		logrus.Warningln("file client.go will be overwrited")
	}
	if f, err = os.Create(clientfile); err != nil {
		panic(err)
	}
	defer f.Close()

	err = copier.DeepCopy(ic.Interfaces[0], &meta)
	if err != nil {
		panic(err)
	}

	modfile = filepath.Join(dir, "go.mod")
	if modf, err = os.Open(modfile); err != nil {
		panic(err)
	}
	reader := bufio.NewReader(modf)
	if firstLine, err = reader.ReadString('\n'); err != nil {
		panic(err)
	}
	modName = strings.TrimSpace(strings.TrimPrefix(firstLine, "module"))

	funcMap := make(map[string]interface{})
	funcMap["toLowerCamel"] = strcase.ToLowerCamel
	funcMap["toCamel"] = strcase.ToCamel
	funcMap["httpMethod"] = httpMethod
	funcMap["pattern"] = pattern
	funcMap["lower"] = strings.ToLower
	funcMap["contains"] = strings.Contains
	funcMap["isBuiltin"] = v3.IsBuiltin
	funcMap["restyMethod"] = restyMethod
	funcMap["toUpper"] = strings.ToUpper
	funcMap["noSplitPattern"] = noSplitPattern
	funcMap["isOptional"] = v3.IsOptional
	funcMap["convertCase"] = caseconvertor
	funcMap["isSlice"] = v3.IsSlice
	funcMap["isVarargs"] = v3.IsVarargs
	if tpl, err = template.New("client.go.tmpl").Funcs(funcMap).Parse(clientTmpl); err != nil {
		panic(err)
	}
	if err = tpl.Execute(&sqlBuf, struct {
		VoPackage            string
		Meta                 astutils.InterfaceMeta
		Env                  string
		RoutePatternStrategy int
	}{
		VoPackage:            modName + "/vo",
		Meta:                 meta,
		Env:                  env,
		RoutePatternStrategy: routePatternStrategy,
	}); err != nil {
		panic(err)
	}

	source = strings.TrimSpace(sqlBuf.String())
	astutils.FixImport([]byte(source), clientfile)
}
