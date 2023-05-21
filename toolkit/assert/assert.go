package assert

import (
	"github.com/unionj-cloud/go-doudou/v2/toolkit/stringutils"
	"github.com/unionj-cloud/go-doudou/v2/toolkit/zlogger"
)

func NotNil(input interface{}, format string, v ...interface{}) {
	if input == nil {
		zlogger.Panic().Msgf(format, v...)
	}
}

func NotTrue(input bool, format string, v ...interface{}) {
	if input {
		zlogger.Panic().Msgf(format, v...)
	}
}

func True(input bool, format string, v ...interface{}) {
	if !input {
		zlogger.Panic().Msgf(format, v...)
	}
}

func NotEmpty(input string, format string, v ...interface{}) {
	if stringutils.IsEmpty(input) {
		zlogger.Panic().Msgf(format, v...)
	}
}

func Empty(input string, format string, v ...interface{}) {
	if stringutils.IsNotEmpty(input) {
		zlogger.Panic().Msgf(format, v...)
	}
}
