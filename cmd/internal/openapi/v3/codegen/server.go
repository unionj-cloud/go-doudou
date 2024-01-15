package codegen

import (
	"bytes"
	"fmt"
	"github.com/iancoleman/strcase"
	"github.com/pkg/errors"
	"github.com/unionj-cloud/go-doudou/v2/cmd/internal/templates"
	"github.com/unionj-cloud/go-doudou/v2/toolkit/astutils"
	v3 "github.com/unionj-cloud/go-doudou/v2/toolkit/openapi/v3"
	"github.com/unionj-cloud/go-doudou/v2/toolkit/sliceutils"
	"github.com/unionj-cloud/go-doudou/v2/toolkit/stringutils"
	"github.com/unionj-cloud/go-doudou/v2/version"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"
)

type IOperationConverter interface {
	ConvertOperation(endpoint, httpMethod string, operation *v3.Operation, gparams []v3.Parameter) (astutils.MethodMeta, error)
}

type ServerOperationConverter struct {
	_         [0]int
	Generator *OpenAPICodeGenerator
}

// /shelves/{shelf}/books/{book}
// [shelves,{shelf},books,{book}]
// shelves_{shelf}_books_{book}
// Shelves_ShelfBooks_Book
func Pattern2Method(pattern string) string {
	pattern = strings.TrimSuffix(strings.TrimPrefix(pattern, "/"), "/")
	partials := strings.Split(pattern, "/")
	var converted []string
	nosymbolreg := regexp.MustCompile(`[^a-zA-Z0-9_]`)
	for _, item := range partials {
		if strings.HasPrefix(item, "{") && strings.HasSuffix(item, "}") {
			item = strings.TrimSuffix(strings.TrimPrefix(item, "{"), "}")
			item = nosymbolreg.ReplaceAllLiteralString(item, "")
			item = cases.Title(language.English).String(strings.ToLower(item))
			converted = append(converted, "_"+item)
			continue
		}
		item = nosymbolreg.ReplaceAllLiteralString(item, "")
		item = cases.Title(language.English).String(strings.ToLower(item))
		converted = append(converted, item)
	}
	return strings.Join(converted, "")
}

func (receiver *ServerOperationConverter) globalParams(gparams []v3.Parameter) (queryVars, pathVars, headerVars []astutils.FieldMeta) {
	for _, item := range gparams {
		switch item.In {
		case v3.InQuery:
			queryVars = append(queryVars, receiver.Generator.parameter2Field(item))
		case v3.InPath:
			pathVar := receiver.Generator.parameter2Field(item)
			pathVar.Name = strings.ToLower(pathVar.Name)
			pathVars = append(pathVars, pathVar)
		case v3.InHeader:
			headerVars = append(headerVars, receiver.Generator.parameter2Field(item))
		default:
			panic(fmt.Errorf("not support %s parameter yet", item.In))
		}
	}
	return
}

func (receiver *ServerOperationConverter) operationParams(parameters []v3.Parameter, queryVars, pathVars, headerVars *[]astutils.FieldMeta) {
	for _, item := range parameters {
		switch item.In {
		case v3.InQuery:
			*queryVars = append(*queryVars, receiver.Generator.parameter2Field(item))
		case v3.InPath:
			pathVar := receiver.Generator.parameter2Field(item)
			pathVar.Name = strings.ToLower(pathVar.Name)
			*pathVars = append(*pathVars, pathVar)
		case v3.InHeader:
			*headerVars = append(*headerVars, receiver.Generator.parameter2Field(item))
		default:
			panic(fmt.Errorf("not support %s parameter yet", item.In))
		}
	}
}

func (receiver *ServerOperationConverter) parseForm(form *v3.MediaType) (bodyParams []astutils.FieldMeta) {
	schema := *form.Schema
	if stringutils.IsNotEmpty(schema.Ref) {
		schema = receiver.Generator.Schemas[strings.TrimPrefix(form.Schema.Ref, "#/components/schemas/")]
	}
	for k, v := range schema.Properties {
		var gotype string
		if v.Type == v3.StringT && v.Format == v3.BinaryF {
			gotype = "v3.FileModel"
		} else if v.Type == v3.ArrayT && v.Items.Type == v3.StringT && v.Items.Format == v3.BinaryF {
			gotype = "[]v3.FileModel"
		}
		if stringutils.IsNotEmpty(gotype) {
			continue
		}
		field := receiver.Generator.schema2Field(v, k, nil, v3.UNKNOWN_EXAMPLE)
		if !sliceutils.StringContains(schema.Required, k) {
			field.Type = v3.ToOptional(field.Type)
		}
		bodyParams = append(bodyParams, *field)
	}
	return
}

