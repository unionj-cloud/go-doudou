package db

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/iancoleman/strcase"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
	"github.com/unionj-cloud/go-doudou/kit/pathutils"
	"log"
)

type dbConfig struct {
	Host    string
	Port    string
	User    string
	Passwd  string
	Schema  string
	Charset string
}

var db *sqlx.DB

func init() {
	err := godotenv.Load(pathutils.Abs("../.env"))
	if err != nil {
		logrus.Fatal("Error loading .env file", err)
	}
	var conf dbConfig
	err = envconfig.Process("db", &conf)
	if err != nil {
		logrus.Fatal("Error processing env", err)
	}

	conn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s",
		conf.User,
		conf.Passwd,
		conf.Host,
		conf.Port,
		conf.Schema,
		conf.Charset)
	conn += `&loc=Asia%2FShanghai&parseTime=True`

	db, err = sqlx.Connect("mysql", conn)
	if err != nil {
		log.Fatalln(err)
	}
	db.MapperFunc(strcase.ToSnake)
}

func Db() *sqlx.DB {
	return db
}
