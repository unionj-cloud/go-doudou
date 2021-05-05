package codegen

import (
	"bytes"
	"github.com/iancoleman/strcase"
	"github.com/sirupsen/logrus"
	"github.com/unionj-cloud/go-doudou/astutils"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

var httpHandlerTmpl = `package httpsrv

import (
	"net/http"
)

type {{.Name}}Handler interface {
{{- range $m := .Methods }}
	{{$m.Name}}(w http.ResponseWriter, r *http.Request)
{{- end }}
}

type route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

func routes(handler {{.Name}}Handler) []route {
	return []route{
		{{- range $m := .Methods }}
		{
			"{{$m.Name | routeName}}",
			"{{$m.Name | httpMethod}}",
			"/{{$m.Name | pattern}}",
			handler.{{$m.Name}},
		},
		{{- end }}
	}
}
`

func pattern(method string) string {
	httpMethods := []string{"GET", "POST", "PUT", "DELETE"}
	snake := strcase.ToSnake(method)
	splits := strings.Split(snake, "_")
	head := strings.ToUpper(splits[0])
	for _, m := range httpMethods {
		if head == m {
			return strings.ToLower(method[len(m):])
		}
	}
	return strings.ToLower(method)
}

func routeName(method string) string {
	httpMethods := []string{"GET", "POST", "PUT", "DELETE"}
	snake := strcase.ToSnake(method)
	splits := strings.Split(snake, "_")
	head := strings.ToUpper(splits[0])
	for _, m := range httpMethods {
		if head == m {
			return method[len(m):]
		}
	}
	return method
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

func GenHttpHandler(dir string, ic astutils.InterfaceCollector) {
	var (
		err         error
		handlerfile string
		f           *os.File
		tpl         *template.Template
		httpDir     string
		source      string
		sqlBuf      bytes.Buffer
	)
	httpDir = filepath.Join(dir, "transport/httpsrv")
	if err = os.MkdirAll(httpDir, os.ModePerm); err != nil {
		panic(err)
	}

	handlerfile = filepath.Join(httpDir, "handler.go")
	if _, err = os.Stat(handlerfile); os.IsNotExist(err) {
		if f, err = os.Create(handlerfile); err != nil {
			panic(err)
		}
		defer f.Close()

		funcMap := make(map[string]interface{})
		funcMap["httpMethod"] = httpMethod
		funcMap["routeName"] = routeName
		funcMap["pattern"] = pattern
		if tpl, err = template.New("handler.go.tmpl").Funcs(funcMap).Parse(httpHandlerTmpl); err != nil {
			panic(err)
		}
		if err = tpl.Execute(&sqlBuf, ic.Interfaces[0]); err != nil {
			panic(err)
		}
		source = strings.TrimSpace(sqlBuf.String())
		astutils.FixImport([]byte(source), handlerfile)
	} else {
		logrus.Warnf("file %s already exists", handlerfile)
	}
}
