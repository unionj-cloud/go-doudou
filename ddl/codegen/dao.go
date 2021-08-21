package codegen

import (
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/unionj-cloud/go-doudou/ddl/table"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

var daotmpl = `package dao

type {{.DomainName}}Dao interface {
	Base
}`

func GenDaoGo(domainpath string, t table.Table, folder ...string) error {
	var (
		err     error
		daopath string
		f       *os.File
		tpl     *template.Template
		df      string
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

		if tpl, err = template.New("dao.go.tmpl").Parse(daotmpl); err != nil {
			return errors.Wrap(err, "error")
		}
		if err = tpl.Execute(f, struct {
			DomainName string
		}{
			DomainName: t.Meta.Name,
		}); err != nil {
			return errors.Wrap(err, "error")
		}
	} else {
		log.Warnf("file %s already exists", daofile)
	}
	return nil
}
