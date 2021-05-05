package config

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
)

type Dotenv struct {
	Fp   string // absolute config file path
	Conf Config
}

func (d *Dotenv) GetConf() Config {
	return d.Conf
}

func (d *Dotenv) GetLogLevel() logrus.Level {
	switch d.Conf.AppConf.LogLevel {
	case "panic":
		return logrus.PanicLevel
	case "fatal":
		return logrus.FatalLevel
	case "error":
		return logrus.ErrorLevel
	case "warn":
		return logrus.WarnLevel
	case "debug":
		return logrus.DebugLevel
	case "trace":
		return logrus.TraceLevel
	default:
		return logrus.InfoLevel
	}
}

func (d *Dotenv) Load() {
	err := godotenv.Load(d.Fp)
	if err != nil {
		logrus.Fatal("Error loading .env file", err)
	}
	var appconf AppConfig
	err = envconfig.Process("app", &appconf)
	if err != nil {
		logrus.Fatal("Error processing env", err)
	}
	var dbconf DbConfig
	err = envconfig.Process("db", &dbconf)
	if err != nil {
		logrus.Fatal("Error processing env", err)
	}
	var srvconf HttpConfig
	err = envconfig.Process("srv", &srvconf)
	if err != nil {
		logrus.Fatal("Error processing env", err)
	}
	var svcconf SvcConfig
	err = envconfig.Process("svc", &svcconf)
	if err != nil {
		logrus.Fatal("Error processing env", err)
	}
	d.Conf = Config{
		dbconf,
		srvconf,
		svcconf,
		appconf,
	}
}

func NewDotenv(fp string) Configurator {
	env := &Dotenv{
		Fp: fp,
	}
	env.Load()
	return env
}
