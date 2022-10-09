package router

import (
	"github.com/ucarion/urlpath"
	"github.com/valyala/fasthttp"
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
