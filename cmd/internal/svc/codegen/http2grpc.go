package codegen

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/iancoleman/strcase"
	"github.com/sirupsen/logrus"
	"github.com/unionj-cloud/go-doudou/v2/cmd/internal/templates"
	"github.com/unionj-cloud/toolkit/astutils"
	"github.com/unionj-cloud/toolkit/copier"
	"github.com/unionj-cloud/go-doudou/v2/version"
)

var appendHttp2GrpcTmpl = `
{{- range $m := .Meta.Methods }}
	func (receiver *{{$.Meta.Name}}Http2Grpc) {{$m.Name}}(_writer http.ResponseWriter, _req *http.Request) {
    	var (
			{{- range $p := $m.Params }}
			{{- if eq $p.Type "context.Context"}}
			{{ $p.Name }} {{ $p.Type }}
			{{- else }}
			{{$p.Name}} = new({{ trimLeft $p.Type "*" | replacePkg }})
			{{- end }}
			{{- end }}
			{{- range $r := $m.Results }}
			{{- if eq $r.Type "error"}}
			{{ $r.Name }} {{ $r.Type }}
			{{- else }}
			{{ $r.Name }} = new({{ trimLeft $r.Type "*" | replacePkg }})
			{{- end }}
			{{- end }}
		)
		{{- range $p := $m.Params }}
		{{- if eq $p.Type "context.Context" }}
		{{$p.Name}} = _req.Context()
		{{- else if and (eq $m.HttpMethod "GET") (not $.Config.AllowGetWithReqBody)}}
		if _err := _req.ParseForm(); _err != nil {
			rest.PanicBadRequestErr(_err)
		}
		if _err := rest.DecodeForm({{ $p.Name }}, _req.Form); _err != nil {
			rest.PanicBadRequestErr(_err)
		}
		if _err := rest.ValidateStruct({{ $p.Name }}); _err != nil {
			rest.PanicBadRequestErr(_err)
		}
		{{- else }}
		if _err := json.NewDecoder(_req.Body).Decode({{$p.Name}}); _err != nil {
			rest.PanicBadRequestErr(_err)	
		}
		if _err := rest.ValidateStruct({{$p.Name}}); _err != nil {
			rest.PanicBadRequestErr(_err)
		}
		{{- end }}
		{{- end }}
		{{- if eq (len $m.Results) 1 }}
		_, {{ range $i, $r := $m.Results }}{{- $r.Name }}{{- end }} = receiver.{{$.Meta.Name | toLowerCamel}}.{{$m.Name}}Rpc(
		{{- else }}
		{{ range $i, $r := $m.Results }}{{- if $i}},{{- end}}{{- $r.Name }}{{- end }} = receiver.{{$.Meta.Name | toLowerCamel}}.{{$m.Name}}Rpc(
		{{- end }}
			{{- if eq (len $m.Params) 1 }}
			{{- range $p := $m.Params }}
			{{ $p.Name }},
			{{- end }}
			new(emptypb.Empty),
			{{- else }}
			{{- range $p := $m.Params }}
			{{- if eq $p.Type "context.Context"}}
			{{ $p.Name }},
			{{- else }}
			{{$p.Name}},
			{{- end }}
			{{- end }}
			{{- end }}
		)
		{{- range $r := $m.Results }}
			{{- if eq $r.Type "error" }}
				if {{ $r.Name }} != nil {
					panic({{ $r.Name }})
				}
			{{- end }}
		{{- end }}
		_writer.Header().Set("Content-Type", "application/json; charset=UTF-8")
		{{- if eq (len $m.Results) 1 }}
			if _err := json.NewEncoder(_writer).Encode(struct {}{}); _err != nil {
				rest.HandleInternalServerError(_err)
			}
		{{- else }}
			if _err := json.NewEncoder(_writer).Encode(struct {
				{{- range $r := $m.Results }}
				{{- if ne $r.Type "error" }}
				{{ $r.Name | toCamel }} *{{ trimLeft $r.Type "*" | replacePkg }} ` + "`" + `json:"{{ $r.Name | convertCase }}{{if $.Config.Omitempty}},omitempty{{end}}"` + "`" + `
				{{- end }}
				{{- end }}
			}{
				{{- range $r := $m.Results }}
				{{- if ne $r.Type "error" }}
				{{ $r.Name | toCamel }}: {{ $r.Name }},
				{{- end }}
				{{- end }}
			}); _err != nil {
				rest.HandleInternalServerError(_err)
			}
		{{- end }}
    }
{{- end }}
`

