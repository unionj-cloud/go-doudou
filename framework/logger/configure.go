package logger

import (
	"github.com/mattn/go-colorable"
	"github.com/sirupsen/logrus"
	"github.com/unionj-cloud/go-doudou/framework/internal/config"
	"github.com/unionj-cloud/go-doudou/toolkit/stringutils"
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
		tf.ForceColors = true
		formatter = tf
	default:
	}
	return formatter
}

// LogLevel alias for logrus.Level
type LogLevel logrus.Level

// Decode decodes value to LogLevel
func (ll *LogLevel) Decode(value string) error {
	//if stringutils.IsEmpty(value) {
	//	value = DefaultGddLogLevel
	//}
	switch value {
	case "panic":
		*ll = LogLevel(logrus.PanicLevel)
	case "fatal":
		*ll = LogLevel(logrus.FatalLevel)
	case "error":
		*ll = LogLevel(logrus.ErrorLevel)
	case "warn":
		*ll = LogLevel(logrus.WarnLevel)
	case "debug":
		*ll = LogLevel(logrus.DebugLevel)
	case "trace":
		*ll = LogLevel(logrus.TraceLevel)
	default:
		*ll = LogLevel(logrus.InfoLevel)
	}
	return nil
}

func Init(opts ...LoggerOption) {
	var loglevel LogLevel
	(&loglevel).Decode(config.GddLogLevel.Load())

	logger := logrus.StandardLogger()
	logger.SetFormatter(defaultFormatter())
	logger.SetLevel(logrus.Level(loglevel))
	logrus.SetOutput(colorable.NewColorableStdout())

	for _, opt := range opts {
		opt(logger)
	}
}
