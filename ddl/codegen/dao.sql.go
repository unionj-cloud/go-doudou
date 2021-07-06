package codegen

import (
	"github.com/iancoleman/strcase"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/unionj-cloud/go-doudou/ddl/table"
	"github.com/unionj-cloud/go-doudou/pathutils"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

func GenDaoSql(domainpath string, t table.Table, folder ...string) error {
	var (
		err      error
		tplpath  string
		daopath  string
		f        *os.File
		funcMap  map[string]interface{}
		tpl      *template.Template
		iColumns []table.Column
		uColumns []table.Column
		df       string
	)
	df = "dao"
	if len(folder) > 0 {
		df = folder[0]
	}
	daopath = filepath.Join(filepath.Dir(domainpath), df)
	if err = os.MkdirAll(daopath, 0644); err != nil {
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

		var pkColumn table.Column
		for _, co := range t.Columns {
			if co.Pk {
				pkColumn = co
				break
			}
		}

		if err = tpl.Execute(f, struct {
			Schema        string
			TableName     string
			DomainName    string
			InsertColumns []table.Column
			UpdateColumns []table.Column
			Pk            table.Column
		}{
			Schema:        os.Getenv("DB_SCHEMA"),
			TableName:     t.Name,
			DomainName:    t.Meta.Name,
			InsertColumns: iColumns,
			UpdateColumns: uColumns,
			Pk:            pkColumn,
		}); err != nil {
			return errors.Wrap(err, "error")
		}
	} else {
		log.Warnf("file %s already exists", daofile)
	}
	return nil
}
