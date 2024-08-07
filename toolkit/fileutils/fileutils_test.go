package fileutils

import (
	"testing"

	"github.com/unionj-cloud/go-doudou/v2/toolkit/pathutils"
)

func TestCreateDirectory(t *testing.T) {
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
				dir: pathutils.Abs("testfiles"),
			},
			// it should have error because testfiles has already existed as a file not a directory
			wantErr: true,
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
			//defer os.RemoveAll(dir)
		})
	}
}
