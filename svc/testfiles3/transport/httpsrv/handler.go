package httpsrv

import (
	"net/http"

	ddhttp "github.com/unionj-cloud/go-doudou/svc/http"
)

type Testfiles3Handler interface {
	PageUsers(w http.ResponseWriter, r *http.Request)
}

func Routes(handler Testfiles3Handler) []ddhttp.Route {
	return []ddhttp.Route{
		{
			"PageUsers",
			"POST",
			"/testfiles3/pageusers",
			handler.PageUsers,
		},
	}
}
