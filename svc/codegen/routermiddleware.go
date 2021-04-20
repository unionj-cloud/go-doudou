package codegen

import (
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"text/template"
)

var routerMwTmpl = `package router

import (
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

func logger(inner http.Handler, name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		inner.ServeHTTP(w, r)
		logrus.Infof(
			"%s\t%s\t%s\t%s\n",
			r.Method,
			r.RequestURI,
			name,
			time.Since(start),
		)
	})
}

func rest(inner http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		inner.ServeHTTP(w, r)
	})
}
`

func GenRouterMiddleware(dir string) {
	var (
		err       error
		mwfile    string
		f         *os.File
		tpl       *template.Template
		routerDir string
	)
	routerDir = filepath.Join(dir, "router")
	if err = os.MkdirAll(routerDir, os.ModePerm); err != nil {
		panic(err)
	}

	mwfile = filepath.Join(routerDir, "middleware.go")
	if _, err = os.Stat(mwfile); os.IsNotExist(err) {
		if f, err = os.Create(mwfile); err != nil {
			panic(err)
		}
		defer f.Close()

		if tpl, err = template.New("middleware.go.tmpl").Parse(routerMwTmpl); err != nil {
			panic(err)
		}
		if err = tpl.Execute(f, nil); err != nil {
			panic(err)
		}
	} else {
		logrus.Warnf("file %s already exists", mwfile)
	}
}
