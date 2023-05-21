package errorx

import (
	"github.com/unionj-cloud/go-doudou/v2/toolkit/zlogger"
)

func Panic(msg string) {
	zlogger.Panic().Msg(msg)
}

func Wrap(err error) error {
	if err == nil {
		return nil
	}
	zlogger.Error().Err(err).Msg(err.Error())
	return err
}
