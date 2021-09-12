package codegen

import (
	log "github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"text/template"
)

var basetmpl = `package dao

import (
	"context"
	"github.com/unionj-cloud/go-doudou/ddl/query"
)

type Base interface {
	Insert(ctx context.Context, data interface{}) (int64, error)
	Upsert(ctx context.Context, data interface{}) (int64, error)
	UpsertNoneZero(ctx context.Context, data interface{}) (int64, error)
	DeleteMany(ctx context.Context, where query.Q) (int64, error)
	Update(ctx context.Context, data interface{}) (int64, error)
	UpdateNoneZero(ctx context.Context, data interface{}) (int64, error)
	UpdateMany(ctx context.Context, data interface{}, where query.Q) (int64, error)
	UpdateManyNoneZero(ctx context.Context, data interface{}, where query.Q) (int64, error)
	Get(ctx context.Context, id interface{}) (interface{}, error)
	SelectMany(ctx context.Context, where ...query.Q) (interface{}, error)
	CountMany(ctx context.Context, where ...query.Q) (int, error)
	PageMany(ctx context.Context, page query.Page, where ...query.Q) (query.PageRet, error)
}
`

// GenBaseGo generates Base interface code
// Base interface wraps some common CRUD operations for convenient use
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
	_ = os.MkdirAll(daopath, os.ModePerm)
	basefile := filepath.Join(daopath, "base.go")
	if _, err = os.Stat(basefile); os.IsNotExist(err) {
		f, _ = os.Create(basefile)
		defer f.Close()
		tpl, _ = template.New("base.go.tmpl").Parse(basetmpl)
		_ = tpl.Execute(f, nil)
	} else {
		log.Warnf("file %s already exists", basefile)
	}
	return nil
}
