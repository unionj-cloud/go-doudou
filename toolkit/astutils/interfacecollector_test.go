package astutils

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/unionj-cloud/go-doudou/v2/toolkit/pathutils"
)

func TestInterfaceCollector(t *testing.T) {
	file := pathutils.Abs("testdata/svc.go")
	fset := token.NewFileSet()
	root, err := parser.ParseFile(fset, file, nil, parser.ParseComments)
	if err != nil {
		panic(err)
	}
	//spew.Dump(root)
	sc := NewInterfaceCollector(ExprString)
	sc.cmap = ast.NewCommentMap(fset, root, root.Comments)
	ast.Walk(sc, root)
	fmt.Println(sc)
}

func TestBuildInterfaceCollector(t *testing.T) {
	file := pathutils.Abs("testdata/svc.go")
	ic := BuildInterfaceCollector(file, ExprString)
	assert.NotNil(t, ic)
}

func Test_pattern(t *testing.T) {
	type args struct {
		method string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "1",
			args: args{
				method: "GetBooks",
			},
			want: "books",
		},
		{
			name: "2",
			args: args{
				method: "PageUsers",
			},
			want: "page/users",
		},
		{
			name: "",
			args: args{
				method: "Get",
			},
			want: "",
		},
		{
			name: "",
			args: args{
				method: "GetShelves_ShelfBooks_Book",
			},
			want: "shelves/:shelf/books/:book",
		},
		{
			name: "",
			args: args{
				method: "Goodfood_BigappleBooks_Mybird",
			},
			want: "goodfood/:bigapple/books/:mybird",
		},
		{
			name: "",
			args: args{
				method: "ApiV1Query_range",
			},
			want: "api/v1/query_range",
		},
		{
			name: "",
			args: args{
				method: "GetQuery_range",
			},
			want: "query_range",
		},
		{
			name: "",
			args: args{
				method: "GetQuery",
			},
			want: "query",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, endpoint := Pattern(tt.args.method)
			assert.Equalf(t, tt.want, endpoint, "pattern(%v)", tt.args.method)
		})
	}
}
