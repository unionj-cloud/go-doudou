package service

import (
	"context"
	"testsvc/config"
	"testsvc/vo"
)

var _ Testsvc = (*TestsvcImpl)(nil)

type TestsvcImpl struct {
	conf *config.Config
}

func (receiver *TestsvcImpl) PageUsers(ctx context.Context, query vo.PageQuery) (code int, data vo.PageRet, err error) {
	panic("implement me")
}

func NewTestsvc(conf *config.Config) Testsvc {
	return &TestsvcImpl{
		conf,
	}
}
