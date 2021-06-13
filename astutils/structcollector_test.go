package astutils

import (
	"fmt"
	"github.com/unionj-cloud/go-doudou/pathutils"
	"go/ast"
	"go/parser"
	"go/token"
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
