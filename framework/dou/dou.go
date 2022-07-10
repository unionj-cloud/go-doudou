package dou

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ascarter/requestid"
	"github.com/common-nighthawk/go-figure"
	"github.com/go-playground/form/v4"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/iancoleman/strcase"
	"github.com/olekukonko/tablewriter"
	"github.com/pkg/errors"
	"github.com/rs/cors"
	ddhttp "github.com/unionj-cloud/go-doudou/framework/http"
	configui "github.com/unionj-cloud/go-doudou/framework/http/config"
	"github.com/unionj-cloud/go-doudou/framework/http/model"
	"github.com/unionj-cloud/go-doudou/framework/http/onlinedoc"
	"github.com/unionj-cloud/go-doudou/framework/http/prometheus"
	"github.com/unionj-cloud/go-doudou/framework/http/registry"
	"github.com/unionj-cloud/go-doudou/framework/internal/config"
	"github.com/unionj-cloud/go-doudou/framework/logger"
	"github.com/unionj-cloud/go-doudou/toolkit/cast"
	"github.com/unionj-cloud/go-doudou/toolkit/sliceutils"
	"github.com/unionj-cloud/go-doudou/toolkit/stringutils"
	"io"
	"mime"
	"net/http"
	"os"
	"os/signal"
	"path"
	"reflect"
	"strconv"
	"strings"
	"time"
)

var bizRouter *mux.Router
var rootRouter *mux.Router
var gddRoutes []model.Route
var bizRoutes []model.Route
var middlewares []mux.MiddlewareFunc
var decoder *form.Decoder

func RegisterCustomTypeFunc(fn form.DecodeCustomTypeFunc, types ...interface{}) {
	decoder.RegisterCustomTypeFunc(fn, types...)
}

func init() {
	rr := config.DefaultGddRouteRootPath
	if stringutils.IsNotEmpty(config.GddRouteRootPath.Load()) {
		rr = config.GddRouteRootPath.Load()
	}
	rootRouter = mux.NewRouter().StrictSlash(true)
	bizRouter = rootRouter.PathPrefix(rr).Subrouter().StrictSlash(true)
	middlewares = make([]mux.MiddlewareFunc, 0)
	middlewares = append(middlewares,
		tracing,
		metrics,
	)
	if cast.ToBoolOrDefault(config.GddEnableResponseGzip.Load(), config.DefaultGddEnableResponseGzip) {
		middlewares = append(middlewares, handlers.CompressHandler)
	}
	if cast.ToBoolOrDefault(config.GddLogReqEnable.Load(), config.DefaultGddLogReqEnable) {
		middlewares = append(middlewares, log)
	}
	middlewares = append(middlewares,
		requestid.RequestIDHandler,
		handlers.ProxyHeaders,
	)
	appType := config.GddAppType.LoadOrDefault(config.DefaultGddAppType)
	switch strings.TrimSpace(appType) {
	case "rest":
		middlewares = append(middlewares, rest)
	}
	decoder = form.NewDecoder()
}

func pattern(method string) string {
	httpMethods := []string{"GET", "POST", "PUT", "DELETE"}
	snake := strcase.ToSnake(strings.ReplaceAll(method, "_", "."))
	splits := strings.Split(snake, "_")
	head := strings.ToUpper(splits[0])
	if sliceutils.StringContains(httpMethods, head) {
		splits = splits[1:]
	}
	clean := sliceutils.StringFilter(splits, func(item string) bool {
		return stringutils.IsNotEmpty(item)
	})
	return strings.Join(clean, "/")
}

func buildRoutes(service interface{}) []model.Route {
	routes := make([]model.Route, 0)
	svcType := reflect.TypeOf(service)
	for i := 0; i < svcType.NumMethod(); i++ {
		m := svcType.Method(i)
		routes = append(routes, model.Route{
			Name:        m.Name,
			Method:      "POST",
			Pattern:     fmt.Sprintf("/%s", pattern(m.Name)),
			HandlerFunc: buildHandler(m, reflect.ValueOf(service)),
		})
	}
	return routes
}

