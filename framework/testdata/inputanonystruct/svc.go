package service

import (
	"context"
	"github.com/unionj-cloud/go-doudou/v2/framework/testdata/vo"
)

// 用户服务接口
// v1版本
type Usersvc interface {
	// You can define your service methods as your need. Below is an example.
	PageUsers(ctx context.Context, query struct {
		Filter vo.PageFilter
		Page   vo.Page
	}) (code int, data vo.PageRet, msg error)
}
