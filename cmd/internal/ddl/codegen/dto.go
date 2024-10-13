package codegen

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"text/template"

	"github.com/sirupsen/logrus"
	"github.com/unionj-cloud/go-doudou/v2/cmd/internal/templates"
	"github.com/unionj-cloud/toolkit/astutils"
	"github.com/unionj-cloud/go-doudou/v2/version"
)

var appendDtoTmpl = `
{{- range $m := .Entities}}

type {{$m.Name}} struct {
{{- range $f := $m.Fields }}
	{{$f.Name}} {{$f.Type}}
{{- end }}
}
{{- end }}
`

var initDtoTmpl = templates.EditableHeaderTmpl + `package dto

import (
	"time"

	"github.com/shopspring/decimal"
)

//go:generate go-doudou name --file $GOFILE

` + appendDtoTmpl

// GenDto generates structs code in dto pkg from database tables
func GenDto(dir string, entities []astutils.StructMeta) {
	var (
		err  error
		f    *os.File
		fi   os.FileInfo
		tmpl string
		tpl  *template.Template
		buf  bytes.Buffer
	)
	dir = filepath.Join(filepath.Dir(dir), "dto")
	_ = os.MkdirAll(dir, os.ModePerm)
	dtoFile := filepath.Join(dir, "dto_gen.go")
	fi, err = os.Stat(dtoFile)
	if err != nil && !os.IsNotExist(err) {
		panic(err)
	}
	if fi != nil {
		logrus.Warningln("New content will be append to dto_gen.go file")
		if f, err = os.OpenFile(dtoFile, os.O_APPEND, os.ModePerm); err != nil {
			panic(err)
		}
		defer f.Close()
		tmpl = appendDtoTmpl
	} else {
		if f, err = os.Create(dtoFile); err != nil {
			panic(err)
		}
		defer f.Close()
		tmpl = initDtoTmpl
	}
	if tpl, err = template.New("dto_gen.go.tmpl").Parse(tmpl); err != nil {
		panic(err)
	}
	if err = tpl.Execute(&buf, struct {
		Entities []astutils.StructMeta
		Version  string
	}{
		Entities: entities,
		Version:  version.Release,
	}); err != nil {
		panic(err)
	}
	original, err := ioutil.ReadAll(f)
	if err != nil {
		panic(err)
	}
	original = append(original, buf.Bytes()...)
	astutils.FixImport(original, dtoFile)
}
