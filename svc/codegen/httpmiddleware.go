package codegen

import (
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"text/template"
)

var httpMwTmpl = `package httpsrv

import (
	"github.com/felixge/httpsnoop"
	"github.com/sirupsen/logrus"
	"github.com/unionj-cloud/go-doudou/stringutils"
	"net/http"
)

func Logger(inner http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		m := httpsnoop.CaptureMetrics(inner, w, r)
		logrus.Printf(
			"%s\t%s\t%s\t%d\t%d\t%s\n",
			r.RemoteAddr,
			r.Method,
			r.URL,
			m.Code,
			m.Written,
			m.Duration,
		)
	})
}

func Rest(inner http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if stringutils.IsEmpty(w.Header().Get("Content-Type")) {
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		}
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
