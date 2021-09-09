package astutils

import (
	"fmt"
	"github.com/iancoleman/strcase"
	"github.com/unionj-cloud/go-doudou/pathutils"
)

func ExampleRewriteJsonTag() {
	file := pathutils.Abs("testdata/rewritejsontag.go")
	result, err := RewriteJsonTag(file, true, strcase.ToLowerCamel)
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

func ExampleRewriteJsonTag1() {
	file := pathutils.Abs("testdata/rewritejsontag.go")
	result, err := RewriteJsonTag(file, false, strcase.ToLowerCamel)
	if err != nil {
		panic(err)
	}
	fmt.Println(result)
	// Output:
	//package main
	//
	//type base struct {
	//	Index string `json:"index"`
	//	Type  string `json:"type"`
	//}
	//
	//type struct1 struct {
	//	base
	//	Name       string `json:"name"`
	//	StructType int    `json:"structType" dd:"awesomtag"`
	//	Format     string `dd:"anothertag" json:"format"`
	//	Pos        int    `json:"pos"`
	//}

}
