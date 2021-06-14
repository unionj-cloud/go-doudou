package astutils

import (
	"fmt"
	"github.com/unionj-cloud/go-doudou/pathutils"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
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
	// [{PageFilter [{Name string  [真实姓名，前缀匹配]} {Dept int  [所属部门ID]}] [筛选条件] []} {Order [{Col string  []} {Sort string  []}] [排序条件] []} {PageQuery [{Filter PageFilter  []} {Page Page  []}] [分页筛选条件] []} {PageRet [{Items interface{}  []} {PageNo int  []} {PageSize int  []} {Total int  []} {HasNext bool  []}] [] []}]
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

func TestGetImportPath(t *testing.T) {
	err := os.Chdir("../example")
	if err != nil {
		panic(err)
	}
	defer os.Chdir("../astutils")
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
				file: pathutils.Abs("../example/ddl/domain"),
			},
			want: "example/ddl/domain",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetImportPath(tt.args.file); got != tt.want {
				t.Errorf("GetImportPath() = %v, want %v", got, tt.want)
			}
		})
	}
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
