package config

// DbConfig store database connection parameters
type DbConfig struct {
	Host    string
	Port    string
	User    string
	Passwd  string
	Schema  string
	Charset string
}
