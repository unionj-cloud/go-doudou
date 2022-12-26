package server

import (
	"github.com/iancoleman/strcase"
	"github.com/unionj-cloud/go-doudou/v2/cmd/internal/openapi/v3/codegen"
	v3 "github.com/unionj-cloud/go-doudou/v2/toolkit/openapi/v3"
	"github.com/unionj-cloud/go-doudou/v2/toolkit/stringutils"
	"os"
	"path/filepath"
)

var dtoTmpl = `package {{.Pkg}}

//go:generate go-doudou name --file $GOFILE

{{- range $k, $v := .Schemas }}
{{ toComment $v.Description ($k | toCamel)}}
type {{$k | toCamel}} struct {
{{- range $pk, $pv := $v.Properties }}
	{{ $pv.Description | toComment }}
	{{- if stringContains $v.Required $pk }}
	// required
	{{ $pk | toCamel}} {{$pv | toGoType }}
	{{- else }}
	{{ $pk | toCamel}} {{$pv | toOptionalGoType }}
	{{- end }}
{{- end }}
}
{{- end }}
`

func GenSvcGo(dir string, docPath string) {
	if stringutils.IsEmpty(docPath) {
		matches, _ := filepath.Glob(filepath.Join(dir, "*_openapi3.json"))
		if len(matches) > 0 {
			docPath = matches[0]
		}
	}
	if stringutils.IsEmpty(docPath) {
		return
	}
	api := v3.LoadAPI(docPath)
	generator := &codegen.OpenAPICodeGenerator{
		Schemas:       api.Components.Schemas,
		RequestBodies: api.Components.RequestBodies,
		Responses:     api.Components.Responses,
	}
	var (
		f   *os.File
		err error
	)
	dtoFile := filepath.Join(dir, "dto", "dto.go")
	if f, err = os.Create(dtoFile); err != nil {
		panic(err)
	}
	defer f.Close()
	generator.GenGoDto(api.Components.Schemas, dtoFile, "dto", dtoTmpl)

	methodPaths := make(map[string]codegen.ApiPath)
	for endpoint, path := range api.Paths {
		method := codegen.Pattern2Method(endpoint)
		if path.Get != nil {
			method = "Get" + method
			methodPaths[method] = codegen.ApiPath{
				Operation:        path.Get,
				GlobalParameters: path.Parameters,
			}
		}
		if path.Post != nil {
			methodPaths[method] = codegen.ApiPath{
				Operation:        path.Post,
				GlobalParameters: path.Parameters,
			}
		}
		if path.Put != nil {
			method = "Put" + method
			methodPaths[method] = codegen.ApiPath{
				Operation:        path.Put,
				GlobalParameters: path.Parameters,
			}
		}
		if path.Delete != nil {
			method = "Delete" + method
			methodPaths[method] = codegen.ApiPath{
				Operation:        path.Delete,
				GlobalParameters: path.Parameters,
			}
		}
	}
	svcName := strcase.ToCamel(filepath.Base(dir))
	svcFile := filepath.Join(dir, "svc.go")
	if f, err = os.Create(svcFile); err != nil {
		panic(err)
	}
	defer f.Close()
	generator.GenGoInterface(methodPaths, svcName, svcFile)
}
