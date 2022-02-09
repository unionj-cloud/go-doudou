package logger

import (
	"github.com/sirupsen/logrus"
	"github.com/unionj-cloud/go-doudou/stringutils"
	"github.com/unionj-cloud/go-doudou/svc/config"
	"io"
)

type LoggerOption func(*logrus.Logger)

func WithWritter(writer io.Writer) LoggerOption {
	return func(log *logrus.Logger) {
		log.SetOutput(writer)
	}
}

func WithFormatter(formatter logrus.Formatter) LoggerOption {
	return func(log *logrus.Logger) {
		log.SetFormatter(formatter)
	}
}

func WithReportCaller(reportCaller bool) LoggerOption {
	return func(log *logrus.Logger) {
		log.SetReportCaller(reportCaller)
	}
}

func defaultFormatter() logrus.Formatter {
	format := config.DefaultGddLogFormat
	if stringutils.IsNotEmpty(config.GddLogFormat.Load()) {
		format = config.GddLogFormat.Load()
	}
	var formatter logrus.Formatter
	switch format {
	case "json":
		jf := new(logrus.JSONFormatter)
		jf.TimestampFormat = "2006-01-02 15:04:05"
		jf.DisableHTMLEscape = true
		formatter = jf
	case "text":
		tf := new(logrus.TextFormatter)
		tf.TimestampFormat = "2006-01-02 15:04:05"
		tf.FullTimestamp = true
		formatter = tf
	default:
	}
	return formatter
}

func Init(opts ...LoggerOption) {
	var loglevel config.LogLevel
	(&loglevel).Decode(config.GddLogLevel.Load())

	logger := logrus.StandardLogger()
	logger.SetFormatter(defaultFormatter())
	logger.SetLevel(logrus.Level(loglevel))

	for _, opt := range opts {
		opt(logger)
	}
}
