package router

import (
	"github.com/ucarion/urlpath"
	"github.com/unionj-cloud/go-doudou/toolkit/stringutils"
	"strings"

	"github.com/savsgio/gotils/strconv"
	"github.com/valyala/fasthttp"
)

var (
	routeNameKey = []byte("routeNameKey")
	httpMethods  = []string{fasthttp.MethodGet, fasthttp.MethodHead, fasthttp.MethodPost,
		fasthttp.MethodPut, fasthttp.MethodPatch, fasthttp.MethodDelete, fasthttp.MethodConnect,
		fasthttp.MethodOptions, fasthttp.MethodTrace}
)

// Router is a fasthttp.RequestHandler which can be used to dispatch requests to different
// handler functions via configurable routes
type Router struct {
	registeredPaths map[string][]string

	// If enabled, adds the matched route path onto the ctx.UserValue context
	// before invoking the handler.
	// The matched route path is only added to handlers of routes that were
	// registered when this option was enabled.
	SaveMatchedRoutePath bool

	// If enabled, the router checks if another method is allowed for the
	// current route, if the current request can not be routed.
	// If this is the case, the request is answered with 'Method Not Allowed'
	// and HTTP status code 405.
	// If no other Method is allowed, the request is delegated to the NotFound
	// handler.
	HandleMethodNotAllowed bool

	// If enabled, the router automatically replies to OPTIONS requests.
	// Custom OPTIONS handlers take priority over automatic replies.
	HandleOPTIONS bool

	// An optional fasthttp.RequestHandler that is called on automatic OPTIONS requests.
	// The handler is only called if HandleOPTIONS is true and no OPTIONS
	// handler for the specific path was set.
	// The "Allowed" header is set before calling the handler.
	GlobalOPTIONS fasthttp.RequestHandler

	// Configurable fasthttp.RequestHandler which is called when no matching route is
	// found. If it is not set, default NotFound is used.
	NotFound fasthttp.RequestHandler

	// Configurable fasthttp.RequestHandler which is called when a request
	// cannot be routed and HandleMethodNotAllowed is true.
	// If it is not set, ctx.Error with fasthttp.StatusMethodNotAllowed is used.
	// The "Allow" header with allowed request methods is set before the handler
	// is called.
	MethodNotAllowed fasthttp.RequestHandler

	// Function to handle panics recovered from http handlers.
	// It should be used to generate a error page and return the http error code
	// 500 (Internal Server Error).
	// The handler can be used to keep your server from crashing because of
	// unrecovered panics.
	PanicHandler func(*fasthttp.RequestCtx, interface{})

	// Cached value of global (*) allowed methods
	globalAllowed string

	handlers        map[string]fasthttp.RequestHandler
	dynamicHandlers []map[*urlpath.Path]fasthttp.RequestHandler
}

// Group is a sub-router to group paths
type Group struct {
	router *Router
	prefix string
}

func (r *Router) saveMatchedRoutePath(name string, handler fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		ctx.SetUserValueBytes(routeNameKey, name)
		handler(ctx)
	}
}

// New returns a new router.
// Path auto-correction, including trailing slashes, is enabled by default.
func New() *Router {
	r := &Router{
		registeredPaths:        make(map[string][]string),
		HandleMethodNotAllowed: true,
		HandleOPTIONS:          true,
		SaveMatchedRoutePath:   true,
		handlers:               make(map[string]fasthttp.RequestHandler),
		dynamicHandlers:        make([]map[*urlpath.Path]fasthttp.RequestHandler, len(httpMethods)),
	}
	for i := range httpMethods {
		r.dynamicHandlers[i] = make(map[*urlpath.Path]fasthttp.RequestHandler)
	}
	return r
}

// Group returns a new group.
// Path auto-correction, including trailing slashes, is enabled by default.
func (r *Router) Group(path string) *Group {
	validatePath(path)
	return &Group{
		router: r,
		prefix: path,
	}
}

func (r *Router) methodIndexOf(method string) int {
	switch method {
	case fasthttp.MethodGet:
		return 0
	case fasthttp.MethodHead:
		return 1
	case fasthttp.MethodPost:
		return 2
	case fasthttp.MethodPut:
		return 3
	case fasthttp.MethodPatch:
		return 4
	case fasthttp.MethodDelete:
		return 5
	case fasthttp.MethodConnect:
		return 6
	case fasthttp.MethodOptions:
		return 7
	case fasthttp.MethodTrace:
		return 8
	}
	return -1
}

// List returns all registered routes grouped by method
func (r *Router) List() map[string][]string {
	return r.registeredPaths
}

func path2key(method, path string) string {
	var sb strings.Builder
	sb.WriteString(method)
	sb.WriteString(":")
	sb.WriteString(path)
	return sb.String()
}

