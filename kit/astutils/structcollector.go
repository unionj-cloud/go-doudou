package astutils

import (
	"go/ast"
	"go/token"
	"log"
)

type PackageMeta struct {
	Name string
}

type FieldMeta struct {
	Name string
	Type string
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
				i := field.Type.(*ast.Ident)
				fieldType := i.Name
				for _, name := range field.Names {
					log.Printf("\tField: name=%s type=%s\n", name.Name, fieldType)
					fields = append(fields, FieldMeta{
						Name: name.Name,
						Type: fieldType,
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
