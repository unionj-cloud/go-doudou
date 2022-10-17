package service

import (
	"context"
	"github.com/unionj-cloud/go-doudou/v2/framework/testdata/vo"
)

// 用户服务接口
// v1版本
type Usersvc interface {
	// You can define your service methods as your need. Below is an example.
	PageUsers(ctx context.Context, query vo.PageQuery) (code int, data struct {
		Items    interface{}
		PageNo   int
		PageSize int
		Total    int
		HasNext  bool
	}, msg error)
}
