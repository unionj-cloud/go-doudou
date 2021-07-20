package onlinedoc

import (
	"net/http"
)

var Oas string

type OnlineDocHandlerImpl struct {
}

func (receiver *OnlineDocHandlerImpl) GetOpenAPI(_writer http.ResponseWriter, _req *http.Request) {
	_writer.Write([]byte(Oas))
}

func (receiver *OnlineDocHandlerImpl) GetDoc(_writer http.ResponseWriter, _req *http.Request) {
	_writer.Write([]byte("Greeting from doudou"))
}

func NewOnlineDocHandler() OnlineDocHandler {
	return &OnlineDocHandlerImpl{}
}
