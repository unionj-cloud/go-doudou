package logger

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/unionj-cloud/go-doudou/pathutils"
	"github.com/unionj-cloud/go-doudou/stringutils"
	"github.com/unionj-cloud/go-doudou/svc/config"
	"io"
	"os"
	"path/filepath"
)

func Init() {
	var loglevel config.LogLevel
	(&loglevel).Decode(config.GddLogLevel.Load())

	formatter := new(logrus.TextFormatter)
	formatter.TimestampFormat = "2006-01-02 15:04:05"
	formatter.FullTimestamp = true

	logger := logrus.StandardLogger()
	logger.SetFormatter(formatter)
	logger.SetLevel(logrus.Level(loglevel))
}

func PersistLogToDisk() *os.File {
	var (
		logpath string
		err     error
		logFile *os.File
	)
	if logpath, err = pathutils.FixPath(config.GddLogPath.Load(), ""); err != nil {
		logrus.Panic(fmt.Sprintf("%+v\n", err))
	}
	if stringutils.IsNotEmpty(logpath) {
		if err = os.MkdirAll(logpath, os.ModePerm); err != nil {
			logrus.Panic(fmt.Sprintf("%+v\n", err))
		}
	}
	if logFile, err = os.OpenFile(filepath.Join(logpath, "app.log"), os.O_APPEND|os.O_CREATE|os.O_RDWR, os.ModePerm); err != nil {
		logrus.Panic(fmt.Sprintf("%+v\n", err))
	}
	logrus.SetOutput(io.MultiWriter(os.Stdout, logFile))
	return logFile
}
