package gorilla

import (
	"context"
	"fmt"
	"github.com/arl/statsviz"
	"github.com/ascarter/requestid"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/klauspost/compress/gzhttp"
	"github.com/olekukonko/tablewriter"
	"github.com/rs/cors"
	"github.com/unionj-cloud/go-doudou/v2/framework"
	"github.com/unionj-cloud/go-doudou/v2/framework/banner"
	"github.com/unionj-cloud/go-doudou/v2/framework/config"
	register "github.com/unionj-cloud/go-doudou/v2/framework/registry"
	"github.com/unionj-cloud/go-doudou/v2/framework/registry/constants"
	"github.com/unionj-cloud/go-doudou/v2/framework/rest"
	"github.com/unionj-cloud/go-doudou/v2/toolkit/cast"
	"github.com/unionj-cloud/go-doudou/v2/toolkit/stringutils"
	logger "github.com/unionj-cloud/go-doudou/v2/toolkit/zlogger"
	"net/http"
	"net/http/pprof"
	"os"
	"os/signal"
	"path"
	"strconv"
	"strings"
	"time"
)

var startAt time.Time

func init() {
	startAt = time.Now()
}

type common struct {
	gddRoutes   []rest.Route
	debugRoutes []rest.Route
	bizRoutes   []rest.Route
	Middlewares []mux.MiddlewareFunc
	data        map[string]interface{}
}

// RestServer wraps gorilla mux router
type RestServer struct {
	*mux.Router
	rootRouter *mux.Router
	common
}

const gddPathPrefix = "/go-doudou/"
const debugPathPrefix = "/debug/"

var contentTypeShouldbeGzip []string

func init() {
	contentTypeShouldbeGzip = []string{
		"text/html",
		"text/css",
		"text/plain",
		"text/xml",
		"text/x-component",
		"text/javascript",
		"application/x-javascript",
		"application/javascript",
		"application/json",
		"application/manifest+json",
		"application/vnd.api+json",
		"application/xml",
		"application/xhtml+xml",
		"application/rss+xml",
		"application/atom+xml",
		"application/vnd.ms-fontobject",
		"application/x-font-ttf",
		"application/x-font-opentype",
		"application/x-font-truetype",
		"image/svg+xml",
		"image/x-icon",
		"image/vnd.microsoft.icon",
		"font/ttf",
		"font/eot",
		"font/otf",
		"font/opentype",
	}
}

// mux.MiddlewareFunc is alias for func(http.Handler) http.Handler
func toMiddlewareFunc(m func(http.Handler) http.HandlerFunc) mux.MiddlewareFunc {
	return func(handler http.Handler) http.Handler {
		return m(handler)
	}
}

// NewRestServer create a RestServer instance
func NewRestServer(data ...map[string]interface{}) *RestServer {
	rr := config.DefaultGddRouteRootPath
	if stringutils.IsNotEmpty(config.GddRouteRootPath.Load()) {
		rr = config.GddRouteRootPath.Load()
	}
	rootRouter := mux.NewRouter().StrictSlash(true)
	srv := &RestServer{
		Router:     rootRouter.PathPrefix(rr).Subrouter().StrictSlash(true),
		rootRouter: rootRouter,
	}
	srv.Middlewares = append(srv.Middlewares,
		rest.Tracing,
		rest.Metrics,
	)
	if cast.ToBoolOrDefault(config.GddEnableResponseGzip.Load(), config.DefaultGddEnableResponseGzip) {
		gzipMiddleware, err := gzhttp.NewWrapper(gzhttp.ContentTypes(contentTypeShouldbeGzip))
		if err != nil {
			panic(err)
		}
		srv.Middlewares = append(srv.Middlewares, toMiddlewareFunc(gzipMiddleware))
	}
	if cast.ToBoolOrDefault(config.GddLogReqEnable.Load(), config.DefaultGddLogReqEnable) {
		srv.Middlewares = append(srv.Middlewares, rest.Log)
	}
	srv.Middlewares = append(srv.Middlewares,
		requestid.RequestIDHandler,
		handlers.ProxyHeaders,
		rest.FallbackContentType(config.GddFallbackContentType.LoadOrDefault(config.DefaultGddFallbackContentType)),
	)
	if len(data) > 0 {
		srv.data = data[0]
	}
	return srv
}

// AddRoute adds routes to router
func (srv *RestServer) AddRoute(route ...rest.Route) {
	srv.bizRoutes = append(srv.bizRoutes, route...)
}

