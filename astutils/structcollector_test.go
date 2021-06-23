package astutils

import (
	"encoding/json"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"
	"github.com/unionj-cloud/go-doudou/pathutils"
	"go/ast"
	"go/parser"
	"go/token"
	"regexp"
	"testing"
)

func ExampleStruct() {
	file := pathutils.Abs("testfiles/vo.go")
	fset := token.NewFileSet()
	root, err := parser.ParseFile(fset, file, nil, parser.ParseComments)
	if err != nil {
		panic(err)
	}
	sc := NewStructCollector(ExprString)
	ast.Walk(sc, root)
	fmt.Println(sc.Structs)
	// Output:
	// [{PageFilter [{Name string  [真实姓名，前缀匹配] true Name} {Dept int  [所属部门ID] true Dept}] [筛选条件] [] true} {Order [{Col string  [] true Col} {Sort string  [] true Sort}] [排序条件] [] true} {PageQuery [{Filter PageFilter  [] true Filter} {Page Page  [] true Page}] [分页筛选条件] [] true} {PageRet [{Items interface{}  [] true Items} {PageNo int  [] true PageNo} {PageSize int  [] true PageSize} {Total int  [] true Total} {HasNext bool  [] true HasNext}] [] [] true} {Base [{Index string  [] true Index} {Type string  [] true Type}] [] [] true} {QueryCond [{Pair map[string][]interface{}  [] true Pair} {QueryLogic queryLogic  [] true QueryLogic} {QueryType queryType  [] true QueryType} {Children []QueryCond  [] true Children}] [] [] true} {Sort [{Field string  [] true Field} {Ascending bool  [] true Ascending}] [] [] true} {Paging [{StartDate string  [] true StartDate} {EndDate string  [] true EndDate} {DateField string  [] true DateField} {QueryConds []QueryCond  [] true QueryConds} {Skip int  [] true Skip} {Limit int  [] true Limit} {Sortby []Sort  [] true Sortby}] [] [] true} {BulkSavePayload [{Base embed:Base  [] true Base} {Docs []map[string]interface{}  [] true Docs}] [] [] true} {SavePayload [{Base embed:Base  [] true Base} {Doc map[string]interface{}  [] true Doc}] [] [] true} {BulkDeletePayload [{Base embed:Base  [] true Base} {DocIds []string  [] true DocIds}] [] [] true} {PagePayload [{Base embed:Base  [] true Base} {Paging embed:Paging  [] true Paging}] [] [] true} {PageResult [{Page int  [] true Page} {PageSize int  [] true PageSize} {Total int  [] true Total} {Docs []map[string]interface{}  [] true Docs} {HasNextPage bool  [] true HasNextPage}] [] [] true} {StatPayload [{Base embed:Base  [] true Base} {Paging embed:Paging  [] true Paging} {Aggr interface{}  [] true Aggr}] [] [] true} {RandomPayload [{Base embed:Base  [] true Base} {Paging embed:Paging  [] true Paging}] [] [] true} {CountPayload [{Base embed:Base  [] true Base} {Paging embed:Paging  [] true Paging}] [] [] true} {Field [{Name string  [] true Name} {Type esFieldType  [] true Type} {Format string  [] true Format}] [] [] true} {MappingPayload [{Base embed:Base  [] true Base} {Fields []Field  [] true Fields}] [] [] true}]
}

func ExampleInter() {
	file := pathutils.Abs("testfiles/svc.go")
	fset := token.NewFileSet()
	root, err := parser.ParseFile(fset, file, nil, parser.ParseComments)
	if err != nil {
		panic(err)
	}
	sc := NewInterfaceCollector(ExprString)
	ast.Walk(sc, root)
	fmt.Println(sc.Interfaces)
	// Output:
	// [{Usersvc [func PageUsers(ctx context.Context, query PageQuery) (code int, data PageRet, msg error) func GetUser(ctx context.Context, userId string, photo string) (code int, data string, msg error) func SignUp(ctx context.Context, username string, password int, actived bool, score float64) (code int, data string, msg error) func UploadAvatar(pc context.Context, pf []*multipart.FileHeader, ps string) (ri int, rs string, re error) func DownloadAvatar(ctx context.Context, userId string) (rf *os.File, re error)] [用户服务接口 v1版本]}]
}

func TestStructFuncDecl(t *testing.T) {
	file := pathutils.Abs("testfiles/cat.go")
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

func ExampleRegex() {
	re := regexp.MustCompile(`anonystruct«(.*)»`)
	a := `[]anonystruct«{"Name":"","Fields":[{"Name":"Name","Type":"string","Tag":"","Comments":null,"IsExport":true,"DocName":"Name"},{"Name":"Addr","Type":"anonystruct«{\"Name\":\"\",\"Fields\":[{\"Name\":\"Zip\",\"Type\":\"string\",\"Tag\":\"\",\"Comments\":null,\"IsExport\":true,\"DocName\":\"Zip\"},{\"Name\":\"Block\",\"Type\":\"string\",\"Tag\":\"\",\"Comments\":null,\"IsExport\":true,\"DocName\":\"Block\"},{\"Name\":\"Full\",\"Type\":\"string\",\"Tag\":\"\",\"Comments\":null,\"IsExport\":true,\"DocName\":\"Full\"}],\"Comments\":null,\"Methods\":null,\"IsExport\":false}»","Tag":"","Comments":null,"IsExport":true,"DocName":"Addr"}],"Comments":null,"Methods":null,"IsExport":false}»`
	result := re.FindStringSubmatch(a)
	fmt.Println(result[1])

	j := result[1]
	var structmeta StructMeta
	json.Unmarshal([]byte(j), &structmeta)
	// Output:
	// {"Name":"","Fields":[{"Name":"Name","Type":"string","Tag":"","Comments":null,"IsExport":true,"DocName":"Name"},{"Name":"Addr","Type":"anonystruct«{\"Name\":\"\",\"Fields\":[{\"Name\":\"Zip\",\"Type\":\"string\",\"Tag\":\"\",\"Comments\":null,\"IsExport\":true,\"DocName\":\"Zip\"},{\"Name\":\"Block\",\"Type\":\"string\",\"Tag\":\"\",\"Comments\":null,\"IsExport\":true,\"DocName\":\"Block\"},{\"Name\":\"Full\",\"Type\":\"string\",\"Tag\":\"\",\"Comments\":null,\"IsExport\":true,\"DocName\":\"Full\"}],\"Comments\":null,\"Methods\":null,\"IsExport\":false}»","Tag":"","Comments":null,"IsExport":true,"DocName":"Addr"}],"Comments":null,"Methods":null,"IsExport":false}
}

func TestStructCollector_DocFlatEmbed(t *testing.T) {
	file := pathutils.Abs("testfiles/embed.go")
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
	file := pathutils.Abs("testfiles/embed1.go")
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
	file := pathutils.Abs("testfiles/embed2.go")
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
	file := pathutils.Abs("testfiles/embed3.go")
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
	file := pathutils.Abs("testfiles/embed4.go")
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
	file := pathutils.Abs("testfiles/embed5.go")
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
	file := pathutils.Abs("testfiles/alias.go")
	fset := token.NewFileSet()
	root, err := parser.ParseFile(fset, file, nil, parser.ParseComments)
	if err != nil {
		panic(err)
	}
	spew.Dump(root)
	sc := NewStructCollector(ExprString)
	ast.Walk(sc, root)
	structs := sc.DocFlatEmbed()
	for _, item := range structs {
		if item.Name == "TestAlias" {
			fmt.Println(item)
		}
	}
}
