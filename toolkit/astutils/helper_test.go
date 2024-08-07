package astutils

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetPkgPath(t *testing.T) {
	workDir, _ := filepath.Abs("./testdata")
	oldDir, _ := os.Getwd()
	os.Chdir(workDir)
	defer func() {
		os.Chdir(oldDir)
	}()
	type args struct {
		filePath string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "",
			args: args{
				filePath: filepath.Join(workDir, "demo", "config"),
			},
			want: "testdata/demo/config",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, GetPkgPath(tt.args.filePath), "GetPkgPath(%v)", tt.args.filePath)
		})
	}
}
