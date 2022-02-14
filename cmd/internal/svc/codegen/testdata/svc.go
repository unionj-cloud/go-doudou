package service

import (
	"context"
	v3 "github.com/unionj-cloud/go-doudou/toolkit/openapi/v3"
	"mime/multipart"
	"os"
	"testdata/vo"
)

// 用户服务接口
// v1版本
type Usersvc interface {
	// You can define your service methods as your need. Below is an example.
	PageUsers(ctx context.Context, query vo.PageQuery) (code int, data vo.PageRet, msg error)

	// comment1
	// comment2
	GetUser(ctx context.Context,
		// 用户ID
		userId,
		// 图片地址
		photo string,
	) (code int, data string, msg error)

	// comment3
	SignUp(ctx context.Context, username string, password int, actived bool, score []int) (code int, data string, msg error)

	// comment4
	UploadAvatar(context.Context, []v3.FileModel, string, v3.FileModel, *multipart.FileHeader, []*multipart.FileHeader) (int, interface{}, error)

	// comment5
	DownloadAvatar(ctx context.Context, userId interface{}, userAttrs ...string) (*os.File, error)
}
