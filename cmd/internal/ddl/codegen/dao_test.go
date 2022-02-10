package codegen

import (
	"github.com/sirupsen/logrus"
	"github.com/unionj-cloud/go-doudou/cmd/internal/astutils"
	"github.com/unionj-cloud/go-doudou/cmd/internal/ddl/ddlast"
	"github.com/unionj-cloud/go-doudou/cmd/internal/ddl/table"
	"github.com/unionj-cloud/go-doudou/toolkit/pathutils"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"testing"
)

func TestGenDaoGo(t *testing.T) {
	domain := "../testdata/domain"

	sc := astutils.NewStructCollector(astutils.ExprString)

	usergo := pathutils.Abs(domain + "/user.go")
	fset := token.NewFileSet()
	root, err := parser.ParseFile(fset, usergo, nil, parser.ParseComments)
	if err != nil {
		logrus.Panicln(err)
	}
	ast.Walk(sc, root)

	basego := pathutils.Abs(domain + "/base.go")
	fset = token.NewFileSet()
	root, err = parser.ParseFile(fset, basego, nil, parser.ParseComments)
	if err != nil {
		logrus.Panicln(err)
	}
	ast.Walk(sc, root)

	var tables []table.Table
	flattened := ddlast.FlatEmbed(sc.Structs)
	for _, sm := range flattened {
		tables = append(tables, table.NewTableFromStruct(sm, ""))
	}
	type args struct {
		domainpath string
		t          table.Table
		folder     []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "1",
			args: args{
				domainpath: pathutils.Abs(domain),
				t:          tables[0],
				folder:     nil,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := GenDaoGo(tt.args.domainpath, tt.args.t, tt.args.folder...); (err != nil) != tt.wantErr {
				t.Errorf("GenDaoGo() error = %v, wantErr %v", err, tt.wantErr)
			}
			defer os.RemoveAll(pathutils.Abs("../testdata/dao"))
			expect := `package dao

type UserDao interface {
	Base
}`
			daofile := pathutils.Abs("../testdata/dao/userdao.go")
			f, err := os.Open(daofile)
			if err != nil {
				t.Fatal(err)
			}
			content, err := ioutil.ReadAll(f)
			if err != nil {
				t.Fatal(err)
			}
			if string(content) != expect {
				t.Errorf("want %s, got %s\n", expect, string(content))
			}
		})
	}
}