// Handle registers a new request handler with the given path and method.
//
// For GET, POST, PUT, PATCH and DELETE requests the respective shortcut
// functions can be used.
//
// This function is intended for bulk loading and to allow the usage of less
// frequently used, non-standardized or custom methods (e.g. for internal
// communication with a proxy).
func (r *Router) Handle(method, path string, handler fasthttp.RequestHandler, name ...string) {
	switch {
	case len(method) == 0:
		panic("method must not be empty")
	case len(path) < 1 || path[0] != '/':
		panic("path must begin with '/' in path '" + path + "'")
	case handler == nil:
		panic("handler must not be nil")
	}
	idx := r.methodIndexOf(method)
	if idx < 0 {
		panic("unknown http method")
	}
	_, f := r.registeredPaths[method]
	r.registeredPaths[method] = append(r.registeredPaths[method], path)
	if !f {
		r.globalAllowed = r.allowed("*", "")
	}
	if r.SaveMatchedRoutePath {
		if len(name) == 0 {
			panic("route name must not be nil")
		}
		handler = r.saveMatchedRoutePath(name[0], handler)
	}
	if strings.Contains(path, "*") {
		pt := urlpath.New(path)
		r.dynamicHandlers[idx][&pt] = handler
	} else {
		r.handlers[path2key(method, path)] = handler
	}
}

func (r *Router) recv(ctx *fasthttp.RequestCtx) {
	if rcv := recover(); rcv != nil {
		r.PanicHandler(ctx, rcv)
	}
}

func (r *Router) search(method, path string) fasthttp.RequestHandler {
	idx := r.methodIndexOf(method)
	if idx < 0 {
		return nil
	}
	for k := range r.dynamicHandlers[idx] {
		if _, ok := k.Match(path); !ok {
			continue
		}
		return r.dynamicHandlers[idx][k]
	}
	return nil
}

func (r *Router) allowed(path, reqMethod string) (allow string) {
	allowed := make([]string, 0, 9)

	if path == "*" || path == "/*" { // server-wide{ // server-wide
		// empty method is used for internal calls to refresh the cache
		if reqMethod == "" {
			for method := range r.registeredPaths {
				if method == fasthttp.MethodOptions {
					continue
				}
				// Add request method to list of allowed methods
				allowed = append(allowed, method)
			}
		} else {
			return r.globalAllowed
		}
	} else { // specific path
		for method := range r.registeredPaths {
			// Skip the requested method - we already tried this one
			if method == reqMethod || method == fasthttp.MethodOptions {
				continue
			}

			handle, _ := r.handlers[path2key(method, path)]
			if handle == nil {
				handle = r.search(method, path)
			}
			if handle != nil {
				// Add request method to list of allowed methods
				allowed = append(allowed, method)
			}
		}
	}

	if len(allowed) > 0 {
		// Add request method to list of allowed methods
		allowed = append(allowed, fasthttp.MethodOptions)

		// Sort allowed methods.
		// sort.Strings(allowed) unfortunately causes unnecessary allocations
		// due to allowed being moved to the heap and interface conversion
		for i, l := 1, len(allowed); i < l; i++ {
			for j := i; j > 0 && allowed[j] < allowed[j-1]; j-- {
				allowed[j], allowed[j-1] = allowed[j-1], allowed[j]
			}
		}

		// return as comma separated list
		return strings.Join(allowed, ", ")
	}
	return
}

// Handler makes the router implement the http.Handler interface.
func (r *Router) Handler(ctx *fasthttp.RequestCtx) {
	if r.PanicHandler != nil {
		defer r.recv(ctx)
	}

	path := strconv.B2S(ctx.Path())
	method := strconv.B2S(ctx.Request.Header.Method())
	methodIndex := r.methodIndexOf(method)

	if methodIndex > -1 {
		handler := r.handlers[path2key(method, path)]
		if handler == nil {
			handler = r.search(method, path)
		}
		if handler != nil {
			handler(ctx)
			return
		}
	}

	if r.HandleOPTIONS && method == fasthttp.MethodOptions {
		// Handle OPTIONS requests

		if allow := r.allowed(path, fasthttp.MethodOptions); allow != "" {
			ctx.Response.Header.Set("Allow", allow)
			if r.GlobalOPTIONS != nil {
				r.GlobalOPTIONS(ctx)
			}
			return
		}
	} else if r.HandleMethodNotAllowed {
		// Handle 405

		if allow := r.allowed(path, method); allow != "" {
			ctx.Response.Header.Set("Allow", allow)
			if r.MethodNotAllowed != nil {
				r.MethodNotAllowed(ctx)
			} else {
				ctx.SetStatusCode(fasthttp.StatusMethodNotAllowed)
				ctx.SetBodyString(fasthttp.StatusMessage(fasthttp.StatusMethodNotAllowed))
			}
			return
		}
	}

	// Handle 404
	if r.NotFound != nil {
		r.NotFound(ctx)
	} else {
		ctx.Error(fasthttp.StatusMessage(fasthttp.StatusNotFound), fasthttp.StatusNotFound)
	}
}


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

