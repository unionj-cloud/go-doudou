package codegen

import (
	"bufio"
	"bytes"
	"github.com/iancoleman/strcase"
	"github.com/sirupsen/logrus"
	"github.com/unionj-cloud/go-doudou/cmd/internal/astutils"
	"github.com/unionj-cloud/go-doudou/toolkit/copier"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

var svcimportTmpl = `
	"context"
	"{{.ConfigPackage}}"
	"{{.VoPackage}}"
	"github.com/jmoiron/sqlx"
	"github.com/brianvoe/gofakeit/v6"
`

var appendPart = `{{- range $m := .Meta.Methods }}
	func (receiver *{{$.Meta.Name}}Impl) {{$m.Name}}({{- range $i, $p := $m.Params}}
    {{- if $i}},{{end}}
    {{- $p.Name}} {{$p.Type}}
    {{- end }}) ({{- range $i, $r := $m.Results}}
                     {{- if $i}},{{end}}
                     {{- $r.Name}} {{$r.Type}}
                     {{- end }}) {
    	var _result struct{
			{{- range $r := $m.Results }}
			{{- if ne $r.Type "error" }}
			{{ $r.Name | toCamel }} {{ $r.Type }}
			{{- end }}
			{{- end }}
		}
		_ = gofakeit.Struct(&_result)
		return {{range $i, $r := $m.Results }}{{- if $i}},{{end}}{{ if eq $r.Type "error" }}nil{{else}}_result.{{ $r.Name | toCamel }}{{end}}{{- end }}
    }
{{- end }}`

var svcimplTmpl = `package {{.SvcPackage}}

import ()

type {{.Meta.Name}}Impl struct {
	conf *config.Config
}

` + appendPart + `

func New{{.Meta.Name}}(conf *config.Config) {{.Meta.Name}} {
	return &{{.Meta.Name}}Impl{
		conf,
	}
}
`

// GenSvcImpl generates service implementation
func GenSvcImpl(dir string, ic astutils.InterfaceCollector) {
	var (
		err         error
		modfile     string
		modName     string
		svcimplfile string
		firstLine   string
		f           *os.File
		tpl         *template.Template
		buf         bytes.Buffer
		meta        astutils.InterfaceMeta
		tmpl        string
		importBuf   bytes.Buffer
	)
	svcimplfile = filepath.Join(dir, "svcimpl.go")
	err = copier.DeepCopy(ic.Interfaces[0], &meta)
	if err != nil {
		panic(err)
	}
	modfile = filepath.Join(dir, "go.mod")
	if f, err = os.Open(modfile); err != nil {
		panic(err)
	}
	reader := bufio.NewReader(f)
	firstLine, _ = reader.ReadString('\n')
	modName = strings.TrimSpace(strings.TrimPrefix(firstLine, "module"))
	if _, err = os.Stat(svcimplfile); os.IsNotExist(err) {
		if f, err = os.Create(svcimplfile); err != nil {
			panic(err)
		}
		defer f.Close()
		tmpl = svcimplTmpl
	} else {
		logrus.Warningln("New content will be append to file svcimpl.go")
		if f, err = os.OpenFile(svcimplfile, os.O_APPEND, os.ModePerm); err != nil {
			panic(err)
		}
		defer f.Close()
		tmpl = appendPart

		fset := token.NewFileSet()
		root, err := parser.ParseFile(fset, svcimplfile, nil, parser.ParseComments)
		if err != nil {
			panic(err)
		}
		sc := astutils.NewStructCollector(astutils.ExprString)
		ast.Walk(sc, root)
		if implementations, exists := sc.Methods[meta.Name+"Impl"]; exists {
			var notimplemented []astutils.MethodMeta
			for _, item := range meta.Methods {
				for _, implemented := range implementations {
					if item.Name == implemented.Name {
						goto L
					}
				}
				notimplemented = append(notimplemented, item)

			L:
			}

			meta.Methods = notimplemented
		}
	}

	funcMap := make(map[string]interface{})
	funcMap["toCamel"] = strcase.ToCamel
	if tpl, err = template.New("svcimpl.go.tmpl").Funcs(funcMap).Parse(tmpl); err != nil {
		panic(err)
	}
	if err = tpl.Execute(&buf, struct {
		ConfigPackage string
		VoPackage     string
		SvcPackage    string
		Meta          astutils.InterfaceMeta
	}{
		VoPackage:     modName + "/vo",
		ConfigPackage: modName + "/config",
		SvcPackage:    ic.Package.Name,
		Meta:          meta,
	}); err != nil {
		panic(err)
	}

	original, err := ioutil.ReadAll(f)
	if err != nil {
		panic(err)
	}

	original = append(original, buf.Bytes()...)
	if tpl, err = template.New("simportimpl.go.tmpl").Parse(svcimportTmpl); err != nil {
		panic(err)
	}
	if err = tpl.Execute(&importBuf, struct {
		ConfigPackage string
		VoPackage     string
	}{
		VoPackage:     modName + "/vo",
		ConfigPackage: modName + "/config",
	}); err != nil {
		panic(err)
	}
	original = astutils.AppendImportStatements(original, importBuf.Bytes())
	//fmt.Println(string(original))
	astutils.FixImport(original, svcimplfile)
}
