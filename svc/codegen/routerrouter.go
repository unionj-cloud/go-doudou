package codegen

import (
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"text/template"
)

var routerRouterTmpl = `package router

import (
	"github.com/gorilla/mux"
	"net/http"
)

type route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

func NewRouter() *mux.Router {
	rous := routes()
	router := mux.NewRouter().StrictSlash(true)
	for _, r := range rous {
		var handler http.Handler

		handler = r.HandlerFunc
		handler = logger(handler, r.Name)
		handler = rest(handler)

		router.
			Methods(r.Method).
			Path(r.Pattern).
			Name(r.Name).
			Handler(handler)
	}
	return router
}
`

func GenRouterRouter(dir string) {
	var (
		err        error
		routerfile string
		f          *os.File
		tpl        *template.Template
		routerDir  string
	)
	routerDir = filepath.Join(dir, "router")
	if err = os.MkdirAll(routerDir, os.ModePerm); err != nil {
		panic(err)
	}

	routerfile = filepath.Join(routerDir, "router.go")
	if _, err = os.Stat(routerfile); os.IsNotExist(err) {
		if f, err = os.Create(routerfile); err != nil {
			panic(err)
		}
		defer f.Close()

		if tpl, err = template.New("router.go.tmpl").Parse(routerRouterTmpl); err != nil {
			panic(err)
		}
		if err = tpl.Execute(f, nil); err != nil {
			panic(err)
		}
	} else {
		logrus.Warnf("file %s already exists", routerfile)
	}
}
