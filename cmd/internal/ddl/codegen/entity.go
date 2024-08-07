package codegen

import (
	"os"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/unionj-cloud/go-doudou/v2/cmd/internal/templates"
	"github.com/unionj-cloud/go-doudou/v2/toolkit/astutils"
	"github.com/unionj-cloud/go-doudou/v2/toolkit/templateutils"
	"github.com/unionj-cloud/go-doudou/v2/version"
)

var entityTmpl = templates.EditableHeaderTmpl + `package entity

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

// GenEntityGo generates structs code in entity pkg from database tables
func GenEntityGo(dpath string, entity astutils.StructMeta) {
	var (
		err error
		f   *os.File
	)
	_ = os.MkdirAll(dpath, os.ModePerm)
	dfile := filepath.Join(dpath, strings.ToLower(entity.Name)+".go")
	if _, err = os.Stat(dfile); os.IsNotExist(err) {
		f, _ = os.Create(dfile)
		defer f.Close()
		var source string
		entity.Version = version.Release
		source, _ = templateutils.String("entity.go.tmpl", entityTmpl, entity)
		astutils.FixImport([]byte(source), dfile)
	} else {
		log.Warnf("file %s already exists", dfile)
	}
}
