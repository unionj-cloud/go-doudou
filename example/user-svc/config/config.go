package config

import "github.com/sirupsen/logrus"

type Configurator interface {
	Load()
	GetConf() Config
	GetLogLevel() logrus.Level
}

type Config struct {
	DbConf   DbConfig
	HttpConf HttpConfig
	SvcConf  SvcConfig
	AppConf  AppConfig
}

type DbConfig struct {
	Driver  string
	Host    string
	Port    string
	User    string
	Passwd  string
	Schema  string
	Charset string
}

type HttpConfig struct {
	Host         string
	Port         string
	WriteTimeout string
	ReadTimeout  string
	IdleTimeout  string
}

type SvcConfig struct {
}

type AppConfig struct {
	Logo         string
	LogLevel     string
	GraceTimeout string
}
