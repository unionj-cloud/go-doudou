package astutils

import (
	"github.com/stretchr/testify/require"
	"github.com/unionj-cloud/go-doudou/toolkit/pathutils"
	"go/ast"
	"go/parser"
	"go/token"
	"testing"
)

func TestEnum(t *testing.T) {
	file := pathutils.Abs("testdata/enum.go")
	fset := token.NewFileSet()
	root, err := parser.ParseFile(fset, file, nil, parser.ParseComments)
	if err != nil {
		panic(err)
	}
	sc := NewEnumCollector(ExprString)
	ast.Walk(sc, root)
	for k, v := range sc.Methods {
		if IsEnum(v) {
			em := EnumMeta{
				Name:   k,
				Values: sc.Consts[k],
			}
			sc.Enums[k] = em
		}
	}
	require.Equal(t, 1, len(sc.Enums))
}
