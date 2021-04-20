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

var svcimplTmpl = `package {{.Ic.Package.Name}}

{{- range $interface := .Ic.Interfaces }}
type {{$interface.Name}}Impl struct{}

{{- range $m := $interface.Methods }}
	func (receiver {{$interface.Name}}Impl) {{$m.Name}}({{- range $i, $p := $m.Params}}
    {{- if $i}},{{end}}
    {{- $p.Name}} {{$p.Type}}
    {{- end }}) ({{- range $i, $r := $m.Results}}
                     {{- if $i}},{{end}}
                     {{- $r.Name}} {{$r.Type}}
                     {{- end }}) {
    	panic("implement me")
    }
{{- end }}

func New{{$interface.Name}}() {{$interface.Name}} {
	return {{$interface.Name}}Impl{}
}
{{- end }}
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
			VoPackage string
			Ic        astutils.InterfaceCollector
		}{
			VoPackage: modName + "/vo",
			Ic:        ic,
		}); err != nil {
			panic(err)
		}

		source = strings.TrimSpace(sqlBuf.String())
		astutils.FixImport([]byte(source), svcimplfile)
	} else {
		logrus.Warnf("file %s already exists.", svcimplfile)
	}
}
