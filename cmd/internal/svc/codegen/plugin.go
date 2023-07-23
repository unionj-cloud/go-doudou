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

func genPlugin(dir string, ic astutils.InterfaceCollector) {
	var (
		err      error
		mainfile string
		f        *os.File
		tpl      *template.Template
		cmdDir   string
	)
	cmdDir = filepath.Join(dir, "cmd")
	if err = MkdirAll(cmdDir, os.ModePerm); err != nil {
		panic(err)
	}
	mainfile = filepath.Join(cmdDir, "main.go")
	if _, err = Stat(mainfile); os.IsNotExist(err) {
		if f, err = Create(mainfile); err != nil {
			panic(err)
		}
		defer f.Close()

		if tpl, err = template.New(templates.MainModuleTmpl).Parse(templates.MainModuleTmpl); err != nil {
			panic(err)
		}

		servicePkg := astutils.GetPkgPath(dir)
		cfgPkg := astutils.GetPkgPath(filepath.Join(dir, "config"))
		httpsrvPkg := astutils.GetPkgPath(filepath.Join(dir, "transport", "httpsrv"))
		transGrpcPkg := astutils.GetPkgPath(filepath.Join(dir, "transport", "grpc"))
		svcName := ic.Interfaces[0].Name
		alias := ic.Package.Name
		if err = tpl.Execute(f, struct {
			ServicePackage       string
			ConfigPackage        string
			TransportGrpcPackage string
			TransportHttpPackage string
			ServiceAlias         string
			SvcName              string
			Version              string
		}{
			ServicePackage:       servicePkg,
			ConfigPackage:        cfgPkg,
			TransportGrpcPackage: transGrpcPkg,
			TransportHttpPackage: httpsrvPkg,
			ServiceAlias:         alias,
			SvcName:              svcName,
			Version:              version.Release,
		}); err != nil {
			panic(err)
		}
	} else {
		logrus.Warnf("file %s already exists", mainfile)
	}
}
