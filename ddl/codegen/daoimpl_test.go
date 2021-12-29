package codegen

import (
	"github.com/sirupsen/logrus"
	"github.com/unionj-cloud/go-doudou/astutils"
	"github.com/unionj-cloud/go-doudou/ddl/ddlast"
	"github.com/unionj-cloud/go-doudou/ddl/table"
	"github.com/unionj-cloud/go-doudou/pathutils"
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
	type args struct {
		domainpath string
		t          table.Table
		folder     []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "",
			args: args{
				domainpath: dir,
				t:          tables[0],
				folder:     nil,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := GenDaoGo(tt.args.domainpath, tt.args.t, tt.args.folder...); (err != nil) != tt.wantErr {
				t.Errorf("GenDaoGo() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err := GenDaoImplGo(tt.args.domainpath, tt.args.t, tt.args.folder...); (err != nil) != tt.wantErr {
				t.Errorf("GenDaoGo() error = %v, wantErr %v", err, tt.wantErr)
			}
			defer os.RemoveAll(filepath.Join(dir, "../dao"))
			expect := `package dao

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/pkg/errors"
	"testdata/domain"
	"github.com/unionj-cloud/go-doudou/ddl/query"
	"github.com/unionj-cloud/go-doudou/ddl/wrapper"
	"github.com/unionj-cloud/go-doudou/reflectutils"
	"github.com/unionj-cloud/go-doudou/stringutils"
	"github.com/unionj-cloud/go-doudou/templateutils"
	"strings"
	"math"
)

type UserDaoImpl struct {
	db wrapper.Querier
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
	)
	if statement, err = templateutils.BlockMysql("userdao.sql", userdaosql, "InsertUser", nil); err != nil {
		return 0, errors.Wrap(err, "")
	}
	if result, err = receiver.db.NamedExecContext(ctx, statement, data); err != nil {
		return 0, errors.Wrap(err, "")
	}
	if lastInsertID, err = result.LastInsertId(); err != nil {
		return 0, errors.Wrap(err, "")
	}
	if lastInsertID > 0 {
		if user, ok := data.(*domain.User); ok {
			user.ID = int(lastInsertID)
		}
	}
	return result.RowsAffected()
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
	)
	if statement, err = templateutils.BlockMysql("userdao.sql", userdaosql, "UpsertUser", nil); err != nil {
		return 0, errors.Wrap(err, "")
	}
	if result, err = receiver.db.NamedExecContext(ctx, statement, data); err != nil {
		return 0, errors.Wrap(err, "")
	}
	if lastInsertID, err = result.LastInsertId(); err != nil {
		return 0, errors.Wrap(err, "")
	}
	if lastInsertID > 0 {
		if user, ok := data.(*domain.User); ok {
			user.ID = int(lastInsertID)
		}
	}
	return result.RowsAffected()
}

func (receiver UserDaoImpl) UpsertNoneZero(ctx context.Context, data interface{}) (int64, error) {
	var (
		statement    string
		err          error
		result       sql.Result
		lastInsertID int64
	)
	value := reflectutils.ValueOf(data).Interface()
	if _, ok := value.(domain.User); !ok {
		return 0, errors.New("underlying type of data should be domain.User")
	}
	if statement, err = templateutils.BlockMysql("userdao.sql", userdaosql, "UpsertUserNoneZero", data); err != nil {
		return 0, errors.Wrap(err, "")
	}
	if result, err = receiver.db.NamedExecContext(ctx, statement, data); err != nil {
		return 0, errors.Wrap(err, "")
	}
	if lastInsertID, err = result.LastInsertId(); err != nil {
		return 0, errors.Wrap(err, "")
	}
	if lastInsertID > 0 {
		if user, ok := data.(*domain.User); ok {
			user.ID = int(lastInsertID)
		}
	}
	return result.RowsAffected()
}

func (receiver UserDaoImpl) DeleteMany(ctx context.Context, where query.Q) (int64, error) {
	var (
		err    error
		result sql.Result
		w      string
		args   []interface{}
	)
	w, args = where.Sql()
	if result, err = receiver.db.ExecContext(ctx, receiver.db.Rebind(fmt.Sprintf("delete from user where %s;", w)), args...); err != nil {
		return 0, errors.Wrap(err, "")
	}
	return result.RowsAffected()
}

func (receiver UserDaoImpl) Update(ctx context.Context, data interface{}) (int64, error) {
	var (
		statement string
		err       error
		result    sql.Result
	)
	if statement, err = templateutils.BlockMysql("userdao.sql", userdaosql, "UpdateUser", nil); err != nil {
		return 0, errors.Wrap(err, "")
	}
	if result, err = receiver.db.NamedExecContext(ctx, statement, data); err != nil {
		return 0, errors.Wrap(err, "")
	}
	return result.RowsAffected()
}

func (receiver UserDaoImpl) UpdateNoneZero(ctx context.Context, data interface{}) (int64, error) {
	var (
		statement string
		err       error
		result    sql.Result
	)
	value := reflectutils.ValueOf(data).Interface()
	if _, ok := value.(domain.User); !ok {
		return 0, errors.New("underlying type of data should be domain.User")
	}
	if statement, err = templateutils.BlockMysql("userdao.sql", userdaosql, "UpdateUserNoneZero", data); err != nil {
		return 0, errors.Wrap(err, "")
	}
	if result, err = receiver.db.NamedExecContext(ctx, statement, data); err != nil {
		return 0, errors.Wrap(err, "")
	}
	return result.RowsAffected()
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
	)
	if statement, err = templateutils.BlockMysql("userdao.sql", userdaosql, "UpdateUsers", nil); err != nil {
		return 0, errors.Wrap(err, "")
	}
	if q, args, err = receiver.db.BindNamed(statement, data); err != nil {
		return 0, errors.Wrap(err, "")
	}
	w, wargs = where.Sql()
	if stringutils.IsNotEmpty(w) {
		q += " where " + w
	}
	args = append(args, wargs...)
	if result, err = receiver.db.ExecContext(ctx, receiver.db.Rebind(q), args...); err != nil {
		return 0, errors.Wrap(err, "")
	}
	return result.RowsAffected()
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
	)
	value := reflectutils.ValueOf(data).Interface()
	if _, ok := value.(domain.User); !ok {
		return 0, errors.New("underlying type of data should be domain.User")
	}
	if statement, err = templateutils.BlockMysql("userdao.sql", userdaosql, "UpdateUsersNoneZero", data); err != nil {
		return 0, errors.Wrap(err, "")
	}
	if q, args, err = receiver.db.BindNamed(statement, data); err != nil {
		return 0, errors.Wrap(err, "")
	}
	w, wargs = where.Sql()
	if stringutils.IsNotEmpty(w) {
		q += " where " + w
	}
	args = append(args, wargs...)
	if result, err = receiver.db.ExecContext(ctx, receiver.db.Rebind(q), args...); err != nil {
		return 0, errors.Wrap(err, "")
	}
	return result.RowsAffected()
}

func (receiver UserDaoImpl) Get(ctx context.Context, id interface{}) (interface{}, error) {
	var (
		statement string
		err       error
		user      domain.User
	)
	if statement, err = templateutils.BlockMysql("userdao.sql", userdaosql, "GetUser", nil); err != nil {
		return domain.User{}, errors.Wrap(err, "")
	}
	if err = receiver.db.GetContext(ctx, &user, receiver.db.Rebind(statement), id); err != nil {
		return domain.User{}, errors.Wrap(err, "")
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
    statements = append(statements, "select * from user")
    if len(where) > 0 {
        statements = append(statements, "where")
        for _, item := range where {
			q, wargs := item.Sql()
			statements = append(statements, q)
			args = append(args, wargs...)
		}
    }
	if err = receiver.db.SelectContext(ctx, &users, strings.Join(statements, " "), args...); err != nil {
		return nil, errors.Wrap(err, "")
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
	statements = append(statements, "select count(1) from user")
    if len(where) > 0 {
        statements = append(statements, "where")
        for _, item := range where {
			q, wargs := item.Sql()
			statements = append(statements, q)
			args = append(args, wargs...)
		}
    }
	if err = receiver.db.GetContext(ctx, &total, strings.Join(statements, " "), args...); err != nil {
		return 0, errors.Wrap(err, "")
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
	q := strings.Join(statements, " ")
	if err = receiver.db.SelectContext(ctx, &users, q, args...); err != nil {
		return query.PageRet{}, errors.Wrap(err, "")
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
	if err = receiver.db.GetContext(ctx, &total, strings.Join(statements, " "), args...); err != nil {
		return query.PageRet{}, errors.Wrap(err, "")
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
			content, err := ioutil.ReadAll(f)
			if err != nil {
				t.Fatal(err)
			}
			if string(content) != expect {
				t.Errorf("want %s, got %s\n", expect, string(content))
			}
		})
	}
}
