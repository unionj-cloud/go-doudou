package codegen

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestGenConfig(t *testing.T) {
	dir := testDir + "config"
	InitSvc(dir)
	defer os.RemoveAll(dir)
	GenConfig(dir)
	expect := `package config

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
	configfile := filepath.Join(dir, "config", "config.go")
	f, err := os.Open(configfile)
	if err != nil {
		t.Fatal(err)
	}
	content, err := ioutil.ReadAll(f)
	if err != nil {
		t.Fatal(err)
	}
	if string(content) != expect {
		t.Errorf("want %s, got %s\n", expect, string(content))
	}
}

func TestGenConfig1(t *testing.T) {
	GenConfig(filepath.Join(testDir))
}
