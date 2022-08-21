package astutils

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/unionj-cloud/go-doudou/toolkit/pathutils"
	"go/ast"
	"go/parser"
	"go/token"
	"regexp"
	"testing"
)

func TestStruct(t *testing.T) {
	file := pathutils.Abs("testdata/vo.go")
	fset := token.NewFileSet()
	root, err := parser.ParseFile(fset, file, nil, parser.ParseComments)
	if err != nil {
		panic(err)
	}
	sc := NewStructCollector(ExprString)
	ast.Walk(sc, root)
	fmt.Println(sc.Structs)
}

func TestInter(t *testing.T) {
	file := pathutils.Abs("testdata/svc.go")
	fset := token.NewFileSet()
	root, err := parser.ParseFile(fset, file, nil, parser.ParseComments)
	if err != nil {
		panic(err)
	}
	sc := NewInterfaceCollector(ExprString)
	ast.Walk(sc, root)
	fmt.Println(sc.Interfaces)
}

func TestStructFuncDecl(t *testing.T) {
	file := pathutils.Abs("testdata/cat.go")
	fset := token.NewFileSet()
	root, err := parser.ParseFile(fset, file, nil, parser.ParseComments)
	if err != nil {
		panic(err)
	}
	sc := NewStructCollector(ExprString)
	ast.Walk(sc, root)
	methods, exists := sc.Methods["Cat"]
	if !exists {
		t.Error("Cat should has methods")
	}
	if len(methods) != 1 {
		t.Error("Cat should has only one method")
	}
}

func TestRegex(t *testing.T) {
	re := regexp.MustCompile(`anonystruct«(.*)»`)
	a := `[]anonystruct«{"Name":"","Fields":[{"Name":"Name","Type":"string","Tag":"","Comments":null,"IsExport":true,"DocName":"Name"},{"Name":"Addr","Type":"anonystruct«{\"Name\":\"\",\"Fields\":[{\"Name\":\"Zip\",\"Type\":\"string\",\"Tag\":\"\",\"Comments\":null,\"IsExport\":true,\"DocName\":\"Zip\"},{\"Name\":\"Block\",\"Type\":\"string\",\"Tag\":\"\",\"Comments\":null,\"IsExport\":true,\"DocName\":\"Block\"},{\"Name\":\"Full\",\"Type\":\"string\",\"Tag\":\"\",\"Comments\":null,\"IsExport\":true,\"DocName\":\"Full\"}],\"Comments\":null,\"Methods\":null,\"IsExport\":false}»","Tag":"","Comments":null,"IsExport":true,"DocName":"Addr"}],"Comments":null,"Methods":null,"IsExport":false}»`
	result := re.FindStringSubmatch(a)
	fmt.Println(result[1])

	j := result[1]
	var structmeta StructMeta
	json.Unmarshal([]byte(j), &structmeta)
}

func TestStructCollector_DocFlatEmbed(t *testing.T) {
	file := pathutils.Abs("testdata/embed.go")
	fset := token.NewFileSet()
	root, err := parser.ParseFile(fset, file, nil, parser.ParseComments)
	if err != nil {
		panic(err)
	}
	sc := NewStructCollector(ExprString)
	ast.Walk(sc, root)
	structs := sc.DocFlatEmbed()
	for _, item := range structs {
		if item.Name == "TestEmbed" {
			var fields []string
			var docFields []string
			for _, fieldMeta := range item.Fields {
				fields = append(fields, fieldMeta.Name)
				docFields = append(docFields, fieldMeta.DocName)
			}
			assert.ElementsMatch(t, fields, []string{"Fields", "Type", "Index"})
			assert.ElementsMatch(t, docFields, []string{"fields", "type", "index"})
		}
	}
}

func TestStructCollector_DocFlatEmbed1(t *testing.T) {
	file := pathutils.Abs("testdata/embed1.go")
	fset := token.NewFileSet()
	root, err := parser.ParseFile(fset, file, nil, parser.ParseComments)
	if err != nil {
		panic(err)
	}
	sc := NewStructCollector(ExprString)
	ast.Walk(sc, root)
	structs := sc.DocFlatEmbed()
	for _, item := range structs {
		if item.Name == "TestEmbed1" {
			var fields []string
			var docFields []string
			for _, fieldMeta := range item.Fields {
				fields = append(fields, fieldMeta.Name)
				docFields = append(docFields, fieldMeta.DocName)
			}
			assert.ElementsMatch(t, fields, []string{"Fields", "TestBase1"})
			assert.ElementsMatch(t, docFields, []string{"fields", "test_base_1"})
		}
	}
}

