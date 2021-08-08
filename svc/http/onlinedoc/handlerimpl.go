package onlinedoc

import (
	"bytes"
	"encoding/json"
	"net/http"
	"text/template"
)

var Oas string

type OnlineDocHandlerImpl struct {
}

func (receiver *OnlineDocHandlerImpl) GetOpenAPI(_writer http.ResponseWriter, _req *http.Request) {
	_writer.Write([]byte(Oas))
}

func (receiver *OnlineDocHandlerImpl) GetDoc(_writer http.ResponseWriter, _req *http.Request) {
	var (
		tpl *template.Template
		err error
		buf bytes.Buffer
	)
	if tpl, err = template.New("handlerimpl.go.tmpl").Parse(indexTmpl); err != nil {
		panic(err)
	}
	doc, _ := json.Marshal(Oas)
	if err = tpl.Execute(&buf, struct {
		Doc string
	}{
		Doc: string(doc),
	}); err != nil {
		panic(err)
	}
	_writer.WriteHeader(http.StatusOK)
	_writer.Header().Set("Content-Type", "text/html; charset=utf-8")
	_writer.Write(buf.Bytes())
}

func NewOnlineDocHandler() OnlineDocHandler {
	return &OnlineDocHandlerImpl{}
}
