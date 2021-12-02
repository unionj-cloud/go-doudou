package service

import (
	"context"
	"svcimplappend/config"
	"svcimplappend/vo"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/jmoiron/sqlx"
)

type SvcimplappendImpl struct {
	conf *config.Config
}

func (receiver *SvcimplappendImpl) PageUsers(ctx context.Context, query vo.PageQuery) (code int, data vo.PageRet, err error) {
	var _result struct {
		Code int
		Data vo.PageRet
	}
	_ = gofakeit.Struct(&_result)
	return _result.Code, _result.Data, nil
}

func NewSvcimplappend(conf *config.Config, db *sqlx.DB) Svcimplappend {
	return &SvcimplappendImpl{
		conf,
	}
}
