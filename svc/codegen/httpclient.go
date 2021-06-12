package codegen

import (
	"bufio"
	"bytes"
	"github.com/iancoleman/strcase"
	"github.com/sirupsen/logrus"
	"github.com/unionj-cloud/go-doudou/astutils"
	"github.com/unionj-cloud/go-doudou/copier"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

var tmpl = `package client

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
	{{.ServiceAlias}} "{{.ServicePackage}}"
	"{{.VoPackage}}"
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
		{{- else if not (isSimple $p)}}
		_req.SetBody({{$p.Name}})
		{{- else if contains $p.Type "["}}
		for _, _item := range {{$p.Name}} {
			_urlValues.Add("{{$p.Name}}", _item)
		}
		{{- else }}
		_urlValues.Set("{{$p.Name}}", {{$p.Name}})
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
				_output := os.Getenv("OUTPUT")
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

func New{{.Meta.Name}}(opts ...{{.Meta.Name}}ClientOption) service.{{.Meta.Name}} {
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

func restyMethod(method string) string {
	return strings.Title(strings.ToLower(httpMethod(method)))
}

func GenGoClient(dir string, ic astutils.InterfaceCollector) {
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
	funcMap["isSimple"] = IsSimple
	funcMap["restyMethod"] = restyMethod
	if tpl, err = template.New("client.go.tmpl").Funcs(funcMap).Parse(tmpl); err != nil {
		panic(err)
	}
	if err = tpl.Execute(&sqlBuf, struct {
		ServicePackage string
		ServiceAlias   string
		VoPackage      string
		Meta           astutils.InterfaceMeta
	}{
		ServicePackage: modName,
		ServiceAlias:   ic.Package.Name,
		VoPackage:      modName + "/vo",
		Meta:           meta,
	}); err != nil {
		panic(err)
	}

	source = strings.TrimSpace(sqlBuf.String())
	astutils.FixImport([]byte(source), clientfile)
}
