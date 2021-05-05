package codegen

import (
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"text/template"
)

var configTmpl = `package config

import (
	"github.com/sirupsen/logrus"
	"time"
)

type Configurator interface {
	Load()
	Get() Config
}

type Config struct {
	DbConf   DbConfig
	HttpConf HttpConfig
	SvcConf  SvcConfig
	AppConf  AppConfig
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

type HttpConfig struct {
	Host         string        ` + "`" + `default:"0.0.0.0"` + "`" + `
	Port         string        ` + "`" + `default:"8080"` + "`" + `
	WriteTimeout time.Duration ` + "`" + `default:"15s"` + "`" + `
	ReadTimeout  time.Duration ` + "`" + `default:"15s"` + "`" + `
	IdleTimeout  time.Duration ` + "`" + `default:"60s"` + "`" + `
}

type SvcConfig struct {
}

type AppConfig struct {
	Banner       Switch   ` + "`" + `default:"on"` + "`" + `
	BannerText   string   ` + "`" + `default:"Go-doudou"` + "`" + `
	LogLevel     LogLevel ` + "`" + `default:"info"` + "`" + `
	LogPath      *string
	GraceTimeout time.Duration ` + "`" + `default:"15s"` + "`" + `
}

type Switch bool

func (s *Switch) Decode(value string) error {
	if value == "on" {
		*s = true
	}
	return nil
}

type LogLevel logrus.Level

func (ll *LogLevel) Decode(value string) error {
	switch value {
	case "panic":
		*ll = LogLevel(logrus.PanicLevel)
	case "fatal":
		*ll = LogLevel(logrus.FatalLevel)
	case "error":
		*ll = LogLevel(logrus.ErrorLevel)
	case "warn":
		*ll = LogLevel(logrus.WarnLevel)
	case "debug":
		*ll = LogLevel(logrus.DebugLevel)
	case "trace":
		*ll = LogLevel(logrus.TraceLevel)
	default:
		*ll = LogLevel(logrus.InfoLevel)
	}
	return nil
}
`

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
