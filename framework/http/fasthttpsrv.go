package ddhttp

import (
	"fmt"
	configui "github.com/unionj-cloud/go-doudou/framework/http/config"
	"github.com/unionj-cloud/go-doudou/framework/http/fast"
	"github.com/unionj-cloud/go-doudou/framework/http/fast/cors"
	"github.com/unionj-cloud/go-doudou/framework/http/model"
	"github.com/unionj-cloud/go-doudou/framework/http/onlinedoc"
	"github.com/unionj-cloud/go-doudou/framework/http/prefork"
	"github.com/unionj-cloud/go-doudou/framework/http/prometheus"
	"github.com/unionj-cloud/go-doudou/framework/http/registry"
	"github.com/unionj-cloud/go-doudou/framework/http/router"
	"github.com/unionj-cloud/go-doudou/framework/internal/banner"
	"github.com/unionj-cloud/go-doudou/framework/internal/config"
	"github.com/unionj-cloud/go-doudou/toolkit/cast"
	"github.com/unionj-cloud/go-doudou/toolkit/stringutils"
	logger "github.com/unionj-cloud/go-doudou/toolkit/zlogger"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttpadaptor"
	"github.com/valyala/fasthttp/pprofhandler"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"time"
)

// FastHttpSrv wraps fasthttp router
type FastHttpSrv struct {
	Router     *router.Group
	rootRouter *router.Router
	common
}

// NewFastHttpSrv create a FastHttpSrv instance
func NewFastHttpSrv() *FastHttpSrv {
	rr := config.DefaultGddRouteRootPath
	if stringutils.IsNotEmpty(config.GddRouteRootPath.Load()) {
		rr = config.GddRouteRootPath.Load()
	}
	if stringutils.IsEmpty(rr) {
		rr = "/"
	}
	rootRouter := router.New()
	srv := &FastHttpSrv{
		Router:     rootRouter.Group(rr),
		rootRouter: rootRouter,
	}
	srv.FastMiddlewares = append(srv.FastMiddlewares,
		fast.Tracing,
		fast.Metrics,
	)
	if cast.ToBoolOrDefault(config.GddEnableResponseGzip.Load(), config.DefaultGddEnableResponseGzip) {
		srv.FastMiddlewares = append(srv.FastMiddlewares, fasthttp.CompressHandler)
	}
	if cast.ToBoolOrDefault(config.GddLogReqEnable.Load(), config.DefaultGddLogReqEnable) {
		srv.FastMiddlewares = append(srv.FastMiddlewares, fast.ReqLog)
	}
	srv.FastMiddlewares = append(srv.FastMiddlewares,
		fast.RequestIDHandler,
		fast.RealIp,
		fast.FallbackContentType(config.GddFallbackContentType.LoadOrDefault(config.DefaultGddFallbackContentType)),
	)
	return srv
}

// AddRoute adds routes to router
func (srv *FastHttpSrv) AddRoute(route ...model.Route) {
	srv.bizRoutes = append(srv.bizRoutes, route...)
}

// AddMiddleware adds middlewares to the end of chain
func (srv *FastHttpSrv) AddMiddleware(mwf ...func(inner fasthttp.RequestHandler) fasthttp.RequestHandler) {
	for _, item := range mwf {
		srv.FastMiddlewares = append(srv.FastMiddlewares, item)
	}
}

// PreMiddleware adds middlewares to the head of chain
func (srv *FastHttpSrv) PreMiddleware(mwf ...func(inner fasthttp.RequestHandler) fasthttp.RequestHandler) {
	var middlewares []fast.MiddlewareFunc
	for _, item := range mwf {
		middlewares = append(middlewares, item)
	}
	srv.FastMiddlewares = append(middlewares, srv.FastMiddlewares...)
}

// RootRouter returns pointer type of httprouter.Router for directly putting into http.ListenAndServe as http.Handler implementation
func (srv *FastHttpSrv) RootRouter() *router.Router {
	return srv.rootRouter
}

func (srv *FastHttpSrv) newFastHttpServer() *fasthttp.Server {
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
	addr := strings.Join([]string{httpHost, httpPort}, ":")
	httpServer := &fasthttp.Server{
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout:                       write,
		ReadTimeout:                        read,
		IdleTimeout:                        idle,
		NoDefaultServerHeader:              true,
		NoDefaultDate:                      true,
		DisableHeaderNamesNormalizing:      true,
		NoDefaultContentType:               true,
		Logger:                             &logger.Logger,
		TLSConfig:                          nil,             // TODO Plan to make it configurable with environment variable
		WriteBufferSize:                    4096,            // TODO 4096 is default value. Plan to make it configurable with environment variable
		ReadBufferSize:                     4096,            // TODO 4096 is default value. Plan to make it configurable with environment variable
		MaxRequestBodySize:                 4 * 1024 * 1024, // TODO 4 * 1024 * 1024 is default value. Plan to make it configurable with environment variable
		DisableKeepalive:                   false,           // TODO not tcp keep alive. false is default value. Plan to make it configurable with environment variable
		ReduceMemoryUsage:                  false,           // TODO false is default value. Plan to make it configurable with environment variable
		GetOnly:                            false,           // TODO false is default value. Plan to make it configurable with environment variable
		LogAllErrors:                       false,           // TODO false is default value. Plan to make it configurable with environment variable
		Concurrency:                        256 * 1024,      // TODO 256 * 1024 is default value. Plan to make it configurable with environment variable
		SleepWhenConcurrencyLimitsExceeded: 0,               // TODO false is default value. Plan to make it configurable with environment variable
		Handler:                            srv.rootRouter.Handler,
	}

	// Run our server in a goroutine so that it doesn't block.
	go func() {
		if cast.ToBoolOrDefault(config.GddPreforkEnable.Load(), config.DefaultGddPreforkEnable) {
			preforkServer := prefork.New(httpServer, addr)
			if err := preforkServer.ListenAndServe(); err != nil {
				logger.Error().Err(err).Msg("")
			}
		} else {
			logger.Info().Msgf("Http server is listening at %v", addr)
			logger.Info().Msgf("Http server started in %s", time.Since(startAt))
			if err := httpServer.ListenAndServe(addr); err != nil {
				logger.Error().Err(err).Msg("")
			}
		}
	}()

	return httpServer
}

