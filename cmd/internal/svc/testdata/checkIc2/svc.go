package service

import (
	"context"
	"mime/multipart"

	v3 "github.com/unionj-cloud/go-doudou/v2/toolkit/openapi/v3"
)

// 用户服务接口
// v1版本
type Usersvc interface {
	// comment4
	UploadAvatar(context.Context, *multipart.FileHeader, v3.FileModel, string, []*multipart.FileHeader, []*multipart.FileHeader) (int, string, error)
}
