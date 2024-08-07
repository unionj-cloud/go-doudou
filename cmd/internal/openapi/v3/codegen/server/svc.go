package server

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/unionj-cloud/go-doudou/v2/cmd/internal/openapi/v3/codegen"
	"github.com/unionj-cloud/go-doudou/v2/cmd/internal/templates"
	"github.com/unionj-cloud/go-doudou/v2/toolkit/assert"
	v3 "github.com/unionj-cloud/go-doudou/v2/toolkit/openapi/v3"
)

var dtoTmpl = templates.EditableHeaderTmpl + `package {{.Pkg}}

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

// GenSvcGo may panic if docPath is empty string
func GenSvcGo(dir string, docPath string) {
	assert.NotEmpty(docPath, "docPath should not be empty")
	var (
		f   *os.File
		err error
	)
	svcName := strcase.ToCamel(filepath.Base(dir))
	svcFile := filepath.Join(dir, "svc.go")
	if f, err = os.Create(svcFile); err != nil {
		panic(err)
	}
	defer f.Close()
	modfile := filepath.Join(dir, "go.mod")
	var modf *os.File
	if modf, err = os.Open(modfile); err != nil {
		panic(err)
	}
	defer modf.Close()
	reader := bufio.NewReader(modf)
	firstLine, _ := reader.ReadString('\n')
	modName := strings.TrimSpace(strings.TrimPrefix(firstLine, "module"))
	api := v3.LoadAPI(docPath)
	generator := &codegen.OpenAPICodeGenerator{
		Schemas:       api.Components.Schemas,
		RequestBodies: api.Components.RequestBodies,
		Responses:     api.Components.Responses,
		SvcName:       svcName,
		ModName:       modName,
		ApiInfo:       api.Info,
	}
	operationConverter := &codegen.ServerOperationConverter{
		Generator: generator,
	}
	dtoFile := filepath.Join(dir, "dto", "dto.go")
	if f, err = os.Create(dtoFile); err != nil {
		panic(err)
	}
	defer f.Close()
	generator.GenGoDto(api.Components.Schemas, dtoFile, "dto", dtoTmpl)
	generator.DtoPkg = "dto"
	generator.GenGoInterface(svcFile, api.Paths, operationConverter)
}
