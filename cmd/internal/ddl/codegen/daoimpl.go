package codegen

import (
	"github.com/iancoleman/strcase"
	log "github.com/sirupsen/logrus"
	"github.com/unionj-cloud/go-doudou/cmd/internal/astutils"
	"github.com/unionj-cloud/go-doudou/cmd/internal/ddl/table"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

var daoimpltmpl = `package dao

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/pkg/errors"
	"{{.DomainPackage}}"
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

type {{.DomainName}}DaoImpl struct {
	db wrapper.Querier
}

func (receiver {{.DomainName}}DaoImpl) BeforeSaveHook(ctx context.Context, data interface{}) {
	// implement your business logic
}

func (receiver {{.DomainName}}DaoImpl) AfterSaveHook(ctx context.Context, data interface{}, lastInsertID int64, affected int64) {
	// implement your business logic
}

func (receiver {{.DomainName}}DaoImpl) BeforeUpdateManyHook(ctx context.Context, data interface{}, where query.Q) {
	// implement your business logic
}

func (receiver {{.DomainName}}DaoImpl) AfterUpdateManyHook(ctx context.Context, data interface{}, where query.Q, affected int64) {
	// implement your business logic
}

func (receiver {{.DomainName}}DaoImpl) BeforeDeleteManyHook(ctx context.Context, data interface{}, where query.Q) {
	// implement your business logic
}

func (receiver {{.DomainName}}DaoImpl) AfterDeleteManyHook(ctx context.Context, data interface{}, where query.Q, affected int64) {
	// implement your business logic
}

func (receiver {{.DomainName}}DaoImpl) BeforeReadManyHook(ctx context.Context, page *query.Page, where ...query.Q) {
	// implement your business logic
}

func (receiver {{.DomainName}}DaoImpl) DeleteManySoft(ctx context.Context, where query.Q) (int64, error) {
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
	if result, err = receiver.db.ExecContext(ctx, receiver.db.Rebind(fmt.Sprintf("update {{.TableName}} set delete_at=? where %s;", w)), args...); err != nil {
		return 0, errors.Wrap(err, caller.NewCaller().String())
	}
	if affected, err = result.RowsAffected(); err == nil {
		receiver.AfterDeleteManyHook(ctx, nil, where, affected)
	}
	return affected, err
}

func New{{.DomainName}}Dao(querier wrapper.Querier) {{.DomainName}}Dao {
	return {{.DomainName}}DaoImpl{
		db: querier,
	}
}

func (receiver {{.DomainName}}DaoImpl) Insert(ctx context.Context, data interface{}) (int64, error) {
	var (
		statement    string
		err          error
		result       sql.Result
		{{- if .PkCol.Autoincrement }}
		lastInsertID int64
		{{- end }}
		affected     int64
	)
	receiver.BeforeSaveHook(ctx, data)
	if statement, err = templateutils.BlockMysql("{{.DomainName | ToLower}}dao.sql", {{.DomainName | ToLower}}daosql, "Insert{{.DomainName}}", nil); err != nil {
		return 0, errors.Wrap(err, caller.NewCaller().String())
	}
	if result, err = receiver.db.NamedExecContext(ctx, statement, data); err != nil {
		return 0, errors.Wrap(err, caller.NewCaller().String())
	}
	{{- if .PkCol.Autoincrement }}
	if lastInsertID, err = result.LastInsertId(); err != nil {
		return 0, errors.Wrap(err, caller.NewCaller().String())
	}
	if lastInsertID > 0 {
		if {{.DomainName | ToLower}}, ok := data.(*domain.{{.DomainName}}); ok {
			{{- if eq .PkField.Type "int64"}}
			{{.DomainName | ToLower}}.{{.PkField.Name}} = lastInsertID
			{{- else }}
			{{.DomainName | ToLower}}.{{.PkField.Name}} = {{.PkField.Type}}(lastInsertID)
			{{- end }}
		}
	}
	{{- end }}
	if affected, err = result.RowsAffected(); err == nil {
		{{- if .PkCol.Autoincrement }}
		receiver.AfterSaveHook(ctx, data, lastInsertID, affected)
		{{- else }}
		receiver.AfterSaveHook(ctx, data, 0, affected)
		{{- end }}
	}
	return affected, err
}

// Upsert With ON DUPLICATE KEY UPDATE, the affected-rows value per row is 1 if the row is inserted as a new row,
// 2 if an existing row is updated, and 0 if an existing row is set to its current values.
// If you specify the CLIENT_FOUND_ROWS flag to the mysql_real_connect() C API function when connecting to mysqld,
// the affected-rows value is 1 (not 0) if an existing row is set to its current values.
// https://dev.mysql.com/doc/refman/5.7/en/insert-on-duplicate.html
func (receiver {{.DomainName}}DaoImpl) Upsert(ctx context.Context, data interface{}) (int64, error) {
	var (
		statement    string
		err          error
		result       sql.Result
		{{- if .PkCol.Autoincrement }}
		lastInsertID int64
		{{- end }}
		affected     int64
	)
	receiver.BeforeSaveHook(ctx, data)
	if statement, err = templateutils.BlockMysql("{{.DomainName | ToLower}}dao.sql", {{.DomainName | ToLower}}daosql, "Upsert{{.DomainName}}", nil); err != nil {
		return 0, errors.Wrap(err, caller.NewCaller().String())
	}
	if result, err = receiver.db.NamedExecContext(ctx, statement, data); err != nil {
		return 0, errors.Wrap(err, caller.NewCaller().String())
	}
	{{- if .PkCol.Autoincrement }}
	if lastInsertID, err = result.LastInsertId(); err != nil {
		return 0, errors.Wrap(err, caller.NewCaller().String())
	}
	if lastInsertID > 0 {
		if {{.DomainName | ToLower}}, ok := data.(*domain.{{.DomainName}}); ok {
			{{- if eq .PkField.Type "int64"}}
			{{.DomainName | ToLower}}.{{.PkField.Name}} = lastInsertID
			{{- else }}
			{{.DomainName | ToLower}}.{{.PkField.Name}} = {{.PkField.Type}}(lastInsertID)
			{{- end }}
		}
	}
	{{- end }}
	if affected, err = result.RowsAffected(); err == nil {
		{{- if .PkCol.Autoincrement }}
		receiver.AfterSaveHook(ctx, data, lastInsertID, affected)
		{{- else }}
		receiver.AfterSaveHook(ctx, data, 0, affected)
		{{- end }}
	}
	return affected, err
}

func (receiver {{.DomainName}}DaoImpl) UpsertNoneZero(ctx context.Context, data interface{}) (int64, error) {
	var (
		statement    string
		err          error
		result       sql.Result
		{{- if .PkCol.Autoincrement }}
		lastInsertID int64
		{{- end }}
		affected     int64
	)
	receiver.BeforeSaveHook(ctx, data)
	value := reflectutils.ValueOf(data).Interface()
	if _, ok := value.(domain.{{.DomainName}}); !ok {
		return 0, errors.New("underlying type of data should be domain.{{.DomainName}}")
	}
	if statement, err = templateutils.BlockMysql("{{.DomainName | ToLower}}dao.sql", {{.DomainName | ToLower}}daosql, "Upsert{{.DomainName}}NoneZero", data); err != nil {
		return 0, errors.Wrap(err, caller.NewCaller().String())
	}
	if result, err = receiver.db.NamedExecContext(ctx, statement, data); err != nil {
		return 0, errors.Wrap(err, caller.NewCaller().String())
	}
	{{- if .PkCol.Autoincrement }}
	if lastInsertID, err = result.LastInsertId(); err != nil {
		return 0, errors.Wrap(err, caller.NewCaller().String())
	}
	if lastInsertID > 0 {
		if {{.DomainName | ToLower}}, ok := data.(*domain.{{.DomainName}}); ok {
			{{- if eq .PkField.Type "int64"}}
			{{.DomainName | ToLower}}.{{.PkField.Name}} = lastInsertID
			{{- else }}
			{{.DomainName | ToLower}}.{{.PkField.Name}} = {{.PkField.Type}}(lastInsertID)
			{{- end }}
		}
	}
	{{- end }}
	if affected, err = result.RowsAffected(); err == nil {
		{{- if .PkCol.Autoincrement }}
		receiver.AfterSaveHook(ctx, data, lastInsertID, affected)
		{{- else }}
		receiver.AfterSaveHook(ctx, data, 0, affected)
		{{- end }}
	}
	return affected, err
}

func (receiver {{.DomainName}}DaoImpl) DeleteMany(ctx context.Context, where query.Q) (int64, error) {
	var (
		err    error
		result sql.Result
		w      string
		args   []interface{}
		affected int64
	)
	receiver.BeforeDeleteManyHook(ctx, nil, where)
	w, args = where.Sql()
	if result, err = receiver.db.ExecContext(ctx, receiver.db.Rebind(fmt.Sprintf("delete from {{.TableName}} where %s;", w)), args...); err != nil {
		return 0, errors.Wrap(err, caller.NewCaller().String())
	}
	if affected, err = result.RowsAffected(); err == nil {
		receiver.AfterDeleteManyHook(ctx, nil, where, affected)
	}
	return affected, err
}

func (receiver {{.DomainName}}DaoImpl) Update(ctx context.Context, data interface{}) (int64, error) {
	var (
		statement string
		err       error
		result    sql.Result
		affected  int64
	)
	receiver.BeforeSaveHook(ctx, data)
	if statement, err = templateutils.BlockMysql("{{.DomainName | ToLower}}dao.sql", {{.DomainName | ToLower}}daosql, "Update{{.DomainName}}", nil); err != nil {
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

func (receiver {{.DomainName}}DaoImpl) UpdateNoneZero(ctx context.Context, data interface{}) (int64, error) {
	var (
		statement string
		err       error
		result    sql.Result
		affected  int64
	)
	receiver.BeforeSaveHook(ctx, data)
	value := reflectutils.ValueOf(data).Interface()
	if _, ok := value.(domain.{{.DomainName}}); !ok {
		return 0, errors.New("underlying type of data should be domain.{{.DomainName}}")
	}
	if statement, err = templateutils.BlockMysql("{{.DomainName | ToLower}}dao.sql", {{.DomainName | ToLower}}daosql, "Update{{.DomainName}}NoneZero", data); err != nil {
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

func (receiver {{.DomainName}}DaoImpl) UpdateMany(ctx context.Context, data interface{}, where query.Q) (int64, error) {
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
	if statement, err = templateutils.BlockMysql("{{.DomainName | ToLower}}dao.sql", {{.DomainName | ToLower}}daosql, "Update{{.DomainName}}s", nil); err != nil {
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

func (receiver {{.DomainName}}DaoImpl) UpdateManyNoneZero(ctx context.Context, data interface{}, where query.Q) (int64, error) {
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
	if _, ok := value.(domain.{{.DomainName}}); !ok {
		return 0, errors.New("underlying type of data should be domain.{{.DomainName}}")
	}
	if statement, err = templateutils.BlockMysql("{{.DomainName | ToLower}}dao.sql", {{.DomainName | ToLower}}daosql, "Update{{.DomainName}}sNoneZero", data); err != nil {
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

func (receiver {{.DomainName}}DaoImpl) Get(ctx context.Context, id interface{}) (interface{}, error) {
	var (
		statement string
		err       error
		{{.DomainName | ToLower}}      domain.{{.DomainName}}
	)
	if statement, err = templateutils.BlockMysql("{{.DomainName | ToLower}}dao.sql", {{.DomainName | ToLower}}daosql, "Get{{.DomainName}}", nil); err != nil {
		return domain.{{.DomainName}}{}, errors.Wrap(err, caller.NewCaller().String())
	}
	if err = receiver.db.GetContext(ctx, &{{.DomainName | ToLower}}, receiver.db.Rebind(statement), id); err != nil {
		return domain.{{.DomainName}}{}, errors.Wrap(err, caller.NewCaller().String())
	}
	return {{.DomainName | ToLower}}, nil
}

func (receiver {{.DomainName}}DaoImpl) SelectMany(ctx context.Context, where ...query.Q) (interface{}, error) {
	var (
		statements []string
		err       error
		{{.DomainName | ToLower}}s     []domain.{{.DomainName}}
		args       []interface{}
	)
	receiver.BeforeReadManyHook(ctx, nil, where...)
    statements = append(statements, "select * from {{.TableName}}")
    if len(where) > 0 {
        statements = append(statements, "where")
        for _, item := range where {
			q, wargs := item.Sql()
			statements = append(statements, q)
			args = append(args, wargs...)
		}
    }
	sqlStr := strings.TrimSpace(strings.TrimSuffix(strings.TrimSpace(strings.Join(statements, " ")), "where"))
	if err = receiver.db.SelectContext(ctx, &{{.DomainName | ToLower}}s, receiver.db.Rebind(sqlStr), args...); err != nil {
		return nil, errors.Wrap(err, caller.NewCaller().String())
	}
	return {{.DomainName | ToLower}}s, nil
}

func (receiver {{.DomainName}}DaoImpl) CountMany(ctx context.Context, where ...query.Q) (int, error) {
	var (
		statements []string
		err       error
		total     int
		args       []interface{}
	)
	receiver.BeforeReadManyHook(ctx, nil, where...)
	statements = append(statements, "select count(1) from {{.TableName}}")
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

func (receiver {{.DomainName}}DaoImpl) PageMany(ctx context.Context, page query.Page, where ...query.Q) (query.PageRet, error) {
	var (
		statements []string
		err       error
		{{.DomainName | ToLower}}s     []domain.{{.DomainName}}
		total     int
		args       []interface{}
	)
	receiver.BeforeReadManyHook(ctx, &page, where...)
	statements = append(statements, "select * from {{.TableName}}")
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
	if err = receiver.db.SelectContext(ctx, &{{.DomainName | ToLower}}s, receiver.db.Rebind(sqlStr), args...); err != nil {
		return query.PageRet{}, errors.Wrap(err, caller.NewCaller().String())
	}
	
	args = nil
    statements = nil
	statements = append(statements, "select count(1) from {{.TableName}}")
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
		return query.PageRet{}, errors.Wrap(err, caller.NewCaller().String())
	}

	pageRet := query.NewPageRet(page)
	pageRet.Items = {{.DomainName | ToLower}}s
	pageRet.Total = total
	
	if pageRet.PageSize > 0 && math.Ceil(float64(total)/float64(pageRet.PageSize)) > float64(pageRet.PageNo) {
		pageRet.HasNext = true
	}

	return pageRet, nil
}`

// GenDaoImplGo generates dao layer implementation code
func GenDaoImplGo(domainpath string, t table.Table, folder ...string) error {
	var (
		err      error
		dpkg     string
		daopath  string
		f        *os.File
		funcMap  map[string]interface{}
		tpl      *template.Template
		pkColumn table.Column
		df       string
	)
	df = "dao"
	if len(folder) > 0 {
		df = folder[0]
	}
	daopath = filepath.Join(filepath.Dir(domainpath), df)
	_ = os.MkdirAll(daopath, os.ModePerm)

	daofile := filepath.Join(daopath, strings.ToLower(t.Meta.Name)+"daoimpl.go")
	if _, err = os.Stat(daofile); os.IsNotExist(err) {
		f, _ = os.Create(daofile)
		defer f.Close()

		dpkg = astutils.GetImportPath(domainpath)
		funcMap = make(map[string]interface{})
		funcMap["ToLower"] = strings.ToLower
		funcMap["ToSnake"] = strcase.ToSnake
		tpl, _ = template.New("daoimpl.go.tmpl").Funcs(funcMap).Parse(daoimpltmpl)
		for _, column := range t.Columns {
			if column.Pk {
				pkColumn = column
				break
			}
		}
		_ = tpl.Execute(f, struct {
			DomainPackage string
			DomainName    string
			TableName     string
			PkField       astutils.FieldMeta
			PkCol         table.Column
		}{
			DomainPackage: dpkg,
			DomainName:    t.Meta.Name,
			TableName:     t.Name,
			PkField:       pkColumn.Meta,
			PkCol:         pkColumn,
		})
	} else {
		log.Warnf("file %s already exists", daofile)
	}
	return nil
}
