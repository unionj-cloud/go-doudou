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
	file := pathutils.Abs("testdata/vo.go")
	fset := token.NewFileSet()
	root, err := parser.ParseFile(fset, file, nil, parser.ParseComments)
	if err != nil {
		panic(err)
	}
	sc := NewStructCollector(ExprString)
	ast.Walk(sc, root)
	fmt.Println(sc.Structs)
	// Output:
	// [{PageFilter [{Name string json:"name,omitempty" [真实姓名，前缀匹配] true name} {Dept int json:"dept,omitempty" [所属部门ID] true dept}] [筛选条件] [] true} {Order [{Col string json:"col,omitempty" [] true col} {Sort string  [] true Sort} {Name string  [] true Name} {Banana string  [] true Banana}] [排序条件] [] true} {PageQuery [{Filter PageFilter json:"filter,omitempty" [] true filter} {Page Page json:"page,omitempty" [] true page}] [分页筛选条件] [] true} {PageRet [{Items interface{} json:"items,omitempty" [] true items} {PageNo int json:"pageNo,omitempty" [] true pageNo} {PageSize int json:"pageSize,omitempty" [] true pageSize} {Total int json:"total,omitempty" [] true total} {HasNext bool json:"hasNext,omitempty" [] true hasNext}] [] [] true} {Base [{Index string json:"index,omitempty" [] true index} {Type string json:"type,omitempty" [] true type}] [] [] true} {QueryCond [{Pair map[string][]interface{} json:"pair,omitempty" [] true pair} {QueryLogic queryLogic json:"queryLogic,omitempty" [] true queryLogic} {QueryType queryType json:"queryType,omitempty" [] true queryType} {Children []QueryCond json:"children,omitempty" [] true children}] [] [] true} {Sort [{Field string json:"field,omitempty" [] true field} {Ascending bool json:"ascending,omitempty" [] true ascending}] [] [] true} {Paging [{StartDate string json:"startDate,omitempty" [] true startDate} {EndDate string json:"endDate,omitempty" [] true endDate} {DateField string json:"dateField,omitempty" [] true dateField} {QueryConds []QueryCond json:"queryConds,omitempty" [] true queryConds} {Skip int json:"skip,omitempty" [] true skip} {Limit int json:"limit,omitempty" [] true limit} {Sortby []Sort json:"sortby,omitempty" [] true sortby}] [] [] true} {BulkSavePayload [{Base embed:Base  [] true Base} {Docs []map[string]interface{} json:"docs,omitempty" [] true docs}] [] [] true} {SavePayload [{Base embed:Base  [] true Base} {Doc map[string]interface{} json:"doc,omitempty" [] true doc}] [] [] true} {BulkDeletePayload [{Base embed:Base  [] true Base} {DocIds []string json:"docIds,omitempty" [] true docIds}] [] [] true} {PagePayload [{Base embed:Base  [] true Base} {Paging embed:Paging  [] true Paging}] [] [] true} {PageResult [{Page int json:"page,omitempty" [] true page} {PageSize int json:"pageSize,omitempty" [] true pageSize} {Total int json:"total,omitempty" [] true total} {Docs []map[string]interface{} json:"docs,omitempty" [] true docs} {HasNextPage bool json:"hasNextPage,omitempty" [] true hasNextPage}] [] [] true} {StatPayload [{Base embed:Base  [] true Base} {Paging embed:Paging  [] true Paging} {Aggr interface{} json:"aggr,omitempty" [] true aggr}] [] [] true} {RandomPayload [{Base embed:Base  [] true Base} {Paging embed:Paging  [] true Paging}] [] [] true} {CountPayload [{Base embed:Base  [] true Base} {Paging embed:Paging  [] true Paging}] [] [] true} {Field [{Name string json:"name,omitempty" [] true name} {Type esFieldType json:"type,omitempty" [] true type} {Format string json:"format,omitempty" [] true format}] [] [] true} {MappingPayload [{Base embed:Base  [] true Base} {Fields []Field json:"fields,omitempty" [] true fields}] [] [] true}]
}

func ExampleInter() {
	file := pathutils.Abs("testdata/svc.go")
	fset := token.NewFileSet()
	root, err := parser.ParseFile(fset, file, nil, parser.ParseComments)
	if err != nil {
		panic(err)
	}
	sc := NewInterfaceCollector(ExprString)
	ast.Walk(sc, root)
	fmt.Println(sc.Interfaces)
	// Output:
	// [{Usersvc [func PageUsers(ctx context.Context, query PageQuery) (code int, data PageRet, msg error) func GetUser(ctx context.Context, userId string, photo string) (code int, data string, msg error) func SignUp(ctx context.Context, username string, password int, actived bool, score float64) (code int, data string, msg error) func UploadAvatar(pc context.Context, pf []*multipart.FileHeader, ps string) (ri int, rs string, re error) func DownloadAvatar(ctx context.Context, userId string, userType string, userNo string) (a string, b string) func BulkSaveOrUpdate(pc context.Context, pi int) (re error)] [用户服务接口 v1版本]}]
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

func TestStructCollector_Domain(t *testing.T) {
	file := pathutils.Abs("testdata/domain/purchase.go")
	fset := token.NewFileSet()
	root, err := parser.ParseFile(fset, file, nil, parser.ParseComments)
	if err != nil {
		panic(err)
	}
	spew.Dump(root)
	sc := NewStructCollector(ExprString)
	ast.Walk(sc, root)
	fmt.Println(sc)
}
