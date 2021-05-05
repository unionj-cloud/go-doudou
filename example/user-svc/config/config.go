package config

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
	Driver  string `default:"mysql"`
	Host    string `default:"localhost"`
	Port    string `default:"3306"`
	User    string
	Passwd  string
	Schema  string
	Charset string `default:"utf8mb4"`
}

type HttpConfig struct {
	Host         string        `default:"0.0.0.0"`
	Port         string        `default:"8080"`
	WriteTimeout time.Duration `default:"15s"`
	ReadTimeout  time.Duration `default:"15s"`
	IdleTimeout  time.Duration `default:"60s"`
}

type SvcConfig struct {
}

type AppConfig struct {
	Banner       Switch   `default:"on"`
	BannerText   string   `default:"Go-doudou"`
	LogLevel     LogLevel `default:"info"`
	LogPath      *string
	GraceTimeout time.Duration `default:"15s"`
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
