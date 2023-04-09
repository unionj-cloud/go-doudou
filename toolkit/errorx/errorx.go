package errorx

import (
	"github.com/pkg/errors"
	"github.com/unionj-cloud/go-doudou/v2/toolkit/caller"
	"github.com/unionj-cloud/go-doudou/v2/toolkit/zlogger"
)

func Handle(err error) error {
	return errors.Wrap(err, caller.NewCaller().String())
}

func Panic(msg string) {
	zlogger.Panic().Msg(msg)
}

func Log(err error) {
	zlogger.Error().Err(err).Msg(err.Error())
}
