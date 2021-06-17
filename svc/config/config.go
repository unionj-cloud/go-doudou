package config

import (
	"github.com/sirupsen/logrus"
	"os"
)

type envVariable string

const (
	SvcName     envVariable = "SVC_NAME"
	SvcHostname envVariable = "SVC_HOSTNAME"
	SvcPort     envVariable = "SVC_PORT"
	SvcMemPort  envVariable = "SVC_MEM_PORT"
	SvcBaseUrl  envVariable = "SVC_BASE_URL"
	SvcSeed     envVariable = "SVC_SEED"
	// Accept 'mono' for monolith mode or 'micro' for microservice mode
	SvcMode envVariable = "SVC_MODE"
)

func (receiver envVariable) Load() string {
	return os.Getenv(string(receiver))
}

type Switch bool

func (s *Switch) Decode(value string) error {
	if value == "on" {
		*s = true
	}
	return nil
}

type LogLevel logrus.Level

func (ll *LogLevel) Decode(value string) error {
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
