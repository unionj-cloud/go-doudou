package dao

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"example/user-svc/domain"
	"github.com/unionj-cloud/go-doudou/ddl/query"
	"github.com/unionj-cloud/go-doudou/pathutils"
	"github.com/unionj-cloud/go-doudou/templateutils"
	"math"
)

type UserDaoImpl struct {
	db *sqlx.DB
}

func NewUserDao(db *sqlx.DB) UserDao {
	return UserDaoImpl{
		db: db,
	}
}

func (receiver UserDaoImpl) InsertUser(ctx context.Context, user *domain.User) (int64, error) {
	var (
		statement    string
		err          error
		result       sql.Result
		lastInsertID int64
	)
	if statement, err = templateutils.StringBlockMysql(pathutils.Abs("userdao.sql"), "InsertUser", nil); err != nil {
		return 0, err
	}
	if result, err = receiver.db.NamedExecContext(ctx, statement, user); err != nil {
		return 0, errors.Wrap(err, "error returned from calling db.Exec")
	}
	if lastInsertID, err = result.LastInsertId(); err != nil {
		return 0, errors.Wrap(err, "error returned from calling result.LastInsertId")
	}
	if lastInsertID > 0 {
		user.Id = int(lastInsertID)
	}
	return result.RowsAffected()
}

// With ON DUPLICATE KEY UPDATE, the affected-rows value per row is 1 if the row is inserted as a new row,
// 2 if an existing row is updated, and 0 if an existing row is set to its current values.
// If you specify the CLIENT_FOUND_ROWS flag to the mysql_real_connect() C API function when connecting to mysqld,
// the affected-rows value is 1 (not 0) if an existing row is set to its current values.
// https://dev.mysql.com/doc/refman/5.7/en/insert-on-duplicate.html
func (receiver UserDaoImpl) UpsertUser(ctx context.Context, user *domain.User) (int64, error) {
	var (
		statement    string
		err          error
		result       sql.Result
		lastInsertID int64
	)
	if statement, err = templateutils.StringBlockMysql(pathutils.Abs("userdao.sql"), "UpsertUser", nil); err != nil {
		return 0, err
	}
	if result, err = receiver.db.NamedExecContext(ctx, statement, user); err != nil {
		return 0, errors.Wrap(err, "error returned from calling db.Exec")
	}
	if lastInsertID, err = result.LastInsertId(); err != nil {
		return 0, errors.Wrap(err, "error returned from calling result.LastInsertId")
	}
	if lastInsertID > 0 {
		user.Id = int(lastInsertID)
	}
	return result.RowsAffected()
}

func (receiver UserDaoImpl) UpsertUserNoneZero(ctx context.Context, user *domain.User) (int64, error) {
	var (
		statement    string
		err          error
		result       sql.Result
		lastInsertID int64
	)
	if statement, err = templateutils.StringBlockMysql(pathutils.Abs("userdao.sql"), "UpsertUserNoneZero", user); err != nil {
		return 0, err
	}
	if result, err = receiver.db.ExecContext(ctx, statement); err != nil {
		return 0, errors.Wrap(err, "error returned from calling db.Exec")
	}
	if lastInsertID, err = result.LastInsertId(); err != nil {
		return 0, errors.Wrap(err, "error returned from calling result.LastInsertId")
	}
	if lastInsertID > 0 {
		user.Id = int(lastInsertID)
	}
	return result.RowsAffected()
}

func (receiver UserDaoImpl) DeleteUsers(ctx context.Context, where query.Q) (int64, error) {
	var (
		statement string
		err       error
		result    sql.Result
	)
	statement = fmt.Sprintf("delete from dd_user where %s;", where.Sql())
	if result, err = receiver.db.ExecContext(ctx, statement); err != nil {
		return 0, errors.Wrap(err, "error returned from calling db.ExecContext")
	}
	return result.RowsAffected()
}

func (receiver UserDaoImpl) UpdateUser(ctx context.Context, user *domain.User) (int64, error) {
	var (
		statement string
		err       error
		result    sql.Result
	)
	if statement, err = templateutils.StringBlockMysql(pathutils.Abs("userdao.sql"), "UpdateUser", nil); err != nil {
		return 0, err
	}
	if result, err = receiver.db.NamedExecContext(ctx, statement, user); err != nil {
		return 0, errors.Wrap(err, "error returned from calling db.Exec")
	}
	return result.RowsAffected()
}

