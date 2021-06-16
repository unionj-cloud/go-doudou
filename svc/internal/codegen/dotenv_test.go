package codegen

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestGenDotenv(t *testing.T) {
	dir := testDir + "dotenv"
	InitSvc(dir)
	defer os.RemoveAll(dir)
	GenDotenv(dir)
	expect := `package config

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
)

type Dotenv struct {
	Fp   string // absolute config file path
	Conf Config
}

func (d *Dotenv) Get() Config {
	return d.Conf
}

func (d *Dotenv) Load() {
	err := godotenv.Load(d.Fp)
	if err != nil {
		logrus.Fatal("Error loading .env file", err)
	}
	var dbconf DbConfig
	err = envconfig.Process("db", &dbconf)
	if err != nil {
		logrus.Fatal("Error processing env", err)
	}
	d.Conf = Config{
		dbconf,
	}
}

func NewDotenv(fp string) Configurator {
	env := &Dotenv{
		Fp: fp,
	}
	env.Load()
	return env
}
`
	configfile := dir + "/config/dotenv.go"
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
