package codegen

import (
	"bytes"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/unionj-cloud/go-doudou/astutils"
	"github.com/unionj-cloud/go-doudou/pathutils"
	"github.com/unionj-cloud/go-doudou/templateutils"
	"golang.org/x/tools/imports"
	"io/ioutil"
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

	if err = os.MkdirAll(dpath, os.ModePerm); err != nil {
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

		var src, res []byte
		src = []byte(source)

		if res, err = imports.Process(dfile, src, &imports.Options{
			TabWidth:  8,
			TabIndent: true,
			Comments:  true,
			Fragment:  true,
		}); err != nil {
			return errors.Wrap(err, "error")
		}

		if !bytes.Equal(src, res) {
			// On Windows, we need to re-set the permissions from the file. See golang/go#38225.
			var perms os.FileMode
			if fi, err := os.Stat(dfile); err == nil {
				perms = fi.Mode() & os.ModePerm
			}
			err = ioutil.WriteFile(dfile, res, perms)
			if err != nil {
				return err
			}
		}
	} else {
		log.Warnf("file %s already exists", dfile)
	}
	return nil
}
