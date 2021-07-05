package service

import (
	"context"
	"testfiles3/vo"
)

type Testfiles3 interface {
	// You can define your service methods as your need. Below is an example.
	PageUsers(ctx context.Context, query vo.PageQuery) (code int, data vo.PageRet, msg error)
}
