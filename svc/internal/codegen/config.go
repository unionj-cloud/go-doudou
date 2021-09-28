package codegen

import (
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"text/template"
)

var configTmpl = `package config

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
)

type Config struct {
	DbConf   DbConfig
}

type DbConfig struct {
	Driver  string ` + "`" + `default:"mysql"` + "`" + `
	Host    string ` + "`" + `default:"localhost"` + "`" + `
	Port    string ` + "`" + `default:"3306"` + "`" + `
	User    string
	Passwd  string
	Schema  string
	Charset string ` + "`" + `default:"utf8mb4"` + "`" + `
}

func LoadFromEnv() *Config {
	var dbconf DbConfig
	err := envconfig.Process("db", &dbconf)
	if err != nil {
		logrus.Panicln("Error processing env", err)
	}
	return &Config{
		dbconf,
	}
}
`

//GenConfig generates config file
func GenConfig(dir string) {
	var (
		err        error
		configfile string
		f          *os.File
		tpl        *template.Template
		configDir  string
	)
	configDir = filepath.Join(dir, "config")
	if err = os.MkdirAll(configDir, os.ModePerm); err != nil {
		panic(err)
	}

	configfile = filepath.Join(configDir, "config.go")
	if _, err = os.Stat(configfile); os.IsNotExist(err) {
		if f, err = os.Create(configfile); err != nil {
			panic(err)
		}
		defer f.Close()

		if tpl, err = template.New("config.go.tmpl").Parse(configTmpl); err != nil {
			panic(err)
		}
		if err = tpl.Execute(f, nil); err != nil {
			panic(err)
		}
	} else {
		logrus.Warnf("file %s already exists", configfile)
	}
}
