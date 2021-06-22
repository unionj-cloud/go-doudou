package astutils

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/unionj-cloud/go-doudou/pathutils"
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
	"testing"
)

func ExampleStruct() {
	file := pathutils.Abs("testfiles/vo.go")
	fset := token.NewFileSet()
	root, err := parser.ParseFile(fset, file, nil, parser.ParseComments)
	if err != nil {
		panic(err)
	}
	sc := NewStructCollector()
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
	var sc InterfaceCollector
	ast.Walk(&sc, root)
	fmt.Println(sc.Interfaces)
	// Output:
	// [{Usersvc [{PageUsers [{ctx context.Context  [] false } {query PageQuery  [] false }] [{code int  [] false } {data PageRet  [] false } {msg error  [] false }] [You can define your service methods as your need. Below is an example.]} {GetUser [{ctx context.Context  [] false } {userId string  [] false } {photo string  [] false }] [{code int  [] false } {data string  [] false } {msg error  [] false }] [comment1 comment2]} {SignUp [{ctx context.Context  [] false } {username string  [] false } {password int  [] false } {actived bool  [] false } {score float64  [] false }] [{code int  [] false } {data string  [] false } {msg error  [] false }] [comment3]} {UploadAvatar [{pc context.Context  [] false } {pf []*multipart.FileHeader  [] false } {ps string  [] false }] [{ri int  [] false } {rs string  [] false } {re error  [] false }] [comment4]} {DownloadAvatar [{ctx context.Context  [] false } {userId string  [] false }] [{rf *os.File  [] false } {re error  [] false }] [comment5]}] [用户服务接口 v1版本]}]
}

func TestStructFuncDecl(t *testing.T) {
	file := pathutils.Abs("testfiles/cat.go")
	fset := token.NewFileSet()
	root, err := parser.ParseFile(fset, file, nil, parser.ParseComments)
	if err != nil {
		panic(err)
	}
	sc := NewStructCollector()
	ast.Walk(sc, root)
	methods, exists := sc.Methods["Cat"]
	if !exists {
		t.Error("Cat should has methods")
	}
	if len(methods) != 1 {
		t.Error("Cat should has only one method")
	}
}

func TestNewTableFromStruct(t *testing.T) {
	var files []string
	var err error
	testDir := pathutils.Abs("testfiles/domain")
	err = filepath.Walk(testDir, Visit(&files))
	if err != nil {
		panic(err)
	}
	var sc StructCollector
	for _, file := range files {
		fset := token.NewFileSet()
		root, err := parser.ParseFile(fset, file, nil, parser.ParseComments)
		if err != nil {
			panic(err)
		}
		ast.Walk(&sc, root)
	}
	flattened := sc.FlatEmbed()

	for _, sm := range flattened {
		if len(sm.Fields) != 10 {
			t.Errorf("want 10, got %d\n", len(sm.Fields))
		}
	}
}

func TestStructCollector_DocFlatEmbed(t *testing.T) {
	file := pathutils.Abs("testfiles/embed.go")
	fset := token.NewFileSet()
	root, err := parser.ParseFile(fset, file, nil, parser.ParseComments)
	if err != nil {
		panic(err)
	}
	sc := NewStructCollector()
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
	sc := NewStructCollector()
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
	sc := NewStructCollector()
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
	sc := NewStructCollector()
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
	sc := NewStructCollector()
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
	sc := NewStructCollector()
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
