package httpsrv

import (
	"example/user-svc/config"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/olekukonko/tablewriter"
	"github.com/sirupsen/logrus"
	"github.com/unionj-cloud/go-doudou/stringutils"
	"github.com/unionj-cloud/go-doudou/timeutils"
	"net/http"
	"strings"
	"time"
)

func newRouter(handler UserHandler) *mux.Router {
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

func Run(conf config.HttpConfig, handler UserHandler) *http.Server {
	host := conf.Host
	if stringutils.IsEmpty(host) {
		host = "0.0.0.0"
	}
	port := conf.Port
	if stringutils.IsEmpty(port) {
		port = "8080"
	}
	wt, err := timeutils.Parse(conf.WriteTimeout, time.Second*15)
	if err != nil {
		logrus.Errorln(fmt.Sprintf("%+v", err))
	}

	rt, err := parseTimeout(conf.ReadTimeout, time.Second*15)
	if err != nil {
		logrus.Errorln(fmt.Sprintf("%+v", err))
	}

	it, err := parseTimeout(conf.IdleTimeout, time.Second*60)
	if err != nil {
		logrus.Errorln(fmt.Sprintf("%+v", err))
	}

	srv := &http.Server{
		Addr: strings.Join([]string{host, port}, ":"),
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: wt,
		ReadTimeout:  rt,
		IdleTimeout:  it,
		Handler:      newRouter(handler), // Pass our instance of gorilla/mux in.
	}

	// Run our server in a goroutine so that it doesn't block.
	go func() {
		logrus.Infof("server is listening on %s\n", srv.Addr)
		if err := srv.ListenAndServe(); err != nil {
			logrus.Println(err)
		}
	}()

	return srv
}
