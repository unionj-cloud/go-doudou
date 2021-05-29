package codegen

import (
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"text/template"
)

var httpMwTmpl = `package httpsrv

import (
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

func Logger(inner http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		inner.ServeHTTP(w, r)
		logrus.Infof(
			"%s\t%s\t%s\n",
			r.Method,
			r.RequestURI,
			time.Since(start),
		)
	})
}

func Rest(inner http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		inner.ServeHTTP(w, r)
	})
}
`

func GenHttpMiddleware(dir string) {
	var (
		err     error
		mwfile  string
		f       *os.File
		tpl     *template.Template
		httpDir string
	)
	httpDir = filepath.Join(dir, "transport/httpsrv")
	if err = os.MkdirAll(httpDir, os.ModePerm); err != nil {
		panic(err)
	}

	mwfile = filepath.Join(httpDir, "middleware.go")
	if _, err = os.Stat(mwfile); os.IsNotExist(err) {
		if f, err = os.Create(mwfile); err != nil {
			panic(err)
		}
		defer f.Close()

		if tpl, err = template.New("middleware.go.tmpl").Parse(httpMwTmpl); err != nil {
			panic(err)
		}
		if err = tpl.Execute(f, nil); err != nil {
			panic(err)
		}
	} else {
		logrus.Warnf("file %s already exists", mwfile)
	}
}
