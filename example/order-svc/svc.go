package service

import (
	"context"
	"github.com/unionj-cloud/go-doudou/ddl/query"
	"order-svc/vo"
)

type OrderSvc interface {
	// You can define your service methods as your need. Below is an example.
	PageUsers(context.Context, vo.PageQuery) (query.PageRet, error)
	GetUser(ctx context.Context, id int) (vo.UserVo, error)
}
