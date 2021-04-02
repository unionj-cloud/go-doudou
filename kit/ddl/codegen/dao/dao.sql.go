package dao

import (
	"github.com/iancoleman/strcase"
	"github.com/pkg/errors"
	"github.com/unionj-cloud/go-doudou/kit/ddl/table"
	"github.com/unionj-cloud/go-doudou/kit/pathutils"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

func GenDaoSql(domainpath string, t table.Table) error {
	var (
		err      error
		tplpath  string
		daopath  string
		f        *os.File
		funcMap  map[string]interface{}
		tpl      *template.Template
		iColumns []table.Column
		uColumns []table.Column
	)
	daopath = filepath.Join(filepath.Dir(domainpath), "dao")
	if err = os.MkdirAll(daopath, os.ModePerm); err != nil {
		return errors.Wrap(err, "error")
	}

	daofile := filepath.Join(daopath, strings.ToLower(t.Meta.Name)+"dao.sql")
	if _, err = os.Stat(daofile); os.IsNotExist(err) {
		if f, err = os.Create(daofile); err != nil {
			return errors.Wrap(err, "error")
		}
		defer f.Close()

		tplpath = pathutils.Abs("dao.sql.tmpl")
		funcMap = make(map[string]interface{})
		funcMap["ToSnake"] = strcase.ToSnake
		if tpl, err = template.New("dao.sql.tmpl").Funcs(funcMap).ParseFiles(tplpath); err != nil {
			return errors.Wrap(err, "error")
		}

		for _, co := range t.Columns {
			if !co.AutoSet {
				iColumns = append(iColumns, co)
			}
			if !co.AutoSet && !co.Pk {
				uColumns = append(uColumns, co)
			}
		}

		if err = tpl.Execute(f, struct {
			Schema        string
			DomainName    string
			InsertColumns []table.Column
			UpdateColumns []table.Column
			Pk            string
		}{
			Schema:        os.Getenv("DB_SCHEMA"),
			DomainName:    t.Meta.Name,
			InsertColumns: iColumns,
			UpdateColumns: uColumns,
			Pk:            t.Pk,
		}); err != nil {
			return errors.Wrap(err, "error")
		}
	}
	return nil
}