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
	"github.com/unionj-cloud/go-doudou/toolkit/sqlext/query"
	"github.com/unionj-cloud/go-doudou/toolkit/sqlext/wrapper"
	"github.com/unionj-cloud/go-doudou/toolkit/reflectutils"
	"github.com/unionj-cloud/go-doudou/toolkit/stringutils"
	"github.com/unionj-cloud/go-doudou/toolkit/templateutils"
	"strings"
	"math"
)

type {{.DomainName}}DaoImpl struct {
	db wrapper.Querier
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
	)
	if statement, err = templateutils.BlockMysql("{{.DomainName | ToLower}}dao.sql", {{.DomainName | ToLower}}daosql, "Insert{{.DomainName}}", nil); err != nil {
		return 0, errors.Wrap(err, "")
	}
	if result, err = receiver.db.NamedExecContext(ctx, statement, data); err != nil {
		return 0, errors.Wrap(err, "")
	}
	{{- if .PkCol.Autoincrement }}
	if lastInsertID, err = result.LastInsertId(); err != nil {
		return 0, errors.Wrap(err, "")
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
	return result.RowsAffected()
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
	)
	if statement, err = templateutils.BlockMysql("{{.DomainName | ToLower}}dao.sql", {{.DomainName | ToLower}}daosql, "Upsert{{.DomainName}}", nil); err != nil {
		return 0, errors.Wrap(err, "")
	}
	if result, err = receiver.db.NamedExecContext(ctx, statement, data); err != nil {
		return 0, errors.Wrap(err, "")
	}
	{{- if .PkCol.Autoincrement }}
	if lastInsertID, err = result.LastInsertId(); err != nil {
		return 0, errors.Wrap(err, "")
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
	return result.RowsAffected()
}

func (receiver {{.DomainName}}DaoImpl) UpsertNoneZero(ctx context.Context, data interface{}) (int64, error) {
	var (
		statement    string
		err          error
		result       sql.Result
		{{- if .PkCol.Autoincrement }}
		lastInsertID int64
		{{- end }}
	)
	value := reflectutils.ValueOf(data).Interface()
	if _, ok := value.(domain.{{.DomainName}}); !ok {
		return 0, errors.New("underlying type of data should be domain.{{.DomainName}}")
	}
	if statement, err = templateutils.BlockMysql("{{.DomainName | ToLower}}dao.sql", {{.DomainName | ToLower}}daosql, "Upsert{{.DomainName}}NoneZero", data); err != nil {
		return 0, errors.Wrap(err, "")
	}
	if result, err = receiver.db.NamedExecContext(ctx, statement, data); err != nil {
		return 0, errors.Wrap(err, "")
	}
	{{- if .PkCol.Autoincrement }}
	if lastInsertID, err = result.LastInsertId(); err != nil {
		return 0, errors.Wrap(err, "")
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
	return result.RowsAffected()
}

func (receiver {{.DomainName}}DaoImpl) DeleteMany(ctx context.Context, where query.Q) (int64, error) {
	var (
		err    error
		result sql.Result
		w      string
		args   []interface{}
	)
	w, args = where.Sql()
	if result, err = receiver.db.ExecContext(ctx, receiver.db.Rebind(fmt.Sprintf("delete from {{.TableName}} where %s;", w)), args...); err != nil {
		return 0, errors.Wrap(err, "")
	}
	return result.RowsAffected()
}

func (receiver {{.DomainName}}DaoImpl) Update(ctx context.Context, data interface{}) (int64, error) {
	var (
		statement string
		err       error
		result    sql.Result
	)
	if statement, err = templateutils.BlockMysql("{{.DomainName | ToLower}}dao.sql", {{.DomainName | ToLower}}daosql, "Update{{.DomainName}}", nil); err != nil {
		return 0, errors.Wrap(err, "")
	}
	if result, err = receiver.db.NamedExecContext(ctx, statement, data); err != nil {
		return 0, errors.Wrap(err, "")
	}
	return result.RowsAffected()
}

func (receiver {{.DomainName}}DaoImpl) UpdateNoneZero(ctx context.Context, data interface{}) (int64, error) {
	var (
		statement string
		err       error
		result    sql.Result
	)
	value := reflectutils.ValueOf(data).Interface()
	if _, ok := value.(domain.{{.DomainName}}); !ok {
		return 0, errors.New("underlying type of data should be domain.{{.DomainName}}")
	}
	if statement, err = templateutils.BlockMysql("{{.DomainName | ToLower}}dao.sql", {{.DomainName | ToLower}}daosql, "Update{{.DomainName}}NoneZero", data); err != nil {
		return 0, errors.Wrap(err, "")
	}
	if result, err = receiver.db.NamedExecContext(ctx, statement, data); err != nil {
		return 0, errors.Wrap(err, "")
	}
	return result.RowsAffected()
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
	)
	if statement, err = templateutils.BlockMysql("{{.DomainName | ToLower}}dao.sql", {{.DomainName | ToLower}}daosql, "Update{{.DomainName}}s", nil); err != nil {
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

func (receiver {{.DomainName}}DaoImpl) UpdateManyNoneZero(ctx context.Context, data interface{}, where query.Q) (int64, error) {
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
	if _, ok := value.(domain.{{.DomainName}}); !ok {
		return 0, errors.New("underlying type of data should be domain.{{.DomainName}}")
	}
	if statement, err = templateutils.BlockMysql("{{.DomainName | ToLower}}dao.sql", {{.DomainName | ToLower}}daosql, "Update{{.DomainName}}sNoneZero", data); err != nil {
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

func (receiver {{.DomainName}}DaoImpl) Get(ctx context.Context, id interface{}) (interface{}, error) {
	var (
		statement string
		err       error
		{{.DomainName | ToLower}}      domain.{{.DomainName}}
	)
	if statement, err = templateutils.BlockMysql("{{.DomainName | ToLower}}dao.sql", {{.DomainName | ToLower}}daosql, "Get{{.DomainName}}", nil); err != nil {
		return domain.{{.DomainName}}{}, errors.Wrap(err, "")
	}
	if err = receiver.db.GetContext(ctx, &{{.DomainName | ToLower}}, receiver.db.Rebind(statement), id); err != nil {
		return domain.{{.DomainName}}{}, errors.Wrap(err, "")
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
    statements = append(statements, "select * from {{.TableName}}")
    if len(where) > 0 {
        statements = append(statements, "where")
        for _, item := range where {
			q, wargs := item.Sql()
			statements = append(statements, q)
			args = append(args, wargs...)
		}
    }
	if err = receiver.db.SelectContext(ctx, &{{.DomainName | ToLower}}s, strings.Join(statements, " "), args...); err != nil {
		return nil, errors.Wrap(err, "")
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
	statements = append(statements, "select count(1) from {{.TableName}}")
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

func (receiver {{.DomainName}}DaoImpl) PageMany(ctx context.Context, page query.Page, where ...query.Q) (query.PageRet, error) {
	var (
		statements []string
		err       error
		{{.DomainName | ToLower}}s     []domain.{{.DomainName}}
		total     int
		args       []interface{}
	)
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
	q := strings.Join(statements, " ")
	if err = receiver.db.SelectContext(ctx, &{{.DomainName | ToLower}}s, q, args...); err != nil {
		return query.PageRet{}, errors.Wrap(err, "")
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
	if err = receiver.db.GetContext(ctx, &total, strings.Join(statements, " "), args...); err != nil {
		return query.PageRet{}, errors.Wrap(err, "")
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