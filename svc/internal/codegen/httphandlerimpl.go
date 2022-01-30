package codegen

import (
	"bufio"
	"bytes"
	"github.com/iancoleman/strcase"
	"github.com/sirupsen/logrus"
	"github.com/unionj-cloud/go-doudou/astutils"
	"github.com/unionj-cloud/go-doudou/copier"
	v3 "github.com/unionj-cloud/go-doudou/openapi/v3"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

var httpHandlerImpl = `package httpsrv

import (
	{{.ServiceAlias}} "{{.ServicePackage}}"
	"net/http"
)

type {{.Meta.Name}}HandlerImpl struct{
	{{.Meta.Name | toLowerCamel}} {{.ServiceAlias}}.{{.Meta.Name}}
}

{{- range $m := .Meta.Methods }}
	func (receiver *{{$.Meta.Name}}HandlerImpl) {{$m.Name}}(w http.ResponseWriter, r *http.Request) {
    	panic("implement me")
    }
{{- end }}

func New{{.Meta.Name}}Handler({{.Meta.Name | toLowerCamel}} {{.ServiceAlias}}.{{.Meta.Name}}) {{.Meta.Name}}Handler {
	return &{{.Meta.Name}}HandlerImpl{
		{{.Meta.Name | toLowerCamel}},
	}
}
`

// GenHttpHandlerImpl generates http handler implementations
func GenHttpHandlerImpl(dir string, ic astutils.InterfaceCollector) {
	var (
		err             error
		modfile         string
		modName         string
		firstLine       string
		handlerimplfile string
		f               *os.File
		tpl             *template.Template
		source          string
		buf             bytes.Buffer
		httpDir         string
	)
	httpDir = filepath.Join(dir, "transport/httpsrv")
	if err = os.MkdirAll(httpDir, os.ModePerm); err != nil {
		panic(err)
	}

	handlerimplfile = filepath.Join(httpDir, "handlerimpl.go")
	if _, err = os.Stat(handlerimplfile); os.IsNotExist(err) {
		modfile = filepath.Join(dir, "go.mod")
		if f, err = os.Open(modfile); err != nil {
			panic(err)
		}
		reader := bufio.NewReader(f)
		if firstLine, err = reader.ReadString('\n'); err != nil {
			panic(err)
		}
		modName = strings.TrimSpace(strings.TrimPrefix(firstLine, "module"))

		if f, err = os.Create(handlerimplfile); err != nil {
			panic(err)
		}
		defer f.Close()

		funcMap := make(map[string]interface{})
		funcMap["toLowerCamel"] = strcase.ToLowerCamel
		funcMap["toCamel"] = strcase.ToCamel
		if tpl, err = template.New("handlerimpl.go.tmpl").Funcs(funcMap).Parse(httpHandlerImpl); err != nil {
			panic(err)
		}
		if err = tpl.Execute(&buf, struct {
			ServicePackage string
			ServiceAlias   string
			VoPackage      string
			Meta           astutils.InterfaceMeta
		}{
			ServicePackage: modName,
			ServiceAlias:   ic.Package.Name,
			VoPackage:      modName + "/vo",
			Meta:           ic.Interfaces[0],
		}); err != nil {
			panic(err)
		}

		source = strings.TrimSpace(buf.String())
		astutils.FixImport([]byte(source), handlerimplfile)
	} else {
		logrus.Warnf("file %s already exists.", handlerimplfile)
	}
}

var appendHttpHandlerImplTmpl = `
{{- range $m := .Meta.Methods }}
	func (receiver *{{$.Meta.Name}}HandlerImpl) {{$m.Name}}(_writer http.ResponseWriter, _req *http.Request) {
    	var (
			{{- range $p := $m.Params }}
			{{ $p.Name }} {{ $p.Type }}
			{{- end }}
			{{- range $r := $m.Results }}
			{{ $r.Name }} {{ $r.Type }}
			{{- end }}
		)
		{{- $multipartFormParsed := false }}
		{{- $formParsed := false }}
		{{- range $p := $m.Params }}
		{{- if contains $p.Type "*multipart.FileHeader" }}
		{{- if not $multipartFormParsed }}
		if err := _req.ParseMultipartForm(32 << 20); err != nil {
			http.Error(_writer, err.Error(), http.StatusBadRequest)
			return
		}
		{{- $multipartFormParsed = true }}
		{{- end }}
		{{- if contains $p.Type "["}}
		{{$p.Name}} = _req.MultipartForm.File["{{$p.Name}}"]
		{{- else}}
		{{$p.Name}}Files := _req.MultipartForm.File["{{$p.Name}}"]
		if len({{$p.Name}}Files) > 0 {
			{{$p.Name}} = {{$p.Name}}Files[0]
		}
		{{- end}}
		{{- else if contains $p.Type "*v3.FileModel" }}
		{{- if not $multipartFormParsed }}
		if err := _req.ParseMultipartForm(32 << 20); err != nil {
			http.Error(_writer, err.Error(), http.StatusBadRequest)
			return
		}
		{{- $multipartFormParsed = true }}
		{{- end }}
		{{$p.Name}}FileHeaders := _req.MultipartForm.File["{{$p.Name}}"]
		{{- if contains $p.Type "["}}
		for _, _fh :=range {{$p.Name}}FileHeaders {
			_f, err := _fh.Open()
			if err != nil {
				http.Error(_writer, err.Error(), http.StatusBadRequest)
				return
			}
			{{$p.Name}} = append({{$p.Name}}, &v3.FileModel{
				Filename: _fh.Filename,
				Reader: _f,
			})
		}
		{{- else}}
		if len({{$p.Name}}FileHeaders) > 0 {
			_fh := {{$p.Name}}FileHeaders[0]
			_f, err := _fh.Open()
			if err != nil {
				http.Error(_writer, err.Error(), http.StatusBadRequest)
				return
			}
			{{$p.Name}} = &v3.FileModel{
				Filename: _fh.Filename,
				Reader: _f,
			}
		}
		{{- end}}
		{{- else if eq $p.Type "context.Context" }}
		{{$p.Name}} = _req.Context()
		{{- else if not (isBuiltin $p)}}
		if err := json.NewDecoder(_req.Body).Decode(&{{$p.Name}}); err != nil {
			http.Error(_writer, err.Error(), http.StatusBadRequest)
			return
		}
		defer _req.Body.Close()
		{{- else if contains $p.Type "["}}
		{{- if not $formParsed }}
		if err := _req.ParseForm(); err != nil {
			http.Error(_writer, err.Error(), http.StatusBadRequest)
			return
		}
		{{- $formParsed = true }}
		{{- end }}
		if _, exists := _req.Form["{{$p.Name}}"]; exists {
			{{- if $p.Type | isSupport }}
			if casted, err := cast.{{$p.Type | castFunc}}E(_req.Form["{{$p.Name}}"]); err != nil {
				http.Error(_writer, err.Error(), http.StatusBadRequest)
				return
			} else {
				{{$p.Name}} = casted
			}
			{{- else }}
			{{$p.Name}} = _req.Form["{{$p.Name}}"]
			{{- end }}
		} else {
			if _, exists := _req.Form["{{$p.Name}}[]"]; exists {
				{{- if $p.Type | isSupport }}
				if casted, err := cast.{{$p.Type | castFunc}}E(_req.Form["{{$p.Name}}[]"]); err != nil {
					http.Error(_writer, err.Error(), http.StatusBadRequest)
					return
				} else {
					{{$p.Name}} = casted
				}
				{{- else }}
				{{$p.Name}} = _req.Form["{{$p.Name}}[]"]
				{{- end }}
			}
		}
		{{- else }}
		{{- if not $formParsed }}
		if err := _req.ParseForm(); err != nil {
			http.Error(_writer, err.Error(), http.StatusBadRequest)
			return
		}
		{{- $formParsed = true }}
		{{- end }}
		if _, exists := _req.Form["{{$p.Name}}"]; exists {
			{{- if $p.Type | isSupport }}
			if casted, err := cast.{{$p.Type | castFunc}}E(_req.FormValue("{{$p.Name}}")); err != nil {
				http.Error(_writer, err.Error(), http.StatusBadRequest)
				return
			} else {
				{{$p.Name}} = casted
			}
			{{- else }}
			{{$p.Name}} = _req.FormValue("{{$p.Name}}")
			{{- end }}
		}
		{{- end }}
		{{- end }}
		{{ range $i, $r := $m.Results }}{{- if $i}},{{- end}}{{- $r.Name }}{{- end }} = receiver.{{$.Meta.Name | toLowerCamel}}.{{$m.Name}}(
			{{- range $p := $m.Params }}
			{{ $p.Name }},
			{{- end }}
		)
		{{- range $r := $m.Results }}
			{{- if eq $r.Type "error" }}
				if {{ $r.Name }} != nil {
					if errors.Is({{ $r.Name }}, context.Canceled) {
						http.Error(_writer, {{ $r.Name }}.Error(), http.StatusBadRequest)
					} else {
						http.Error(_writer, {{ $r.Name }}.Error(), http.StatusInternalServerError)
					}
					return
				}
			{{- end }}
		{{- end }}
		{{- $done := false }}
		{{- range $r := $m.Results }}
			{{- if eq $r.Type "*os.File" }}
				if {{$r.Name}} == nil {
					http.Error(_writer, "No file returned", http.StatusInternalServerError)
					return
				}
				defer {{$r.Name}}.Close()
				var _fi os.FileInfo
				_fi, _err := {{$r.Name}}.Stat()
				if _err != nil {
					http.Error(_writer, _err.Error(), http.StatusInternalServerError)
					return
				}
				_writer.Header().Set("Content-Disposition", "attachment; filename="+_fi.Name())
				_writer.Header().Set("Content-Type", "application/octet-stream")
				_writer.Header().Set("Content-Length", fmt.Sprintf("%d", _fi.Size()))
				io.Copy(_writer, {{$r.Name}})
				{{- $done = true }}	
			{{- end }}
		{{- end }}
		{{- if not $done }}
			if err := json.NewEncoder(_writer).Encode(struct{
				{{- range $r := $m.Results }}
				{{- if ne $r.Type "error" }}
				{{ $r.Name | toCamel }} {{ $r.Type }} ` + "`" + `json:"{{ $r.Name | convertCase }}{{if $.Omitempty}},omitempty{{end}}"` + "`" + `
				{{- end }}
				{{- end }}
			}{
				{{- range $r := $m.Results }}
				{{- if ne $r.Type "error" }}
				{{ $r.Name | toCamel }}: {{ $r.Name }},
				{{- end }}
				{{- end }}
			}); err != nil {
				http.Error(_writer, err.Error(), http.StatusInternalServerError)
				return
			}
		{{- end }}
    }
{{- end }}
`

var initHttpHandlerImplTmpl = `package httpsrv

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/unionj-cloud/cast"
	{{.ServiceAlias}} "{{.ServicePackage}}"
	"net/http"
	"{{.VoPackage}}"
	"github.com/pkg/errors"
)

type {{.Meta.Name}}HandlerImpl struct{
	{{.Meta.Name | toLowerCamel}} {{.ServiceAlias}}.{{.Meta.Name}}
}

` + appendHttpHandlerImplTmpl + `

func New{{.Meta.Name}}Handler({{.Meta.Name | toLowerCamel}} {{.ServiceAlias}}.{{.Meta.Name}}) {{.Meta.Name}}Handler {
	return &{{.Meta.Name}}HandlerImpl{
		{{.Meta.Name | toLowerCamel}},
	}
}
`
var castFuncMap = map[string]string{
	"bool":          "ToBool",
	"float64":       "ToFloat64",
	"float32":       "ToFloat32",
	"int64":         "ToInt64",
	"int32":         "ToInt32",
	"int16":         "ToInt16",
	"int8":          "ToInt8",
	"int":           "ToInt",
	"uint":          "ToUint",
	"uint8":         "ToUint8",
	"uint16":        "ToUint16",
	"uint32":        "ToUint32",
	"uint64":        "ToUint64",
	"[]interface{}": "ToSlice",
	"[]bool":        "ToBoolSlice",
	"[]string":      "ToStringSlice",
	"[]int":         "ToIntSlice",
}

func isSupport(t string) bool {
	_, exists := castFuncMap[strings.TrimLeft(t, "*")]
	return exists
}

func castFunc(t string) string {
	return castFuncMap[strings.TrimLeft(t, "*")]
}

// GenHttpHandlerImplWithImpl generates http handler implementation
// Parsed value from query string parameters or application/x-www-form-urlencoded form will be string type.
// You may need to convert the type by yourself.
func GenHttpHandlerImplWithImpl(dir string, ic astutils.InterfaceCollector, omitempty bool, caseconvertor func(string) string) {
	var (
		err             error
		modfile         string
		modName         string
		firstLine       string
		handlerimplfile string
		f               *os.File
		modf            *os.File
		tpl             *template.Template
		buf             bytes.Buffer
		httpDir         string
		fi              os.FileInfo
		tmpl            string
		meta            astutils.InterfaceMeta
	)
	httpDir = filepath.Join(dir, "transport/httpsrv")
	if err = os.MkdirAll(httpDir, os.ModePerm); err != nil {
		panic(err)
	}

	handlerimplfile = filepath.Join(httpDir, "handlerimpl.go")
	fi, err = os.Stat(handlerimplfile)
	if err != nil && !os.IsNotExist(err) {
		panic(err)
	}

	err = copier.DeepCopy(ic.Interfaces[0], &meta)
	if err != nil {
		panic(err)
	}
	if fi != nil {
		logrus.Warningln("New content will be append to handlerimpl.go file")
		if f, err = os.OpenFile(handlerimplfile, os.O_APPEND, os.ModePerm); err != nil {
			panic(err)
		}
		defer f.Close()
		tmpl = appendHttpHandlerImplTmpl

		unimplementedMethods(&meta, handlerimplfile)
	} else {
		if f, err = os.Create(handlerimplfile); err != nil {
			panic(err)
		}
		defer f.Close()
		tmpl = initHttpHandlerImplTmpl
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
	funcMap["contains"] = strings.Contains
	funcMap["isBuiltin"] = v3.IsBuiltin
	funcMap["isSupport"] = isSupport
	funcMap["castFunc"] = castFunc
	funcMap["convertCase"] = caseconvertor
	if tpl, err = template.New("handlerimpl.go.tmpl").Funcs(funcMap).Parse(tmpl); err != nil {
		panic(err)
	}
	if err = tpl.Execute(&buf, struct {
		ServicePackage string
		ServiceAlias   string
		VoPackage      string
		Meta           astutils.InterfaceMeta
		Omitempty      bool
	}{
		ServicePackage: modName,
		ServiceAlias:   ic.Package.Name,
		VoPackage:      modName + "/vo",
		Meta:           meta,
		Omitempty:      omitempty,
	}); err != nil {
		panic(err)
	}

	original, err := ioutil.ReadAll(f)
	if err != nil {
		panic(err)
	}

	original = append(original, buf.Bytes()...)
	astutils.FixImport(original, handlerimplfile)
}

func unimplementedMethods(meta *astutils.InterfaceMeta, handlerimplfile string) {
	fset := token.NewFileSet()
	root, err := parser.ParseFile(fset, handlerimplfile, nil, parser.ParseComments)
	if err != nil {
		panic(err)
	}
	sc := astutils.NewStructCollector(astutils.ExprString)
	ast.Walk(sc, root)
	if handlers, exists := sc.Methods[meta.Name+"HandlerImpl"]; exists {
		var notimplemented []astutils.MethodMeta
		for _, item := range meta.Methods {
			for _, handler := range handlers {
				if item.Name == handler.Name {
					goto L
				}
			}
			notimplemented = append(notimplemented, item)

		L:
		}

		meta.Methods = notimplemented
	}
}
