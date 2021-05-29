package codegen

import (
	"bufio"
	"bytes"
	"github.com/sirupsen/logrus"
	"github.com/unionj-cloud/go-doudou/astutils"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

var svcimplTmpl = `package {{.SvcPackage}}

import (
	"context"
	"{{.ConfigPackage}}"
	"{{.VoPackage}}"
	"github.com/jmoiron/sqlx"
)

type {{.Meta.Name}}Impl struct {
	conf config.Config
}

{{- range $m := .Meta.Methods }}
	func (receiver *{{$.Meta.Name}}Impl) {{$m.Name}}({{- range $i, $p := $m.Params}}
    {{- if $i}},{{end}}
    {{- $p.Name}} {{$p.Type}}
    {{- end }}) ({{- range $i, $r := $m.Results}}
                     {{- if $i}},{{end}}
                     {{- $r.Name}} {{$r.Type}}
                     {{- end }}) {
    	panic("implement me")
    }
{{- end }}

func New{{.Meta.Name}}(conf config.Config, db *sqlx.DB) {{.Meta.Name}} {
	return &{{.Meta.Name}}Impl{
		conf,
	}
}
`

func GenSvcImpl(dir string, ic astutils.InterfaceCollector) {
	var (
		err         error
		modfile     string
		modName     string
		svcimplfile string
		firstLine   string
		f           *os.File
		tpl         *template.Template
		source      string
		sqlBuf      bytes.Buffer
	)
	svcimplfile = filepath.Join(dir, "svcimpl.go")
	if _, err = os.Stat(svcimplfile); os.IsNotExist(err) {
		modfile = filepath.Join(dir, "go.mod")
		if f, err = os.Open(modfile); err != nil {
			panic(err)
		}
		reader := bufio.NewReader(f)
		if firstLine, err = reader.ReadString('\n'); err != nil {
			panic(err)
		}
		modName = strings.TrimSpace(strings.TrimPrefix(firstLine, "module"))

		if f, err = os.Create(svcimplfile); err != nil {
			panic(err)
		}
		defer f.Close()

		if tpl, err = template.New("svcimpl.go.tmpl").Parse(svcimplTmpl); err != nil {
			panic(err)
		}
		if err = tpl.Execute(&sqlBuf, struct {
			ConfigPackage string
			VoPackage     string
			SvcPackage    string
			Meta          astutils.InterfaceMeta
		}{
			VoPackage:     modName + "/vo",
			ConfigPackage: modName + "/config",
			SvcPackage:    ic.Package.Name,
			Meta:          ic.Interfaces[0],
		}); err != nil {
			panic(err)
		}

		source = strings.TrimSpace(sqlBuf.String())
		astutils.FixImport([]byte(source), svcimplfile)
	} else {
		logrus.Warnf("file %s already exists.", svcimplfile)
	}
}
