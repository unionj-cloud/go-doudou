package cmd

import (
	"github.com/unionj-cloud/go-doudou/astutils"
	"github.com/unionj-cloud/go-doudou/svc"
	"os"
	"testing"
)

func TestNameCmd(t *testing.T) {
	dir := testDir + "namecmd"
	receiver := svc.Svc{
		Dir: dir,
	}
	receiver.Init()
	defer os.RemoveAll(dir)
	err := os.Chdir(dir)
	if err != nil {
		t.Fatal(err)
	}
	// go-doudou name -f /Users/wubin1989/workspace/chengdutreeyee/team3-cloud-analyse/vo/vo.go -o
	_, _, err = ExecuteCommandC(rootCmd, []string{"name", "-f", dir + "/vo/vo.go", "-o"}...)
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetImportPath(t *testing.T) {
	dir := testDir + "importpath"
	receiver := svc.Svc{
		Dir: dir,
	}
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
				file: dir + "/domain",
			},
			want: "testfilesimportpath/domain",
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
