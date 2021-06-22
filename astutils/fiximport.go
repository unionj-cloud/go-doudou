package astutils

import (
	"go/ast"
	"golang.org/x/tools/imports"
	"io/ioutil"
	"os"
)

func FixImport(src []byte, file string) {
	var (
		res []byte
		err error
	)
	if res, err = imports.Process(file, src, &imports.Options{
		TabWidth:  8,
		TabIndent: true,
		Comments:  true,
		Fragment:  true,
	}); err != nil {
		panic(err)
	}

	// On Windows, we need to re-set the permissions from the file. See golang/go#38225.
	var perms os.FileMode
	var fi os.FileInfo
	if fi, err = os.Stat(file); err == nil {
		perms = fi.Mode() & os.ModePerm
	}
	err = ioutil.WriteFile(file, res, perms)
	if err != nil {
		panic(err)
	}
}

func getMethodMeta(spec *ast.FuncDecl) MethodMeta {
	methodName := exprString(spec.Name)
	ft := spec.Type
	var params, results []FieldMeta
	for _, param := range ft.Params.List {
		var pn string
		if len(param.Names) > 0 {
			pn = param.Names[0].Name
		}
		pt := exprString(param.Type)
		var pComments []string
		if param.Comment != nil {
			for _, comment := range param.Comment.List {
				pComments = append(pComments, comment.Text)
			}
		}
		params = append(params, FieldMeta{
			Name:     pn,
			Type:     pt,
			Tag:      "",
			Comments: pComments,
		})
	}
	if ft.Results != nil {
		for _, result := range ft.Results.List {
			var rn string
			if len(result.Names) > 0 {
				rn = result.Names[0].Name
			}
			rt := exprString(result.Type)
			var rComments []string
			if result.Comment != nil {
				for _, comment := range result.Comment.List {
					rComments = append(rComments, comment.Text)
				}
			}
			results = append(results, FieldMeta{
				Name:     rn,
				Type:     rt,
				Tag:      "",
				Comments: rComments,
			})
		}
	}
	return MethodMeta{
		Name:     methodName,
		Params:   params,
		Results:  results,
		Comments: nil,
	}
}
