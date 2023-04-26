package errorx

import (
	"github.com/unionj-cloud/go-doudou/v2/toolkit/zlogger"
)

func Panic(msg string) {
	zlogger.Panic().Msg(msg)
}

func Log(err error) {
	zlogger.Error().Err(err).Msg(err.Error())
}