func (receiver *ServerOperationConverter) form(operation *v3.Operation) (bodyParams []astutils.FieldMeta) {
	receiver.Generator.resolveSchemaFromRef(operation)
	content := operation.RequestBody.Content
	if content.JSON != nil {
	} else if content.FormURL != nil {
		bodyParams = receiver.parseForm(content.FormURL)
	} else if content.FormData != nil {
		bodyParams = receiver.parseForm(content.FormData)
	}
	return
}

func (receiver *ServerOperationConverter) ConvertOperation(endpoint, httpMethod string, operation *v3.Operation, gparams []v3.Parameter) (astutils.MethodMeta, error) {
	var files, params, bodyParams []astutils.FieldMeta
	var bodyJSON *astutils.FieldMeta
	comments := commentLines(operation)
	queryVars, pathvars, headervars := receiver.globalParams(gparams)
	receiver.operationParams(operation.Parameters, &queryVars, &pathvars, &headervars)

	if httpMethod != "Get" && operation.RequestBody != nil {
		bodyJSON, files = receiver.Generator.requestBody(operation)
		bodyParams = receiver.form(operation)
	}

	if operation.Responses == nil {
		return astutils.MethodMeta{}, errors.Errorf("response definition not found in api %s %s", httpMethod, endpoint)
	}

	if operation.Responses.Resp200 == nil {
		return astutils.MethodMeta{}, errors.Errorf("200 response definition not found in api %s %s", httpMethod, endpoint)
	}

	results, err := receiver.Generator.responseBody(endpoint, httpMethod, operation)
	if err != nil {
		return astutils.MethodMeta{}, err
	}

	params = append(params, queryVars...)
	params = append(params, pathvars...)
	params = append(params, headervars...)
	params = append(params, bodyParams...)
	params = append(params, files...)

	if bodyJSON != nil {
		params = append(params, *bodyJSON)
	}

	if httpMethod == "Post" {
		httpMethod = ""
	}

	ret := astutils.MethodMeta{
		Name:       httpMethod + Pattern2Method(endpoint),
		Params:     params,
		Results:    results,
		PathVars:   pathvars,
		HeaderVars: headervars,
		BodyJSON:   bodyJSON,
		Files:      files,
		Comments:   comments,
		Path:       endpoint,
	}
	return ret, nil
}

var svcTmpl = templates.EditableHeaderTmpl + `package service

import (
	"context"
	"{{.DtoPackage}}"
	v3 "github.com/unionj-cloud/go-doudou/v2/toolkit/openapi/v3"
)

//go:generate go-doudou svc http
//go:generate go-doudou svc grpc

{{ range $i, $c := .Meta.Comments }}
{{- if eq $i 0}}
// {{$.Meta.Name}} {{$c}}
{{- else}}
// {{$c}}
{{- end}}
{{- end }}
type {{.Meta.Name}} interface {
	{{ range $m := .Meta.Methods }}
	{{- range $i, $c := $m.Comments }}
	{{- if eq $i 0}}
	// {{$m.Name}} {{$c}}
	{{- else}}
	// {{$c}}
	{{- end}}
	{{- end }}
	{{$m.Name}}(ctx context.Context{{if $m.Params}}, {{end}}{{- range $i, $p := $m.Params}}
	{{- if $i}},{{end}}
	{{- range $c := $p.Comments }}
	// {{$c}}
	{{- end }}
	{{ $p.Name}} {{$p.Type}}
	{{- end }}) ({{- range $i, $r := $m.Results}}
					 {{- if $i}},{{end}}
					 {{- range $c := $r.Comments }}
					 // {{$c}}
					 {{- end }}
					 {{ $r.Name}} {{$r.Type}}
					 {{- end }}{{if $m.Results}}, {{end}}err error)

	{{ end }}
}
`

func (receiver *OpenAPICodeGenerator) GenGoInterface(output string, paths map[string]v3.Path, operationConverter IOperationConverter) {
	_ = os.MkdirAll(filepath.Dir(output), os.ModePerm)
	funcMap := make(map[string]interface{})
	funcMap["toCamel"] = strcase.ToCamel
	funcMap["contains"] = strings.Contains
	funcMap["restyMethod"] = RestyMethod
	funcMap["toUpper"] = strings.ToUpper
	funcMap["isOptional"] = IsOptional
	tpl, _ := template.New(svcTmpl).Funcs(funcMap).Parse(svcTmpl)
	var sqlBuf bytes.Buffer
	_ = tpl.Execute(&sqlBuf, struct {
		Meta       astutils.InterfaceMeta
		DtoPackage string
		Version    string
	}{
		Meta:       receiver.Api2Interface(paths, receiver.SvcName, operationConverter),
		DtoPackage: receiver.ModName + "/dto",
		Version:    version.Release,
	})
	source := strings.TrimSpace(sqlBuf.String())
	astutils.FixImport([]byte(source), output)
}
