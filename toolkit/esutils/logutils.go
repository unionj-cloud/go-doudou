package esutils

import (
	"github.com/sirupsen/logrus"
)

func newLogger(level logrus.Level) *logrus.Logger {
	logger := logrus.New()
	formatter := new(logrus.TextFormatter)
	formatter.TimestampFormat = "2006-01-02 15:04:05"
	formatter.FullTimestamp = true
	logger.SetFormatter(formatter)
	logger.SetLevel(level)
	return logger
}