func TestStructCollector_DocFlatEmbed_ExcludeUnexportedFields(t *testing.T) {
	file := pathutils.Abs("testdata/embed2.go")
	fset := token.NewFileSet()
	root, err := parser.ParseFile(fset, file, nil, parser.ParseComments)
	if err != nil {
		panic(err)
	}
	sc := NewStructCollector(ExprString)
	ast.Walk(sc, root)
	structs := sc.DocFlatEmbed()
	for _, item := range structs {
		if item.Name == "TestEmbed2" {
			var fields []string
			var docFields []string
			for _, fieldMeta := range item.Fields {
				fields = append(fields, fieldMeta.Name)
				docFields = append(docFields, fieldMeta.DocName)
			}
			assert.ElementsMatch(t, fields, []string{"Fields", "Index", "Type"})
			assert.ElementsMatch(t, docFields, []string{"fields", "index", "type"})
		}
	}
}

func TestStructCollector_DocFlatEmbed_ExcludeUnexportedFields2(t *testing.T) {
	file := pathutils.Abs("testdata/embed3.go")
	fset := token.NewFileSet()
	root, err := parser.ParseFile(fset, file, nil, parser.ParseComments)
	if err != nil {
		panic(err)
	}
	sc := NewStructCollector(ExprString)
	ast.Walk(sc, root)
	structs := sc.DocFlatEmbed()
	for _, item := range structs {
		if item.Name == "TestEmbed3" {
			var fields []string
			var docFields []string
			for _, fieldMeta := range item.Fields {
				fields = append(fields, fieldMeta.Name)
				docFields = append(docFields, fieldMeta.DocName)
			}
			assert.ElementsMatch(t, fields, []string{"Fields", "Type"})
			assert.ElementsMatch(t, docFields, []string{"fields", "type"})
		}
	}
}

func TestStructCollector_DocFlatEmbed_ExcludeUnexportedFields3(t *testing.T) {
	file := pathutils.Abs("testdata/embed4.go")
	fset := token.NewFileSet()
	root, err := parser.ParseFile(fset, file, nil, parser.ParseComments)
	if err != nil {
		panic(err)
	}
	sc := NewStructCollector(ExprString)
	ast.Walk(sc, root)
	structs := sc.DocFlatEmbed()
	for _, item := range structs {
		if item.Name == "TestEmbed4" {
			var fields []string
			var docFields []string
			for _, fieldMeta := range item.Fields {
				fields = append(fields, fieldMeta.Name)
				docFields = append(docFields, fieldMeta.DocName)
			}
			assert.ElementsMatch(t, fields, []string{"Fields", "TestBase1"})
			assert.ElementsMatch(t, docFields, []string{"fields", "test_base_1"})
		}
		if item.Name == "TestBase4" {
			var fields []string
			var docFields []string
			for _, fieldMeta := range item.Fields {
				fields = append(fields, fieldMeta.Name)
				docFields = append(docFields, fieldMeta.DocName)
			}
			assert.ElementsMatch(t, fields, []string{"Type"})
			assert.ElementsMatch(t, docFields, []string{"Type"})
		}
	}
}

func TestStructCollector_DocFlatEmbed_ExcludeUnexportedFields4(t *testing.T) {
	file := pathutils.Abs("testdata/embed5.go")
	fset := token.NewFileSet()
	root, err := parser.ParseFile(fset, file, nil, parser.ParseComments)
	if err != nil {
		panic(err)
	}
	sc := NewStructCollector(ExprString)
	ast.Walk(sc, root)
	structs := sc.DocFlatEmbed()
	for _, item := range structs {
		if item.Name == "TestEmbed5" {
			var fields []string
			var docFields []string
			for _, fieldMeta := range item.Fields {
				fields = append(fields, fieldMeta.Name)
				docFields = append(docFields, fieldMeta.DocName)
			}
			assert.ElementsMatch(t, fields, []string{"Fields", "testBase5"})
			assert.ElementsMatch(t, docFields, []string{"fields", "testBase"})
		}
	}
}

func TestStructCollector_Alias(t *testing.T) {
	file := pathutils.Abs("testdata/alias.go")
	fset := token.NewFileSet()
	root, err := parser.ParseFile(fset, file, nil, parser.ParseComments)
	if err != nil {
		panic(err)
	}
	//spew.Dump(root)
	sc := NewStructCollector(ExprString)
	ast.Walk(sc, root)
	structs := sc.DocFlatEmbed()
	for _, item := range structs {
		if item.Name == "TestAlias" {
			fmt.Println(item)
		}
	}
}

func TestStructCollector_Domain(t *testing.T) {
	file := pathutils.Abs("testdata/domain/purchase.go")
	fset := token.NewFileSet()
	root, err := parser.ParseFile(fset, file, nil, parser.ParseComments)
	if err != nil {
		panic(err)
	}
	//spew.Dump(root)
	sc := NewStructCollector(ExprString)
	ast.Walk(sc, root)
	fmt.Println(sc)
}
