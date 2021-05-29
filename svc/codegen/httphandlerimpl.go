package codegen

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/iancoleman/strcase"
	"github.com/sirupsen/logrus"
	"github.com/unionj-cloud/go-doudou/astutils"
	"github.com/unionj-cloud/go-doudou/copier"
	"github.com/unionj-cloud/go-doudou/templateutils"
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
		funcMap["hasPrefix"] = templateutils.HasPrefix
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

{{- range $m := .Meta.Methods }}
	func (receiver *{{$.Meta.Name}}HandlerImpl) {{$m.Name}}(w http.ResponseWriter, r *http.Request) {
    	var (
			{{- range $p := $m.Params }}
			{{ $p.Name }} {{ $p.Type }}
			{{- end }}
			{{- range $r := $m.Results }}
			{{ $r.Name }} {{ $r.Type }}
			{{- end }}
		)
		{{- range $p := $m.Params }}
		{{ if or (hasPrefix $p.Type "vo.") (hasPrefix $p.Type "*vo.") (hasPrefix $p.Type "[]vo.") (hasPrefix $p.Type "[]*vo.") (contains $p.Type "map[")}}
		if err := json.NewDecoder(r.Body).Decode(&{{$p.Name}}); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		defer r.Body.Close()
		{{ else if eq $p.Type "context.Context" }}
		{{$p.Name}} = context.Background()
		{{ else if contains $p.Type "*multipart.FileHeader" }}
		if err := r.ParseMultipartForm(32 << 20); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		files := r.MultipartForm.File["{{$p.Name}}"]
		if len(files) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("no file received"))
			return
		}
		{{if contains $p.Type "["}}
		{{$p.Name}} = files
		{{else}}
		{{$p.Name}} = files[0]
		{{end}}
		{{ else if contains $p.Type "["}}
		{{$p.Name}} = r.Form("{{$p.Name}}")
		{{ else }}
		{{$p.Name}} = r.FormValue("{{$p.Name}}")
		{{ end }}
		{{- end }}
		{{range $i, $r := $m.Results }}{{- if $i}},{{end}}{{ $r.Name }}{{- end }} = receiver.{{$.Meta.Name | toLowerCamel}}.{{$m.Name}}(
			{{- range $p := $m.Params }}
			{{ $p.Name }},
			{{- end }}
		)
		{{- range $r := $m.Results }}
			{{ if eq $r.Type "error" }}
				if {{ $r.Name }} != nil {
					w.WriteHeader(http.StatusInternalServerError)
				}
			{{ end }}
		{{- end }}
		if err := json.NewEncoder(w).Encode(struct{
			{{- range $r := $m.Results }}
			{{ $r.Name | toCamel }} {{ $r.Type }} ` + "`" + `json:"{{ $r.Name | toLowerCamel }}"` + "`" + `
			{{- end }}
		}{
			{{- range $r := $m.Results }}
			{{ $r.Name | toCamel }}: {{ $r.Name }},
			{{- end }}
		}); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
		}
    }
{{- end }}

func New{{.Meta.Name}}Handler({{.Meta.Name | toLowerCamel}} {{.ServiceAlias}}.{{.Meta.Name}}) {{.Meta.Name}}Handler {
	return &{{.Meta.Name}}HandlerImpl{
		{{.Meta.Name | toLowerCamel}},
	}
}
`

var appendHttpHandlerImplTmpl = `
{{- range $m := .Meta.Methods }}
	func (receiver *{{$.Meta.Name}}HandlerImpl) {{$m.Name}}(w http.ResponseWriter, r *http.Request) {
    	var (
			{{- range $p := $m.Params }}
			{{ $p.Name }} {{ $p.Type }}
			{{- end }}
			{{- range $r := $m.Results }}
			{{ $r.Name }} {{ $r.Type }}
			{{- end }}
		)
		{{- range $p := $m.Params }}
		{{ if or (hasPrefix $p.Type "vo.") (hasPrefix $p.Type "*vo.") (hasPrefix $p.Type "[]vo.") (hasPrefix $p.Type "[]*vo.") (contains $p.Type "map[")}}
		if err := json.NewDecoder(r.Body).Decode(&{{$p.Name}}); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		defer r.Body.Close()
		{{ else if eq $p.Type "context.Context" }}
		{{$p.Name}} = context.Background()
		{{ else if contains $p.Type "*multipart.FileHeader" }}
		if err := r.ParseMultipartForm(32 << 20); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		files := r.MultipartForm.File["{{$p.Name}}"]
		if len(files) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("no file received"))
			return
		}
		{{if contains $p.Type "["}}
		{{$p.Name}} = files
		{{else}}
		{{$p.Name}} = files[0]
		{{end}}
		{{ else if contains $p.Type "["}}
		{{$p.Name}} = r.Form("{{$p.Name}}")
		{{ else }}
		{{$p.Name}} = r.FormValue("{{$p.Name}}")
		{{ end }}
		{{- end }}
		{{range $i, $r := $m.Results }}{{- if $i}},{{end}}{{ $r.Name }}{{- end }} = receiver.{{$.Meta.Name | toLowerCamel}}.{{$m.Name}}(
			{{- range $p := $m.Params }}
			{{ $p.Name }},
			{{- end }}
		)
		{{- range $r := $m.Results }}
			{{ if eq $r.Type "error" }}
				if {{ $r.Name }} != nil {
					w.WriteHeader(http.StatusInternalServerError)
				}
			{{ end }}
		{{- end }}
		if err := json.NewEncoder(w).Encode(struct{
			{{- range $r := $m.Results }}
			{{ $r.Name | toCamel }} {{ $r.Type }} ` + "`" + `json:"{{ $r.Name | toLowerCamel }}"` + "`" + `
			{{- end }}
		}{
			{{- range $r := $m.Results }}
			{{ $r.Name | toCamel }}: {{ $r.Name }},
			{{- end }}
		}); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
		}
    }
{{- end }}
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
		root, err := parser.ParseFile(fset, handlerimplfile, nil, 0)
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
	funcMap["hasPrefix"] = templateutils.HasPrefix
	funcMap["contains"] = strings.Contains
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
	astutils.FixImport(original, handlerimplfile)
}
