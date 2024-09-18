package service

import (
	"context"
	"github.com/unionj-cloud/go-doudou/v2/cmd/internal/svc/testdata/outputanonystruct/dto"
)

// 用户服务接口
// v1版本
type Usersvc interface {
	// You can define your service methods as your need. Below is an example.
	PageUsers(ctx context.Context, query dto.PageQuery) (page dto.Page, msg error)
}
