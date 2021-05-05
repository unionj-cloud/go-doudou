package main

import (
	"context"
	service "example/user-svc"
	"example/user-svc/config"
	"example/user-svc/db"
	"example/user-svc/transport/httpsrv"
	"fmt"
	"github.com/common-nighthawk/go-figure"
	"github.com/sirupsen/logrus"
	"github.com/unionj-cloud/go-doudou/pathutils"
	"github.com/unionj-cloud/go-doudou/stringutils"
	"io"
	"os"
	"os/signal"
	"path/filepath"
	"time"
)

func configureLogger(logger *logrus.Logger, logptr *string, level logrus.Level) *os.File {
	formatter := new(logrus.TextFormatter)
	formatter.TimestampFormat = "2006-01-02 15:04:05"
	formatter.FullTimestamp = true
	logger.SetFormatter(formatter)
	logger.SetLevel(level)

	if logptr != nil {
		var (
			err error
			f   *os.File
		)
		logpath := *logptr
		logpath, err = pathutils.FixPath(logpath, "")
		if err != nil {
			logger.Errorln(fmt.Sprintf("%+v\n", err))
		}
		if stringutils.IsNotEmpty(logpath) {
			err = os.MkdirAll(logpath, os.ModePerm)
			if err != nil {
				logger.Errorln(err)
				return nil
			}
		}
		f, err = os.OpenFile(filepath.Join(logpath, "app.log"), os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
		if err != nil {
			logger.Errorf("error opening file: %v\n", err)
			return nil
		}
		mw := io.MultiWriter(os.Stdout, f)
		logger.SetOutput(mw)
		return f
	}

	return nil
}

func main() {
	start := time.Now()
	env := config.NewDotenv(pathutils.Abs("../.env"))
	conf := env.Get()

	logFile := configureLogger(logrus.StandardLogger(), conf.AppConf.LogPath, logrus.Level(conf.AppConf.LogLevel))
	defer func() {
		if logFile != nil {
			logFile.Close()
		}
	}()

	if conf.AppConf.Banner {
		figure.NewColorFigure(conf.AppConf.BannerText, "doom", "green", true).Print()
	}

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

	svc := service.NewUserService(conf.SvcConf, db)
	handler := httpsrv.NewUserHandler(svc)

	srv := httpsrv.Run(conf.HttpConf, handler)

	logrus.Infof("Started in %s\n", time.Since(start))

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), conf.AppConf.GraceTimeout)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	srv.Shutdown(ctx)

	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	logrus.Infoln("shutting down")
}
