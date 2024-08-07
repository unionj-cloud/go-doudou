package astutils

import (
	"fmt"
	"testing"

	"github.com/iancoleman/strcase"
	"github.com/unionj-cloud/go-doudou/v2/toolkit/pathutils"
)

func ExampleRewriteTag() {
	file := pathutils.Abs("testdata/rewritejsontag.go")
	config := RewriteTagConfig{
		File:        file,
		Omitempty:   true,
		ConvertFunc: strcase.ToLowerCamel,
		Form:        false,
	}
	result, err := RewriteTag(config)
	if err != nil {
		panic(err)
	}
	fmt.Println(result)
	// Output:
	//package main
	//
	//type base struct {
	//	Index string `json:"index,omitempty"`
	//	Type  string `json:"type,omitempty"`
	//}
	//
	//type struct1 struct {
	//	base
	//	Name       string `json:"name,omitempty"`
	//	StructType int    `json:"structType,omitempty" dd:"awesomtag"`
	//	Format     string `dd:"anothertag" json:"format,omitempty"`
	//	Pos        int    `json:"pos,omitempty"`
	//}
}

func Test_isExport(t *testing.T) {
	type args struct {
		field string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "",
			args: args{
				field: "unExportField",
			},
			want: false,
		},
		{
			name: "",
			args: args{
				field: "ExportField",
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isExport(tt.args.field); got != tt.want {
				t.Errorf("isExport() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_extractJsonPropName(t *testing.T) {
	type args struct {
		tag string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "",
			args: args{
				tag: `json:"name, omitempty"`,
			},
			want: "name",
		},
		{
			name: "",
			args: args{
				tag: `json:"name"`,
			},
			want: "name",
		},
		{
			name: "",
			args: args{
				tag: `json:"-, omitempty"`,
			},
			want: "-",
		},
		{
			name: "",
			args: args{
				tag: `json:"-"`,
			},
			want: "-",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := extractJsonPropName(tt.args.tag); got != tt.want {
				t.Errorf("extractJsonPropName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func ExampleRewriteTagForm() {
	file := pathutils.Abs("testdata/rewritejsontag.go")
	config := RewriteTagConfig{
		File:        file,
		Omitempty:   true,
		ConvertFunc: strcase.ToLowerCamel,
		Form:        true,
	}
	result, err := RewriteTag(config)
	if err != nil {
		panic(err)
	}
	fmt.Println(result)
	// Output:
	//package main
	//
	//type base struct {
	//	Index string `json:"index,omitempty" form:"index,omitempty"`
	//	Type  string `json:"type,omitempty" form:"type,omitempty"`
	//}
	//
	//type struct1 struct {
	//	base
	//	Name       string `json:"name,omitempty" form:"name,omitempty"`
	//	StructType int    `json:"structType,omitempty" dd:"awesomtag" form:"structType,omitempty"`
	//	Format     string `dd:"anothertag" json:"format,omitempty" form:"format,omitempty"`
	//	Pos        int    `json:"pos,omitempty" form:"pos,omitempty"`
	//}
}
