package onlinedoc

import (
	"net/http"
)

type OnlineDocHandlerImpl struct {
}

func (receiver *OnlineDocHandlerImpl) GetDoc(_writer http.ResponseWriter, _req *http.Request) {
	_writer.Write([]byte("Greeting from doudou"))
}

func NewOnlineDocHandler() OnlineDocHandler {
	return &OnlineDocHandlerImpl{}
}
