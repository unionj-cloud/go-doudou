package errorx

import (
	"github.com/pkg/errors"
	"github.com/unionj-cloud/go-doudou/v2/toolkit/caller"
)

func Handle(err error) error {
	return errors.Wrap(err, caller.NewCaller().String())
}
