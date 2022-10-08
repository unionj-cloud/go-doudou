package fast

import (
	"context"
	"encoding/base64"
	"fmt"
	scrypt "github.com/elithrar/simple-scrypt"
	realip "github.com/ferluci/fast-realip"
	"github.com/google/uuid"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"github.com/slok/goresilience"
	"github.com/slok/goresilience/bulkhead"
	"github.com/uber/jaeger-client-go"
	"github.com/unionj-cloud/go-doudou/framework/http/fast/trace"
	"github.com/unionj-cloud/go-doudou/framework/http/model"
	"github.com/unionj-cloud/go-doudou/framework/internal/config"
	"github.com/unionj-cloud/go-doudou/toolkit/stringutils"
	logger "github.com/unionj-cloud/go-doudou/toolkit/zlogger"
	"github.com/valyala/fasthttp"
	"net/http"
	"net/url"
	"runtime/debug"
	"strings"
	"time"
)

type MiddlewareFunc func(inner fasthttp.RequestHandler) fasthttp.RequestHandler

// Middleware allows MiddlewareFunc to implement the middleware interface.
func (mw MiddlewareFunc) Middleware(handler fasthttp.RequestHandler) fasthttp.RequestHandler {
	return mw(handler)
}

var realIpKey = []byte("realIpKey")

func ContextWithRealIp(ctx *fasthttp.RequestCtx, ip string) {
	ctx.SetUserValueBytes(realIpKey, ip)
}

func RealIpFromContext(ctx *fasthttp.RequestCtx) string {
	return ctx.UserValueBytes(realIpKey).(string)
}

func RealIp(inner fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		ContextWithRealIp(ctx, realip.FromRequest(ctx))
		inner(ctx)
	}
}

// Metrics logs some metrics for http request
func Metrics(inner fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		start := time.Now()
		inner(ctx)
		logger.Info().
			Str("ip", RealIpFromContext(ctx)).
			Str("method", string(ctx.Method())).
			Str("url", string(ctx.RequestURI())).
			Int("code", ctx.Response.StatusCode()).
			Int("size", ctx.Response.Header.ContentLength()).
			Str("duration", time.Since(start).String())
	}
}

// ReqLog logs http request body and response body for debugging
func ReqLog(inner fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		start := time.Now()
		inner(ctx)
		elapsed := time.Since(start)
		reqQuery := string(ctx.Request.URI().QueryString())
		if unescape, err := url.QueryUnescape(reqQuery); err == nil {
			reqQuery = unescape
		}
		var traceId string
		span := trace.SpanFromContext(ctx)
		if jspan, ok := span.(*jaeger.Span); ok {
			traceId = jspan.SpanContext().TraceID().String()
		}
		fields := map[string]interface{}{
			"remoteAddr":        RealIpFromContext(ctx),
			"requestId":         RidFromContext(ctx),
			"httpMethod":        ctx.Method(),
			"requestUrl":        string(ctx.RequestURI()),
			"host":              ctx.Host(),
			"reqContentLength":  ctx.Request.Header.ContentLength(),
			"reqHeader":         string(ctx.Request.Header.RawHeaders()),
			"reqQuery":          reqQuery,
			"reqBody":           string(ctx.PostBody()),
			"respBody":          string(ctx.Response.Body()),
			"statusCode":        ctx.Response.StatusCode(),
			"respHeader":        string(ctx.Response.Header.Header()),
			"respContentLength": ctx.Response.Header.ContentLength(),
			"elapsedTime":       elapsed.String(),
			"elapsed":           elapsed.Milliseconds(),
			"span":              span,
			"traceId":           traceId,
		}
		var reqLog string
		var err error
		if reqLog, err = model.JsonMarshalIndent(fields, "", "    ", true); err != nil {
			reqLog = fmt.Sprintf("call jsonMarshalIndent(fields, \"\", \"    \", true) error: %s", err)
		}
		logger.Info().Fields(fields).Msg(reqLog)
	}
}

// Rest set Content-Type to application/json
func Rest(inner fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		if stringutils.IsEmpty(string(ctx.Response.Header.ContentType())) {
			ctx.SetContentType("application/json; charset=UTF-8")
		}
		inner(ctx)
	}
}

// FallbackContentType set fallback response Content-Type to contentType
func FallbackContentType(contentType string) func(inner fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(inner fasthttp.RequestHandler) fasthttp.RequestHandler {
		return func(ctx *fasthttp.RequestCtx) {
			if stringutils.IsEmpty(string(ctx.Response.Header.ContentType())) {
				ctx.SetContentType(contentType)
			}
			inner(ctx)
		}
	}
}

// basicAuth returns the username and password provided in the request's
// Authorization header, if the request uses HTTP Basic Authentication.
// See RFC 2617, Section 2.
func peekBasicAuth(ctx *fasthttp.RequestCtx) (username, password string, ok bool) {
	auth := ctx.Request.Header.Peek("Authorization")
	if auth == nil {
		return
	}
	return parseBasicAuth(string(auth))
}

// parseBasicAuth parses an HTTP Basic Authentication string.
// "Basic QWxhZGRpbjpvcGVuIHNlc2FtZQ==" returns ("Aladdin", "open sesame", true).
func parseBasicAuth(auth string) (username, password string, ok bool) {
	const prefix = "Basic "
	if !strings.HasPrefix(auth, prefix) {
		return
	}
	c, err := base64.StdEncoding.DecodeString(auth[len(prefix):])
	if err != nil {
		return
	}
	cs := string(c)
	s := strings.IndexByte(cs, ':')
	if s < 0 {
		return
	}
	return cs[:s], cs[s+1:], true
}

