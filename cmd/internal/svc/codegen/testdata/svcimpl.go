package service

import (
	"context"
	"mime/multipart"
	"os"
	"testdata/config"
	"testdata/vo"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/jmoiron/sqlx"
	v3 "github.com/unionj-cloud/go-doudou/v2/toolkit/openapi/v3"
)

type UsersvcImpl struct {
	conf *config.Config
}

func (receiver *UsersvcImpl) PageUsers(ctx context.Context, query vo.PageQuery) (code int, data vo.PageRet, msg error) {
	var _result struct {
		Code int
		Data vo.PageRet
	}
	_ = gofakeit.Struct(&_result)
	return _result.Code, _result.Data, nil
}
func (receiver *UsersvcImpl) GetUser(ctx context.Context, userId string, photo string) (code int, data string, msg error) {
	var _result struct {
		Code int
		Data string
	}
	_ = gofakeit.Struct(&_result)
	return _result.Code, _result.Data, nil
}
func (receiver *UsersvcImpl) SignUp(ctx context.Context, username string, password int, actived bool, score []int) (code int, data string, msg error) {
	var _result struct {
		Code int
		Data string
	}
	_ = gofakeit.Struct(&_result)
	return _result.Code, _result.Data, nil
}
func (receiver *UsersvcImpl) UploadAvatar(pc context.Context, pf []v3.FileModel, ps string, pf2 v3.FileModel, pf3 *multipart.FileHeader, pf4 []*multipart.FileHeader) (ri int, rs string, re error) {
	var _result struct {
		Ri int
		Rs string
	}
	_ = gofakeit.Struct(&_result)
	return _result.Ri, _result.Rs, nil
}
func (receiver *UsersvcImpl) DownloadAvatar(ctx context.Context, userId string, userAttrs ...string) (rf *os.File, re error) {
	var _result struct {
		Rf *os.File
	}
	_ = gofakeit.Struct(&_result)
	return _result.Rf, nil
}

func NewUsersvc(conf *config.Config, db *sqlx.DB) Usersvc {
	return &UsersvcImpl{
		conf,
	}
}
