package codegen

import (
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/unionj-cloud/go-doudou/astutils"
	"github.com/unionj-cloud/go-doudou/templateutils"
	"os"
	"path/filepath"
	"strings"
)

var domaintmpl = `package domain

import (
	"time"

	"github.com/shopspring/decimal"
)

//dd:table
type {{.Name}} struct {
{{- range $f := .Fields }}
	{{$f.Name}} {{$f.Type}} ` + "`" + `{{$f.Tag}}` + "`" + `
{{- end }}
}`

func GenDomainGo(dpath string, domain astutils.StructMeta) error {
	var (
		err error
		f   *os.File
	)

	if err = os.MkdirAll(dpath, os.ModePerm); err != nil {
		return errors.Wrap(err, "error")
	}

	dfile := filepath.Join(dpath, strings.ToLower(domain.Name)+".go")
	if _, err = os.Stat(dfile); os.IsNotExist(err) {
		if f, err = os.Create(dfile); err != nil {
			return errors.Wrap(err, "error")
		}
		defer f.Close()

		var source string
		if source, err = templateutils.String("domain.go.tmpl", domaintmpl, domain); err != nil {
			return errors.Wrap(err, "error")
		}

		astutils.FixImport([]byte(source), dfile)
	} else {
		log.Warnf("file %s already exists", dfile)
	}
	return nil
}
