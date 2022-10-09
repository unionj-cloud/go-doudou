package router

import (
	"strings"

	"github.com/savsgio/gotils/strconv"
	"github.com/valyala/fasthttp"
)

var (
	defaultContentType = []byte("text/plain; charset=utf-8")
	questionMark       = byte('?')
	routeNameKey       = []byte("routeNameKey")
)

func (r *Router) saveMatchedRoutePath(name string, handler fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		ctx.SetUserValueBytes(routeNameKey, name)
		handler(ctx)
	}
}

// New returns a new router.
// Path auto-correction, including trailing slashes, is enabled by default.
func New() *Router {
	return &Router{
		registeredPaths:        make(map[string][]string),
		RedirectTrailingSlash:  true,
		RedirectFixedPath:      true,
		HandleMethodNotAllowed: true,
		HandleOPTIONS:          true,
		SaveMatchedRoutePath:   true,
		handlerStore:           make(map[string]fasthttp.RequestHandler),
	}
}

// Group returns a new group.
// Path auto-correction, including trailing slashes, is enabled by default.
func (r *Router) Group(path string) *Group {
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

// GET is a shortcut for router.Handle(fasthttp.MethodGet, path, handler)
func (r *Router) GET(path string, handler fasthttp.RequestHandler, name ...string) {
	r.Handle(fasthttp.MethodGet, path, handler, name...)
}

// HEAD is a shortcut for router.Handle(fasthttp.MethodHead, path, handler)
func (r *Router) HEAD(path string, handler fasthttp.RequestHandler, name ...string) {
	r.Handle(fasthttp.MethodHead, path, handler, name...)
}

// POST is a shortcut for router.Handle(fasthttp.MethodPost, path, handler)
func (r *Router) POST(path string, handler fasthttp.RequestHandler, name ...string) {
	r.Handle(fasthttp.MethodPost, path, handler, name...)
}

// PUT is a shortcut for router.Handle(fasthttp.MethodPut, path, handler)
func (r *Router) PUT(path string, handler fasthttp.RequestHandler, name ...string) {
	r.Handle(fasthttp.MethodPut, path, handler, name...)
}

// PATCH is a shortcut for router.Handle(fasthttp.MethodPatch, path, handler)
func (r *Router) PATCH(path string, handler fasthttp.RequestHandler, name ...string) {
	r.Handle(fasthttp.MethodPatch, path, handler, name...)
}

// DELETE is a shortcut for router.Handle(fasthttp.MethodDelete, path, handler)
func (r *Router) DELETE(path string, handler fasthttp.RequestHandler, name ...string) {
	r.Handle(fasthttp.MethodDelete, path, handler, name...)
}

// CONNECT is a shortcut for router.Handle(fasthttp.MethodConnect, path, handler)
func (r *Router) CONNECT(path string, handler fasthttp.RequestHandler, name ...string) {
	r.Handle(fasthttp.MethodConnect, path, handler, name...)
}

// OPTIONS is a shortcut for router.Handle(fasthttp.MethodOptions, path, handler)
func (r *Router) OPTIONS(path string, handler fasthttp.RequestHandler, name ...string) {
	r.Handle(fasthttp.MethodOptions, path, handler, name...)
}

// TRACE is a shortcut for router.Handle(fasthttp.MethodTrace, path, handler)
func (r *Router) TRACE(path string, handler fasthttp.RequestHandler, name ...string) {
	r.Handle(fasthttp.MethodTrace, path, handler, name...)
}

// ServeFiles serves files from the given file system root.
// The path must end with "/{filepath:*}", files are then served from the local
// path /defined/root/dir/{filepath:*}.
// For example if root is "/etc" and {filepath:*} is "passwd", the local file
// "/etc/passwd" would be served.
// Internally a fasthttp.FSHandler is used, therefore fasthttp.NotFound is used instead
// Use:
//     router.ServeFiles("/src/{filepath:*}", "./")
func (r *Router) ServeFiles(path string, rootPath string) {
	r.ServeFilesCustom(path, &fasthttp.FS{
		Root:               rootPath,
		IndexNames:         []string{"index.html"},
		GenerateIndexPages: true,
		AcceptByteRange:    true,
	})
}

// ServeFilesCustom serves files from the given file system settings.
// The path must end with "/{filepath:*}", files are then served from the local
// path /defined/root/dir/{filepath:*}.
// For example if root is "/etc" and {filepath:*} is "passwd", the local file
// "/etc/passwd" would be served.
// Internally a fasthttp.FSHandler is used, therefore http.NotFound is used instead
// of the Router's NotFound handler.
// Use:
//     router.ServeFilesCustom("/src/{filepath:*}", *customFS)
func (r *Router) ServeFilesCustom(path string, fs *fasthttp.FS) {
	suffix := "/{filepath:*}"

	if !strings.HasSuffix(path, suffix) {
		panic("path must end with " + suffix + " in path '" + path + "'")
	}

	prefix := path[:len(path)-len(suffix)]
	stripSlashes := strings.Count(prefix, "/")

	if fs.PathRewrite == nil && stripSlashes > 0 {
		fs.PathRewrite = fasthttp.NewPathSlashesStripper(stripSlashes)
	}
	fileHandler := fs.NewRequestHandler()

	r.GET(path, fileHandler)
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
	case len(name) == 0:
		panic("route name must not be nil")
	}
	if r.SaveMatchedRoutePath {
		handler = r.saveMatchedRoutePath(name[0], handler)
	}
	r.handlerStore[path2key(method, path)] = handler
	_, f := r.registeredPaths[method]
	r.registeredPaths[method] = append(r.registeredPaths[method], path)
	if !f {
		r.globalAllowed = r.allowed("*", "")
	}
}

func (r *Router) recv(ctx *fasthttp.RequestCtx) {
	if rcv := recover(); rcv != nil {
		r.PanicHandler(ctx, rcv)
	}
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

			handle, _ := r.handlerStore[path2key(method, path)]
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

	path := strconv.B2S(ctx.Request.URI().Path())
	method := strconv.B2S(ctx.Request.Header.Method())
	methodIndex := r.methodIndexOf(method)

	if methodIndex > -1 {
		if handler := r.handlerStore[path2key(method, path)]; handler != nil {
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
