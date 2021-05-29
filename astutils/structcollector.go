package astutils

import (
	"bufio"
	"go/ast"
	"go/token"
	"os"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/unionj-cloud/go-doudou/sliceutils"
)

type PackageMeta struct {
	Name string
}

type FieldMeta struct {
	Name     string
	Type     string
	Tag      string
	Comments []string
}

type StructMeta struct {
	Name     string
	Fields   []FieldMeta
	Comments []string
	Methods  []MethodMeta
}

type StructCollector struct {
	Structs []StructMeta
	Methods map[string][]MethodMeta
	Package PackageMeta
}

func (sc *StructCollector) Visit(n ast.Node) ast.Visitor {
	return sc.Collect(n)
}

func exprString(expr ast.Expr) string {
	switch _expr := expr.(type) {
	case *ast.Ident:
		return _expr.Name
	case *ast.StarExpr:
		return "*" + exprString(_expr.X)
	case *ast.SelectorExpr:
		return exprString(_expr.X) + "." + _expr.Sel.Name
	case *ast.InterfaceType:
		return "interface{}"
	case *ast.ArrayType:
		if _expr.Len == nil {
			return "[]" + exprString(_expr.Elt)
		} else {
			return "[" + exprString(_expr.Len) + "]" + exprString(_expr.Elt)
		}
	case *ast.BasicLit:
		if _expr.Kind == token.INT {
			return _expr.Value
		}
	case *ast.MapType:
		return "map[" + exprString(_expr.Key) + "]" + exprString(_expr.Value)
	case *ast.ChanType: // TODO
	case *ast.FuncType: // TODO
	}
	return ""
}

func (sc *StructCollector) Collect(n ast.Node) ast.Visitor {
	switch spec := n.(type) {
	case *ast.Package:
		return sc
	case *ast.File: // actually it is package name
		logrus.Printf("File: name=%s\n", spec.Name)
		sc.Package = PackageMeta{
			Name: spec.Name.Name,
		}
		return sc
	case *ast.FuncDecl:
		if spec.Recv != nil {
			structName := strings.TrimPrefix(exprString(spec.Recv.List[0].Type), "*")
			methods, _ := sc.Methods[structName]
			methods = append(methods, getMethodMeta(spec))
			sc.Methods[structName] = methods
		}
	case *ast.GenDecl:
		if spec.Tok == token.TYPE {
			var comments []string
			if spec.Doc != nil {
				for _, comment := range spec.Doc.List {
					comments = append(comments, comment.Text)
				}
			}
			for _, item := range spec.Specs {
				typeSpec := item.(*ast.TypeSpec)
				typeName := typeSpec.Name.Name
				logrus.Printf("Type: name=%s\n", typeName)
				switch specType := typeSpec.Type.(type) {
				case *ast.StructType:
					var fields []FieldMeta
					for _, field := range specType.Fields.List {
						var tag string
						if field.Tag != nil {
							tag = strings.Trim(field.Tag.Value, "`")
						}

						var fieldComments []string
						if field.Comment != nil {
							for _, comment := range field.Comment.List {
								fieldComments = append(fieldComments, comment.Text)
							}
						}

						var names []string
						fieldType := exprString(field.Type)

						if field.Names != nil {
							for _, name := range field.Names {
								names = append(names, name.Name)
							}
						} else {
							splits := strings.Split(fieldType, ".")
							names = append(names, splits[len(splits)-1])
							fieldType = "embed:" + fieldType
						}

						for _, name := range names {
							logrus.Printf("\tField: name=%s type=%s tag=%s\n", name, fieldType, tag)
							fields = append(fields, FieldMeta{
								Name:     name,
								Type:     fieldType,
								Tag:      tag,
								Comments: fieldComments,
							})
						}
					}

					sc.Structs = append(sc.Structs, StructMeta{
						Name:     typeName,
						Fields:   fields,
						Comments: comments,
					})
				}
			}
		}
	}
	return nil
}

// Only for ddl tool
func (sc *StructCollector) FlatEmbed() []StructMeta {
	structMap := make(map[string]StructMeta)
	for _, structMeta := range sc.Structs {
		if _, exists := structMap[structMeta.Name]; !exists {
			structMap[structMeta.Name] = structMeta
		}
	}
	var result []StructMeta
	for _, structMeta := range sc.Structs {
		if sliceutils.IsEmpty(structMeta.Comments) {
			continue
		}
		if structMeta.Comments[0] != "//dd:table" {
			continue
		}
		_structMeta := StructMeta{
			Name:     structMeta.Name,
			Fields:   make([]FieldMeta, 0),
			Comments: make([]string, len(structMeta.Comments)),
		}
		copy(_structMeta.Comments, structMeta.Comments)

		fieldMap := make(map[string]FieldMeta)
		embedFieldMap := make(map[string]FieldMeta)
		for _, fieldMeta := range structMeta.Fields {
			if strings.HasPrefix(fieldMeta.Type, "embed") {
				if embeded, exists := structMap[fieldMeta.Name]; exists {
					for _, field := range embeded.Fields {
						embedFieldMap[field.Name] = field
					}
				}
			} else {
				_structMeta.Fields = append(_structMeta.Fields, fieldMeta)
				fieldMap[fieldMeta.Name] = fieldMeta
			}
		}

		for key, field := range embedFieldMap {
			if _, exists := fieldMap[key]; !exists {
				_structMeta.Fields = append(_structMeta.Fields, field)
			}
		}
		result = append(result, _structMeta)
	}

	return result
}

func Visit(files *[]string) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			logrus.Panicln(err)
		}
		if !info.IsDir() {
			*files = append(*files, path)
		}
		return nil
	}
}

func GetMod() string {
	var (
		f         *os.File
		err       error
		firstLine string
	)
	dir, _ := os.Getwd()
	mod := filepath.Join(dir, "go.mod")
	if f, err = os.Open(mod); err != nil {
		panic(err)
	}
	reader := bufio.NewReader(f)
	if firstLine, err = reader.ReadString('\n'); err != nil {
		panic(err)
	}
	return strings.TrimSpace(strings.TrimPrefix(firstLine, "module"))
}

func GetImportPath(file string) string {
	dir, _ := os.Getwd()
	return GetMod() + strings.TrimPrefix(file, dir)
}

func NewStructCollector() *StructCollector {
	return &StructCollector{
		Structs: nil,
		Methods: make(map[string][]MethodMeta),
		Package: PackageMeta{},
	}
}
