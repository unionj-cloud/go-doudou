package service

import (
	"context"
	"example/user-svc/config"
	"example/user-svc/vo"
	"github.com/jmoiron/sqlx"
	"github.com/unionj-cloud/go-doudou/ddl/query"
)

type userServiceImpl struct {
	conf config.SvcConfig
}

func (u *userServiceImpl) PostSignUp(ctx context.Context, form vo.SignUpForm) (int, error) {
	panic("implement me")
}

func (u *userServiceImpl) PostLogIn(ctx context.Context, form vo.LogInForm) (vo.Auth, error) {
	panic("implement me")
}

func (u *userServiceImpl) GetUser(ctx context.Context, id int) (vo.UserVo, error) {
	panic("implement me")
}

func (u *userServiceImpl) PostPageUsers(ctx context.Context, query vo.PageQuery) (query.PageRet, error) {
	panic("implement me")
}

func NewUserService(conf config.SvcConfig, db *sqlx.DB) UserService {
	return &userServiceImpl{
		conf,
	}
}
