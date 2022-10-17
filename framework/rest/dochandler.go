package rest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
)

// Oas store OpenAPI3.0 description json string
var Oas string

func docRoutes() []Route {
	return []Route{
		{
			Name:    "GetDoc",
			Method:  "GET",
			Pattern: "/go-doudou/doc",
			HandlerFunc: func(_writer http.ResponseWriter, _req *http.Request) {
				var (
					tpl    *template.Template
					err    error
					buf    bytes.Buffer
					scheme string
				)
				if tpl, err = template.New("onlinedoc.tmpl").Parse(indexTmpl); err != nil {
					panic(err)
				}
				doc, _ := json.Marshal(Oas)
				if _req.TLS == nil {
					scheme = "http"
				} else {
					scheme = "https"
				}
				if err = tpl.Execute(&buf, struct {
					Doc    string
					DocUrl string
				}{
					Doc:    string(doc),
					DocUrl: fmt.Sprintf("%s://%s/go-doudou/openapi.json", scheme, _req.Host),
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
