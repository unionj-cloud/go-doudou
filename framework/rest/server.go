package rest

import (
	"context"
	"github.com/arl/statsviz"
	"github.com/ascarter/requestid"
	mapset "github.com/deckarep/golang-set/v2"
	"github.com/gorilla/handlers"
	"github.com/klauspost/compress/gzhttp"
	"github.com/olekukonko/tablewriter"
	"github.com/samber/lo"
	"github.com/unionj-cloud/go-doudou/v2/framework"
	"github.com/unionj-cloud/go-doudou/v2/framework/config"
	register "github.com/unionj-cloud/go-doudou/v2/framework/registry"
	"github.com/unionj-cloud/go-doudou/v2/framework/registry/constants"
	"github.com/unionj-cloud/go-doudou/v2/framework/rest/httprouter"
	"github.com/unionj-cloud/go-doudou/v2/toolkit/stringutils"
	logger "github.com/unionj-cloud/go-doudou/v2/toolkit/zlogger"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path"
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
	bizRouter    *httprouter.RouteGroup
	rootRouter   *httprouter.Router
	gddRoutes    []Route
	debugRoutes  []Route
	bizRoutes    []Route
	middlewares  []MiddlewareFunc
	data         map[string]interface{}
	panicHandler func(inner http.Handler) http.Handler
	*http.Server
}

func path2key(method, path string) string {
	var sb strings.Builder
	sb.WriteString(method)
	sb.WriteString(":")
	sb.WriteString(path)
	return sb.String()
}

func (srv *RestServer) printRoutes() {
	if !config.CheckDev() {
		return
	}
	logger.Info().Msg("================ Registered Routes ================")
	data := [][]string{}
	var all []Route
	all = append(all, srv.bizRoutes...)
	all = append(all, srv.gddRoutes...)
	all = append(all, srv.debugRoutes...)
	routes := make(map[string]struct{})
	for _, r := range all {
		key := path2key(r.Method, r.Pattern)
		_, ok := routes[key]
		if !ok {
			routes[key] = struct{}{}
			data = append(data, []string{r.Name, r.Method, path.Clean(r.Pattern)})
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
	var options []ServerOption
	if len(data) > 0 {
		options = append(options, WithUserData(data[0]))
	}
	return NewRestServerWithOptions(options...)
}

type ServerOption func(server *RestServer)

func WithPanicHandler(panicHandler func(inner http.Handler) http.Handler) ServerOption {
	return func(server *RestServer) {
		server.panicHandler = panicHandler
	}
}

func WithUserData(userData map[string]interface{}) ServerOption {
	return func(server *RestServer) {
		server.data = userData
	}
}

// NewRestServerWithOptions create a RestServer instance with options
func NewRestServerWithOptions(options ...ServerOption) *RestServer {
	rootRouter := httprouter.New()
	rootRouter.SaveMatchedRoutePath = true
	srv := &RestServer{
		bizRouter:    rootRouter.NewGroup(config.GddConfig.RouteRootPath),
		rootRouter:   rootRouter,
		panicHandler: recovery,
		Server: &http.Server{
			// Good practice to set timeouts to avoid Slowloris attacks.
			WriteTimeout: config.GddConfig.WriteTimeout,
			ReadTimeout:  config.GddConfig.ReadTimeout,
			IdleTimeout:  config.GddConfig.IdleTimeout,
			Handler:      rootRouter, // Pass our instance of httprouter.Router in.
		},
	}
	for _, fn := range options {
		fn(srv)
	}
	srv.middlewares = append(srv.middlewares,
		srv.panicHandler,
		tracing,
		metrics,
		gzipBody,
	)
	if config.GddConfig.EnableResponseGzip {
		gzipMiddleware, err := gzhttp.NewWrapper(gzhttp.ContentTypes(contentTypeShouldbeGzip))
		if err != nil {
			panic(err)
		}
		srv.middlewares = append(srv.middlewares, toMiddlewareFunc(gzipMiddleware))
	}
	if config.GddConfig.LogReqEnable {
		srv.middlewares = append(srv.middlewares, log)
	}
	srv.middlewares = append(srv.middlewares,
		requestid.RequestIDHandler,
		handlers.ProxyHeaders,
	)
	if config.GddConfig.ManageEnable {
		srv.middlewares = append([]MiddlewareFunc{PrometheusMiddleware}, srv.middlewares...)
		basicAuthMiddle := MiddlewareFunc(basicAuth())
		gddmiddlewares := []MiddlewareFunc{metrics, basicAuthMiddle}
		srv.gddRoutes = append(srv.gddRoutes, promRoutes()...)
		srv.gddRoutes = append(srv.gddRoutes, configRoutes()...)
		if _, ok := config.ServiceDiscoveryMap()[constants.SD_MEMBERLIST]; ok {
			srv.gddRoutes = append(srv.gddRoutes, MemberlistUIRoutes()...)
		}
		freq, err := time.ParseDuration(config.GddStatsFreq.Load())
		if err != nil {
			logger.Debug().Msgf("Parse %s %s as time.Duration failed: %s, use default %s instead.\n", string(config.GddStatsFreq),
				config.GddStatsFreq.Load(), err.Error(), config.DefaultGddStatsFreq)
			freq, _ = time.ParseDuration(config.DefaultGddStatsFreq)
		}
		statsvizServer, _ := statsviz.NewServer(statsviz.Root(srv.bizRouter.SubPath(gddPathPrefix+"statsviz/")), statsviz.SendFrequency(freq))
		ws := statsvizServer.Ws()
		index := statsvizServer.Index()
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
						ws(writer, request)
						return
					}
					index(writer, request)
				},
			},
		}...)
		for k, item := range srv.gddRoutes {
			if item.HandlerFunc != nil {
				h := http.Handler(item.HandlerFunc)
				for i := len(gddmiddlewares) - 1; i >= 0; i-- {
					h = gddmiddlewares[i].Middleware(h)
				}
				srv.bizRouter.Handler(item.Method, item.Pattern, h, item.Name)
			}
			item.Pattern = srv.bizRouter.SubPath(item.Pattern)
			srv.gddRoutes[k] = item
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
						Cmdline(writer, request)
						return
					case "/profile":
						Profile(writer, request)
						return
					case "/symbol":
						Symbol(writer, request)
						return
					case "/trace":
						Trace(writer, request)
						return
					}
					Index(writer, request)
				},
			},
		}...)
		for k, item := range srv.debugRoutes {
			if item.HandlerFunc != nil {
				h := http.Handler(item.HandlerFunc)
				for i := len(gddmiddlewares) - 1; i >= 0; i-- {
					h = gddmiddlewares[i].Middleware(h)
				}
				srv.bizRouter.Handler(item.Method, item.Pattern, h, item.Name)
			}
			item.Pattern = srv.bizRouter.SubPath(item.Pattern)
			srv.debugRoutes[k] = item
		}
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
	return srv
}

