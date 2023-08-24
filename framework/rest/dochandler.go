package rest

import (
	"bytes"
	"html/template"
	"net/http"
)

// Oas store OpenAPI3.0 description json string
var Oas string
var DocRoutes = docRoutes

func docRoutes() []Route {
	return []Route{
		{
			Name:    "GetDoc",
			Method:  "GET",
			Pattern: "/go-doudou/doc",
			HandlerFunc: func(_writer http.ResponseWriter, _req *http.Request) {
				var (
					tpl *template.Template
					err error
					buf bytes.Buffer
				)
				if tpl, err = template.New("onlinedoc.tmpl").Parse(onlinedocTmpl); err != nil {
					panic(err)
				}
				doc := Oas
				docUrl := "openapi.json"
				if err = tpl.Execute(&buf, struct {
					Doc    string
					DocUrl string
				}{
					Doc:    doc,
					DocUrl: docUrl,
				}); err != nil {
					panic(err)
				}
				_writer.Header().Set("Content-Type", "text/html; charset=utf-8")
				_writer.Write(buf.Bytes())
			},
		},
		{
			Name:    "GetOpenAPI",
			Method:  "GET",
			Pattern: "/go-doudou/openapi.json",
			HandlerFunc: func(_writer http.ResponseWriter, _req *http.Request) {
				_writer.Write([]byte(Oas))
			},
		},
	}
}
