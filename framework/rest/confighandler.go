package rest

import (
	"fmt"
	"github.com/unionj-cloud/go-doudou/v2/toolkit/stringutils"
	"net/http"
	"os"
	"strings"
)

var ConfigRoutes = configRoutes

func configRoutes() []Route {
	return []Route{
		{
			Name:    "GetConfig",
			Method:  "GET",
			Pattern: "/go-doudou/config",
			HandlerFunc: func(_writer http.ResponseWriter, _req *http.Request) {
				pre := _req.FormValue("pre")
				var builder strings.Builder
				for _, pair := range os.Environ() {
					if stringutils.IsEmpty(pre) || strings.HasPrefix(pair, pre) {
						builder.WriteString(fmt.Sprintf("%s\n", pair))
					}
				}
				_writer.Header().Set("Content-Type", "text/plain; charset=utf-8")
				_writer.Write([]byte(builder.String()))
			},
		},
	}
}
