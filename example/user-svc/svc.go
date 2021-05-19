package service

import (
	"context"
	"example/user-svc/vo"
)

type UserService interface {
	PostSignUp(ctx context.Context, form vo.SignUpForm) (int, error)
	PostLogIn(ctx context.Context, form vo.LogInForm) (vo.Auth, error)
	GetUser(ctx context.Context, id int) (vo.UserVo, error)
	PostPageUsers(ctx context.Context, query vo.PageQuery) (vo.PageRet, error)
}
