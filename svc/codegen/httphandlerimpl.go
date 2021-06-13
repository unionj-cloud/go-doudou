package codegen

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/iancoleman/strcase"
	"github.com/sirupsen/logrus"
	"github.com/unionj-cloud/go-doudou/astutils"
	"github.com/unionj-cloud/go-doudou/copier"
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
		sqlBuf          bytes.Buffer
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
		if err = tpl.Execute(&sqlBuf, struct {
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

		source = strings.TrimSpace(sqlBuf.String())
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
		{{- range $p := $m.Params }}
		{{- if contains $p.Type "*multipart.FileHeader" }}
		if err := _req.ParseMultipartForm(32 << 20); err != nil {
			http.Error(_writer, err.Error(), http.StatusBadRequest)
			return
		}
		{{$p.Name}}Files := _req.MultipartForm.File["{{$p.Name}}"]
		{{- if contains $p.Type "["}}
		{{$p.Name}} = {{$p.Name}}Files
		{{- else}}
		if len({{$p.Name}}Files) > 0 {
			{{$p.Name}} = {{$p.Name}}Files[0]
		}
		{{- end}}
		{{- else if eq $p.Type "context.Context" }}
		{{$p.Name}} = _req.Context()
		{{- else if not (isSimple $p)}}
		if err := json.NewDecoder(_req.Body).Decode(&{{$p.Name}}); err != nil {
			http.Error(_writer, err.Error(), http.StatusBadRequest)
			return
		}
		defer _req.Body.Close()
		{{- else if contains $p.Type "["}}
		if err := _req.ParseForm(); err != nil {
			http.Error(_writer, err.Error(), http.StatusBadRequest)
			return
		}
		{{$p.Name}} = _req.Form["{{$p.Name}}"]
		{{- else }}
		{{$p.Name}} = _req.FormValue("{{$p.Name}}")
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
					if {{ $r.Name }} == context.Canceled {
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
				{{ $r.Name | toCamel }} {{ $r.Type }} ` + "`" + `json:"{{ $r.Name | toLowerCamel }},omitempty"` + "`" + `
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
	{{.ServiceAlias}} "{{.ServicePackage}}"
	"net/http"
	"{{.VoPackage}}"
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

// Parsed value from query string parameters or application/x-www-form-urlencoded form will be string type.
// You may need to convert the type by yourself.
func GenHttpHandlerImplWithImpl(dir string, ic astutils.InterfaceCollector) {
	var (
		err             error
		modfile         string
		modName         string
		firstLine       string
		handlerimplfile string
		f               *os.File
		modf            *os.File
		tpl             *template.Template
		sqlBuf          bytes.Buffer
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
		logrus.Warningln("New content will be append to file handlerimpl.go")
		if f, err = os.OpenFile(handlerimplfile, os.O_APPEND, 0666); err != nil {
			panic(err)
		}
		defer f.Close()
		tmpl = appendHttpHandlerImplTmpl

		fset := token.NewFileSet()
		root, err := parser.ParseFile(fset, handlerimplfile, nil, parser.ParseComments)
		if err != nil {
			panic(err)
		}
		sc := astutils.NewStructCollector()
		ast.Walk(sc, root)
		fmt.Println(sc.Structs)

		if handlers, exists := sc.Methods[meta.Name+"HandlerImpl"]; exists {
			var notimplemented []astutils.MethodMeta
			for _, item := range meta.Methods {
				for _, handler := range handlers {
					if len(handler.Params) != 2 {
						continue
					}
					if handler.Params[0].Type == "http.ResponseWriter" && handler.Params[1].Type == "*http.Request" {
						if item.Name == handler.Name {
							goto L
						}
					}
				}
				notimplemented = append(notimplemented, item)

			L:
			}

			meta.Methods = notimplemented
		}
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
	funcMap["isSimple"] = IsSimple
	if tpl, err = template.New("handlerimpl.go.tmpl").Funcs(funcMap).Parse(tmpl); err != nil {
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

	original, err := ioutil.ReadAll(f)
	if err != nil {
		panic(err)
	}

	original = append(original, sqlBuf.Bytes()...)
	//fmt.Println(string(original))
	astutils.FixImport(original, handlerimplfile)
}
