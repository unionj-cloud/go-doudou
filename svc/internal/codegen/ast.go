package codegen

import (
	"encoding/json"
	"fmt"
	. "github.com/unionj-cloud/go-doudou/astutils"
	"go/ast"
	"strings"
)

// Support all built-in type referenced from https://golang.org/pkg/builtin/
// Support map with string key
// Support structs of vo package
// Support slice of types mentioned above
// Not support alias type (all alias type fields of a struct will be outputed as v3.Any in openapi 3.0 json document)
// TODO support anonymous struct type
// as struct field type in vo package
// or as parameter type in method signature in svc.go file besides context.Context, multipart.FileHeader, os.File
// when go-doudou command line flag doc is true
func ExprStringP(expr ast.Expr) string {
	switch _expr := expr.(type) {
	case *ast.Ident:
		return _expr.Name
	case *ast.StarExpr:
		return "*" + ExprStringP(_expr.X)
	case *ast.SelectorExpr:
		result := ExprStringP(_expr.X) + "." + _expr.Sel.Name
		if !strings.HasPrefix(result, "vo.") &&
			result != "context.Context" &&
			result != "multipart.FileHeader" &&
			result != "os.File" {
			panic(fmt.Errorf("not support %s in svc.go file and vo package", result))
		}
		return result
	case *ast.InterfaceType:
		return "interface{}"
	case *ast.ArrayType:
		if _expr.Len == nil {
			return "[]" + ExprStringP(_expr.Elt)
		} else {
			return "[" + ExprStringP(_expr.Len) + "]" + ExprStringP(_expr.Elt)
		}
	case *ast.BasicLit:
		return _expr.Value
	case *ast.MapType:
		if ExprStringP(_expr.Key) != "string" {
			panic("support string map key only in svc.go file and vo package")
		}
		return "map[string]" + ExprStringP(_expr.Value)
	case *ast.StructType:
		structmeta := NewStructMeta(_expr, ExprStringP)
		b, _ := json.Marshal(structmeta)
		return "anonystruct«" + string(b) + "»"
	case *ast.FuncType:
		panic("not support function as struct field type in vo package and as parameter in method signature in svc.go file")
	case *ast.ChanType:
		panic("not support channel as struct field type in vo package and as parameter in method signature in svc.go file")
	default:
		panic(fmt.Errorf("not support expression as struct field type in vo package and in method signature in svc.go file: %+v", expr))
	}
}
