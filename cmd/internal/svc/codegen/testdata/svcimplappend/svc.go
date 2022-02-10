package service

import (
	"context"
	"svcimplappend/vo"
)

type Svcimplappend interface {
	// You can define your service methods as your need. Below is an example.
	PageUsers(ctx context.Context, query vo.PageQuery) (code int, data vo.PageRet, err error)

	// comment1
	// comment2
	GetUser(ctx context.Context,
		// 用户ID
		userId string,
		// 图片地址
		photo string,
	) (code int, data string, msg error)
}
