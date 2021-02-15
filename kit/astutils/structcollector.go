package astutils

import (
	"go/ast"
	"go/token"
	"log"
	"strings"
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
}

type StructCollector struct {
	Structs []StructMeta
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
	}
	return ""
}

func (sc *StructCollector) Collect(n ast.Node) ast.Visitor {
	switch spec := n.(type) {
	case *ast.Package:
		return sc
	case *ast.File: // actually it is package name
		log.Printf("File: name=%s\n", spec.Name)
		sc.Package = PackageMeta{
			Name: spec.Name.Name,
		}
		return sc
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
				log.Printf("Type: name=%s\n", typeName)
				switch specType := typeSpec.Type.(type) {
				case *ast.StructType:
					var fields []FieldMeta
					for _, field := range specType.Fields.List {
						var tag string
						if field.Tag != nil {
							tag = field.Tag.Value
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
							fieldType = "embed"
						}

						for _, name := range names {
							log.Printf("\tField: name=%s type=%s tag=%s\n", name, fieldType, tag)
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
