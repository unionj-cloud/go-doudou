package dao

import (
	"context"
	"database/sql"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/unionj-cloud/go-doudou/kit/ddl/example/domain"
	"github.com/unionj-cloud/go-doudou/kit/pathutils"
	"github.com/unionj-cloud/go-doudou/kit/templateutils"
)

type UserDaoGen struct {
}

// With ON DUPLICATE KEY UPDATE, the affected-rows value per row is 1 if the row is inserted as a new row,
// 2 if an existing row is updated, and 0 if an existing row is set to its current values.
// If you specify the CLIENT_FOUND_ROWS flag to the mysql_real_connect() C API function when connecting to mysqld,
// the affected-rows value is 1 (not 0) if an existing row is set to its current values.
// https://dev.mysql.com/doc/refman/5.7/en/insert-on-duplicate.html
func (u UserDaoGen) UpsertUser(ctx context.Context, db *sqlx.DB, user *domain.User) (int64, error) {
	var (
		statement    string
		err          error
		result       sql.Result
		lastInsertID int64
	)
	if statement, err = templateutils.StringBlock(pathutils.Abs("userdao_gen.sql"), "UpsertUser", nil); err != nil {
		return 0, err
	}
	if result, err = db.NamedExecContext(ctx, statement, user); err != nil {
		return 0, errors.Wrap(err, "error returned from calling db.Exec")
	}
	if lastInsertID, err = result.LastInsertId(); err != nil {
		return 0, errors.Wrap(err, "error returned from calling result.LastInsertId")
	}
	if lastInsertID > 0 {
		user.ID = int(lastInsertID)
	}
	return result.RowsAffected()
}

func (u UserDaoGen) GetUser(ctx context.Context, db *sqlx.DB, id int) (domain.User, error) {
	var (
		statement string
		err       error
		user      domain.User
	)
	if statement, err = templateutils.StringBlock(pathutils.Abs("userdao_gen.sql"), "GetUser", nil); err != nil {
		return domain.User{}, err
	}
	if err = db.Select(&user, db.Rebind(statement), id); err != nil {
		return domain.User{}, errors.Wrap(err, "error returned from calling db.Select")
	}
	return user, nil
}