func (srv *common) printRoutes() {
	if !framework.CheckDev() {
		return
	}
	logger.Info().Msg("================ Registered Routes ================")
	data := [][]string{}
	rr := config.DefaultGddRouteRootPath
	if stringutils.IsNotEmpty(config.GddRouteRootPath.Load()) {
		rr = config.GddRouteRootPath.Load()
	}
	var all []rest.Route
	all = append(all, srv.bizRoutes...)
	all = append(all, srv.gddRoutes...)
	all = append(all, srv.debugRoutes...)
	for _, r := range all {
		if strings.HasPrefix(r.Pattern, gddPathPrefix) || strings.HasPrefix(r.Pattern, debugPathPrefix) {
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
		logger.Info().Msg(row)
	}
	logger.Info().Msg("===================================================")
}

// AddMiddleware adds middlewares to the end of chain
func (srv *RestServer) AddMiddleware(mwf ...func(http.Handler) http.Handler) {
	for _, item := range mwf {
		srv.Middlewares = append(srv.Middlewares, item)
	}
}

// PreMiddleware adds middlewares to the head of chain
func (srv *RestServer) PreMiddleware(mwf ...func(http.Handler) http.Handler) {
	var middlewares []mux.MiddlewareFunc
	for _, item := range mwf {
		middlewares = append(middlewares, item)
	}
	srv.Middlewares = append(middlewares, srv.Middlewares...)
}

func (srv *RestServer) newHttpServer() *http.Server {
	write, err := time.ParseDuration(config.GddWriteTimeout.Load())
	if err != nil {
		logger.Debug().Msgf("Parse %s %s as time.Duration failed: %s, use default %s instead.\n", string(config.GddWriteTimeout),
			config.GddWriteTimeout.Load(), err.Error(), config.DefaultGddWriteTimeout)
		write, _ = time.ParseDuration(config.DefaultGddWriteTimeout)
	}

	read, err := time.ParseDuration(config.GddReadTimeout.Load())
	if err != nil {
		logger.Debug().Msgf("Parse %s %s as time.Duration failed: %s, use default %s instead.\n", string(config.GddReadTimeout),
			config.GddReadTimeout.Load(), err.Error(), config.DefaultGddReadTimeout)
		read, _ = time.ParseDuration(config.DefaultGddReadTimeout)
	}

	idle, err := time.ParseDuration(config.GddIdleTimeout.Load())
	if err != nil {
		logger.Debug().Msgf("Parse %s %s as time.Duration failed: %s, use default %s instead.\n", string(config.GddIdleTimeout),
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
		logger.Info().Msgf("Http server is listening at %v", httpServer.Addr)
		logger.Info().Msgf("Http server started in %s", time.Since(startAt))
		if err := httpServer.ListenAndServe(); err != nil {
			logger.Error().Err(err).Msg("")
		}
	}()

	return httpServer
}

// Run runs http server
func (srv *RestServer) Run() {
	banner.Print()
	register.NewRest(srv.data)
	manage := cast.ToBoolOrDefault(config.GddManage.Load(), config.DefaultGddManage)
	if manage {
		srv.Middlewares = append([]mux.MiddlewareFunc{rest.PrometheusMiddleware}, srv.Middlewares...)
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
		gddRouter.Use(rest.Metrics)
		gddRouter.Use(corsOpts.Handler)
		gddRouter.Use(rest.BasicAuth())
		srv.gddRoutes = append(srv.gddRoutes, rest.DocRoutes()...)
		srv.gddRoutes = append(srv.gddRoutes, rest.PromRoutes()...)
		srv.gddRoutes = append(srv.gddRoutes, rest.ConfigRoutes()...)
		if _, ok := config.ServiceDiscoveryMap()[constants.SD_MEMBERLIST]; ok {
			srv.gddRoutes = append(srv.gddRoutes, rest.MemberlistUIRoutes()...)
		}
		for _, item := range srv.gddRoutes {
			gddRouter.
				Methods(item.Method, http.MethodOptions).
				Path("/" + strings.TrimPrefix(item.Pattern, gddPathPrefix)).
				Name(item.Name).
				Handler(item.HandlerFunc)
		}
		freq, err := time.ParseDuration(config.GddStatsFreq.Load())
		if err != nil {
			logger.Debug().Msgf("Parse %s %s as time.Duration failed: %s, use default %s instead.\n", string(config.GddStatsFreq),
				config.GddStatsFreq.Load(), err.Error(), config.DefaultGddStatsFreq)
			freq, _ = time.ParseDuration(config.DefaultGddStatsFreq)
		}
		srv.gddRoutes = append(srv.gddRoutes, []rest.Route{
			{
				Name:    "GetStatsvizWs",
				Method:  "GET",
				Pattern: gddPathPrefix + "statsviz/ws",
			},
			{
				Name:    "GetStatsviz",
				Method:  "GET",
				Pattern: gddPathPrefix + "statsviz/",
			},
		}...)
		gddRouter.
			Methods(http.MethodGet).
			Path("/statsviz/ws").
			Name("GetStatsvizWs").
			HandlerFunc(statsviz.NewWsHandler(freq))
		gddRouter.
			Methods(http.MethodGet).
			PathPrefix("/statsviz/").
			Name("GetStatsviz").
			Handler(statsviz.IndexAtRoot(gddPathPrefix + "statsviz/"))
		srv.debugRoutes = append(srv.debugRoutes, []rest.Route{
			{
				Name:    "GetDebugPprofCmdline",
				Method:  "GET",
				Pattern: debugPathPrefix + "pprof/cmdline",
			},
			{
				Name:    "GetDebugPprofProfile",
				Method:  "GET",
				Pattern: debugPathPrefix + "pprof/profile",
			},
			{
				Name:    "GetDebugPprofSymbol",
				Method:  "GET",
				Pattern: debugPathPrefix + "pprof/symbol",
			},
			{
				Name:    "GetDebugPprofTrace",
				Method:  "GET",
				Pattern: debugPathPrefix + "pprof/trace",
			},
			{
				Name:    "GetDebugPprofIndex",
				Method:  "GET",
				Pattern: debugPathPrefix + "pprof/",
			},
		}...)
		debugRouter := srv.rootRouter.PathPrefix(debugPathPrefix).Subrouter().StrictSlash(true)
		debugRouter.Use(rest.Metrics)
		debugRouter.Use(corsOpts.Handler)
		debugRouter.Use(rest.BasicAuth())
		debugRouter.Methods(http.MethodGet).Path("/pprof/cmdline").Name("GetDebugPprofCmdline").HandlerFunc(pprof.Cmdline)
		debugRouter.Methods(http.MethodGet).Path("/pprof/profile").Name("GetDebugPprofProfile").HandlerFunc(pprof.Profile)
		debugRouter.Methods(http.MethodGet).Path("/pprof/symbol").Name("GetDebugPprofSymbol").HandlerFunc(pprof.Symbol)
		debugRouter.Methods(http.MethodGet).Path("/pprof/trace").Name("GetDebugPprofTrace").HandlerFunc(pprof.Trace)
		debugRouter.Methods(http.MethodGet).PathPrefix("/pprof/").Name("GetDebugPprofIndex").HandlerFunc(pprof.Index)
	}
	srv.Middlewares = append(srv.Middlewares, rest.Recovery)
	srv.Use(srv.Middlewares...)
	for _, item := range srv.bizRoutes {
		srv.
			Methods(item.Method, http.MethodOptions).
			Path(item.Pattern).
			Name(item.Name).
			Handler(item.HandlerFunc)
	}
	srv.rootRouter.NotFoundHandler = srv.rootRouter.NewRoute().BuildOnly().HandlerFunc(http.NotFound).GetHandler()
	srv.rootRouter.MethodNotAllowedHandler = srv.rootRouter.NewRoute().BuildOnly().HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("405 method not allowed"))
	}).GetHandler()
	for i := len(srv.Middlewares) - 1; i >= 0; i-- {
		srv.rootRouter.NotFoundHandler = srv.Middlewares[i].Middleware(srv.rootRouter.NotFoundHandler)
		srv.rootRouter.MethodNotAllowedHandler = srv.Middlewares[i].Middleware(srv.rootRouter.MethodNotAllowedHandler)
	}
	srv.printRoutes()
	httpServer := srv.newHttpServer()
	defer func() {
		register.ShutdownRest()
		grace, err := time.ParseDuration(config.GddGraceTimeout.Load())
		if err != nil {
			logger.Debug().Msgf("Parse %s %s as time.Duration failed: %s, use default %s instead.\n", string(config.GddGraceTimeout),
				config.GddGraceTimeout.Load(), err.Error(), config.DefaultGddGraceTimeout)
			grace, _ = time.ParseDuration(config.DefaultGddGraceTimeout)
		}
		logger.Info().Msgf("Http server is gracefully shutting down in %s", grace)

		ctx, cancel := context.WithTimeout(context.Background(), grace)
		defer cancel()
		// Doesn't block if no connections, but will otherwise wait
		// until the timeout deadline.
		httpServer.Shutdown(ctx)
	}()

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c
}
