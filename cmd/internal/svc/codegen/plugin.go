package codegen

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/unionj-cloud/go-doudou/v2/cmd/internal/templates"
	"github.com/unionj-cloud/toolkit/astutils"
	"github.com/unionj-cloud/go-doudou/v2/version"
)

type CodeGenConfig struct {
	ProjectType string
}

func genPlugin(dir string, ic astutils.InterfaceCollector, conf CodeGenConfig) {
	var (
		err        error
		pluginFile string
		f          *os.File
		tpl        *template.Template
		pluginDir  string
		buf        bytes.Buffer
	)
	pluginDir = filepath.Join(dir, "plugin")
	if err = MkdirAll(pluginDir, os.ModePerm); err != nil {
		panic(err)
	}
	pluginFile = filepath.Join(pluginDir, "plugin.go")
	if f, err = Create(pluginFile); err != nil {
		panic(err)
	}
	defer f.Close()

	funcMap := make(map[string]interface{})
	funcMap["toLower"] = strings.ToLower
	if tpl, err = template.New(templates.PluginTmpl).Funcs(funcMap).Parse(templates.PluginTmpl); err != nil {
		panic(err)
	}

	servicePkg := astutils.GetPkgPath(dir)
	cfgPkg := astutils.GetPkgPath(filepath.Join(dir, "config"))
	httpsrvPkg := astutils.GetPkgPath(filepath.Join(dir, "transport", "httpsrv"))
	transGrpcPkg := astutils.GetPkgPath(filepath.Join(dir, "transport", "grpc"))
	svcName := ic.Interfaces[0].Name
	alias := ic.Package.Name
	if err = tpl.Execute(&buf, struct {
		CodeGenConfig
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
		CodeGenConfig:        conf,
	}); err != nil {
		panic(err)
	}
	astutils.FixImport([]byte(strings.TrimSpace(buf.String())), pluginFile)
}
