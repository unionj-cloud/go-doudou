package dou

import (
	"context"
	"github.com/goccy/go-json"
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
	"github.com/unionj-cloud/go-doudou/toolkit/reflectutils"
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
var fileType reflect.Type
var errorType reflect.Type

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
	fileType = reflect.TypeOf((*os.File)(nil))
	errorType = reflect.TypeOf((*error)(nil)).Elem()
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

func handleDownload(w http.ResponseWriter, respValue reflect.Value, fileFieldIndex int) {
	out := respValue.Field(fileFieldIndex).Interface().(*os.File)
	if out == nil {
		http.Error(w, "no file returned", http.StatusInternalServerError)
		return
	}
	defer out.Close()
	fi, err := out.Stat()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Disposition", "attachment; filename="+fi.Name())
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", fi.Size()))

	dispositionValue := respValue.FieldByName("Disposition")
	if dispositionValue.IsValid() && !dispositionValue.IsZero() {
		vov := reflectutils.ValueOfValue(dispositionValue)
		if vov.IsValid() && !vov.IsZero() {
			disposition, ok := vov.Interface().(string)
			if ok && stringutils.IsNotEmpty(disposition) {
				w.Header().Set("Content-Disposition", disposition)
			}
		}
	}
	typeValue := respValue.FieldByName("Type")
	if typeValue.IsValid() && !typeValue.IsZero() {
		vov := reflectutils.ValueOfValue(typeValue)
		if vov.IsValid() && !vov.IsZero() {
			contentType, ok := vov.Interface().(string)
			if ok && stringutils.IsNotEmpty(contentType) {
				w.Header().Set("Content-Type", contentType)
			}
		}
	}
	lengthValue := respValue.FieldByName("Length")
	if lengthValue.IsValid() && !lengthValue.IsZero() {
		vov := reflectutils.ValueOfValue(lengthValue)
		if vov.IsValid() && !vov.IsZero() {
			length, ok := vov.Interface().(string)
			if ok && stringutils.IsNotEmpty(length) {
				w.Header().Set("Content-Length", length)
			}
		}
	}
	io.Copy(w, out)
}

func buildHandler(method reflect.Method, svc reflect.Value) http.HandlerFunc {
	if method.Type.NumIn() <= 1 {
		panic("service method must have context.Context as the first input parameter")
	}
	if method.Type.NumIn() > 3 {
		panic("only support up to 2 input parameters including context.Context")
	}
	var bodyType reflect.Type
	inType := method.Type.In(1)
	ctxInterface := reflect.TypeOf((*context.Context)(nil)).Elem()
	if !inType.Implements(ctxInterface) {
		panic("service method must have context.Context as the first input parameter")
	}
	if method.Type.NumIn() > 2 {
		inType = method.Type.In(2)
		if inType.Kind() == reflect.Ptr && inType.Elem().Kind() == reflect.Struct {
			// pointer of struct
			bodyType = inType
		} else if inType.Kind() == reflect.Struct {
			// struct
			bodyType = inType
		} else {
			panic("only support struct type, pointer of struct type as the second input parameter")
		}
	}
	if method.Type.NumOut() > 2 {
		panic("only support up to 2 output parameters including error")
	}
	outTypes := make([]reflect.Type, 0)
	for i := 0; i < method.Type.NumOut(); i++ {
		outType := method.Type.Out(i)
		if outType.Implements(errorType) {
			continue
		}
		outTypes = append(outTypes, outType)
	}
	if len(outTypes) > 0 {
		ot := outTypes[0]
		if ot.Kind() == reflect.Ptr {
			ot = ot.Elem()
		}
		if ot.Kind() != reflect.Struct {
			panic("only support struct type, pointer of struct type and error as output parameter")
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
				http.Error(w, fmt.Sprintf("incorrect Content-Type header %s, only accept application/json", ct), http.StatusBadRequest)
				return
			}
			if httpMethod == POST || httpMethod == PUT {
				switch ct {
				case "multipart/form-data":
					// TODO add maxMemory config
					if err = r.ParseMultipartForm(32 << 20); err != nil {
						http.Error(w, err.Error(), http.StatusBadRequest)
						return
					}
					err = decoder.Decode(pqPtr.Interface(), r.MultipartForm.Value)
					for i := 0; i < bodyType.NumField(); i++ {
						field := bodyType.Field(i)
						fieldType := field.Type
						fieldName := field.Name
						formTag := field.Tag.Get("form")
						var formFieldName string
						if stringutils.IsNotEmpty(formTag) {
							formFieldName = strings.Split(formTag, ",")[0]
							if formFieldName == "-" {
								formFieldName = ""
							}
						}
						if stringutils.IsEmpty(formFieldName) {
							formFieldName = strcase.ToLowerCamel(fieldName)
						}
						if fileHeaders, exists := r.MultipartForm.File[formFieldName]; exists {
							if len(fileHeaders) > 0 {
								if reflect.TypeOf(fileHeaders).AssignableTo(fieldType) {
									pqPtr.Field(i).Set(reflect.ValueOf(fileHeaders))
								} else if reflect.TypeOf(fileHeaders[0]).AssignableTo(fieldType) {
									pqPtr.Field(i).Set(reflect.ValueOf(fileHeaders[0]))
								}
							}
						}
					}
					goto VALIDATE
				case "application/json":
					if err = json.NewDecoder(r.Body).Decode(pqPtr.Interface()); err != nil {
						if err != io.EOF {
							http.Error(w, err.Error(), http.StatusBadRequest)
							return
						}
						err = nil
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
			switch out := item.Interface().(type) {
			case error:
				if out != nil {
					if errors.Is(out, context.Canceled) {
						http.Error(w, out.Error(), http.StatusBadRequest)
					} else if _err, ok := out.(*ddhttp.BizError); ok {
						http.Error(w, _err.Error(), _err.StatusCode)
					} else {
						http.Error(w, out.Error(), http.StatusInternalServerError)
					}
					return
				}
			default:
				copyOutValues = append(copyOutValues, item)
			}
		}
		var resp interface{}
		if len(copyOutValues) > 0 {
			outValue := copyOutValues[0]
			respType := reflect.TypeOf(outValue)
			if respType.Kind() == reflect.Ptr {
				respType = respType.Elem()
			}
			var fileFieldIndex *int
			for i := 0; i < respType.NumField(); i++ {
				field := respType.Field(i)
				if field.Type.AssignableTo(fileType) {
					fileFieldIndex = &i
					break
				}
			}
			respValue := reflectutils.ValueOfValue(outValue)
			if fileFieldIndex != nil {
				if respValue.IsValid() && !respValue.IsZero() {
					handleDownload(w, respValue, *fileFieldIndex)
				} else {
					http.Error(w, "empty response", http.StatusInternalServerError)
				}
				return
			}
			if respValue.IsValid() {
				resp = respValue.Interface()
			} else {
				resp = reflect.New(respType).Elem()
			}
		} else {
			resp = make(map[string]interface{})
		}
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
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
