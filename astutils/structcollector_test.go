package astutils

import (
	"fmt"
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
	// [{PageFilter [{Name string  [真实姓名，前缀匹配]} {Dept int  [所属部门ID]}] [筛选条件] []} {Order [{Col string  []} {Sort string  []}] [排序条件] []} {PageQuery [{Filter PageFilter  []} {Page Page  []}] [分页筛选条件] []} {PageRet [{Items interface{}  []} {PageNo int  []} {PageSize int  []} {Total int  []} {HasNext bool  []}] [] []} {Base [{Index string  []} {Type string  []}] [] []} {QueryCond [{Pair map[string][]interface{}  []} {QueryLogic queryLogic  []} {QueryType queryType  []} {Children []QueryCond  []}] [] []} {Sort [{Field string  []} {Ascending bool  []}] [] []} {Paging [{StartDate string  []} {EndDate string  []} {DateField string  []} {QueryConds []QueryCond  []} {Skip int  []} {Limit int  []} {Sortby []Sort  []}] [] []} {BulkSavePayload [{Base embed:Base  []} {Docs []map[string]interface{}  []}] [] []} {SavePayload [{Base embed:Base  []} {Doc map[string]interface{}  []}] [] []} {BulkDeletePayload [{Base embed:Base  []} {DocIds []string  []}] [] []} {PagePayload [{Base embed:Base  []} {Paging embed:Paging  []}] [] []} {PageResult [{Page int  []} {PageSize int  []} {Total int  []} {Docs []map[string]interface{}  []} {HasNextPage bool  []}] [] []} {StatPayload [{Base embed:Base  []} {Paging embed:Paging  []} {Aggr interface{}  []}] [] []} {RandomPayload [{Base embed:Base  []} {Paging embed:Paging  []}] [] []} {CountPayload [{Base embed:Base  []} {Paging embed:Paging  []}] [] []} {Field [{Name string  []} {Type esFieldType  []} {Format string  []}] [] []} {MappingPayload [{Base embed:Base  []} {Fields []Field  []}] [] []}]
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
	// [{Usersvc [{PageUsers [{ctx context.Context  []} {query PageQuery  []}] [{code int  []} {data PageRet  []} {msg error  []}] [You can define your service methods as your need. Below is an example.]} {GetUser [{ctx context.Context  []} {userId string  []} {photo string  []}] [{code int  []} {data string  []} {msg error  []}] [comment1 comment2]} {SignUp [{ctx context.Context  []} {username string  []} {password int  []} {actived bool  []} {score float64  []}] [{code int  []} {data string  []} {msg error  []}] [comment3]} {UploadAvatar [{pc context.Context  []} {pf []*multipart.FileHeader  []} {ps string  []}] [{ri int  []} {rs string  []} {re error  []}] [comment4]} {DownloadAvatar [{ctx context.Context  []} {userId string  []}] [{rf *os.File  []} {re error  []}] [comment5]}] [用户服务接口 v1版本]}]
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