func buildHandler(method reflect.Method, svc reflect.Value) http.HandlerFunc {
	var bodyType reflect.Type
	for i := 1; i < method.Type.NumIn(); i++ {
		inType := method.Type.In(i)
		if inType.String() == "context.Context" {
			continue
		}
		bodyType = inType
		if bodyType.Kind() != reflect.Struct {
			if !(bodyType.Kind() == reflect.Ptr && bodyType.Elem().Kind() == reflect.Struct) {
				panic("only support struct type or pointer of struct type as method input parameter")
			}
		}
	}
	return func(w http.ResponseWriter, r *http.Request) {
		var httpMethod HttpMethod
		httpMethod.StringSetter(r.Method)
		outValues := make([]reflect.Value, 0)
		if bodyType != nil {
			ct := r.Header.Get("Content-Type")
			if ct == "" {
				ct = "application/octet-stream"
			}
			ct, _, err := mime.ParseMediaType(ct)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			pqPtr := reflect.New(bodyType)
			if pqPtr.Elem().Kind() == reflect.Ptr && ct != "application/json" {
				http.Error(w, "only accept application/json", http.StatusBadRequest)
				return
			}
			if httpMethod == POST || httpMethod == PUT {
				if ct == "application/json" {
					if err = json.NewDecoder(r.Body).Decode(pqPtr.Interface()); err != nil {
						if err != io.EOF {
							http.Error(w, err.Error(), http.StatusBadRequest)
							return
						}
						err = nil
					} else {
						goto VALIDATE
					}
				}
			}
			if pqPtr.Elem().Kind() == reflect.Struct {
				if err = r.ParseForm(); err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}
				if err = decoder.Decode(pqPtr.Interface(), r.Form); err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}
			}
		VALIDATE:
			if pqPtr.Elem().Kind() == reflect.Struct || (pqPtr.Elem().Kind() == reflect.Ptr && !pqPtr.Elem().IsNil()) {
				if err = ddhttp.ValidateStruct(pqPtr.Elem().Interface()); err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}
			}
			outValues = method.Func.Call([]reflect.Value{svc, reflect.ValueOf(r.Context()), pqPtr.Elem()})
		} else {
			outValues = method.Func.Call([]reflect.Value{svc, reflect.ValueOf(r.Context())})
		}
		copyOutValues := make([]reflect.Value, 0)
		for _, item := range outValues {
			err, ok := item.Interface().(error)
			if !ok {
				copyOutValues = append(copyOutValues, item)
				continue
			}
			if err != nil {
				if errors.Is(err, context.Canceled) {
					http.Error(w, err.Error(), http.StatusBadRequest)
				} else if _err, ok := err.(*ddhttp.BizError); ok {
					http.Error(w, _err.Error(), _err.StatusCode)
				} else {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
				return
			}
		}
		var resp interface{}
		if len(copyOutValues) > 0 {
			resp = copyOutValues[0].Interface()
		} else {
			resp = make(map[string]interface{})
		}
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

const gddPathPrefix = "/go-doudou/"

// Run runs http server
func Run(services ...interface{}) {
	manage := cast.ToBoolOrDefault(config.GddManage.Load(), config.DefaultGddManage)
	if manage {
		middlewares = append([]mux.MiddlewareFunc{prometheus.PrometheusMiddleware}, middlewares...)
		gddRouter := rootRouter.PathPrefix(gddPathPrefix).Subrouter().StrictSlash(true)
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
		gddRouter.Use(metrics)
		gddRouter.Use(corsOpts.Handler)
		gddRouter.Use(basicAuth())
		gddRoutes = append(gddRoutes, onlinedoc.Routes()...)
		gddRoutes = append(gddRoutes, prometheus.Routes()...)
		gddRoutes = append(gddRoutes, registry.Routes()...)
		gddRoutes = append(gddRoutes, configui.Routes()...)
		for _, item := range gddRoutes {
			gddRouter.
				Methods(item.Method, http.MethodOptions).
				Path("/" + strings.TrimPrefix(item.Pattern, gddPathPrefix)).
				Name(item.Name).
				Handler(item.HandlerFunc)
		}
	}
	middlewares = append(middlewares, recovery)
	bizRouter.Use(middlewares...)
	if len(services) > 0 {
		service := services[0]
		bizRoutes = make([]model.Route, 0)
		bizRoutes = append(bizRoutes, buildRoutes(service)...)
	}
	for _, item := range bizRoutes {
		bizRouter.
			Methods(item.Method, http.MethodOptions).
			Path(item.Pattern).
			Name(item.Name).
			Handler(item.HandlerFunc)
	}
	rootRouter.NotFoundHandler = rootRouter.NewRoute().BuildOnly().HandlerFunc(http.NotFound).GetHandler()
	rootRouter.MethodNotAllowedHandler = rootRouter.NewRoute().BuildOnly().HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("405 method not allowed"))
	}).GetHandler()

	for _, item := range middlewares {
		rootRouter.NotFoundHandler = item.Middleware(rootRouter.NotFoundHandler)
		rootRouter.MethodNotAllowedHandler = item.Middleware(rootRouter.MethodNotAllowedHandler)
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

	printRoutes()
	httpServer := newHttpServer()
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

// AddRoute adds routes to router
func AddRoute(route ...model.Route) {
	bizRoutes = append(bizRoutes, route...)
}

func printRoutes() {
	logger.Infoln("================ Registered Routes ================")
	data := [][]string{}
	rr := config.DefaultGddRouteRootPath
	if stringutils.IsNotEmpty(config.GddRouteRootPath.Load()) {
		rr = config.GddRouteRootPath.Load()
	}
	var all []model.Route
	all = append(all, bizRoutes...)
	all = append(all, gddRoutes...)
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
func AddMiddleware(mwf ...func(http.Handler) http.Handler) {
	for _, item := range mwf {
		middlewares = append(middlewares, item)
	}
}

// PreMiddleware adds middlewares to the head of chain
func PreMiddleware(mwf ...func(http.Handler) http.Handler) {
	var preMiddlewares []mux.MiddlewareFunc
	for _, item := range mwf {
		preMiddlewares = append(preMiddlewares, item)
	}
	middlewares = append(preMiddlewares, middlewares...)
}

func newHttpServer() *http.Server {
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
		Handler:      rootRouter, // Pass our instance of gorilla/mux in.
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
