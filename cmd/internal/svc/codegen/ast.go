package codegen

import (
	"encoding/json"
	"fmt"
	"github.com/unionj-cloud/go-doudou/v2/cmd/internal/astutils"
	"go/ast"
	"strings"
)

// ExprStringP Support all built-in type referenced from https://golang.org/pkg/builtin/
// Support map with string key
// Support structs of vo and dto package
// Support slice of types mentioned above
// Not support alias type (all alias type fields of a struct will be outputted as v3.Any in openapi 3.0 json document)
// Support anonymous struct type
// as struct field type in vo and dto package
// or as parameter type in method signature in svc.go file besides context.Context, multipart.FileHeader, v3.FileModel, os.File
// when go-doudou command line flag doc is true
func ExprStringP(expr ast.Expr) string {
	switch _expr := expr.(type) {
	case *ast.Ident:
		return _expr.Name
	case *ast.StarExpr:
		return "*" + ExprStringP(_expr.X)
	case *ast.SelectorExpr:
		return parseSelectorExpr(_expr)
	case *ast.InterfaceType:
		return "interface{}"
	case *ast.ArrayType:
		if _expr.Len == nil {
			return "[]" + ExprStringP(_expr.Elt)
		}
		return "[" + ExprStringP(_expr.Len) + "]" + ExprStringP(_expr.Elt)
	case *ast.BasicLit:
		return _expr.Value
	case *ast.MapType:
		if ExprStringP(_expr.Key) != "string" {
			panic("support string map key only in svc.go file and vo, dto package")
		}
		return "map[string]" + ExprStringP(_expr.Value)
	case *ast.StructType:
		structmeta := astutils.NewStructMeta(_expr, ExprStringP)
		b, _ := json.Marshal(structmeta)
		return "anonystruct«" + string(b) + "»"
	case *ast.Ellipsis:
		if _expr.Ellipsis.IsValid() {
			return "..." + ExprStringP(_expr.Elt)
		}
		panic(fmt.Sprintf("invalid ellipsis expression: %+v\n", expr))
	case *ast.FuncType:
		panic("not support function as struct field type in vo and dto package and as parameter in method signature in svc.go file")
	case *ast.ChanType:
		panic("not support channel as struct field type in vo and dto package and as parameter in method signature in svc.go file")
	default:
		panic(fmt.Errorf("not support expression as struct field type in vo and dto package and in method signature in svc.go file: %+v", expr))
	}
}

func parseSelectorExpr(expr *ast.SelectorExpr) string {
	result := ExprStringP(expr.X) + "." + expr.Sel.Name
	if !strings.HasPrefix(result, "vo.") &&
		!strings.HasPrefix(result, "dto.") &&
		result != "context.Context" &&
		result != "v3.FileModel" &&
		result != "multipart.FileHeader" &&
		result != "os.File" {
		panic(fmt.Errorf("not support %s in svc.go file and vo, dto package", result))
	}
	return result
}
