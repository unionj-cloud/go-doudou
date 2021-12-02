package codegen

import (
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"text/template"
)

var httpMwTmpl = `package httpsrv`

// GenHttpMiddleware generates http middleware file
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
		tpl, _ = template.New("middleware.go.tmpl").Parse(httpMwTmpl)
		_ = tpl.Execute(f, nil)
	} else {
		logrus.Warnf("file %s already exists", mwfile)
	}
}
