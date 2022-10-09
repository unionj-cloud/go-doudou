package router

import (
	"github.com/unionj-cloud/go-doudou/toolkit/stringutils"
	"github.com/valyala/fasthttp"
)

// Group returns a new group.
// Path auto-correction, including trailing slashes, is enabled by default.
func (g *Group) Group(path string) *Group {
	return g.router.Group(g.subPath(path))
}

func validatePath(path string) {
	if stringutils.IsEmpty(path) {
		panic("path should not be empty")
	}
	if path[0] != '/' {
		panic("path must start with a '/'")
	}
}

func (g *Group) subPath(path string) string {
	validatePath(path)
	//Strip traling / (if present) as all added sub paths must start with a /
	if path[len(path)-1] == '/' {
		path = path[:len(path)-1]
	}
	return g.prefix + path
}

// Handle registers a new request handler with the given path and method.
//
// For GET, POST, PUT, PATCH and DELETE requests the respective shortcut
// functions can be used.
//
// This function is intended for bulk loading and to allow the usage of less
// frequently used, non-standardized or custom methods (e.g. for internal
// communication with a proxy).
func (g *Group) Handle(method, path string, handler fasthttp.RequestHandler, name ...string) {
	g.router.Handle(method, g.subPath(path), handler, name...)
}
