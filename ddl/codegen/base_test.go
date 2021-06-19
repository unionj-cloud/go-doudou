package codegen

import (
	"github.com/unionj-cloud/go-doudou/pathutils"
	"io/ioutil"
	"os"
	"testing"
)

func TestGenBaseGo(t *testing.T) {
	testDir := pathutils.Abs("testfiles")
	dir := testDir + "basego"
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
			name: "1",
			args: args{
				domainpath: dir + "/domain",
				folder:     nil,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := GenBaseGo(tt.args.domainpath, tt.args.folder...); (err != nil) != tt.wantErr {
				t.Errorf("GenBaseGo() error = %v, wantErr %v", err, tt.wantErr)
			}
			defer os.RemoveAll(dir)
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
			basefile := dir + "/dao/base.go"
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
