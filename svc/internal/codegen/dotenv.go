package codegen

import (
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"text/template"
)

var dotenvTmpl = `package config

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
)

type Dotenv struct {
	Fp   string // absolute config file path
	Conf Config
}

func (d *Dotenv) Get() Config {
	return d.Conf
}

func (d *Dotenv) Load() {
	err := godotenv.Load(d.Fp)
	if err != nil {
		logrus.Fatal("Error loading .env file", err)
	}
	var dbconf DbConfig
	err = envconfig.Process("db", &dbconf)
	if err != nil {
		logrus.Fatal("Error processing env", err)
	}
	d.Conf = Config{
		dbconf,
	}
}

func NewDotenv(fp string) Configurator {
	env := &Dotenv{
		Fp: fp,
	}
	env.Load()
	return env
}
`

func GenDotenv(dir string) {
	var (
		err        error
		dotenvfile string
		f          *os.File
		tpl        *template.Template
		configDir  string
	)
	configDir = filepath.Join(dir, "config")
	if err = os.MkdirAll(configDir, 0644); err != nil {
		panic(err)
	}

	dotenvfile = filepath.Join(configDir, "dotenv.go")
	if _, err = os.Stat(dotenvfile); os.IsNotExist(err) {
		if f, err = os.Create(dotenvfile); err != nil {
			panic(err)
		}
		defer f.Close()

		if tpl, err = template.New("dotenv.go.tmpl").Parse(dotenvTmpl); err != nil {
			panic(err)
		}
		if err = tpl.Execute(f, nil); err != nil {
			panic(err)
		}
	} else {
		logrus.Warnf("file %s already exists", dotenvfile)
	}
}
