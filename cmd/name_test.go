package cmd_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/unionj-cloud/go-doudou/v2/cmd"
	"github.com/unionj-cloud/toolkit/astutils"
)

func TestNameCmd(t *testing.T) {
	dir := testDir + "/namecmd"
	receiver := NewMockSvc(dir)
	receiver.Init()
	defer os.RemoveAll(dir)
	err := os.Chdir(dir)
	if err != nil {
		t.Fatal(err)
	}
	_, _, err = ExecuteCommandC(cmd.GetRootCmd(), []string{"name", "-f", filepath.Join(dir, "dto", "dto.go"), "-o"}...)
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetImportPath(t *testing.T) {
	dir := testDir + "/importpath"
	receiver := NewMockSvc(dir)
	receiver.Init()
	defer os.RemoveAll(dir)
	err := os.Chdir(dir)
	if err != nil {
		t.Fatal(err)
	}
	type args struct {
		file string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "1",
			args: args{
				file: filepath.Join(dir, "/entity"),
			},
			want: "importpath/entity",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := astutils.GetImportPath(tt.args.file); got != tt.want {
				t.Errorf("GetImportPath() = %v, want %v", got, tt.want)
			}
		})
	}
}
