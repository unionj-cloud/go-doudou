package httprouter

import (
	"github.com/unionj-cloud/go-doudou/v2/toolkit/stringutils"
	"net/http"
)

type RouteGroup struct {
	r *Router
	p string
}

func validatePath(path string) {
	if stringutils.IsEmpty(path) {
		panic("path should not be empty")
	}
	if path[0] != '/' {
		panic("path must start with a '/'")
	}
}

func newRouteGroup(r *Router, path string) *RouteGroup {
	validatePath(path)

	//Strip traling / (if present) as all added sub paths must start with a /
	if path[len(path)-1] == '/' {
		path = path[:len(path)-1]
	}
	return &RouteGroup{r: r, p: path}
}

func (r *RouteGroup) NewGroup(path string) *RouteGroup {
	return newRouteGroup(r.r, r.SubPath(path))
}

func (r *RouteGroup) Handle(method, path string, handle Handle, name ...string) {
	r.r.Handle(method, r.SubPath(path), handle, name...)
}

func (r *RouteGroup) Handler(method, path string, handler http.Handler, name ...string) {
	r.r.Handler(method, r.SubPath(path), handler, name...)
}

func (r *RouteGroup) HandlerFunc(method, path string, handler http.HandlerFunc, name ...string) {
	r.r.HandlerFunc(method, r.SubPath(path), handler, name...)
}

func (r *RouteGroup) GET(path string, handle Handle) {
	r.Handle("GET", path, handle)
}
func (r *RouteGroup) HEAD(path string, handle Handle) {
	r.Handle("HEAD", path, handle)
}
func (r *RouteGroup) OPTIONS(path string, handle Handle) {
	r.Handle("OPTIONS", path, handle)
}
func (r *RouteGroup) POST(path string, handle Handle) {
	r.Handle("POST", path, handle)
}
func (r *RouteGroup) PUT(path string, handle Handle) {
	r.Handle("PUT", path, handle)
}
func (r *RouteGroup) PATCH(path string, handle Handle) {
	r.Handle("PATCH", path, handle)
}
func (r *RouteGroup) DELETE(path string, handle Handle) {
	r.Handle("DELETE", path, handle)
}

func (r *RouteGroup) SubPath(path string) string {
	result := r.p + path
	if stringutils.IsEmpty(result) {
		return "/"
	}
	return r.p + path
}
