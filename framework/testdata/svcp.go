package service

import (
	"context"
	v3 "github.com/unionj-cloud/go-doudou/v2/toolkit/openapi/v3"
	"github.com/unionj-cloud/go-doudou/v2/framework/testdata/vo"
	"os"
)

// 用户服务接口
// v1版本
type UsersvcPanic interface {
	// You can define your service methods as your need. Below is an example.
	PageUsers(ctx context.Context, query vo.PageQuery) (code int, data vo.PageRet, msg error)

	// comment1
	// comment2
	GetUser(ctx context.Context,
		// 用户ID
		userId string,
		// 图片地址
		photo string,
	) (code int, data string, msg error)

	// comment3
	SignUp(ctx context.Context, username string, password int, actived bool, score float64) (code int, data string, msg error)

	// comment4
	UploadAvatar(context.Context, []v3.FileModel, string) (int, string, error)

	// comment5
	DownloadAvatar(ctx context.Context, userId string) (*os.File, error)

	// Wrong1
	Wrong1(ctx context.Context, userId string, query vo.PageQuery, files []v3.FileModel) (*os.File, error)

	// Wrong2
	Wrong2(ctx context.Context, ch chan string) (*os.File, error)

	// Wrong3
	Wrong3(ctx context.Context, fun func()) (*os.File, error)
}
