package dao

import (
	"context"
	"github.com/unionj-cloud/go-doudou/kit/ddl/example/domain"
	"github.com/unionj-cloud/go-doudou/kit/ddl/query"
)

type IUserDao interface {
	UpsertUser(ctx context.Context, user *domain.User) (int64, error)
	GetUser(ctx context.Context, id int) (domain.User, error)
	DeleteUsers(ctx context.Context, where query.Q) (int64, error)
	PageUsers(ctx context.Context, where query.Q, page query.Page) (query.PageRet, error)
}

type UserDao struct {
	UserDaoGen
}
