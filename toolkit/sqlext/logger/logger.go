package logger

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/unionj-cloud/go-doudou/toolkit/cast"
	"github.com/unionj-cloud/go-doudou/toolkit/reflectutils"
	"os"
	"regexp"
	"strings"
)

type ISqlLogger interface {
	Log(query string, args ...interface{})
	Enable() bool
}

type SqlLogger struct {
	logger logrus.StdLogger
}

func (receiver *SqlLogger) Enable() bool {
	return cast.ToBoolOrDefault(os.Getenv("GDD_SQL_LOG_ENABLE"), false)
}

func (receiver *SqlLogger) SetLogger(logger logrus.FieldLogger) {
	receiver.logger = logger
}

func NewSqlLogger(logger logrus.StdLogger) ISqlLogger {
	return &SqlLogger{logger: logger}
}

func (receiver *SqlLogger) Log(query string, args ...interface{}) {
	if receiver.Enable() {
		copiedArgs := make([]interface{}, len(args))
		copy(copiedArgs, args)
		for i, arg := range copiedArgs {
			if arg == nil {
				continue
			}
			value := reflectutils.ValueOf(arg)
			if value.IsValid() {
				copiedArgs[i] = value.Interface()
			}
		}
		str := strings.ReplaceAll(fmt.Sprintf(strings.ReplaceAll(query, "?", "'%v'"), copiedArgs...), "'<nil>'", "null")
		re := regexp.MustCompile(`limit '\d+'(,'\d+')?`)
		if re.MatchString(str) {
			str = re.ReplaceAllStringFunc(str, func(s string) string {
				return strings.ReplaceAll(s, "'", "")
			})
		}
		receiver.logger.Printf(str)
	}
}