func (srv *RestServer) groupRoutes(routeGroup *httprouter.RouteGroup, routes []Route, mwf ...MiddlewareFunc) {
	for _, item := range routes {
		h := http.Handler(item.HandlerFunc)
		for i := len(mwf) - 1; i >= 0; i-- {
			h = mwf[i].Middleware(h)
		}
		routeGroup.Handler(item.Method, item.Pattern, h, item.Name)
		item.Pattern = routeGroup.SubPath(item.Pattern)
		srv.bizRoutes = append(srv.bizRoutes, item)
	}
}

// AddRoute adds routes to router
func (srv *RestServer) AddRoute(route ...Route) {
	srv.groupRoutes(srv.bizRouter, route, srv.middlewares...)
}

// AddRoutes adds routes to router
func (srv *RestServer) AddRoutes(routes []Route, mwf ...func(http.Handler) http.Handler) {
	m := make([]MiddlewareFunc, 0)
	m = append(m, srv.middlewares...)
	m = append(m, lo.Map(mwf, func(item func(http.Handler) http.Handler, index int) MiddlewareFunc {
		return item
	})...)
	srv.groupRoutes(srv.bizRouter, routes, m...)
}

// GroupRoutes adds routes to router
func (srv *RestServer) GroupRoutes(group string, routes []Route, mwf ...func(http.Handler) http.Handler) {
	m := make([]MiddlewareFunc, 0)
	m = append(m, srv.middlewares...)
	m = append(m, lo.Map(mwf, func(item func(http.Handler) http.Handler, index int) MiddlewareFunc {
		return item
	})...)
	srv.groupRoutes(srv.bizRouter.NewGroup(group), routes, m...)
}

// AddMiddleware adds middlewares to the end of chain
// Deprecated: use Use instead
func (srv *RestServer) AddMiddleware(mwf ...func(http.Handler) http.Handler) {
	for _, item := range mwf {
		srv.middlewares = append(srv.middlewares, item)
	}
}

// Use adds middlewares to the end of chain
func (srv *RestServer) Use(mwf ...func(http.Handler) http.Handler) {
	for _, item := range mwf {
		srv.middlewares = append(srv.middlewares, item)
	}
}

// Run runs http server
func (srv *RestServer) Run() {
	ln, err := net.Listen("tcp", strings.Join([]string{config.GddConfig.Host, config.GddConfig.Port}, ":"))
	if err != nil {
		logger.Panic().Msg(err.Error())
	}
	srv.Serve(ln)
	defer func() {
		logger.Info().Msgf("Grpc server is gracefully shutting down in %s", config.GddConfig.GraceTimeout)
		// Make sure to set a deadline on exiting the process
		// after upg.Exit() is closed. No new upgrades can be
		// performed if the parent doesn't exit.
		time.AfterFunc(config.GddConfig.GraceTimeout, func() {
			logger.Error().Msg("Graceful shutdown timed out")
			os.Exit(1)
		})
		register.ShutdownRest()
		srv.Shutdown(context.Background())
	}()

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c
}

func (srv *RestServer) Serve(ln net.Listener) {
	var all []Route
	all = append(all, srv.bizRoutes...)
	all = append(all, srv.gddRoutes...)
	all = append(all, srv.debugRoutes...)
	docs := mapset.NewSet[DocItem]()
	for _, r := range all {
		if strings.Contains(r.Pattern, gddPathPrefix+"doc") {
			resources := strings.Split(strings.TrimPrefix(strings.TrimSuffix(r.Pattern, gddPathPrefix+"doc"), "/"), "/")
			if len(resources) > 0 {
				module := resources[len(resources)-1]
				if stringutils.IsNotEmpty(module) {
					docs.Add(DocItem{
						Label: module,
						Value: r.Pattern,
					})
				}
			}
		}
	}
	Docs = docs.ToSlice()
	framework.PrintBanner()
	framework.PrintLock.Lock()
	register.NewRest(srv.data)
	srv.printRoutes()
	// Run our server in a goroutine so that it doesn't block.
	go func() {
		if err := srv.Server.Serve(ln); err != http.ErrServerClosed {
			logger.Error().Msgf("HTTP server: %s", err.Error())
		}
	}()
	logger.Info().Msgf("Http server is listening at %v", ln.Addr().String())
	logger.Info().Msgf("Http server started in %s", time.Since(startAt))
	framework.PrintLock.Unlock()
}
