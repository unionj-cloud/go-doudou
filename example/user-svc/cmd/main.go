package main

import (
	"context"
	service "example/user-svc"
	"example/user-svc/config"
	"example/user-svc/dao"
	"example/user-svc/db"
	"example/user-svc/transport/httpsrv"
	"flag"
	"github.com/common-nighthawk/go-figure"
	"github.com/sirupsen/logrus"
	"github.com/unionj-cloud/go-doudou/pathutils"
	"os"
	"os/signal"
	"time"
)

func main() {
	env := config.NewDotenv(pathutils.Abs("../.env"))
	conf := env.GetConf()
	conf.AppConf.GraceTimeout


	if conf.AppConf.Logo != "off" {
		figure.NewColorFigure("Go-doudou", "doom", "green", true).Print()
	}

	customFormatter := new(logrus.TextFormatter)
	customFormatter.TimestampFormat = "2006-01-02 15:04:05"
	customFormatter.FullTimestamp = true
	logrus.SetFormatter(customFormatter)
	logrus.SetLevel(logrus.WarnLevel)



	db, err := db.NewDb(conf.DbConf)
	if err != nil {
		panic(err)
	}
	defer func() {
		if db == nil {
			return
		}
		if err := db.Close(); err == nil {
			logrus.Infoln("Database connection is closed")
		} else {
			logrus.Warnln("Failed to close database connection")
		}
	}()

	userdao := dao.NewUserDao(db)
	svc := service.NewUserService(conf.SvcConf, userdao)
	handler := httpsrv.NewUserHandler(svc)

	srv := httpsrv.Run(conf.HttpConf, handler)

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	srv.Shutdown(ctx)

	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	logrus.Infoln("shutting down")
}
