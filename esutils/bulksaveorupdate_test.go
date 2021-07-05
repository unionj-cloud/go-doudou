package esutils

import (
	"context"
	"github.com/unionj-cloud/go-doudou/constants"
	"github.com/unionj-cloud/go-doudou/test"
	"os"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	var terminator func()
	terminator, esHost, esPort = test.PrepareTestEnvironment()
	code := m.Run()
	terminator()
	os.Exit(code)
}

func TestBulkSaveOrUpdate(t *testing.T) {
	es := setupSubTest("test_bulksaveorupdate")
	type args struct {
		docs []interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "1",
			args: args{
				docs: []interface{}{
					map[string]interface{}{
						"createAt": time.Now().In(constants.Loc).Format(constants.FORMATES),
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := es.BulkSaveOrUpdate(context.Background(), tt.args.docs); (err != nil) != tt.wantErr {
				t.Errorf("BulkSaveOrUpdate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_getId(t *testing.T) {
	type args struct {
		doc interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "",
			args: args{
				doc: map[string]interface{}{
					"id":       "id1",
					"createAt": time.Now().In(constants.Loc).Format(constants.FORMATES),
				},
			},
			want:    "id1",
			wantErr: false,
		},
		{
			name: "",
			args: args{
				doc: struct {
					Id       string
					CreateAt string
				}{
					Id:       "id2",
					CreateAt: time.Now().In(constants.Loc).Format(constants.FORMATES),
				},
			},
			want:    "id2",
			wantErr: false,
		},
		{
			name: "",
			args: args{
				doc: map[string]interface{}{
					"createAt": time.Now().In(constants.Loc).Format(constants.FORMATES),
				},
			},
			want:    "",
			wantErr: false,
		},
		{
			name: "",
			args: args{
				doc: struct {
					CreateAt string
				}{
					CreateAt: time.Now().In(constants.Loc).Format(constants.FORMATES),
				},
			},
			want:    "",
			wantErr: false,
		},
		{
			name: "",
			args: args{
				doc: "random string",
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getId(tt.args.doc)
			if (err != nil) != tt.wantErr {
				t.Errorf("getId() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("getId() got = %v, want %v", got, tt.want)
			}
		})
	}
}
