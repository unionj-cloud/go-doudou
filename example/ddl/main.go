package main

import (
	"context"
	"example/ddl/dao"
	"example/ddl/domain"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/iancoleman/strcase"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
	. "github.com/unionj-cloud/go-doudou/kit/ddl/query"
	"github.com/unionj-cloud/go-doudou/kit/pathutils"
	"log"
)

type DbConfig struct {
	Host    string
	Port    string
	User    string
	Passwd  string
	Schema  string
	Charset string
}

func init() {
	customFormatter := new(logrus.TextFormatter)
	customFormatter.TimestampFormat = "2006-01-02 15:04:05"
	customFormatter.FullTimestamp = true
	logrus.SetFormatter(customFormatter)
}

func main() {
	err := godotenv.Load(pathutils.Abs(".env"))
	if err != nil {
		logrus.Fatal("Error loading .env file", err)
	}
	var dbConfig DbConfig
	err = envconfig.Process("db", &dbConfig)
	if err != nil {
		logrus.Fatal("Error processing env", err)
	}

	conn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s",
		dbConfig.User,
		dbConfig.Passwd,
		dbConfig.Host,
		dbConfig.Port,
		dbConfig.Schema,
		dbConfig.Charset)
	conn += `&loc=Asia%2FShanghai&parseTime=True`
	var db *sqlx.DB
	db, err = sqlx.Connect("mysql", conn)
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()
	db.MapperFunc(strcase.ToSnake)

	u := dao.NewUserDao(db)

	if _, err = u.UpsertUser(context.Background(), &domain.User{
		Name:      "Biden",
		Phone:     "13893997878",
		Age:       70,
		School:    "Harvard Univ.",
		No:        46,
		IsStudent: true,
	}); err != nil {
		logrus.Panicln(err)
	}

	if _, err = u.UpsertUser(context.Background(), &domain.User{
		Name:      "Trump",
		Phone:     "13893997979",
		Age:       12,
		School:    "Harvard Univ.",
		No:        45,
		IsStudent: true,
	}); err != nil {
		logrus.Panicln(err)
	}

	got, err := u.PageUsers(context.TODO(), C().Col("age").Gt(Literal("27")), Page{
		Orders: []Order{
			{
				Col:  "age",
				Sort: "desc",
			},
		},
		Offset: 0,
		Size:   1,
	})
	if err != nil {
		panic(err)
	}
	logrus.Infof("%+v\n", got)
}
