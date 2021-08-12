package ddhttp

import (
	"context"
	"crypto/subtle"
	"github.com/common-nighthawk/go-figure"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/unionj-cloud/go-doudou/stringutils"
	"github.com/unionj-cloud/go-doudou/svc/config"
	"github.com/unionj-cloud/go-doudou/svc/http/model"
	"github.com/unionj-cloud/go-doudou/svc/http/onlinedoc"
	"github.com/unionj-cloud/go-doudou/svc/http/prometheus"
	"github.com/urfave/negroni"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"
)

// gorilla
type DefaultHttpSrv struct {
	*mux.Router
	gddRouter *mux.Router
	routes    []model.Route
}

const gddPathPrefix = "/go-doudou"

func NewDefaultHttpSrv() Srv {
	var gddRouter *mux.Router
	var routes []model.Route
	if config.GddManage.Load() == "true" {
		gddRouter = mux.NewRouter().PathPrefix(config.GddRouteRootPath.Load() + gddPathPrefix).Subrouter().StrictSlash(true)
		var mergedRoutes []model.Route
		mergedRoutes = append(mergedRoutes, onlinedoc.Routes()...)
		mergedRoutes = append(mergedRoutes, prometheus.Routes()...)
		for _, item := range mergedRoutes {
			gddRouter.
				Methods(item.Method).
				Path(strings.TrimPrefix(item.Pattern, gddPathPrefix)).
				Name(item.Name).
				Handler(item.HandlerFunc)
		}
		routes = append(routes, mergedRoutes...)
	}
	return &DefaultHttpSrv{
		mux.NewRouter().PathPrefix(config.GddRouteRootPath.Load()).Subrouter().StrictSlash(true),
		gddRouter,
		routes,
	}
}

func (srv *DefaultHttpSrv) AddRoute(route ...model.Route) {
	var routes []model.Route
	routes = append(routes, route...)
	routes = append(routes, srv.routes...)
	srv.routes = routes[:]
	routes = nil
	for _, item := range route {
		srv.
			Methods(item.Method).
			Path(item.Pattern).
			Name(item.Name).
			Handler(item.HandlerFunc)
	}
}

func BasicAuth(w http.ResponseWriter, r *http.Request) bool {
	username := config.GddManageUser.Load()
	password := config.GddManagePass.Load()
	if stringutils.IsEmpty(username) && stringutils.IsEmpty(password) {
		return true
	}

	user, pass, ok := r.BasicAuth()

	if !ok || subtle.ConstantTimeCompare([]byte(user), []byte(username)) != 1 || subtle.ConstantTimeCompare([]byte(pass), []byte(password)) != 1 {
		w.Header().Set("WWW-Authenticate", `Basic realm="Provide user name and password"`)
		w.WriteHeader(401)
		w.Write([]byte("Unauthorised.\n"))
		return false
	}

	return true
}

func (srv *DefaultHttpSrv) AddMiddleware(mwf ...func(http.Handler) http.Handler) {
	if config.GddManage.Load() == "true" {
		srv.Use(prometheus.PrometheusMiddleware)
	}
	var middlewares []mux.MiddlewareFunc
	for _, item := range mwf {
		middlewares = append(middlewares, item)
	}
	srv.Use(middlewares...)
	if config.GddManage.Load() == "true" {
		srv.PathPrefix(gddPathPrefix).Handler(negroni.New(
			negroni.HandlerFunc(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
				if BasicAuth(w, r) {
					next(w, r)
				}
			}),
			negroni.Wrap(srv.gddRouter),
		))
	}
}

func (srv *DefaultHttpSrv) Run() {
	start := time.Now()
	var bannerSwitch config.Switch
	(&bannerSwitch).Decode(config.GddBanner.Load())
	if bannerSwitch {
		banner := config.GddBannerText.Load()
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
	server.Shutdown(ctx)

	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	logrus.Infoln("shutting down")
}
