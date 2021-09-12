package codegen

import (
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

// GenDomainGo generates structs code in domain pkg from database tables
func GenDomainGo(dpath string, domain astutils.StructMeta) error {
	var (
		err error
		f   *os.File
	)
	_ = os.MkdirAll(dpath, os.ModePerm)
	dfile := filepath.Join(dpath, strings.ToLower(domain.Name)+".go")
	if _, err = os.Stat(dfile); os.IsNotExist(err) {
		f, _ = os.Create(dfile)
		defer f.Close()
		var source string
		source, _ = templateutils.String("domain.go.tmpl", domaintmpl, domain)
		astutils.FixImport([]byte(source), dfile)
	} else {
		log.Warnf("file %s already exists", dfile)
	}
	return nil
}
