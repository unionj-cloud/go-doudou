package logutils

import (
	"github.com/sirupsen/logrus"
	"github.com/unionj-cloud/go-doudou/svc/config"
)

// NewLogger creates a logrus.Logger instance
func NewLogger() *logrus.Logger {
	logger := logrus.New()
	formatter := new(logrus.TextFormatter)
	formatter.TimestampFormat = "2006-01-02 15:04:05"
	formatter.FullTimestamp = true
	logger.SetFormatter(formatter)
	var loglevel config.LogLevel
	(&loglevel).Decode(config.GddLogLevel.Load())
	logger.SetLevel(logrus.Level(loglevel))
	return logger
}
