package dao_bak

import (
	"context"
	"example/domain"
	"github.com/unionj-cloud/go-doudou/kit/ddl/query"
)

type UserDao interface {
	InsertUser(ctx context.Context, user *domain.User) (int64, error)
	UpdateUser(ctx context.Context, user *domain.User) (int64, error)
	UpdateUserNoneZero(ctx context.Context, user *domain.User) (int64, error)
	UpsertUser(ctx context.Context, user *domain.User) (int64, error)
	UpsertUserNoneZero(ctx context.Context, user *domain.User) (int64, error)
	GetUser(ctx context.Context, id int) (domain.User, error)
	SelectUsers(ctx context.Context, where query.Q) ([]domain.User, error)
	UpdateUsers(ctx context.Context, user domain.User, where query.Q) (int64, error)
	UpdateUsersNoneZero(ctx context.Context, user domain.User, where query.Q) (int64, error)
	CountUsers(ctx context.Context, where query.Q) (int, error)
	DeleteUsers(ctx context.Context, where query.Q) (int64, error)
	PageUsers(ctx context.Context, where query.Q, page query.Page) (query.PageRet, error)
}
