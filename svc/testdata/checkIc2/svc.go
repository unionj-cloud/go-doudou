package service

import (
	"context"
	v3 "github.com/unionj-cloud/go-doudou/openapi/v3"
	"mime/multipart"
)

// 用户服务接口
// v1版本
type Usersvc interface {
	// comment4
	UploadAvatar(context.Context, *multipart.FileHeader, v3.FileModel, string, []*multipart.FileHeader, []*multipart.FileHeader) (int, string, error)
}
