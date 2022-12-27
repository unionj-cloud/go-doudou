package main

import (
	"context"
	"mime/multipart"
)

// 用户服务接口
// v1版本
type Usersvc interface {
	// You can define your service methods as your need. Below is an example.
	PageUsers(ctx context.Context, query PageQuery) (code int, data PageRet, msg error)

	// comment1
	// comment2
	GetUser(ctx context.Context,
		// 用户ID
		// @validate(gt=0,lte=10)
		userId,
		// 测试
		school,
		// 图片地址
		photo string,
	) (
		// 业务状态码
		code int,
		// 结果
		data string,
		// 错误
		msg error)

	// comment3
	SignUp(ctx context.Context,
		// @validate(gt=0,lte=10)
		username string, password int, actived bool, score float64) (code int, data string, msg error)

	// comment4
	UploadAvatar(context.Context, []*multipart.FileHeader, string) (int, string, error)

	// comment5
	DownloadAvatar(ctx context.Context, userId string, userAttrs ...int) (a, b string)

	// @role(SUPER_ADMIN)
	BulkSaveOrUpdate(context.Context, int) error
}
