package ddhttp

import (
	"context"
	"fmt"
	"github.com/ascarter/requestid"
	"github.com/common-nighthawk/go-figure"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/olekukonko/tablewriter"
	"github.com/rs/cors"
	configui "github.com/unionj-cloud/go-doudou/framework/http/config"
	"github.com/unionj-cloud/go-doudou/framework/http/model"
	"github.com/unionj-cloud/go-doudou/framework/http/onlinedoc"
	"github.com/unionj-cloud/go-doudou/framework/http/prometheus"
	"github.com/unionj-cloud/go-doudou/framework/http/registry"
	"github.com/unionj-cloud/go-doudou/framework/internal/config"
	"github.com/unionj-cloud/go-doudou/framework/logger"
	"github.com/unionj-cloud/go-doudou/toolkit/cast"
	"github.com/unionj-cloud/go-doudou/toolkit/stringutils"
	"net/http"
	"os"
	"os/signal"
	"path"
	"strconv"
	"strings"
	"time"
)

// DefaultHttpSrv wraps gorilla mux router
type DefaultHttpSrv struct {
	*mux.Router
	rootRouter  *mux.Router
	gddRoutes   []model.Route
	bizRoutes   []model.Route
	middlewares []mux.MiddlewareFunc
}

const gddPathPrefix = "/go-doudou/"

// NewDefaultHttpSrv create a DefaultHttpSrv instance
func NewDefaultHttpSrv() *DefaultHttpSrv {
	rr := config.DefaultGddRouteRootPath
	if stringutils.IsNotEmpty(config.GddRouteRootPath.Load()) {
		rr = config.GddRouteRootPath.Load()
	}
	rootRouter := mux.NewRouter().StrictSlash(true)
	srv := &DefaultHttpSrv{
		Router:     rootRouter.PathPrefix(rr).Subrouter().StrictSlash(true),
		rootRouter: rootRouter,
	}
	srv.middlewares = append(srv.middlewares,
		tracing,
		metrics,
		handlers.CompressHandler,
	)
	logReq := config.DefaultGddLogReqEnable
	if l, err := cast.ToBoolE(config.GddLogReqEnable.Load()); err == nil {
		logReq = l
	}
	if logReq {
		srv.middlewares = append(srv.middlewares, log)
	}
	srv.middlewares = append(srv.middlewares,
		requestid.RequestIDHandler,
		handlers.ProxyHeaders,
		rest,
	)
	return srv
}

// AddRoute adds routes to router
func (srv *DefaultHttpSrv) AddRoute(route ...model.Route) {
	srv.bizRoutes = append(srv.bizRoutes, route...)
}

