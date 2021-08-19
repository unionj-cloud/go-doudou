package codegen

import (
	"github.com/unionj-cloud/go-doudou/astutils"
	v3 "github.com/unionj-cloud/go-doudou/openapi/v3"
	"github.com/unionj-cloud/go-doudou/pathutils"
	"os"
	"path/filepath"
	"testing"
)

var testDir string

func init() {
	testDir = pathutils.Abs("testfiles")
}

func TestGenDoc(t *testing.T) {
	dir := testDir + "doc1"
	InitSvc(dir)
	defer os.RemoveAll(dir)
	type args struct {
		dir string
		ic  astutils.InterfaceCollector
	}
	svcfile := filepath.Join(dir, "svc.go")
	ic := astutils.BuildInterfaceCollector(svcfile, ExprStringP)

	tests := []struct {
		name string
		args args
	}{
		{
			name: "1",
			args: args{
				dir,
				ic,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			GenDoc(tt.args.dir, tt.args.ic, 1)
		})
	}
}

func Test_schemasOf(t *testing.T) {
	type args struct {
		vofile string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "",
			args: args{
				vofile: pathutils.Abs("testfiles") + "/vo/vo.go",
			},
			want: 6,
		},
		{
			name: "",
			args: args{
				vofile: pathutils.Abs("testfiles") + "/vo/vo1.go",
			},
			want: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v3.SchemaNames = getSchemaNames(tt.args.vofile)
			if got := schemasOf(tt.args.vofile); len(got) != tt.want {
				t.Errorf("schemasOf() = %v, want %v", len(got), tt.want)
			}
		})
	}
}
