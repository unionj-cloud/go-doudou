package cors

import (
	"github.com/rs/cors"
	"github.com/unionj-cloud/go-doudou/framework/http/fast"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttpadaptor"
	"net/http"
)

// Options is a configuration container to setup the CORS middleware.
type Options = cors.Options

// corsWrapper is a wrapper of cors.Cors handler which preserves information
// about configured 'optionPassthrough' option.
type corsWrapper struct {
	*cors.Cors
	optionPassthrough bool
}

// build transforms wrapped cors.Cors handler into Gin middleware.
func (c corsWrapper) build() fast.MiddlewareFunc {
	return func(inner fasthttp.RequestHandler) fasthttp.RequestHandler {
		return func(ctx *fasthttp.RequestCtx) {
			fasthttpadaptor.NewFastHTTPHandlerFunc(c.HandlerFunc)(ctx)
			if !c.optionPassthrough &&
				string(ctx.Method()) == http.MethodOptions && len(ctx.Request.Header.Peek("Access-Control-Request-Method")) > 0 {
				return
			}
			inner(ctx)
		}
	}
}

// AllowAll creates a new CORS Gin middleware with permissive configuration
// allowing all origins with all standard methods with any header and
// credentials.
func AllowAll() fast.MiddlewareFunc {
	return corsWrapper{Cors: cors.AllowAll()}.build()
}

// Default creates a new CORS Gin middleware with default options.
func Default() fast.MiddlewareFunc {
	return corsWrapper{Cors: cors.Default()}.build()
}

// New creates a new CORS Gin middleware with the provided options.
func New(options Options) fast.MiddlewareFunc {
	return corsWrapper{cors.New(options), options.OptionsPassthrough}.build()
}
