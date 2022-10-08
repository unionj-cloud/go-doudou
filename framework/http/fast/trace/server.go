//go:build go1.7
// +build go1.7

package trace

import (
	"github.com/d7561985/opentracefasthttp"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/valyala/fasthttp"
)

type mwOptions struct {
	opNameFunc    func(r *fasthttp.Request) string
	spanFilter    func(r *fasthttp.Request) bool
	spanObserver  func(span opentracing.Span, r *fasthttp.Request)
	urlTagFunc    func(u *fasthttp.URI) string
	componentName string
}

// MWOption controls the behavior of the Middleware.
type MWOption func(*mwOptions)

// OperationNameFunc returns a MWOption that uses given function f
// to generate operation name for each server-side span.
func OperationNameFunc(f func(r *fasthttp.Request) string) MWOption {
	return func(options *mwOptions) {
		options.opNameFunc = f
	}
}

// MWComponentName returns a MWOption that sets the component name
// for the server-side span.
func MWComponentName(componentName string) MWOption {
	return func(options *mwOptions) {
		options.componentName = componentName
	}
}

// MWSpanFilter returns a MWOption that filters requests from creating a span
// for the server-side span.
// Span won't be created if it returns false.
func MWSpanFilter(f func(r *fasthttp.Request) bool) MWOption {
	return func(options *mwOptions) {
		options.spanFilter = f
	}
}

// MWSpanObserver returns a MWOption that observe the span
// for the server-side span.
func MWSpanObserver(f func(span opentracing.Span, r *fasthttp.Request)) MWOption {
	return func(options *mwOptions) {
		options.spanObserver = f
	}
}

// MWURLTagFunc returns a MWOption that uses given function f
// to set the span's fasthttp.url tag. Can be used to change the default
// fasthttp.url tag, eg to redact sensitive information.
func MWURLTagFunc(f func(u *fasthttp.URI) string) MWOption {
	return func(options *mwOptions) {
		options.urlTagFunc = f
	}
}

var spanKey = []byte("spanKey")

func ContextWithSpan(ctx *fasthttp.RequestCtx, span opentracing.Span) {
	ctx.SetUserValueBytes(spanKey, span)
}

// SpanFromContext returns the `Span` previously associated with `ctx`, or
// `nil` if no such `Span` could be found.
//
// NOTE: context.Context != SpanContext: the former is Go's intra-process
// context propagation mechanism, and the latter houses OpenTracing's per-Span
// identity and baggage information.
func SpanFromContext(ctx *fasthttp.RequestCtx) opentracing.Span {
	val := ctx.UserValueBytes(spanKey)
	if sp, ok := val.(opentracing.Span); ok {
		return sp
	}
	return nil
}

// Middleware wraps an fasthttp.Handler and traces incoming requests.
// Additionally, it adds the span to the request's context.
//
// By default, the operation name of the spans is set to "HTTP {method}".
// This can be overriden with options.
//
// Example:
// 	 fasthttp.ListenAndServe("localhost:80", netfasthttp.Middleware(tracer, fasthttp.DefaultServeMux))
//
// The options allow fine tuning the behavior of the middleware.
//
// Example:
//   mw := netfasthttp.Middleware(
//      tracer,
//      fasthttp.DefaultServeMux,
//      netfasthttp.OperationNameFunc(func(r *fasthttp.Request) string {
//	        return "HTTP " + r.Method + ":/api/customers"
//      }),
//      netfasthttp.MWSpanObserver(func(sp opentracing.Span, r *fasthttp.Request) {
//			sp.SetTag("fasthttp.uri", r.URL.EscapedPath())
//		}),
//   )
func Middleware(tr opentracing.Tracer, inner fasthttp.RequestHandler, options ...MWOption) fasthttp.RequestHandler {
	opts := mwOptions{
		opNameFunc: func(r *fasthttp.Request) string {
			return "HTTP " + string(r.Header.Method())
		},
		spanFilter:   func(r *fasthttp.Request) bool { return true },
		spanObserver: func(span opentracing.Span, r *fasthttp.Request) {},
		urlTagFunc: func(u *fasthttp.URI) string {
			return u.String()
		},
	}
	for _, opt := range options {
		opt(&opts)
	}
	// set component name, use "net/http" if caller does not specify
	componentName := opts.componentName
	if componentName == "" {
		componentName = defaultComponentName
	}
	return func(ctx *fasthttp.RequestCtx) {
		if !opts.spanFilter(&ctx.Request) {
			inner(ctx)
			return
		}
		clientContext, _ := tr.Extract(opentracing.HTTPHeaders, opentracefasthttp.New(&ctx.Request.Header))
		sp := tr.StartSpan(opts.opNameFunc(&ctx.Request), ext.RPCServerOption(clientContext))
		ext.HTTPMethod.Set(sp, string(ctx.Method()))
		ext.HTTPUrl.Set(sp, opts.urlTagFunc(ctx.URI()))
		ext.Component.Set(sp, componentName)
		opts.spanObserver(sp, &ctx.Request)
		ContextWithSpan(ctx, sp)

		defer func() {
			panicErr := recover()
			didPanic := panicErr != nil

			status := ctx.Response.StatusCode()
			if status == 0 && !didPanic {
				ctx.SetStatusCode(200)
			}
			if status > 0 {
				ext.HTTPStatusCode.Set(sp, uint16(status))
			}
			if status >= fasthttp.StatusInternalServerError || didPanic {
				ext.Error.Set(sp, true)
			}
			sp.Finish()

			if didPanic {
				panic(panicErr)
			}
		}()

		inner(ctx)
	}
}
