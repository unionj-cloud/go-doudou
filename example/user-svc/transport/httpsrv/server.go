package httpsrv

import (
	"usersvc/config"
	"github.com/gorilla/mux"
	"github.com/olekukonko/tablewriter"
	"github.com/sirupsen/logrus"
	"net/http"
	"strings"
)

func newRouter(handler UserServiceHandler) *mux.Router {
	rous := routes(handler)
	router := mux.NewRouter().StrictSlash(true)
	for _, r := range rous {
		var hh http.Handler

		hh = r.HandlerFunc
		hh = logger(hh, r.Name)
		hh = rest(hh)

		router.
			Methods(r.Method).
			Path(r.Pattern).
			Name(r.Name).
			Handler(hh)
	}
	printRoutes(rous)
	return router
}

func printRoutes(rous []route) {
	logrus.Infoln("================ Registered Routes ================")
	data := [][]string{}
	for _, r := range rous {
		data = append(data, []string{r.Name, r.Method, r.Pattern})
	}

	tableString := &strings.Builder{}
	table := tablewriter.NewWriter(tableString)
	table.SetHeader([]string{"Name", "Method", "Pattern"})
	for _, v := range data {
		table.Append(v)
	}
	table.Render() // Send output
	rows := strings.Split(strings.TrimSpace(tableString.String()), "\n")
	for _, row := range rows {
		logrus.Infoln(row)
	}
	logrus.Infoln("===================================================")
}

func Run(conf config.HttpConfig, handler UserServiceHandler) *http.Server {
	srv := &http.Server{
		Addr: strings.Join([]string{conf.Host, conf.Port}, ":"),
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: conf.WriteTimeout,
		ReadTimeout:  conf.ReadTimeout,
		IdleTimeout:  conf.IdleTimeout,
		Handler:      newRouter(handler), // Pass our instance of gorilla/mux in.
	}

	// Run our server in a goroutine so that it doesn't block.
	go func() {
		logrus.Infof("Http server is listening on %s\n", srv.Addr)
		if err := srv.ListenAndServe(); err != nil {
			logrus.Println(err)
		}
	}()

	return srv
}
