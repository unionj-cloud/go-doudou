package codegen

import (
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/unionj-cloud/go-doudou/pathutils"
	"os"
	"path/filepath"
	"text/template"
)

func GenBaseGo(domainpath string, folder ...string) error {
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
	if err = os.MkdirAll(daopath, 0644); err != nil {
		return errors.Wrap(err, "error")
	}

	basefile := filepath.Join(daopath, "base.go")
	if _, err = os.Stat(basefile); os.IsNotExist(err) {
		if f, err = os.Create(basefile); err != nil {
			return errors.Wrap(err, "error")
		}
		defer f.Close()
		if tpl, err = template.New("base.go.tmpl").ParseFiles(pathutils.Abs("base.go.tmpl")); err != nil {
			return errors.Wrap(err, "error")
		}
		if err = tpl.Execute(f, nil); err != nil {
			return errors.Wrap(err, "error")
		}
	} else {
		log.Warnf("file %s already exists", basefile)
	}
	return nil
}
