package service

import (
	"context"
	"mime/multipart"
	"os"
	"testdata/vo"

	"github.com/shopspring/decimal"
	v3 "github.com/unionj-cloud/toolkit/openapi/v3"
)

// 用户服务接口
// v1版本
type Usersvc interface {
	// You can define your service methods as your need. Below is an example.@role(user)
	PageUsers(ctx context.Context, query vo.PageQuery) (code int, data vo.PageRet, msg error)

	// comment1
	// comment2
	// @role(admin)
	GetUser(ctx context.Context,
		// 用户ID
		userId,
		// 图片地址
		photo string,
	) (code int, data string, msg error)

	// comment3
	// @permission(create,update)@role(admin)
	SignUp(ctx context.Context, username string, password int, actived bool, score []int) (code int, data string, msg error)

	// comment4
	// @role(user)
	UploadAvatar(context.Context, []v3.FileModel, string, v3.FileModel, *multipart.FileHeader, []*multipart.FileHeader) (int, interface{}, error)

	// comment5
	DownloadAvatar(ctx context.Context, userId interface{}, data []byte, price decimal.Decimal, userAttrs ...string) (*os.File, error)

	GetQuery_range(ctx context.Context) error
	GetShelves_ShelfBooks_Book(ctx context.Context) error
}
