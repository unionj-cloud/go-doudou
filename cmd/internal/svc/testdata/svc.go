package service

import (
	"context"
	"github.com/unionj-cloud/go-doudou/cmd/internal/svc/testdata/vo"
	"os"
)

// 用户服务接口
// v1版本
type Usersvc interface {
	//// You can define your service methods as your need. Below is an example.
	//PageUsers(ctx context.Context, query vo.PageQuery) (code int, data vo.PageRet, msg error)
	//
	//// comment1
	//// comment2
	//GetUser(ctx context.Context,
	//	// 用户ID
	//	userId string,
	//	// 图片地址
	//	photo string,
	//) (code int, data string, msg error)
	//
	//// comment3
	//SignUp(ctx context.Context, username string, password int, actived bool, score float64) (code int, data string, msg error)
	//
	//// comment4
	//UploadAvatar(context.Context, []*multipart.FileHeader, []*multipart.FileHeader, *multipart.FileHeader, v3.FileModel, string) (int, string, error)

	// comment5
	DownloadAvatar(ctx context.Context, userId string, layout vo.KeyboardLayout) (*os.File, error)
}