var initHttp2GrpcTmpl = templates.EditableHeaderTmpl + `package httpsrv

import ()

var json = sonic.ConfigDefault

type {{.Meta.Name}}Http2Grpc struct{
	{{.Meta.Name | toLowerCamel}} pb.{{.Meta.Name}}ServiceServer
}

` + appendHttp2GrpcTmpl + `

func New{{.Meta.Name}}Http2Grpc({{.Meta.Name | toLowerCamel}} pb.{{.Meta.Name}}ServiceServer) {{.Meta.Name}}Handler {
	return &{{.Meta.Name}}Http2Grpc{
		{{.Meta.Name | toLowerCamel}},
	}
}
`

var importHttp2GrpcTmpl = `
	"context"
	"github.com/bytedance/sonic"
	"github.com/unionj-cloud/go-doudou/v2/framework/rest"
	"net/http"
	pb "{{.TransportGrpcPackage}}"
`

type GenHttp2GrpcConfig struct {
	Omitempty           bool
	AllowGetWithReqBody bool
	CaseConvertor       func(string) string
}

// GenHttp2Grpc generates http handler implementation
// Parsed value from query string parameters or application/x-www-form-urlencoded form will be string type.
// You may need to convert the type by yourself.
func GenHttp2Grpc(dir string, ic astutils.InterfaceCollector, config GenHttp2GrpcConfig) {
	var (
		err             error
		handlerimplfile string
		f               *os.File
		tpl             *template.Template
		buf             bytes.Buffer
		httpDir         string
		fi              os.FileInfo
		tmpl            string
		meta            astutils.InterfaceMeta
		importBuf       bytes.Buffer
	)
	httpDir = filepath.Join(dir, "transport", "httpsrv")
	if err = os.MkdirAll(httpDir, os.ModePerm); err != nil {
		panic(err)
	}
	handlerimplfile = filepath.Join(httpDir, "http2grpc.go")
	fi, err = os.Stat(handlerimplfile)
	if err != nil && !os.IsNotExist(err) {
		panic(err)
	}
	err = copier.DeepCopy(ic.Interfaces[0], &meta)
	if err != nil {
		panic(err)
	}
	unimplementedMethods(&meta, httpDir, meta.Name+"Http2Grpc")
	if fi != nil {
		logrus.Warningln("New content will be append to http2grpc.go file")
		if f, err = os.OpenFile(handlerimplfile, os.O_APPEND, os.ModePerm); err != nil {
			panic(err)
		}
		defer f.Close()
		tmpl = appendHttp2GrpcTmpl
	} else {
		if f, err = os.Create(handlerimplfile); err != nil {
			panic(err)
		}
		defer f.Close()
		tmpl = initHttp2GrpcTmpl
	}

	funcMap := make(map[string]interface{})
	funcMap["toLowerCamel"] = strcase.ToLowerCamel
	funcMap["replacePkg"] = func(input string) string {
		return "pb" + input[strings.LastIndex(input, "."):]
	}
	funcMap["trimLeft"] = strings.TrimLeft
	funcMap["toCamel"] = strcase.ToCamel
	funcMap["convertCase"] = config.CaseConvertor
	if tpl, err = template.New(tmpl).Funcs(funcMap).Parse(tmpl); err != nil {
		panic(err)
	}
	if err = tpl.Execute(&buf, struct {
		Meta    astutils.InterfaceMeta
		Config  GenHttp2GrpcConfig
		Version string
	}{
		Meta:    meta,
		Config:  config,
		Version: version.Release,
	}); err != nil {
		panic(err)
	}
	original, err := ioutil.ReadAll(f)
	if err != nil {
		panic(err)
	}
	original = append(original, buf.Bytes()...)
	if tpl, err = template.New(importHttp2GrpcTmpl).Parse(importHttp2GrpcTmpl); err != nil {
		panic(err)
	}
	transGrpcPkg := astutils.GetPkgPath(filepath.Join(dir, "transport", "grpc"))
	if err = tpl.Execute(&importBuf, struct {
		TransportGrpcPackage string
	}{
		TransportGrpcPackage: transGrpcPkg,
	}); err != nil {
		panic(err)
	}
	original = astutils.AppendImportStatements(original, importBuf.Bytes())
	astutils.FixImport(original, handlerimplfile)
}
