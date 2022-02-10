package codegen

import (
	"bufio"
	"github.com/sirupsen/logrus"
	"github.com/unionj-cloud/go-doudou/cmd/internal/astutils"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

var mainTmpl = `package main

import (
	ddhttp "github.com/unionj-cloud/go-doudou/framework/http"
	{{.ServiceAlias}} "{{.ServicePackage}}"
    "{{.ConfigPackage}}"
	"{{.HttpPackage}}"
)

func main() {
	conf := config.LoadFromEnv()
    svc := {{.ServiceAlias}}.New{{.SvcName}}(conf)
	handler := httpsrv.New{{.SvcName}}Handler(svc)
	srv := ddhttp.NewDefaultHttpSrv()
	srv.AddRoute(httpsrv.Routes(handler)...)
	srv.Run()
}
`

// GenMain generates main function
func GenMain(dir string, ic astutils.InterfaceCollector) {
	var (
		err       error
		modfile   string
		modName   string
		mainfile  string
		firstLine string
		f         *os.File
		tpl       *template.Template
		cmdDir    string
		svcName   string
		alias     string
	)
	cmdDir = filepath.Join(dir, "cmd")
	if err = os.MkdirAll(cmdDir, os.ModePerm); err != nil {
		panic(err)
	}

	svcName = ic.Interfaces[0].Name
	alias = ic.Package.Name
	mainfile = filepath.Join(cmdDir, "main.go")
	if _, err = os.Stat(mainfile); os.IsNotExist(err) {
		modfile = filepath.Join(dir, "go.mod")
		if f, err = os.Open(modfile); err != nil {
			panic(err)
		}
		reader := bufio.NewReader(f)
		if firstLine, err = reader.ReadString('\n'); err != nil {
			panic(err)
		}
		modName = strings.TrimSpace(strings.TrimPrefix(firstLine, "module"))

		if f, err = os.Create(mainfile); err != nil {
			panic(err)
		}
		defer f.Close()

		if tpl, err = template.New("main.go.tmpl").Parse(mainTmpl); err != nil {
			panic(err)
		}
		if err = tpl.Execute(f, struct {
			ServicePackage string
			ConfigPackage  string
			DbPackage      string
			HttpPackage    string
			SvcName        string
			ServiceAlias   string
		}{
			ServicePackage: modName,
			ConfigPackage:  modName + "/config",
			DbPackage:      modName + "/db",
			HttpPackage:    modName + "/transport/httpsrv",
			SvcName:        svcName,
			ServiceAlias:   alias,
		}); err != nil {
			panic(err)
		}
	} else {
		logrus.Warnf("file %s already exists", mainfile)
	}
}
