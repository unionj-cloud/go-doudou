package onlinedoc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/unionj-cloud/go-doudou/svc/config"
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
		scheme string
		host string
	)
	if tpl, err = template.New("handlerimpl.go.tmpl").Parse(indexTmpl); err != nil {
		panic(err)
	}
	doc, _ := json.Marshal(Oas)
	if _req.TLS == nil {
		scheme = "http"
	} else {
		scheme = "https"
	}
	host = fmt.Sprintf("%s://%s%s", scheme, _req.Host, config.GddRouteRootPath.Load())
	if err = tpl.Execute(&buf, struct {
		Doc string
		DocUrl string
	}{
		Doc: string(doc),
		DocUrl: host + "/go-doudou/openapi.json",
	}); err != nil {
		panic(err)
	}
	_writer.Header().Set("Content-Type", "text/html; charset=utf-8")
	_writer.Write(buf.Bytes())
}

func NewOnlineDocHandler() OnlineDocHandler {
	return &OnlineDocHandlerImpl{}
}
