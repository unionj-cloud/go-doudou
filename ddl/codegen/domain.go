package codegen

import (
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/unionj-cloud/go-doudou/astutils"
	"github.com/unionj-cloud/go-doudou/pathutils"
	"github.com/unionj-cloud/go-doudou/templateutils"
	"os"
	"path/filepath"
	"strings"
)

func GenDomainGo(dpath string, domain astutils.StructMeta) error {
	var (
		err     error
		tplpath string
		f       *os.File
	)

	if err = os.MkdirAll(dpath, 0644); err != nil {
		return errors.Wrap(err, "error")
	}

	dfile := filepath.Join(dpath, strings.ToLower(domain.Name)+".go")
	if _, err = os.Stat(dfile); os.IsNotExist(err) {
		if f, err = os.Create(dfile); err != nil {
			return errors.Wrap(err, "error")
		}
		defer f.Close()

		tplpath = pathutils.Abs("domain.go.tmpl")
		var source string
		if source, err = templateutils.String(tplpath, domain); err != nil {
			return errors.Wrap(err, "error")
		}

		astutils.FixImport([]byte(source), dfile)
	} else {
		log.Warnf("file %s already exists", dfile)
	}
	return nil
}
