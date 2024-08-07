package codegen

import (
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/sirupsen/logrus"
	"github.com/unionj-cloud/go-doudou/v2/cmd/internal/templates"
	"github.com/unionj-cloud/go-doudou/v2/toolkit/astutils"
	"github.com/unionj-cloud/go-doudou/v2/version"
)

var configTmpl = templates.EditableHeaderTmpl + `package config

import (
	"github.com/unionj-cloud/go-doudou/v2/toolkit/envconfig"
	"github.com/unionj-cloud/go-doudou/v2/framework/config"
	"github.com/unionj-cloud/go-doudou/v2/toolkit/zlogger"
)

var G_Config *Config

type Config struct {
	Biz struct {
		ApiSecret string ` + "`" + `split_words:"true"` + "`" + `
	}
	config.Config
}

func init() {
	var conf Config
	err := envconfig.Process("{{.ServiceName}}", &conf)
	if err != nil {
		zlogger.Panic().Msgf("Error processing environment variables: %v", err)
	}
	G_Config = &conf
}

func LoadFromEnv() *Config {
	return G_Config
}
`

// GenConfig generates config file
func GenConfig(dir string, ic astutils.InterfaceCollector) {
	var (
		err         error
		configfile  string
		f           *os.File
		tpl         *template.Template
		configDir   string
		serviceName string
	)
	configDir = filepath.Join(dir, "config")
	if err = os.MkdirAll(configDir, os.ModePerm); err != nil {
		panic(err)
	}
	serviceName = strings.ToLower(ic.Interfaces[0].Name)
	configfile = filepath.Join(configDir, "config.go")
	if _, err = os.Stat(configfile); os.IsNotExist(err) {
		if f, err = os.Create(configfile); err != nil {
			panic(err)
		}
		defer f.Close()
		tpl, _ = template.New("config.go.tmpl").Parse(configTmpl)
		_ = tpl.Execute(f, struct {
			Version     string
			ServiceName string
		}{
			Version:     version.Release,
			ServiceName: serviceName,
		})
	} else {
		logrus.Warnf("file %s already exists", configfile)
	}
}
