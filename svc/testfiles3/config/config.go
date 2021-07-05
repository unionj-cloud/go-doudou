package config

type Configurator interface {
	Load()
	Get() Config
}

type Config struct {
	DbConf   DbConfig
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
