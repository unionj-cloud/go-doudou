package fileutils

import (
	"github.com/unionj-cloud/go-doudou/toolkit/pathutils"
	"os"
	"testing"
)

func TestCreateDirectory(t *testing.T) {
	dir := pathutils.Abs("testfiles")
	type args struct {
		dir string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "",
			args: args{
				dir: dir,
			},
			wantErr: false,
		},
		{
			name: "",
			args: args{
				dir: "/TestCreateDirectory",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := CreateDirectory(tt.args.dir); (err != nil) != tt.wantErr {
				t.Errorf("CreateDirectory() error = %v, wantErr %v", err, tt.wantErr)
			}
			defer os.RemoveAll(dir)
		})
	}
}
