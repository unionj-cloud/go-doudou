package codegen

import (
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/unionj-cloud/go-doudou/v2/toolkit/astutils"
	v3helper "github.com/unionj-cloud/go-doudou/v2/toolkit/openapi/v3"
	"github.com/unionj-cloud/go-doudou/v2/toolkit/pathutils"
	"os"
	"path/filepath"
	"testing"
)

var testDir string

func init() {
	testDir = pathutils.Abs("testdata")
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

func TestGenDocUploadFile(t *testing.T) {
	type args struct {
		dir string
		ic  astutils.InterfaceCollector
	}
	svcfile := filepath.Join(testDir, "svc.go")
	ic := astutils.BuildInterfaceCollector(svcfile, ExprStringP)
	ParseDto(testDir, "vo")
	tests := []struct {
		name string
		args args
	}{
		{
			name: "1",
			args: args{
				testDir,
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
				vofile: pathutils.Abs("testdata") + "/vo/vo.go",
			},
			want: 6,
		},
		{
			name: "",
			args: args{
				vofile: pathutils.Abs("testdata") + "/vo/vo1.go",
			},
			want: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v3helper.SchemaNames = getSchemaNames(tt.args.vofile)
			if got := schemasOf(tt.args.vofile); len(got) != tt.want {
				t.Errorf("schemasOf() = %v, want %v", len(got), tt.want)
			}
		})
	}
}

func TestParseVo(t *testing.T) {
	Convey("Test ParseDto", t, func() {
		So(func() {
			ParseDto(testDir, "vo")
		}, ShouldNotPanic)
		So(len(v3helper.Schemas), ShouldNotBeZeroValue)
	})
}

func TestParseVo_Decimal(t *testing.T) {
	Convey("Test Parse decimal.Decimal type", t, func() {
		ParseDto(testDir, "dto")
		schemas := v3helper.Schemas
		laptopSchema := schemas["Laptop"]
		priceSchema := laptopSchema.Properties["Price"]
		So(priceSchema, ShouldResemble, v3helper.Decimal)
	})
}

func Test_pathsOf(t *testing.T) {
	type args struct {
		dir string
		ic  astutils.InterfaceCollector
	}
	svcfile := filepath.Join(testDir, "svc.go")
	ic := astutils.BuildInterfaceCollector(svcfile, ExprStringP)
	ParseDto(testDir, "vo")
	tests := []struct {
		name string
		args args
	}{
		{
			name: "1",
			args: args{
				testDir,
				ic,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			paths := pathsOf(tt.args.ic, 0)
			fmt.Println(paths)
		})
	}
}
