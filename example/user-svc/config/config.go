package config

type Configurator interface {
	NewConf() Config
}

type Config struct {
	DbConf   DbConfig
	HttpConf HttpConfig
	SvcConf  SvcConfig
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

func NewConf(configurator Configurator) Config {
	return configurator.NewConf()
}
