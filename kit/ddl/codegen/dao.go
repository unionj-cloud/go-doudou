package codegen

import (
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/unionj-cloud/go-doudou/kit/astutils"
	"github.com/unionj-cloud/go-doudou/kit/ddl/table"
	"github.com/unionj-cloud/go-doudou/kit/pathutils"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

func GenDaoGo(domainpath string, t table.Table, folder ...string) error {
	var (
		err      error
		dpkg     string
		tplpath  string
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
	if err = os.MkdirAll(daopath, os.ModePerm); err != nil {
		return errors.Wrap(err, "error")
	}

	daofile := filepath.Join(daopath, strings.ToLower(t.Meta.Name)+"dao.go")
	if _, err = os.Stat(daofile); os.IsNotExist(err) {
		if f, err = os.Create(daofile); err != nil {
			return errors.Wrap(err, "error")
		}
		defer f.Close()

		dpkg = astutils.GetImportPath(domainpath)
		tplpath = pathutils.Abs("dao.go.tmpl")
		funcMap = make(map[string]interface{})
		funcMap["ToLower"] = strings.ToLower
		if tpl, err = template.New("dao.go.tmpl").Funcs(funcMap).ParseFiles(tplpath); err != nil {
			return errors.Wrap(err, "error")
		}

		for _, column := range t.Columns {
			if column.Pk {
				pkColumn = column
				break
			}
		}
		if err = tpl.Execute(f, struct {
			DomainPackage string
			DomainName    string
			PkField       astutils.FieldMeta
		}{
			DomainPackage: dpkg,
			DomainName:    t.Meta.Name,
			PkField:       pkColumn.Meta,
		}); err != nil {
			return errors.Wrap(err, "error")
		}
	} else {
		log.Warnf("file %s already exists", daofile)
	}
	return nil
}
