package service

import (
	"context"
	"pushcmd/vo"
)

type Pushcmd interface {
	// You can define your service methods as your need. Below is an example.
	PageUsers(ctx context.Context, query vo.PageQuery) (code int, data vo.PageRet, err error)
}
