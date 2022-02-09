package onlinedoc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/unionj-cloud/go-doudou/stringutils"
	"github.com/unionj-cloud/go-doudou/svc/config"
	"net/http"
	"text/template"
)

// Oas store OpenAPI3.0 description json string
var Oas string

// OnlineDocHandlerImpl define implementation for OnlineDocHandler
type OnlineDocHandlerImpl struct {
}

// GetOpenAPI return OpenAPI3.0 description json string
func (receiver *OnlineDocHandlerImpl) GetOpenAPI(_writer http.ResponseWriter, _req *http.Request) {
	_writer.Write([]byte(Oas))
}

// GetDoc return documentation web UI
func (receiver *OnlineDocHandlerImpl) GetDoc(_writer http.ResponseWriter, _req *http.Request) {
	var (
		tpl    *template.Template
		err    error
		buf    bytes.Buffer
		scheme string
		host   string
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
	rr := config.DefaultGddRouteRootPath
	if stringutils.IsNotEmpty(config.GddRouteRootPath.Load()) {
		rr = config.GddRouteRootPath.Load()
	}
	host = fmt.Sprintf("%s://%s%s", scheme, _req.Host, rr)
	if err = tpl.Execute(&buf, struct {
		Doc    string
		DocUrl string
	}{
		Doc:    string(doc),
		DocUrl: host + "/go-doudou/openapi.json",
	}); err != nil {
		panic(err)
	}
	_writer.Header().Set("Content-Type", "text/html; charset=utf-8")
	_writer.Write(buf.Bytes())
}

// NewOnlineDocHandler creates new instance for OnlineDocHandlerImpl
func NewOnlineDocHandler() OnlineDocHandler {
	return &OnlineDocHandlerImpl{}
}
