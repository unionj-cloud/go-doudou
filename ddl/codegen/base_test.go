package codegen

import (
	"github.com/unionj-cloud/go-doudou/pathutils"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestGenBaseGo(t *testing.T) {
	dir := pathutils.Abs("../testdata")
	type args struct {
		domainpath string
		folder     []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			args: args{
				domainpath: dir + "/domain",
				folder:     nil,
			},
			wantErr: false,
		},
		{
			args: args{
				domainpath: dir + "/domain",
				folder:     []string{"testdao"},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := GenBaseGo(tt.args.domainpath, tt.args.folder...); (err != nil) != tt.wantErr {
				t.Errorf("GenBaseGo() error = %v, wantErr %v", err, tt.wantErr)
			}
			defer func() {
				if len(tt.args.folder) > 0 {
					os.RemoveAll(filepath.Join(dir, tt.args.folder[0]))
				} else {
					os.RemoveAll(filepath.Join(dir, "dao"))
				}
			}()
			expect := `package dao

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
			var basefile string
			if len(tt.args.folder) > 0 {
				basefile = filepath.Join(dir, tt.args.folder[0], "base.go")
			} else {
				basefile = filepath.Join(dir, "dao", "base.go")
			}
			f, err := os.Open(basefile)
			if err != nil {
				t.Fatal(err)
			}
			content, err := ioutil.ReadAll(f)
			if err != nil {
				t.Fatal(err)
			}
			if string(content) != expect {
				t.Errorf("want %s, got %s\n", expect, string(content))
			}
		})
	}
}