func (receiver UserDaoImpl) UpdateUserNoneZero(ctx context.Context, user *domain.User) (int64, error) {
	var (
		statement string
		err       error
		result    sql.Result
	)
	if statement, err = templateutils.StringBlockMysql(pathutils.Abs("userdao.sql"), "UpdateUserNoneZero", user); err != nil {
		return 0, err
	}
	if result, err = receiver.db.ExecContext(ctx, statement); err != nil {
		return 0, errors.Wrap(err, "error returned from calling db.Exec")
	}
	return result.RowsAffected()
}

func (receiver UserDaoImpl) UpdateUsers(ctx context.Context, user domain.User, where query.Q) (int64, error) {
	var (
		statement string
		err       error
		result    sql.Result
	)
	if statement, err = templateutils.StringBlockMysql(pathutils.Abs("userdao.sql"), "UpdateUsers", struct {
		domain.User
		Where string
	}{
		User:  user,
		Where: where.Sql(),
	}); err != nil {
		return 0, err
	}
	if result, err = receiver.db.ExecContext(ctx, statement); err != nil {
		return 0, errors.Wrap(err, "error returned from calling db.Exec")
	}
	return result.RowsAffected()
}

func (receiver UserDaoImpl) UpdateUsersNoneZero(ctx context.Context, user domain.User, where query.Q) (int64, error) {
	var (
		statement string
		err       error
		result    sql.Result
	)
	if statement, err = templateutils.StringBlockMysql(pathutils.Abs("userdao.sql"), "UpdateUsersNoneZero", struct {
		domain.User
		Where string
	}{
		User:  user,
		Where: where.Sql(),
	}); err != nil {
		return 0, err
	}
	if result, err = receiver.db.ExecContext(ctx, statement); err != nil {
		return 0, errors.Wrap(err, "error returned from calling db.Exec")
	}
	return result.RowsAffected()
}

func (receiver UserDaoImpl) GetUser(ctx context.Context, id int) (domain.User, error) {
	var (
		statement string
		err       error
		user      domain.User
	)
	if statement, err = templateutils.StringBlockMysql(pathutils.Abs("userdao.sql"), "GetUser", nil); err != nil {
		return domain.User{}, err
	}
	if err = receiver.db.GetContext(ctx, &user, receiver.db.Rebind(statement), id); err != nil {
		return domain.User{}, errors.Wrap(err, "error returned from calling db.Select")
	}
	return user, nil
}

func (receiver UserDaoImpl) SelectUsers(ctx context.Context, where query.Q) ([]domain.User, error) {
	var (
		statement string
		err       error
		users     []domain.User
	)
	statement = fmt.Sprintf("select * from dd_user where %s", where.Sql())
	if err = receiver.db.SelectContext(ctx, &users, statement); err != nil {
		return nil, errors.Wrap(err, "error returned from calling db.SelectContext")
	}
	return users, nil
}

func (receiver UserDaoImpl) CountUsers(ctx context.Context, where query.Q) (int, error) {
	var (
		statement string
		err       error
		total     int
	)
	statement = fmt.Sprintf("select count(1) from dd_user where %s", where.Sql())
	if err = receiver.db.GetContext(ctx, &total, statement); err != nil {
		return 0, errors.Wrap(err, "error returned from calling db.GetContext")
	}
	return total, nil
}

func (receiver UserDaoImpl) PageUsers(ctx context.Context, where query.Q, page query.Page) (query.PageRet, error) {
	var (
		statement string
		err       error
		users     []domain.User
		total     int
	)
	statement = fmt.Sprintf("select * from dd_user where %s %s;", where.Sql(), page.Sql())
	if err = receiver.db.SelectContext(ctx, &users, statement); err != nil {
		return query.PageRet{}, errors.Wrap(err, "error returned from calling db.SelectContext")
	}

	statement = fmt.Sprintf("select count(1) from dd_user where %s;", where.Sql())
	if err = receiver.db.GetContext(ctx, &total, statement); err != nil {
		return query.PageRet{}, errors.Wrap(err, "error returned from calling db.GetContext")
	}

	pageRet := query.NewPageRet(page)
	pageRet.Items = users
	pageRet.Total = total

	if math.Ceil(float64(total)/float64(pageRet.PageSize)) > float64(pageRet.PageNo) {
		pageRet.HasNext = true
	}

	return pageRet, nil
}
