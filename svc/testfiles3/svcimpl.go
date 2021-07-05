package service

import (
	"context"
	"testfiles3/config"
	"testfiles3/vo"

	"github.com/jmoiron/sqlx"
)

type Testfiles3Impl struct {
	conf config.Config
}

func (receiver *Testfiles3Impl) PageUsers(ctx context.Context, query vo.PageQuery) (code int, data vo.PageRet, msg error) {
	panic("implement me")
}

func NewTestfiles3(conf config.Config, db *sqlx.DB) Testfiles3 {
	return &Testfiles3Impl{
		conf,
	}
}