func (srv *DefaultHttpSrv) printRoutes() {
	logger.Infoln("================ Registered Routes ================")
	data := [][]string{}
	rr := config.DefaultGddRouteRootPath
	if stringutils.IsNotEmpty(config.GddRouteRootPath.Load()) {
		rr = config.GddRouteRootPath.Load()
	}
	var all []model.Route
	all = append(all, srv.bizRoutes...)
	all = append(all, srv.gddRoutes...)
	for _, r := range all {
		if strings.HasPrefix(r.Pattern, gddPathPrefix) {
			data = append(data, []string{r.Name, r.Method, r.Pattern})
		} else {
			data = append(data, []string{r.Name, r.Method, path.Clean(rr + r.Pattern)})
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
		logger.Infoln(row)
	}
	logger.Infoln("===================================================")
}

// AddMiddleware adds middlewares to the end of chain
func (srv *DefaultHttpSrv) AddMiddleware(mwf ...func(http.Handler) http.Handler) {
	for _, item := range mwf {
		srv.middlewares = append(srv.middlewares, item)
	}
}

// PreMiddleware adds middlewares to the head of chain
func (srv *DefaultHttpSrv) PreMiddleware(mwf ...func(http.Handler) http.Handler) {
	var middlewares []mux.MiddlewareFunc
	for _, item := range mwf {
		middlewares = append(middlewares, item)
	}
	srv.middlewares = append(middlewares, srv.middlewares...)
}

func (srv *DefaultHttpSrv) newHttpServer() *http.Server {
	write, err := time.ParseDuration(config.GddWriteTimeout.Load())
	if err != nil {
		logger.Debugf("Parse %s %s as time.Duration failed: %s, use default %s instead.\n", string(config.GddWriteTimeout),
			config.GddWriteTimeout.Load(), err.Error(), config.DefaultGddWriteTimeout)
		write, _ = time.ParseDuration(config.DefaultGddWriteTimeout)
	}

	read, err := time.ParseDuration(config.GddReadTimeout.Load())
	if err != nil {
		logger.Debugf("Parse %s %s as time.Duration failed: %s, use default %s instead.\n", string(config.GddReadTimeout),
			config.GddReadTimeout.Load(), err.Error(), config.DefaultGddReadTimeout)
		read, _ = time.ParseDuration(config.DefaultGddReadTimeout)
	}

	idle, err := time.ParseDuration(config.GddIdleTimeout.Load())
	if err != nil {
		logger.Debugf("Parse %s %s as time.Duration failed: %s, use default %s instead.\n", string(config.GddIdleTimeout),
			config.GddIdleTimeout.Load(), err.Error(), config.DefaultGddIdleTimeout)
		idle, _ = time.ParseDuration(config.DefaultGddIdleTimeout)
	}

	httpPort := strconv.Itoa(config.DefaultGddPort)
	if _, err = cast.ToIntE(config.GddPort.Load()); err == nil {
		httpPort = config.GddPort.Load()
	}
	httpHost := config.DefaultGddHost
	if stringutils.IsNotEmpty(config.GddHost.Load()) {
		httpHost = config.GddHost.Load()
	}
	httpServer := &http.Server{
		Addr: strings.Join([]string{httpHost, httpPort}, ":"),
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: write,
		ReadTimeout:  read,
		IdleTimeout:  idle,
		Handler:      srv.rootRouter, // Pass our instance of gorilla/mux in.
	}

	// Run our server in a goroutine so that it doesn't block.
	go func() {
		logger.Infof("Http server is listening on %s\n", httpServer.Addr)
		if err := httpServer.ListenAndServe(); err != nil {
			logger.Println(err)
		}
	}()

	return httpServer
}

// Run runs http server
func (srv *DefaultHttpSrv) Run() {
	manage := config.DefaultGddManage
	if m, err := cast.ToBoolE(config.GddManage.Load()); err == nil {
		manage = m
	}
	if manage {
		srv.middlewares = append([]mux.MiddlewareFunc{prometheus.PrometheusMiddleware}, srv.middlewares...)
		gddRouter := srv.rootRouter.PathPrefix(gddPathPrefix).Subrouter().StrictSlash(true)
		corsOpts := cors.New(cors.Options{
			AllowedMethods: []string{
				http.MethodGet,
				http.MethodPost,
				http.MethodPut,
				http.MethodPatch,
				http.MethodDelete,
				http.MethodOptions,
				http.MethodHead,
			},

			AllowedHeaders: []string{
				"*",
			},

			AllowOriginRequestFunc: func(r *http.Request, origin string) bool {
				if r.URL.Path == fmt.Sprintf("%sopenapi.json", gddPathPrefix) {
					return true
				}
				return false
			},
		})
		gddRouter.Use(corsOpts.Handler)
		gddRouter.Use(basicAuth)
		srv.gddRoutes = append(srv.gddRoutes, onlinedoc.Routes()...)
		srv.gddRoutes = append(srv.gddRoutes, prometheus.Routes()...)
		srv.gddRoutes = append(srv.gddRoutes, registry.Routes()...)
		srv.gddRoutes = append(srv.gddRoutes, configui.Routes()...)
		for _, item := range srv.gddRoutes {
			gddRouter.
				Methods(item.Method, http.MethodOptions).
				Path("/" + strings.TrimPrefix(item.Pattern, gddPathPrefix)).
				Name(item.Name).
				Handler(item.HandlerFunc)
		}
	}
	srv.middlewares = append(srv.middlewares, recovery)
	srv.Use(srv.middlewares...)
	for _, item := range srv.bizRoutes {
		srv.
			Methods(item.Method, http.MethodOptions).
			Path(item.Pattern).
			Name(item.Name).
			Handler(item.HandlerFunc)
	}

	start := time.Now()
	banner := config.DefaultGddBanner
	if b, err := cast.ToBoolE(config.GddBanner.Load()); err == nil {
		banner = b
	}
	if banner {
		bannerText := config.DefaultGddBannerText
		if stringutils.IsNotEmpty(config.GddBannerText.Load()) {
			bannerText = config.GddBannerText.Load()
		}
		figure.NewColorFigure(bannerText, "doom", "green", true).Print()
	}

	srv.printRoutes()
	httpServer := srv.newHttpServer()
	defer func() {
		logger.Infoln("http server is shutting down...")

		// Create a deadline to wait for.
		grace, err := time.ParseDuration(config.GddGraceTimeout.Load())
		if err != nil {
			logger.Debugf("Parse %s %s as time.Duration failed: %s, use default %s instead.\n", string(config.GddGraceTimeout),
				config.GddGraceTimeout.Load(), err.Error(), config.DefaultGddGraceTimeout)
			grace, _ = time.ParseDuration(config.DefaultGddGraceTimeout)
		}

		ctx, cancel := context.WithTimeout(context.Background(), grace)
		defer cancel()
		// Doesn't block if no connections, but will otherwise wait
		// until the timeout deadline.
		httpServer.Shutdown(ctx)
	}()

	logger.Infof("Started in %s\n", time.Since(start))

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c
}
