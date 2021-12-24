package ddhttp

import (
	"context"
	"fmt"
	"github.com/common-nighthawk/go-figure"
	"github.com/gorilla/mux"
	"github.com/olekukonko/tablewriter"
	"github.com/opentracing-contrib/go-stdlib/nethttp"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"github.com/uber/jaeger-lib/metrics"
	jprom "github.com/uber/jaeger-lib/metrics/prometheus"
	"github.com/unionj-cloud/go-doudou/pathutils"
	"github.com/unionj-cloud/go-doudou/stringutils"
	"github.com/unionj-cloud/go-doudou/svc/config"
	"github.com/unionj-cloud/go-doudou/svc/http/model"
	"github.com/unionj-cloud/go-doudou/svc/http/onlinedoc"
	"github.com/unionj-cloud/go-doudou/svc/http/prometheus"
	"github.com/unionj-cloud/go-doudou/svc/http/registry"
	"github.com/unionj-cloud/go-doudou/svc/tracing"
	"io"
	"net/http"
	"os"
	"os/signal"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// DefaultHttpSrv wraps gorilla mux router
type DefaultHttpSrv struct {
	*mux.Router
	rootRouter    *mux.Router
	routes        []model.Route
	server        *http.Server
	logFile       *os.File
	enableTracing bool
	tracer        opentracing.Tracer
	tracingCloser io.Closer
}

func (srv *DefaultHttpSrv) run() {
	write, err := time.ParseDuration(config.GddWriteTimeout.Load())
	if err != nil {
		logrus.Warnf("Parse %s %s as time.Duration failed: %s, use default 15s instead.\n", "GDD_WRITE_TIMEOUT",
			config.GddWriteTimeout.Load(), err.Error())
		write = 15 * time.Second
	}

	read, err := time.ParseDuration(config.GddReadTimeout.Load())
	if err != nil {
		logrus.Warnf("Parse %s %s as time.Duration failed: %s, use default 15s instead.\n", "GDD_READ_TIMEOUT",
			config.GddReadTimeout.Load(), err.Error())
		read = 15 * time.Second
	}

	idle, err := time.ParseDuration(config.GddIdleTimeout.Load())
	if err != nil {
		logrus.Warnf("Parse %s %s as time.Duration failed: %s, use default 60s instead.\n", "GDD_IDLE_TIMEOUT",
			config.GddIdleTimeout.Load(), err.Error())
		idle = 60 * time.Second
	}

	srv.server = &http.Server{
		Addr: strings.Join([]string{config.GddHost.Load(), config.GddPort.Load()}, ":"),
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: write,
		ReadTimeout:  read,
		IdleTimeout:  idle,
		Handler:      srv.rootRouter, // Pass our instance of gorilla/mux in.
	}

	// Run our server in a goroutine so that it doesn't block.
	go func() {
		logrus.Infof("Http server is listening on %s\n", srv.server.Addr)
		if err := srv.server.ListenAndServe(); err != nil {
			logrus.Println(err)
		}
	}()
}

const gddPathPrefix = "/go-doudou/"

// NewDefaultHttpSrv create a DefaultHttpSrv instance
func NewDefaultHttpSrv() *DefaultHttpSrv {
	rootRouter := mux.NewRouter().StrictSlash(true)
	bizRouter := rootRouter.PathPrefix(config.GddRouteRootPath.Load()).Subrouter().StrictSlash(true)
	var routes []model.Route
	if config.GddManage.Load() == "true" {
		bizRouter.Use(prometheus.PrometheusMiddleware)
		gddRouter := rootRouter.PathPrefix(gddPathPrefix).Subrouter().StrictSlash(true)
		gddRouter.Use(BasicAuth)
		var mergedRoutes []model.Route
		mergedRoutes = append(mergedRoutes, onlinedoc.Routes()...)
		mergedRoutes = append(mergedRoutes, prometheus.Routes()...)
		mergedRoutes = append(mergedRoutes, registry.Routes()...)
		for _, item := range mergedRoutes {
			gddRouter.
				Methods(item.Method).
				Path("/" + strings.TrimPrefix(item.Pattern, gddPathPrefix)).
				Name(item.Name).
				Handler(item.HandlerFunc)
		}
		routes = append(routes, mergedRoutes...)
	}
	srv := &DefaultHttpSrv{
		Router:     bizRouter,
		rootRouter: rootRouter,
		routes:     routes,
	}
	srv.configureLogger()
	enableTracing, _ := strconv.ParseBool(config.GddEnableTracing.Load())
	if enableTracing {
		metricsRoot := config.FrameworkName
		if stringutils.IsNotEmpty(config.GddTracingMetricsRoot.Load()) {
			metricsRoot = config.GddTracingMetricsRoot.Load()
		}
		srv.tracer, srv.tracingCloser = tracing.Init(config.GddServiceName.Load(), jprom.New().Namespace(metrics.NSOptions{Name: metricsRoot, Tags: nil}))
		opentracing.SetGlobalTracer(srv.tracer)
	}
	return srv
}

func (srv *DefaultHttpSrv) configureLogger() {
	var logptr *string
	logpath, isSet := os.LookupEnv(config.GddLogPath.String())
	if isSet {
		logptr = &logpath
	}
	var loglevel config.LogLevel
	(&loglevel).Decode(config.GddLogLevel.Load())

	formatter := new(logrus.TextFormatter)
	formatter.TimestampFormat = "2006-01-02 15:04:05"
	formatter.FullTimestamp = true

	logger := logrus.StandardLogger()
	logger.SetFormatter(formatter)
	logger.SetLevel(logrus.Level(loglevel))

	if logptr != nil {
		var err error
		logpath := *logptr
		logpath, err = pathutils.FixPath(logpath, "")
		if err != nil {
			logger.Errorln(fmt.Sprintf("%+v\n", err))
		}
		if stringutils.IsNotEmpty(logpath) {
			err = os.MkdirAll(logpath, os.ModePerm)
			if err != nil {
				logger.Errorln(err)
			}
		}
		srv.logFile, err = os.OpenFile(filepath.Join(logpath, "app.log"), os.O_APPEND|os.O_CREATE|os.O_RDWR, os.ModePerm)
		if err != nil {
			logger.Errorf("error opening file: %v\n", err)
		}
		mw := io.MultiWriter(os.Stdout, srv.logFile)
		logger.SetOutput(mw)
	}
}

// AddRoute adds routes to router
func (srv *DefaultHttpSrv) AddRoute(route ...model.Route) {
	var routes []model.Route
	routes = append(routes, route...)
	routes = append(routes, srv.routes...)
	srv.routes = routes[:]
	routes = nil
	for _, item := range route {
		var handler http.Handler
		handler = item.HandlerFunc
		if srv.enableTracing {
			handler = nethttp.Middleware(
				opentracing.GlobalTracer(),
				handler,
				nethttp.OperationNameFunc(func(r *http.Request) string {
					return "HTTP " + r.Method + " " + item.Pattern
				}))
		}
		srv.
			Methods(item.Method).
			Path(item.Pattern).
			Name(item.Name).
			Handler(handler)
	}
}

func (srv *DefaultHttpSrv) printRoutes() {
	logrus.Infoln("================ Registered Routes ================")
	data := [][]string{}
	for _, r := range srv.routes {
		if strings.HasPrefix(r.Pattern, gddPathPrefix) {
			data = append(data, []string{r.Name, r.Method, r.Pattern})
		} else {
			data = append(data, []string{r.Name, r.Method, path.Clean(config.GddRouteRootPath.Load() + r.Pattern)})
		}
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

// AddMiddleware adds middlewares to router
func (srv *DefaultHttpSrv) AddMiddleware(mwf ...func(http.Handler) http.Handler) {
	var middlewares []mux.MiddlewareFunc
	for _, item := range mwf {
		middlewares = append(middlewares, item)
	}
	srv.Use(middlewares...)
}

// Run runs http server
func (srv *DefaultHttpSrv) Run() {
	start := time.Now()
	var bannerSwitch config.Switch
	(&bannerSwitch).Decode(config.GddBanner.Load())
	if bannerSwitch {
		banner := config.GddBannerText.Load()
		if stringutils.IsEmpty(banner) {
			banner = config.FrameworkName
		}
		figure.NewColorFigure(banner, "doom", "green", true).Print()
	}

	srv.printRoutes()
	srv.run()

	logrus.Infof("Started in %s\n", time.Since(start))

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c
}

// Shutdown runs http server
func (srv *DefaultHttpSrv) Shutdown() {
	logrus.Infoln("shutting down...")

	if srv.logFile != nil {
		srv.logFile.Close()
	}

	if srv.enableTracing && srv.tracingCloser != nil {
		srv.tracingCloser.Close()
	}

	if srv.server != nil {
		// Create a deadline to wait for.
		grace, err := time.ParseDuration(config.GddGraceTimeout.Load())
		if err != nil {
			logrus.Warnf("Parse %s %s as time.Duration failed: %s, use default 15s instead.\n", "GDD_GRACETIMEOUT",
				config.GddGraceTimeout.Load(), err.Error())
			grace = 15 * time.Second
		}

		ctx, cancel := context.WithTimeout(context.Background(), grace)
		defer cancel()
		// Doesn't block if no connections, but will otherwise wait
		// until the timeout deadline.
		srv.server.Shutdown(ctx)
	}
}
