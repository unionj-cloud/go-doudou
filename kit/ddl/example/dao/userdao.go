package dao

import (
	"context"
	"github.com/jmoiron/sqlx"
	"github.com/unionj-cloud/go-doudou/kit/ddl/example/domain"
)

type IUserDao interface {
	UpsertUser(ctx context.Context, db *sqlx.DB, user *domain.User) (int64, error)
	GetUser(ctx context.Context, db *sqlx.DB, id int) (domain.User, error)
	//DeleteUsers(ids []interface{}) (int, error)
	//PageUsers(ids []interface{}) ([]domain.User, error)
}

type UserDao struct {
	UserDaoGen
}