package httpsrv

import (
	"net/http"
	service "testfiles3"
)

type Testfiles3HandlerImpl struct {
	testfiles3 service.Testfiles3
}

func (receiver *Testfiles3HandlerImpl) PageUsers(w http.ResponseWriter, r *http.Request) {
	panic("implement me")
}

func NewTestfiles3Handler(testfiles3 service.Testfiles3) Testfiles3Handler {
	return &Testfiles3HandlerImpl{
		testfiles3,
	}
}
