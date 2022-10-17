package rest

import (
	"context"
	"fmt"
	"github.com/arl/statsviz"
	"github.com/ascarter/requestid"
	"github.com/gorilla/handlers"
	"github.com/klauspost/compress/gzhttp"
	"github.com/olekukonko/tablewriter"
	"github.com/rs/cors"
	"github.com/unionj-cloud/go-doudou/v2/framework/internal/banner"
	"github.com/unionj-cloud/go-doudou/v2/framework/internal/config"
	register "github.com/unionj-cloud/go-doudou/v2/framework/registry"
	"github.com/unionj-cloud/go-doudou/v2/framework/rest/httprouter"
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

type MiddlewareFunc func(http.Handler) http.Handler

// Middleware allows MiddlewareFunc to implement the middleware interface.
func (mw MiddlewareFunc) Middleware(handler http.Handler) http.Handler {
	return mw(handler)
}

var startAt time.Time

func init() {
	startAt = time.Now()
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

// MiddlewareFunc is alias for func(http.Handler) http.Handler
func toMiddlewareFunc(m func(http.Handler) http.HandlerFunc) MiddlewareFunc {
	return func(handler http.Handler) http.Handler {
		return m(handler)
	}
}

// RestServer wraps httpRouter router
type RestServer struct {
	bizRouter   *httprouter.RouteGroup
	rootRouter  *httprouter.Router
	gddRoutes   []Route
	debugRoutes []Route
	bizRoutes   []Route
	middlewares []MiddlewareFunc
	data        map[string]interface{}
}

func (srv *RestServer) printRoutes() {
	if !config.CheckDev() {
		return
	}
	logger.Info().Msg("================ Registered Routes ================")
	data := [][]string{}
	rr := config.DefaultGddRouteRootPath
	if stringutils.IsNotEmpty(config.GddRouteRootPath.Load()) {
		rr = config.GddRouteRootPath.Load()
	}
	var all []Route
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

// NewRestServer create a RestServer instance
func NewRestServer(data ...map[string]interface{}) *RestServer {
	rr := config.DefaultGddRouteRootPath
	if stringutils.IsNotEmpty(config.GddRouteRootPath.Load()) {
		rr = config.GddRouteRootPath.Load()
	}
	if stringutils.IsEmpty(rr) {
		rr = "/"
	}
	rootRouter := httprouter.New()
	rootRouter.SaveMatchedRoutePath = cast.ToBoolOrDefault(config.GddRouterSaveMatchedRoutePath.Load(), config.DefaultGddRouterSaveMatchedRoutePath)
	srv := &RestServer{
		bizRouter:  rootRouter.NewGroup(rr),
		rootRouter: rootRouter,
	}
	srv.middlewares = append(srv.middlewares,
		tracing,
		metrics,
	)
	if cast.ToBoolOrDefault(config.GddEnableResponseGzip.Load(), config.DefaultGddEnableResponseGzip) {
		gzipMiddleware, err := gzhttp.NewWrapper(gzhttp.ContentTypes(contentTypeShouldbeGzip))
		if err != nil {
			panic(err)
		}
		srv.middlewares = append(srv.middlewares, toMiddlewareFunc(gzipMiddleware))
	}
	if cast.ToBoolOrDefault(config.GddLogReqEnable.Load(), config.DefaultGddLogReqEnable) {
		srv.middlewares = append(srv.middlewares, log)
	}
	srv.middlewares = append(srv.middlewares,
		requestid.RequestIDHandler,
		handlers.ProxyHeaders,
		fallbackContentType(config.GddFallbackContentType.LoadOrDefault(config.DefaultGddFallbackContentType)),
	)
	if len(data) > 0 {
		srv.data = data[0]
	}
	return srv
}

// AddRoute adds routes to router
func (srv *RestServer) AddRoute(route ...Route) {
	srv.bizRoutes = append(srv.bizRoutes, route...)
}

// AddMiddleware adds middlewares to the end of chain
func (srv *RestServer) AddMiddleware(mwf ...func(http.Handler) http.Handler) {
	for _, item := range mwf {
		srv.middlewares = append(srv.middlewares, item)
	}
}

// PreMiddleware adds middlewares to the head of chain
func (srv *RestServer) PreMiddleware(mwf ...func(http.Handler) http.Handler) {
	var middlewares []MiddlewareFunc
	for _, item := range mwf {
		middlewares = append(middlewares, item)
	}
	srv.middlewares = append(middlewares, srv.middlewares...)
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
		Handler:      srv.rootRouter, // Pass our instance of httprouter.Router in.
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
		srv.middlewares = append([]MiddlewareFunc{PrometheusMiddleware}, srv.middlewares...)
		gddRouter := srv.rootRouter.NewGroup(gddPathPrefix)
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
		basicAuthMiddle := MiddlewareFunc(basicAuth())
		gddmiddlewares := []MiddlewareFunc{metrics, corsOpts.Handler, basicAuthMiddle}
		srv.gddRoutes = append(srv.gddRoutes, docRoutes()...)
		srv.gddRoutes = append(srv.gddRoutes, promRoutes()...)
		srv.gddRoutes = append(srv.gddRoutes, configRoutes()...)
		freq, err := time.ParseDuration(config.GddStatsFreq.Load())
		if err != nil {
			logger.Debug().Msgf("Parse %s %s as time.Duration failed: %s, use default %s instead.\n", string(config.GddStatsFreq),
				config.GddStatsFreq.Load(), err.Error(), config.DefaultGddStatsFreq)
			freq, _ = time.ParseDuration(config.DefaultGddStatsFreq)
		}
		_ = freq
		srv.gddRoutes = append(srv.gddRoutes, []Route{
			{
				Name:    "GetStatsvizWs",
				Method:  http.MethodGet,
				Pattern: gddPathPrefix + "statsviz/ws",
			},
			{
				Name:    "GetStatsviz",
				Method:  http.MethodGet,
				Pattern: gddPathPrefix + "statsviz/*",
				HandlerFunc: func(writer http.ResponseWriter, request *http.Request) {
					if strings.HasSuffix(request.URL.Path, "/ws") {
						statsviz.Ws(writer, request)
						return
					}
					statsviz.IndexAtRoot(gddPathPrefix+"statsviz/").ServeHTTP(writer, request)
				},
			},
		}...)
		for _, item := range srv.gddRoutes {
			if item.HandlerFunc == nil {
				continue
			}
			h := http.Handler(item.HandlerFunc)
			for i := len(gddmiddlewares) - 1; i >= 0; i-- {
				h = gddmiddlewares[i].Middleware(h)
			}
			gddRouter.Handler(item.Method, "/"+strings.TrimPrefix(item.Pattern, gddPathPrefix), h, item.Name)
		}
		srv.debugRoutes = append(srv.debugRoutes, []Route{
			{
				Name:    "GetDebugPprofCmdline",
				Method:  http.MethodGet,
				Pattern: debugPathPrefix + "pprof/cmdline",
			},
			{
				Name:    "GetDebugPprofProfile",
				Method:  http.MethodGet,
				Pattern: debugPathPrefix + "pprof/profile",
			},
			{
				Name:    "GetDebugPprofSymbol",
				Method:  http.MethodGet,
				Pattern: debugPathPrefix + "pprof/symbol",
			},
			{
				Name:    "GetDebugPprofTrace",
				Method:  http.MethodGet,
				Pattern: debugPathPrefix + "pprof/trace",
			},
			{
				Name:    "GetDebugPprofIndex",
				Method:  http.MethodGet,
				Pattern: debugPathPrefix + "pprof/*",
				HandlerFunc: func(writer http.ResponseWriter, request *http.Request) {
					lastSegment := request.URL.Path[strings.LastIndex(request.URL.Path, "/"):]
					switch lastSegment {
					case "/cmdline":
						pprof.Cmdline(writer, request)
						return
					case "/profile":
						pprof.Profile(writer, request)
						return
					case "/symbol":
						pprof.Symbol(writer, request)
						return
					case "/trace":
						pprof.Trace(writer, request)
						return
					}
					pprof.Index(writer, request)
				},
			},
		}...)
		debugRouter := srv.rootRouter.NewGroup(debugPathPrefix)
		for _, item := range srv.debugRoutes {
			if item.HandlerFunc == nil {
				continue
			}
			h := http.Handler(item.HandlerFunc)
			for i := len(gddmiddlewares) - 1; i >= 0; i-- {
				h = gddmiddlewares[i].Middleware(h)
			}
			debugRouter.Handler(item.Method, "/"+strings.TrimPrefix(item.Pattern, debugPathPrefix), h, item.Name)
		}
	}
	srv.middlewares = append(srv.middlewares, recovery)
	for _, item := range srv.bizRoutes {
		h := http.Handler(item.HandlerFunc)
		for i := len(srv.middlewares) - 1; i >= 0; i-- {
			h = srv.middlewares[i].Middleware(h)
		}
		srv.bizRouter.Handler(item.Method, item.Pattern, h, item.Name)
	}
	srv.rootRouter.NotFound = http.HandlerFunc(http.NotFound)
	srv.rootRouter.MethodNotAllowed = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("405 method not allowed"))
	})
	for i := len(srv.middlewares) - 1; i >= 0; i-- {
		srv.rootRouter.NotFound = srv.middlewares[i].Middleware(srv.rootRouter.NotFound)
		srv.rootRouter.MethodNotAllowed = srv.middlewares[i].Middleware(srv.rootRouter.MethodNotAllowed)
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
