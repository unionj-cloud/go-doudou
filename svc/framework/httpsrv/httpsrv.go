package httpsrv

import (
	"context"
	"fmt"
	"github.com/common-nighthawk/go-figure"
	"github.com/gorilla/mux"
	"github.com/olekukonko/tablewriter"
	"github.com/sirupsen/logrus"
	"github.com/unionj-cloud/go-doudou/pathutils"
	"github.com/unionj-cloud/go-doudou/stringutils"
	"github.com/unionj-cloud/go-doudou/svc/framework"
	"io"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"time"
)

type Interface interface {
	// Run the service
	Run()
	// Register routes
	Route(route Route)
	// Create http server
	NewServer(router http.Handler) *http.Server
}

type HttpSrv struct {
	*mux.Router
	routes []Route
}

func NewHttpSrv() Interface {
	return &HttpSrv{
		mux.NewRouter().StrictSlash(true),
		[]Route{},
	}
}

func (srv *HttpSrv) Route(route Route) {
	srv.routes = append(srv.routes, route)
	srv.
		Methods(route.Method).
		Path(route.Pattern).
		Name(route.Name).
		Handler(route.HandlerFunc)
}

func (srv *HttpSrv) printRoutes() {
	logrus.Infoln("================ Registered Routes ================")
	data := [][]string{}
	for _, r := range srv.routes {
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

func (srv *HttpSrv) Run() {
	start := time.Now()
	var logptr *string
	logpath, isSet := os.LookupEnv("APP_LOGPATH")
	if isSet {
		logptr = &logpath
	}
	var loglevel framework.LogLevel
	(&loglevel).Decode(os.Getenv("APP_LOGLEVEL"))

	logFile := configureLogger(logrus.StandardLogger(), logptr, logrus.Level(loglevel))
	defer func() {
		if logFile != nil {
			logFile.Close()
		}
	}()

	var bannerSwitch framework.Switch
	(&bannerSwitch).Decode(os.Getenv("APP_BANNER"))
	if bannerSwitch {
		banner := os.Getenv("APP_BANNERTEXT")
		if stringutils.IsEmpty(banner) {
			banner = "Go-doudou"
		}
		figure.NewColorFigure(banner, "doom", "green", true).Print()
	}

	srv.printRoutes()

	server := srv.NewServer(srv.Router)

	logrus.Infof("Started in %s\n", time.Since(start))

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	grace, err := time.ParseDuration(os.Getenv("APP_GRACETIMEOUT"))
	if err != nil {
		logrus.Warnf("Parse %s %s as time.Duration failed: %s, use default 15s instead.\n", "APP_GRACETIMEOUT",
			os.Getenv("APP_GRACETIMEOUT"), err.Error())
		grace = 15 * time.Second
	}

	ctx, cancel := context.WithTimeout(context.Background(), grace)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	server.Shutdown(ctx)

	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	logrus.Infoln("shutting down")
}

func (srv *HttpSrv) NewServer(router http.Handler) *http.Server {
	host := os.Getenv("SRV_HOST")
	port := os.Getenv("SRV_PORT")
	write, err := time.ParseDuration(os.Getenv("SRV_WRITETIMEOUT"))
	if err != nil {
		logrus.Warnf("Parse %s %s as time.Duration failed: %s, use default 15s instead.\n", "SRV_WRITETIMEOUT",
			os.Getenv("SRV_WRITETIMEOUT"), err.Error())
		write = 15 * time.Second
	}

	read, err := time.ParseDuration(os.Getenv("SRV_READTIMEOUT"))
	if err != nil {
		logrus.Warnf("Parse %s %s as time.Duration failed: %s, use default 15s instead.\n", "SRV_READTIMEOUT",
			os.Getenv("SRV_READTIMEOUT"), err.Error())
		read = 15 * time.Second
	}

	idle, err := time.ParseDuration(os.Getenv("SRV_IDLETIMEOUT"))
	if err != nil {
		logrus.Warnf("Parse %s %s as time.Duration failed: %s, use default 60s instead.\n", "SRV_IDLETIMEOUT",
			os.Getenv("SRV_IDLETIMEOUT"), err.Error())
		idle = 60 * time.Second
	}

	server := &http.Server{
		Addr: strings.Join([]string{host, port}, ":"),
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: write,
		ReadTimeout:  read,
		IdleTimeout:  idle,
		Handler:      router, // Pass our instance of gorilla/mux in.
	}

	// Run our server in a goroutine so that it doesn't block.
	go func() {
		logrus.Infof("Http server is listening on %s\n", server.Addr)
		if err := server.ListenAndServe(); err != nil {
			logrus.Println(err)
		}
	}()

	return server
}

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

func configureLogger(logger *logrus.Logger, logptr *string, level logrus.Level) *os.File {
	formatter := new(logrus.TextFormatter)
	formatter.TimestampFormat = "2006-01-02 15:04:05"
	formatter.FullTimestamp = true
	logger.SetFormatter(formatter)
	logger.SetLevel(level)

	if logptr != nil {
		var (
			err error
			f   *os.File
		)
		logpath := *logptr
		logpath, err = pathutils.FixPath(logpath, "")
		if err != nil {
			logger.Errorln(fmt.Sprintf("%+v\n", err))
		}
		if stringutils.IsNotEmpty(logpath) {
			err = os.MkdirAll(logpath, os.ModePerm)
			if err != nil {
				logger.Errorln(err)
				return nil
			}
		}
		f, err = os.OpenFile(filepath.Join(logpath, "app.log"), os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
		if err != nil {
			logger.Errorf("error opening file: %v\n", err)
			return nil
		}
		mw := io.MultiWriter(os.Stdout, f)
		logger.SetOutput(mw)
		return f
	}

	return nil
}
