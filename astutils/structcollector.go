package astutils

import (
	"github.com/sirupsen/logrus"
	"go/ast"
	"go/parser"
	"go/token"
	"regexp"
	"strings"
	"unicode"
)

// StructCollector collect structs by parsing source code
type StructCollector struct {
	Structs          []StructMeta
	Methods          map[string][]MethodMeta
	Package          PackageMeta
	NonStructTypeMap map[string]ast.Expr
	exprString       func(ast.Expr) string
}

// Visit traverse each node from source code
func (sc *StructCollector) Visit(n ast.Node) ast.Visitor {
	return sc.Collect(n)
}

// Collect collects all structs from source code
func (sc *StructCollector) Collect(n ast.Node) ast.Visitor {
	switch spec := n.(type) {
	case *ast.Package:
		return sc
	case *ast.File: // actually it is package name
		sc.Package = PackageMeta{
			Name: spec.Name.Name,
		}
		return sc
	case *ast.FuncDecl:
		if spec.Recv != nil {
			structName := strings.TrimPrefix(sc.exprString(spec.Recv.List[0].Type), "*")
			methods, _ := sc.Methods[structName]
			methods = append(methods, GetMethodMeta(spec))
			if sc.Methods == nil {
				sc.Methods = make(map[string][]MethodMeta)
			}
			sc.Methods[structName] = methods
		}
	case *ast.GenDecl:
		if spec.Tok == token.TYPE {
			var comments []string
			if spec.Doc != nil {
				for _, comment := range spec.Doc.List {
					comments = append(comments, strings.TrimSpace(strings.TrimPrefix(comment.Text, "//")))
				}
			}
			for _, item := range spec.Specs {
				typeSpec := item.(*ast.TypeSpec)
				typeName := typeSpec.Name.Name
				logrus.Printf("Type: name=%s\n", typeName)
				switch specType := typeSpec.Type.(type) {
				case *ast.StructType:
					structmeta := NewStructMeta(specType, sc.exprString)
					structmeta.Name = typeName
					structmeta.Comments = comments
					structmeta.IsExport = unicode.IsUpper(rune(typeName[0]))
					sc.Structs = append(sc.Structs, structmeta)
				default:
					sc.NonStructTypeMap[typeName] = typeSpec.Type
				}
			}
		}
	}
	return nil
}

// DocFlatEmbed flatten embed struct fields
func (sc *StructCollector) DocFlatEmbed() []StructMeta {
	structMap := make(map[string]StructMeta)
	for _, structMeta := range sc.Structs {
		if _, exists := structMap[structMeta.Name]; !exists {
			structMap[structMeta.Name] = structMeta
		}
	}

	var exStructs []StructMeta
	for _, structMeta := range sc.Structs {
		if !structMeta.IsExport {
			continue
		}
		exStructs = append(exStructs, structMeta)
	}

	re := regexp.MustCompile(`json:"(.*?)"`)
	var result []StructMeta
	for _, structMeta := range exStructs {
		_structMeta := StructMeta{
			Name:     structMeta.Name,
			Fields:   make([]FieldMeta, 0),
			Comments: make([]string, len(structMeta.Comments)),
			IsExport: true,
		}
		copy(_structMeta.Comments, structMeta.Comments)
		fieldMap := make(map[string]FieldMeta)
		embedFieldMap := make(map[string]FieldMeta)
		for _, fieldMeta := range structMeta.Fields {
			if strings.HasPrefix(fieldMeta.Type, "embed") {
				if re.MatchString(fieldMeta.Tag) {
					fieldMeta.Type = strings.TrimPrefix(fieldMeta.Type, "embed:")
					_structMeta.Fields = append(_structMeta.Fields, fieldMeta)
					fieldMap[fieldMeta.Name] = fieldMeta
				} else {
					if embeded, exists := structMap[fieldMeta.Name]; exists {
						for _, field := range embeded.Fields {
							if !field.IsExport {
								continue
							}
							embedFieldMap[field.Name] = field
						}
					}
				}
			} else if fieldMeta.IsExport {
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

// NewStructCollector initializes an StructCollector
func NewStructCollector(exprString func(ast.Expr) string) *StructCollector {
	return &StructCollector{
		Structs:          nil,
		Methods:          make(map[string][]MethodMeta),
		Package:          PackageMeta{},
		NonStructTypeMap: make(map[string]ast.Expr),
		exprString:       exprString,
	}
}

// BuildStructCollector initializes an StructCollector and collects structs
func BuildStructCollector(file string, exprString func(ast.Expr) string) StructCollector {
	sc := NewStructCollector(exprString)
	fset := token.NewFileSet()
	root, err := parser.ParseFile(fset, file, nil, parser.ParseComments)
	if err != nil {
		logrus.Panicln(err)
	}
	ast.Walk(sc, root)
	return *sc
}
