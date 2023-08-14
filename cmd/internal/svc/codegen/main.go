package codegen

import (
	"github.com/sirupsen/logrus"
	"github.com/unionj-cloud/go-doudou/v2/cmd/internal/templates"
	"github.com/unionj-cloud/go-doudou/v2/toolkit/astutils"
	"github.com/unionj-cloud/go-doudou/v2/version"
	"os"
	"path/filepath"
	"text/template"
)

var mainTmpl = templates.EditableHeaderTmpl + `package main

import (
	"github.com/unionj-cloud/go-doudou/v2/framework/rest"
	{{.ServiceAlias}} "{{.ServicePackage}}"
    "{{.ConfigPackage}}"
	"{{.HttpPackage}}"
)

func main() {
	conf := config.LoadFromEnv()
    svc := {{.ServiceAlias}}.New{{.SvcName}}(conf)
	handler := httpsrv.New{{.SvcName}}Handler(svc)
	srv := rest.NewRestServer()
	srv.AddRoute(httpsrv.Routes(handler)...)
	srv.Run()
}
`

// GenMain generates main function
func GenMain(dir string, ic astutils.InterfaceCollector) {
	var (
		err      error
		mainfile string
		f        *os.File
		tpl      *template.Template
		cmdDir   string
		svcName  string
		alias    string
	)
	cmdDir = filepath.Join(dir, "cmd")
	if err = MkdirAll(cmdDir, os.ModePerm); err != nil {
		panic(err)
	}

	svcName = ic.Interfaces[0].Name
	alias = ic.Package.Name
	mainfile = filepath.Join(cmdDir, "main.go")
	if _, err = Stat(mainfile); os.IsNotExist(err) {
		servicePkg := astutils.GetPkgPath(dir)
		cfgPkg := astutils.GetPkgPath(filepath.Join(dir, "config"))
		httpsrvPkg := astutils.GetPkgPath(filepath.Join(dir, "transport", "httpsrv"))

		if f, err = Create(mainfile); err != nil {
			panic(err)
		}
		defer f.Close()

		if tpl, err = template.New("main.go.tmpl").Parse(mainTmpl); err != nil {
			panic(err)
		}
		if err = tpl.Execute(f, struct {
			ServicePackage string
			ConfigPackage  string
			HttpPackage    string
			SvcName        string
			ServiceAlias   string
			Version        string
		}{
			ServicePackage: servicePkg,
			ConfigPackage:  cfgPkg,
			HttpPackage:    httpsrvPkg,
			SvcName:        svcName,
			ServiceAlias:   alias,
			Version:        version.Release,
		}); err != nil {
			panic(err)
		}
	} else {
		logrus.Warnf("file %s already exists", mainfile)
	}
}