// BasicAuth adds http basic auth validation
func BasicAuth() func(inner fasthttp.RequestHandler) fasthttp.RequestHandler {
	requiredUser := config.DefaultGddManageUser
	pass := config.DefaultGddManagePass
	return func(inner fasthttp.RequestHandler) fasthttp.RequestHandler {
		return func(ctx *fasthttp.RequestCtx) {
			if stringutils.IsNotEmpty(config.GddManageUser.Load()) {
				requiredUser = config.GddManageUser.Load()
			}
			if stringutils.IsNotEmpty(config.GddManagePass.Load()) {
				pass = config.GddManagePass.Load()
			}
			if requiredPasswordHash, err := scrypt.GenerateFromPassword([]byte(pass), scrypt.DefaultParams); err == nil {
				// Get the Basic Authentication credentials
				user, password, hasAuth := peekBasicAuth(ctx)

				// WARNING:
				// DO NOT use plain-text passwords for real apps.
				// A simple string comparison using == is vulnerable to a timing attack.
				// Instead, use the hash comparison function found in your hash library.
				// This example uses scrypt, which is a solid choice for secure hashing:
				//   go get -u github.com/elithrar/simple-scrypt

				if hasAuth && user == requiredUser {

					// Uses the parameters from the existing derived key. Return an error if they don't match.
					if err = scrypt.CompareHashAndPassword(requiredPasswordHash, []byte(password)); err != nil {

						// log error and request Basic Authentication again below.
						logger.Error().Err(err).Msg("")

					} else {

						// Delegate request to the given handle
						inner(ctx)
						return

					}

				}
			} else {
				logger.Error().Err(err).Msg("")
			}

			// Request Basic Authentication otherwise
			ctx.Error(fasthttp.StatusMessage(fasthttp.StatusUnauthorized), fasthttp.StatusUnauthorized)
			ctx.Response.Header.Set("WWW-Authenticate", "Basic realm=Restricted")
		}
	}
}

// Recovery handles panic from processing incoming http request
func Recovery(inner fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		defer func() {
			if e := recover(); e != nil {
				statusCode := fasthttp.StatusInternalServerError
				respErr := errors.New(fmt.Sprint(e))
				if err, ok := e.(error); ok {
					if errors.Is(err, context.Canceled) {
						statusCode = fasthttp.StatusBadRequest
					} else {
						var bizErr model.BizError
						if errors.As(err, &bizErr) {
							statusCode = bizErr.StatusCode
							if bizErr.Cause != nil {
								e = bizErr.Cause
							}
							respErr = bizErr
						}
					}
				}
				logger.Error().Err(respErr).Msgf("panic: %+v\n\nstacktrace from panic: %s\n", e, string(debug.Stack()))
				ctx.Error(respErr.Error(), statusCode)
				return
			}
		}()
		inner(ctx)
	}
}

// Tracing add jaeger tracing middleware
func Tracing(inner fasthttp.RequestHandler) fasthttp.RequestHandler {
	return trace.Middleware(
		opentracing.GlobalTracer(),
		inner,
		trace.OperationNameFunc(func(r *fasthttp.Request) string {
			return fmt.Sprintf("HTTP %s: %s", r.Header.Method(), r.URI().Path())
		}))
}

var RunnerChain = goresilience.RunnerChain

// BulkHead add bulk head pattern middleware based on https://github.com/slok/goresilience
// workers is the number of workers in the execution pool.
// maxWaitTime is the max time an incoming request will wait to execute before being dropped its execution and return 429 response.
func BulkHead(workers int, maxWaitTime time.Duration) func(inner fasthttp.RequestHandler) fasthttp.RequestHandler {
	runner := RunnerChain(
		bulkhead.NewMiddleware(bulkhead.Config{
			Workers:     workers,
			MaxWaitTime: maxWaitTime,
		}),
	)
	return func(inner fasthttp.RequestHandler) fasthttp.RequestHandler {
		return func(ctx *fasthttp.RequestCtx) {
			err := runner.Run(ctx, func(_ context.Context) error {
				inner(ctx)
				return nil
			})
			if err != nil {
				ctx.Error("too many requests", http.StatusTooManyRequests)
			}
		}
	}
}

var ridKey = []byte("ridKey")

func ContextWithRid(ctx *fasthttp.RequestCtx, rid []byte) {
	ctx.SetUserValueBytes(ridKey, rid)
}

func RidFromContext(ctx *fasthttp.RequestCtx) []byte {
	return ctx.UserValueBytes(ridKey).([]byte)
}

// RequestIDHandler sets unique request id.
// If header `X-Request-ID` is already present in the request, that is considered the
// request id. Otherwise, generates a new unique ID.
func RequestIDHandler(inner fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		rid := ctx.Request.Header.Peek("X-Request-ID")
		if len(rid) == 0 {
			rid = []byte(uuid.New().String())
			ctx.Request.Header.SetBytesV("X-Request-ID", rid)
		}
		ContextWithRid(ctx, rid)
		inner(ctx)
	}
}

func Cors(inner fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {

	}
}
