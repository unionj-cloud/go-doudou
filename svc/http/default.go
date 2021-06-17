package ddhttp

import (
	"context"
	"github.com/common-nighthawk/go-figure"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/unionj-cloud/go-doudou/stringutils"
	"github.com/unionj-cloud/go-doudou/svc/config"
	"net/http"
	"os"
	"os/signal"
	"time"
)

// gorilla
type DefaultHttpSrv struct {
	*mux.Router
	routes []Route
}

func NewDefaultHttpSrv() Srv {
	return &DefaultHttpSrv{
		mux.NewRouter().StrictSlash(true),
		[]Route{},
	}
}

func (srv *DefaultHttpSrv) AddRoute(route ...Route) {
	srv.routes = append(srv.routes, route...)
	for _, item := range route {
		srv.
			Methods(item.Method).
			Path(item.Pattern).
			Name(item.Name).
			Handler(item.HandlerFunc)
	}
}

func (srv *DefaultHttpSrv) AddMiddleware(mwf ...func(http.Handler) http.Handler) {
	var middlewares []mux.MiddlewareFunc
	for _, item := range mwf {
		middlewares = append(middlewares, item)
	}
	srv.Use(middlewares...)
}

func (srv *DefaultHttpSrv) Run() {
	start := time.Now()
	var logptr *string
	logpath, isSet := os.LookupEnv("APP_LOGPATH")
	if isSet {
		logptr = &logpath
	}
	var loglevel config.LogLevel
	(&loglevel).Decode(os.Getenv("APP_LOGLEVEL"))

	logFile := configureLogger(logrus.StandardLogger(), logptr, logrus.Level(loglevel))
	defer func() {
		if logFile != nil {
			logFile.Close()
		}
	}()

	var bannerSwitch config.Switch
	(&bannerSwitch).Decode(os.Getenv("APP_BANNER"))
	if bannerSwitch {
		banner := os.Getenv("APP_BANNERTEXT")
		if stringutils.IsEmpty(banner) {
			banner = "Go-doudou"
		}
		figure.NewColorFigure(banner, "doom", "green", true).Print()
	}

	printRoutes(srv.routes)

	server := newServer(srv)

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
