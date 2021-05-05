package codegen

import (
	"bufio"
	"github.com/sirupsen/logrus"
	"github.com/unionj-cloud/go-doudou/astutils"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

var httpServerTmpl = `package httpsrv

import (
	"{{.ConfigPackage}}"
	"github.com/gorilla/mux"
	"github.com/olekukonko/tablewriter"
	"github.com/sirupsen/logrus"
	"net/http"
	"strings"
)

func newRouter(handler {{.SvcName}}Handler) *mux.Router {
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

func Run(conf config.HttpConfig, handler {{.Meta.Name}}Handler) *http.Server {
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
`

func GenHttpServer(dir string, ic astutils.InterfaceCollector) {
	var (
		err        error
		serverfile string
		f          *os.File
		tpl        *template.Template
		httpDir    string
		modName    string
		modfile    string
		firstLine  string
		svcName    string
	)
	httpDir = filepath.Join(dir, "transport/httpsrv")
	if err = os.MkdirAll(httpDir, os.ModePerm); err != nil {
		panic(err)
	}

	modfile = filepath.Join(dir, "go.mod")
	serverfile = filepath.Join(httpDir, "server.go")
	svcName = ic.Interfaces[0].Name
	if _, err = os.Stat(serverfile); os.IsNotExist(err) {
		if f, err = os.Open(modfile); err != nil {
			panic(err)
		}
		reader := bufio.NewReader(f)
		if firstLine, err = reader.ReadString('\n'); err != nil {
			panic(err)
		}
		modName = strings.TrimSpace(strings.TrimPrefix(firstLine, "module"))

		if f, err = os.Create(serverfile); err != nil {
			panic(err)
		}
		defer f.Close()

		if tpl, err = template.New("server.go.tmpl").Parse(httpServerTmpl); err != nil {
			panic(err)
		}
		if err = tpl.Execute(f, struct {
			ConfigPackage string
			SvcName       string
			Meta          astutils.InterfaceMeta
		}{
			ConfigPackage: modName + "/config",
			SvcName:       svcName,
			Meta:          ic.Interfaces[0],
		}); err != nil {
			panic(err)
		}

	} else {
		logrus.Warnf("file %s already exists", serverfile)
	}
}
