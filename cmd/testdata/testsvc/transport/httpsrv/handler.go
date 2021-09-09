package httpsrv

import (
	"net/http"

	ddmodel "github.com/unionj-cloud/go-doudou/svc/http/model"
)

type TestsvcHandler interface {
	PageUsers(w http.ResponseWriter, r *http.Request)
}

func Routes(handler TestsvcHandler) []ddmodel.Route {
	return []ddmodel.Route{
		{
			"PageUsers",
			"POST",
			"/page/users",
			handler.PageUsers,
		},
	}
}
