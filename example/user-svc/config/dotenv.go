package config

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
)

type Dotenv struct {
	AbsCfg string // absolute config file path
}

func (d Dotenv) NewConf() Config {
	err := godotenv.Load(d.AbsCfg)
	if err != nil {
		logrus.Fatal("Error loading .env file", err)
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
	return Config{
		dbconf,
		srvconf,
	}
}
