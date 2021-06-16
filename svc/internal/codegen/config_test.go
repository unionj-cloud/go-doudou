package codegen

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestGenConfig(t *testing.T) {
	dir := testDir + "config"
	InitSvc(dir)
	defer os.RemoveAll(dir)
	GenConfig(dir)
	expect := `package config

type Configurator interface {
	Load()
	Get() Config
}

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
`
	configfile := dir + "/config/config.go"
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
