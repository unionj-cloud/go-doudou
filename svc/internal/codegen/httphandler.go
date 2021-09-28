package codegen

import (
	"bytes"
	"github.com/unionj-cloud/go-doudou/sliceutils"
	"github.com/unionj-cloud/go-doudou/stringutils"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/iancoleman/strcase"
	"github.com/sirupsen/logrus"
	"github.com/unionj-cloud/go-doudou/astutils"
)

var httpHandlerTmpl = `package httpsrv

import (
	"github.com/unionj-cloud/go-doudou/svc/config"
	ddmodel "github.com/unionj-cloud/go-doudou/svc/http/model"
	"net/http"
	"os"
)

type {{.Meta.Name}}Handler interface {
{{- range $m := .Meta.Methods }}
	{{$m.Name}}(w http.ResponseWriter, r *http.Request)
{{- end }}
}

func Routes(handler {{.Meta.Name}}Handler) []ddmodel.Route {
	return []ddmodel.Route{
		{{- range $m := .Meta.Methods }}
		{
			"{{$m.Name | routeName}}",
			"{{$m.Name | httpMethod}}",
			{{- if eq $.RoutePatternStrategy 1}}
			"/{{$.Meta.Name | lower}}/{{$m.Name | noSplitPattern}}",
			{{- else }}
			"/{{$m.Name | pattern}}",
			{{- end }}
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
	if sliceutils.StringContains(httpMethods, head) {
		splits = splits[1:]
	}
	clean := sliceutils.StringFilter(splits, func(item string) bool {
		return stringutils.IsNotEmpty(item)
	})
	return strings.Join(clean, "/")
}

func noSplitPattern(method string) string {
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

// GenHttpHandler generates http handler interface and routes
func GenHttpHandler(dir string, ic astutils.InterfaceCollector, routePatternStrategy int) {
	var (
		err         error
		handlerfile string
		f           *os.File
		tpl         *template.Template
		httpDir     string
		source      string
		sqlBuf      bytes.Buffer
		fi          os.FileInfo
	)
	httpDir = filepath.Join(dir, "transport/httpsrv")
	if err = os.MkdirAll(httpDir, os.ModePerm); err != nil {
		panic(err)
	}

	handlerfile = filepath.Join(httpDir, "handler.go")
	fi, err = os.Stat(handlerfile)
	if err != nil && !os.IsNotExist(err) {
		panic(err)
	}
	if fi != nil {
		logrus.Warningln("file handler.go will be overwrited")
	}
	if f, err = os.Create(handlerfile); err != nil {
		panic(err)
	}
	defer f.Close()

	funcMap := make(map[string]interface{})
	funcMap["httpMethod"] = httpMethod
	funcMap["routeName"] = routeName
	funcMap["pattern"] = pattern
	funcMap["noSplitPattern"] = noSplitPattern
	funcMap["lower"] = strings.ToLower
	if tpl, err = template.New("handler.go.tmpl").Funcs(funcMap).Parse(httpHandlerTmpl); err != nil {
		panic(err)
	}
	if err = tpl.Execute(&sqlBuf, struct {
		RoutePatternStrategy int
		Meta                 astutils.InterfaceMeta
	}{
		RoutePatternStrategy: routePatternStrategy,
		Meta:                 ic.Interfaces[0],
	}); err != nil {
		panic(err)
	}
	source = strings.TrimSpace(sqlBuf.String())
	astutils.FixImport([]byte(source), handlerfile)
}