// Run runs http server
func (srv *FastHttpSrv) Run() {
	banner.Print()
	manage := cast.ToBoolOrDefault(config.GddManage.Load(), config.DefaultGddManage)
	if manage {
		srv.FastMiddlewares = append([]fast.MiddlewareFunc{prometheus.FastPrometheusMiddleware}, srv.FastMiddlewares...)
		gddRouter := srv.rootRouter.Group(gddPathPrefix)
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
		gddFastMiddlewares := []fast.MiddlewareFunc{fast.Metrics, corsOpts, fast.BasicAuth()}
		srv.gddRoutes = append(srv.gddRoutes, onlinedoc.Routes()...)
		srv.gddRoutes = append(srv.gddRoutes, prometheus.Routes()...)
		srv.gddRoutes = append(srv.gddRoutes, registry.Routes()...)
		srv.gddRoutes = append(srv.gddRoutes, configui.Routes()...)
		for _, item := range srv.gddRoutes {
			if item.HandlerFunc == nil && item.FastHandler == nil {
				continue
			}
			var h fasthttp.RequestHandler
			if item.FastHandler == nil && item.HandlerFunc != nil {
				h = fasthttpadaptor.NewFastHTTPHandlerFunc(item.HandlerFunc)
			} else {
				h = item.FastHandler
			}
			for i := len(gddFastMiddlewares) - 1; i >= 0; i-- {
				h = gddFastMiddlewares[i].Middleware(h)
			}
			gddRouter.Handle(item.Method, "/"+strings.TrimPrefix(item.Pattern, gddPathPrefix), h, item.Name)
		}
		srv.debugRoutes = append(srv.debugRoutes, []model.Route{
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
				Name:        "GetDebugPprofIndex",
				Method:      http.MethodGet,
				Pattern:     debugPathPrefix + "pprof/*filepath",
				FastHandler: pprofhandler.PprofHandler,
			},
		}...)
		debugRouter := srv.rootRouter.Group(debugPathPrefix)
		for _, item := range srv.debugRoutes {
			if item.HandlerFunc == nil && item.FastHandler == nil {
				continue
			}
			var h fasthttp.RequestHandler
			if item.FastHandler == nil && item.HandlerFunc != nil {
				h = fasthttpadaptor.NewFastHTTPHandlerFunc(item.HandlerFunc)
			} else {
				h = item.FastHandler
			}
			for i := len(gddFastMiddlewares) - 1; i >= 0; i-- {
				h = gddFastMiddlewares[i].Middleware(h)
			}
			debugRouter.Handle(item.Method, "/"+strings.TrimPrefix(item.Pattern, debugPathPrefix), h, item.Name)
		}
	}
	srv.FastMiddlewares = append(srv.FastMiddlewares, fast.Recovery)
	for _, item := range srv.bizRoutes {
		if item.HandlerFunc == nil && item.FastHandler == nil {
			continue
		}
		var h fasthttp.RequestHandler
		if item.FastHandler == nil && item.HandlerFunc != nil {
			h = fasthttpadaptor.NewFastHTTPHandlerFunc(item.HandlerFunc)
		} else {
			h = item.FastHandler
		}
		for i := len(srv.FastMiddlewares) - 1; i >= 0; i-- {
			h = srv.FastMiddlewares[i].Middleware(h)
		}
		srv.Router.Handle(item.Method, item.Pattern, h, item.Name)
	}
	srv.rootRouter.NotFound = func(ctx *fasthttp.RequestCtx) {
		ctx.Error(fasthttp.StatusMessage(fasthttp.StatusNotFound), fasthttp.StatusNotFound)
	}
	srv.rootRouter.MethodNotAllowed = func(ctx *fasthttp.RequestCtx) {
		ctx.SetStatusCode(fasthttp.StatusMethodNotAllowed)
		ctx.SetBodyString(fasthttp.StatusMessage(fasthttp.StatusMethodNotAllowed))
	}
	for i := len(srv.FastMiddlewares) - 1; i >= 0; i-- {
		srv.rootRouter.NotFound = srv.FastMiddlewares[i].Middleware(srv.rootRouter.NotFound)
		srv.rootRouter.MethodNotAllowed = srv.FastMiddlewares[i].Middleware(srv.rootRouter.MethodNotAllowed)
	}
	srv.printRoutes()
	httpServer := srv.newFastHttpServer()
	defer func() {
		logger.Info().Msg("Http server is gracefully shutting down...")
		httpServer.Shutdown()
	}()

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c
}
