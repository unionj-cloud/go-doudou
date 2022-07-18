package httpsrv

import (
	"net/http"
	service "testdatahandlerImpl12"
)

type TestdatahandlerImpl12HandlerImpl struct {
	testdatahandlerImpl12 service.TestdatahandlerImpl12
}

func (receiver *TestdatahandlerImpl12HandlerImpl) PageUsers(w http.ResponseWriter, r *http.Request) {
	panic("implement me")
}

func NewTestdatahandlerImpl12Handler(testdatahandlerImpl12 service.TestdatahandlerImpl12) TestdatahandlerImpl12Handler {
	return &TestdatahandlerImpl12HandlerImpl{
		testdatahandlerImpl12,
	}
}
