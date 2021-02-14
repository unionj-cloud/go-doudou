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
	Name string
	Type string
	Tag  string
}

type StructMeta struct {
	Name   string
	Fields []FieldMeta
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
			return sc
		}
	case *ast.TypeSpec:
		log.Printf("Struct: name=%s\n", spec.Name.Name)
		switch specType := spec.Type.(type) {
		case *ast.StructType:
			var fields []FieldMeta
			for _, field := range specType.Fields.List {
				var tag string
				if field.Tag != nil {
					tag = field.Tag.Value
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
						Name: name,
						Type: fieldType,
						Tag:  tag,
					})
				}
			}
			sc.Structs = append(sc.Structs, StructMeta{
				Name:   spec.Name.Name,
				Fields: fields,
			})
		}
	}
	return nil
}
