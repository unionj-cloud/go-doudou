package codegen

import (
	"github.com/sirupsen/logrus"
	"github.com/unionj-cloud/go-doudou/cmd/internal/astutils"
	"github.com/unionj-cloud/go-doudou/cmd/internal/ddl/ddlast"
	"github.com/unionj-cloud/go-doudou/cmd/internal/ddl/table"
	"github.com/unionj-cloud/go-doudou/toolkit/pathutils"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestGenDaoImplGo(t *testing.T) {
	testDir := pathutils.Abs("../testdata")
	err := os.Chdir(testDir)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		err := os.RemoveAll(filepath.Join(testDir, "dao"))
		if err != nil {
			panic(err)
		}
	}()
	dir := testDir + "/domain"
	var files []string
	err = filepath.Walk(dir, astutils.Visit(&files))
	if err != nil {
		logrus.Panicln(err)
	}
	sc := astutils.NewStructCollector(astutils.ExprString)
	usergo := filepath.Join(dir, "user.go")
	fset := token.NewFileSet()
	root, err := parser.ParseFile(fset, usergo, nil, parser.ParseComments)
	if err != nil {
		logrus.Panicln(err)
	}
	ast.Walk(sc, root)

	basego := filepath.Join(dir, "base.go")
	fset = token.NewFileSet()
	root, err = parser.ParseFile(fset, basego, nil, parser.ParseComments)
	if err != nil {
		logrus.Panicln(err)
	}
	ast.Walk(sc, root)

	var tables []table.Table
	flattened := ddlast.FlatEmbed(sc.Structs)
	for _, sm := range flattened {
		tables = append(tables, table.NewTableFromStruct(sm, ""))
	}
	if err := GenBaseGo(dir); (err != nil) != false {
		t.FailNow()
	}
	if err := GenDaoGo(dir, tables[0]); (err != nil) != false {
		t.FailNow()
	}
	if err := GenDaoImplGo(dir, tables[0]); (err != nil) != false {
		t.FailNow()
	}
	expect := `package dao

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/pkg/errors"
	"testdata/domain"
	"github.com/unionj-cloud/go-doudou/toolkit/caller"
	"github.com/unionj-cloud/go-doudou/toolkit/sqlext/query"
	"github.com/unionj-cloud/go-doudou/toolkit/sqlext/wrapper"
	"github.com/unionj-cloud/go-doudou/toolkit/reflectutils"
	"github.com/unionj-cloud/go-doudou/toolkit/stringutils"
	"github.com/unionj-cloud/go-doudou/toolkit/templateutils"
	"strings"
	"math"
	"time"
)

type UserDaoImpl struct {
	db wrapper.Querier
}

func (receiver UserDaoImpl) BeforeSaveHook(ctx context.Context, data interface{}) {
	// implement your business logic
}

func (receiver UserDaoImpl) AfterSaveHook(ctx context.Context, data interface{}, lastInsertID int64, affected int64) {
	// implement your business logic
}

func (receiver UserDaoImpl) BeforeUpdateManyHook(ctx context.Context, data interface{}, where query.Q) {
	// implement your business logic
}

func (receiver UserDaoImpl) AfterUpdateManyHook(ctx context.Context, data interface{}, where query.Q, affected int64) {
	// implement your business logic
}

func (receiver UserDaoImpl) BeforeDeleteManyHook(ctx context.Context, data interface{}, where query.Q) {
	// implement your business logic
}

func (receiver UserDaoImpl) AfterDeleteManyHook(ctx context.Context, data interface{}, where query.Q, affected int64) {
	// implement your business logic
}

func (receiver UserDaoImpl) BeforeReadManyHook(ctx context.Context, page *query.Page, where ...query.Q) {
	// implement your business logic
}

func (receiver UserDaoImpl) DeleteManySoft(ctx context.Context, where query.Q) (int64, error) {
	var (
		err      error
		result   sql.Result
		w        string
		args     []interface{}
		affected int64
	)
	receiver.BeforeDeleteManyHook(ctx, nil, where)
	w, args = where.Sql()
	args = append([]interface{}{time.Now()}, args...)
	if result, err = receiver.db.ExecContext(ctx, receiver.db.Rebind(fmt.Sprintf("update user set delete_at=? where %s;", w)), args...); err != nil {
		return 0, errors.Wrap(err, caller.NewCaller().String())
	}
	if affected, err = result.RowsAffected(); err == nil {
		receiver.AfterDeleteManyHook(ctx, nil, where, affected)
	}
	return affected, err
}

func NewUserDao(querier wrapper.Querier) UserDao {
	return UserDaoImpl{
		db: querier,
	}
}

func (receiver UserDaoImpl) Insert(ctx context.Context, data interface{}) (int64, error) {
	var (
		statement    string
		err          error
		result       sql.Result
		lastInsertID int64
		affected     int64
	)
	receiver.BeforeSaveHook(ctx, data)
	if statement, err = templateutils.BlockMysql("userdao.sql", userdaosql, "InsertUser", nil); err != nil {
		return 0, errors.Wrap(err, caller.NewCaller().String())
	}
	if result, err = receiver.db.NamedExecContext(ctx, statement, data); err != nil {
		return 0, errors.Wrap(err, caller.NewCaller().String())
	}
	if lastInsertID, err = result.LastInsertId(); err != nil {
		return 0, errors.Wrap(err, caller.NewCaller().String())
	}
	if lastInsertID > 0 {
		if user, ok := data.(*domain.User); ok {
			user.ID = int(lastInsertID)
		}
	}
	if affected, err = result.RowsAffected(); err == nil {
		receiver.AfterSaveHook(ctx, data, lastInsertID, affected)
	}
	return affected, err
}

// Upsert With ON DUPLICATE KEY UPDATE, the affected-rows value per row is 1 if the row is inserted as a new row,
// 2 if an existing row is updated, and 0 if an existing row is set to its current values.
// If you specify the CLIENT_FOUND_ROWS flag to the mysql_real_connect() C API function when connecting to mysqld,
// the affected-rows value is 1 (not 0) if an existing row is set to its current values.
// https://dev.mysql.com/doc/refman/5.7/en/insert-on-duplicate.html
func (receiver UserDaoImpl) Upsert(ctx context.Context, data interface{}) (int64, error) {
	var (
		statement    string
		err          error
		result       sql.Result
		lastInsertID int64
		affected     int64
	)
	receiver.BeforeSaveHook(ctx, data)
	if statement, err = templateutils.BlockMysql("userdao.sql", userdaosql, "UpsertUser", nil); err != nil {
		return 0, errors.Wrap(err, caller.NewCaller().String())
	}
	if result, err = receiver.db.NamedExecContext(ctx, statement, data); err != nil {
		return 0, errors.Wrap(err, caller.NewCaller().String())
	}
	if lastInsertID, err = result.LastInsertId(); err != nil {
		return 0, errors.Wrap(err, caller.NewCaller().String())
	}
	if lastInsertID > 0 {
		if user, ok := data.(*domain.User); ok {
			user.ID = int(lastInsertID)
		}
	}
	if affected, err = result.RowsAffected(); err == nil {
		receiver.AfterSaveHook(ctx, data, lastInsertID, affected)
	}
	return affected, err
}

func (receiver UserDaoImpl) UpsertNoneZero(ctx context.Context, data interface{}) (int64, error) {
	var (
		statement    string
		err          error
		result       sql.Result
		lastInsertID int64
		affected     int64
	)
	receiver.BeforeSaveHook(ctx, data)
	value := reflectutils.ValueOf(data).Interface()
	if _, ok := value.(domain.User); !ok {
		return 0, errors.New("underlying type of data should be domain.User")
	}
	if statement, err = templateutils.BlockMysql("userdao.sql", userdaosql, "UpsertUserNoneZero", data); err != nil {
		return 0, errors.Wrap(err, caller.NewCaller().String())
	}
	if result, err = receiver.db.NamedExecContext(ctx, statement, data); err != nil {
		return 0, errors.Wrap(err, caller.NewCaller().String())
	}
	if lastInsertID, err = result.LastInsertId(); err != nil {
		return 0, errors.Wrap(err, caller.NewCaller().String())
	}
	if lastInsertID > 0 {
		if user, ok := data.(*domain.User); ok {
			user.ID = int(lastInsertID)
		}
	}
	if affected, err = result.RowsAffected(); err == nil {
		receiver.AfterSaveHook(ctx, data, lastInsertID, affected)
	}
	return affected, err
}

func (receiver UserDaoImpl) DeleteMany(ctx context.Context, where query.Q) (int64, error) {
	var (
		err    error
		result sql.Result
		w      string
		args   []interface{}
		affected int64
	)
	receiver.BeforeDeleteManyHook(ctx, nil, where)
	w, args = where.Sql()
	if result, err = receiver.db.ExecContext(ctx, receiver.db.Rebind(fmt.Sprintf("delete from user where %s;", w)), args...); err != nil {
		return 0, errors.Wrap(err, caller.NewCaller().String())
	}
	if affected, err = result.RowsAffected(); err == nil {
		receiver.AfterDeleteManyHook(ctx, nil, where, affected)
	}
	return affected, err
}

func (receiver UserDaoImpl) Update(ctx context.Context, data interface{}) (int64, error) {
	var (
		statement string
		err       error
		result    sql.Result
		affected  int64
	)
	receiver.BeforeSaveHook(ctx, data)
	if statement, err = templateutils.BlockMysql("userdao.sql", userdaosql, "UpdateUser", nil); err != nil {
		return 0, errors.Wrap(err, caller.NewCaller().String())
	}
	if result, err = receiver.db.NamedExecContext(ctx, statement, data); err != nil {
		return 0, errors.Wrap(err, caller.NewCaller().String())
	}
	if affected, err = result.RowsAffected(); err == nil {
		receiver.AfterSaveHook(ctx, data, 0, affected)
	}
	return affected, err
}

func (receiver UserDaoImpl) UpdateNoneZero(ctx context.Context, data interface{}) (int64, error) {
	var (
		statement string
		err       error
		result    sql.Result
		affected  int64
	)
	receiver.BeforeSaveHook(ctx, data)
	value := reflectutils.ValueOf(data).Interface()
	if _, ok := value.(domain.User); !ok {
		return 0, errors.New("underlying type of data should be domain.User")
	}
	if statement, err = templateutils.BlockMysql("userdao.sql", userdaosql, "UpdateUserNoneZero", data); err != nil {
		return 0, errors.Wrap(err, caller.NewCaller().String())
	}
	if result, err = receiver.db.NamedExecContext(ctx, statement, data); err != nil {
		return 0, errors.Wrap(err, caller.NewCaller().String())
	}
	if affected, err = result.RowsAffected(); err == nil {
		receiver.AfterSaveHook(ctx, data, 0, affected)
	}
	return affected, err
}

func (receiver UserDaoImpl) UpdateMany(ctx context.Context, data interface{}, where query.Q) (int64, error) {
	var (
		statement string
		err       error
		result    sql.Result
		q         string
		args      []interface{}
		wargs     []interface{}
		w         string
		affected  int64
	)
	receiver.BeforeUpdateManyHook(ctx, data, where)
	if statement, err = templateutils.BlockMysql("userdao.sql", userdaosql, "UpdateUsers", nil); err != nil {
		return 0, errors.Wrap(err, caller.NewCaller().String())
	}
	if q, args, err = receiver.db.BindNamed(statement, data); err != nil {
		return 0, errors.Wrap(err, caller.NewCaller().String())
	}
	w, wargs = where.Sql()
	if stringutils.IsNotEmpty(w) {
		q += " where " + w
	}
	args = append(args, wargs...)
	if result, err = receiver.db.ExecContext(ctx, receiver.db.Rebind(q), args...); err != nil {
		return 0, errors.Wrap(err, caller.NewCaller().String())
	}
	if affected, err = result.RowsAffected(); err == nil {
		receiver.AfterUpdateManyHook(ctx, data, where, affected)
	}
	return affected, err
}

func (receiver UserDaoImpl) UpdateManyNoneZero(ctx context.Context, data interface{}, where query.Q) (int64, error) {
	var (
		statement string
		err       error
		result    sql.Result
		q         string
		args      []interface{}
		wargs     []interface{}
		w         string
		affected  int64
	)
	receiver.BeforeUpdateManyHook(ctx, data, where)
	value := reflectutils.ValueOf(data).Interface()
	if _, ok := value.(domain.User); !ok {
		return 0, errors.New("underlying type of data should be domain.User")
	}
	if statement, err = templateutils.BlockMysql("userdao.sql", userdaosql, "UpdateUsersNoneZero", data); err != nil {
		return 0, errors.Wrap(err, caller.NewCaller().String())
	}
	if q, args, err = receiver.db.BindNamed(statement, data); err != nil {
		return 0, errors.Wrap(err, caller.NewCaller().String())
	}
	w, wargs = where.Sql()
	if stringutils.IsNotEmpty(w) {
		q += " where " + w
	}
	args = append(args, wargs...)
	if result, err = receiver.db.ExecContext(ctx, receiver.db.Rebind(q), args...); err != nil {
		return 0, errors.Wrap(err, caller.NewCaller().String())
	}
	if affected, err = result.RowsAffected(); err == nil {
		receiver.AfterUpdateManyHook(ctx, data, where, affected)
	}
	return affected, err
}

func (receiver UserDaoImpl) Get(ctx context.Context, id interface{}) (interface{}, error) {
	var (
		statement string
		err       error
		user      domain.User
	)
	if statement, err = templateutils.BlockMysql("userdao.sql", userdaosql, "GetUser", nil); err != nil {
		return domain.User{}, errors.Wrap(err, caller.NewCaller().String())
	}
	if err = receiver.db.GetContext(ctx, &user, receiver.db.Rebind(statement), id); err != nil {
		return domain.User{}, errors.Wrap(err, caller.NewCaller().String())
	}
	return user, nil
}

func (receiver UserDaoImpl) SelectMany(ctx context.Context, where ...query.Q) (interface{}, error) {
	var (
		statements []string
		err       error
		users     []domain.User
		args       []interface{}
	)
	receiver.BeforeReadManyHook(ctx, nil, where...)
    statements = append(statements, "select * from user")
    if len(where) > 0 {
        statements = append(statements, "where")
        for _, item := range where {
			q, wargs := item.Sql()
			statements = append(statements, q)
			args = append(args, wargs...)
		}
    }
	sqlStr := strings.TrimSpace(strings.TrimSuffix(strings.TrimSpace(strings.Join(statements, " ")), "where"))
	if err = receiver.db.SelectContext(ctx, &users, receiver.db.Rebind(sqlStr), args...); err != nil {
		return nil, errors.Wrap(err, caller.NewCaller().String())
	}
	return users, nil
}

func (receiver UserDaoImpl) CountMany(ctx context.Context, where ...query.Q) (int, error) {
	var (
		statements []string
		err       error
		total     int
		args       []interface{}
	)
	receiver.BeforeReadManyHook(ctx, nil, where...)
	statements = append(statements, "select count(1) from user")
    if len(where) > 0 {
        statements = append(statements, "where")
        for _, item := range where {
			q, wargs := item.Sql()
			statements = append(statements, q)
			args = append(args, wargs...)
		}
    }
	sqlStr := strings.TrimSpace(strings.TrimSuffix(strings.TrimSpace(strings.Join(statements, " ")), "where"))
	if err = receiver.db.GetContext(ctx, &total, receiver.db.Rebind(sqlStr), args...); err != nil {
		return 0, errors.Wrap(err, caller.NewCaller().String())
	}
	return total, nil
}

func (receiver UserDaoImpl) PageMany(ctx context.Context, page query.Page, where ...query.Q) (query.PageRet, error) {
	var (
		statements []string
		err       error
		users     []domain.User
		total     int
		args       []interface{}
	)
	receiver.BeforeReadManyHook(ctx, &page, where...)
	statements = append(statements, "select * from user")
    if len(where) > 0 {
        statements = append(statements, "where")
        for _, item := range where {
			q, wargs := item.Sql()
			statements = append(statements, q)
			args = append(args, wargs...)
		}
    }
	p, pargs := page.Sql()
	statements = append(statements, p)
	args = append(args, pargs...)
	sqlStr := strings.TrimSpace(strings.TrimSuffix(strings.TrimSpace(strings.Join(statements, " ")), "where"))
	if err = receiver.db.SelectContext(ctx, &users, receiver.db.Rebind(sqlStr), args...); err != nil {
		return query.PageRet{}, errors.Wrap(err, caller.NewCaller().String())
	}
	
	args = nil
    statements = nil
	statements = append(statements, "select count(1) from user")
    if len(where) > 0 {
        statements = append(statements, "where")
        for _, item := range where {
			q, wargs := item.Sql()
			statements = append(statements, q)
			args = append(args, wargs...)
		}
    }
	sqlStr = strings.TrimSpace(strings.TrimSuffix(strings.TrimSpace(strings.Join(statements, " ")), "where"))
	if err = receiver.db.GetContext(ctx, &total, receiver.db.Rebind(sqlStr), args...); err != nil {
		return query.PageRet{}, errors.Wrap(err, caller.NewCaller().String())
	}

	pageRet := query.NewPageRet(page)
	pageRet.Items = users
	pageRet.Total = total
	
	if pageRet.PageSize > 0 && math.Ceil(float64(total)/float64(pageRet.PageSize)) > float64(pageRet.PageNo) {
		pageRet.HasNext = true
	}

	return pageRet, nil
}`
	daofile := filepath.Join(dir, "../dao/userdaoimpl.go")
	f, err := os.Open(daofile)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	content, err := ioutil.ReadAll(f)
	if err != nil {
		t.Fatal(err)
	}
	if string(content) != expect {
		t.Errorf("want %s, got %s\n", expect, string(content))
	}
}
