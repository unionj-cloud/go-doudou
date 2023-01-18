package rest

import (
	"context"
	"crypto/subtle"
	"fmt"
	"github.com/apolloconfig/agollo/v4/storage"
	"github.com/ascarter/requestid"
	"github.com/felixge/httpsnoop"
	"github.com/klauspost/compress/gzip"
	"github.com/opentracing-contrib/go-stdlib/nethttp"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"github.com/slok/goresilience"
	"github.com/slok/goresilience/bulkhead"
	"github.com/uber/jaeger-client-go"
	"github.com/unionj-cloud/go-doudou/v2/framework/configmgr"
	"github.com/unionj-cloud/go-doudou/v2/framework/internal/config"
	"github.com/unionj-cloud/go-doudou/v2/toolkit/stringutils"
	logger "github.com/unionj-cloud/go-doudou/v2/toolkit/zlogger"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"runtime/debug"
	"strings"
	"time"
)

var (
	Tracing             = tracing
	Metrics             = metrics
	Log                 = log
	FallbackContentType = fallbackContentType
	BasicAuth           = basicAuth
	Recovery            = recovery
)

type httpConfigListener struct {
	configmgr.BaseApolloListener
}

func NewHttpConfigListener() *httpConfigListener {
	return &httpConfigListener{}
}

func (c *httpConfigListener) OnChange(event *storage.ChangeEvent) {
	c.Lock.Lock()
	defer c.Lock.Unlock()
	if !c.SkippedFirstEvent {
		c.SkippedFirstEvent = true
		return
	}
	for key, value := range event.Changes {
		upperKey := strings.ToUpper(strings.ReplaceAll(key, ".", "_"))
		if strings.HasPrefix(upperKey, "GDD_MANAGE_") {
			_ = os.Setenv(upperKey, fmt.Sprint(value.NewValue))
		}
	}
}

func CallbackOnChange(listener *httpConfigListener) func(event *configmgr.NacosChangeEvent) {
	return func(event *configmgr.NacosChangeEvent) {
		changes := make(map[string]*storage.ConfigChange)
		for k, v := range event.Changes {
			changes[k] = &storage.ConfigChange{
				OldValue:   v.OldValue,
				NewValue:   v.NewValue,
				ChangeType: storage.ConfigChangeType(v.ChangeType),
			}
		}
		changeEvent := &storage.ChangeEvent{
			Changes: changes,
		}
		listener.OnChange(changeEvent)
	}
}

func InitialiseRemoteConfigListener() {
	listener := &httpConfigListener{}
	configType := config.GddConfigRemoteType.LoadOrDefault(config.DefaultGddConfigRemoteType)
	switch configType {
	case "":
		return
	case config.NacosConfigType:
		dataIdStr := config.GddNacosConfigDataid.LoadOrDefault(config.DefaultGddNacosConfigDataid)
		dataIds := strings.Split(dataIdStr, ",")
		listener.SkippedFirstEvent = true
		for _, dataId := range dataIds {
			configmgr.NacosClient.AddChangeListener(configmgr.NacosConfigListenerParam{
				DataId:   "__" + dataId + "__" + "rest",
				OnChange: CallbackOnChange(listener),
			})
		}
	case config.ApolloConfigType:
		configmgr.ApolloClient.AddChangeListener(listener)
	default:
		logger.Warn().Msgf("[go-doudou] unknown config type: %s\n", configType)
	}
}

func init() {
	InitialiseRemoteConfigListener()
}

// metrics logs some metrics for http request
func metrics(inner http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		m := httpsnoop.CaptureMetrics(inner, w, r)
		logger.Info().
			Msgf("%s\t%s\t%s\t%d\t%d\t%s", r.RemoteAddr,
				r.Method,
				r.URL,
				m.Code,
				m.Written,
				m.Duration.String())

	})
}

// log logs http request body and response body for debugging
func log(inner http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			reqBodyCopy io.ReadCloser
			err         error
			traceId     string
		)
		if reqBodyCopy, r.Body, err = CopyReqBody(r.Body); err != nil {
			logger.Error().Err(err).Msg("call copyReqBody(r.Body) error")
		}

		rec := httptest.NewRecorder()
		start := time.Now()
		inner.ServeHTTP(rec, r)
		elapsed := time.Since(start)
		reqBody := GetReqBody(reqBodyCopy, r)
		rid, _ := requestid.FromContext(r.Context())
		span := opentracing.SpanFromContext(r.Context())
		if jspan, ok := span.(*jaeger.Span); ok {
			traceId = jspan.SpanContext().TraceID().String()
		}
		respBody := GetRespBody(rec)
		reqQuery := r.URL.RawQuery
		if unescape, err := url.QueryUnescape(reqQuery); err == nil {
			reqQuery = unescape
		}
		fields := map[string]interface{}{
			"remoteAddr":        r.RemoteAddr,
			"httpMethod":        r.Method,
			"requestUrl":        r.URL.String(),
			"proto":             r.Proto,
			"host":              r.Host,
			"reqContentLength":  r.ContentLength,
			"reqHeader":         r.Header,
			"requestId":         rid,
			"reqQuery":          reqQuery,
			"reqBody":           reqBody,
			"respBody":          respBody,
			"statusCode":        rec.Result().StatusCode,
			"respHeader":        rec.Result().Header,
			"respContentLength": rec.Body.Len(),
			"elapsedTime":       elapsed.String(),
			"elapsed":           elapsed.Milliseconds(),
			"span":              span,
			"traceId":           traceId,
		}
		var reqLog string
		if reqLog, err = JsonMarshalIndent(fields, "", "    ", true); err != nil {
			reqLog = fmt.Sprintf("call jsonMarshalIndent(fields, \"\", \"    \", true) error: %s", err)
		}
		logger.Info().Fields(fields).Msg(reqLog)
		header := rec.Result().Header
		for k, v := range header {
			w.Header()[k] = v
		}
		w.WriteHeader(rec.Result().StatusCode)
		rec.Body.WriteTo(w)
	})
}

// rest set Content-Type to application/json
func rest(inner http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if stringutils.IsEmpty(w.Header().Get("Content-Type")) {
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		}
		inner.ServeHTTP(w, r)
	})
}

// fallbackContentType set fallback response Content-Type to contentType
func fallbackContentType(contentType string) func(inner http.Handler) http.Handler {
	return func(inner http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if stringutils.IsEmpty(w.Header().Get("Content-Type")) {
				w.Header().Set("Content-Type", contentType)
			}
			inner.ServeHTTP(w, r)
		})
	}
}

// basicAuth adds http basic auth validation
func basicAuth() func(inner http.Handler) http.Handler {
	username := config.DefaultGddManageUser
	password := config.DefaultGddManagePass
	return func(inner http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if stringutils.IsNotEmpty(config.GddManageUser.Load()) {
				username = config.GddManageUser.Load()
			}
			if stringutils.IsNotEmpty(config.GddManagePass.Load()) {
				password = config.GddManagePass.Load()
			}
			user, pass, ok := r.BasicAuth()
			if !ok || subtle.ConstantTimeCompare([]byte(user), []byte(username)) != 1 || subtle.ConstantTimeCompare([]byte(pass), []byte(password)) != 1 {
				w.Header().Set("WWW-Authenticate", `Basic realm="Provide user name and password"`)
				w.WriteHeader(401)
				w.Write([]byte("Unauthorised.\n"))
				return
			}
			inner.ServeHTTP(w, r)
		})
	}
}

// recovery handles panic from processing incoming http request
func recovery(inner http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if e := recover(); e != nil {
				statusCode := http.StatusInternalServerError
				respErr := fmt.Sprintf("%v", e)
				if err, ok := e.(error); ok {
					if errors.Is(err, context.Canceled) {
						statusCode = http.StatusBadRequest
					} else {
						var bizErr BizError
						if errors.As(err, &bizErr) {
							statusCode = bizErr.StatusCode
							if bizErr.Cause != nil {
								e = bizErr.Cause
							}
							respErr = bizErr.Error()
						}
					}
				}
				logger.Error().Msgf("panic: %+v\n\nstacktrace from panic: %s\n", e, string(debug.Stack()))
				http.Error(w, respErr, statusCode)
			}
		}()
		inner.ServeHTTP(w, r)
	})
}

// gzipBody handles gzip-ed request body
func gzipBody(inner http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Handle gzip decoding of the body
		if r.Header.Get("Content-Encoding") == "gzip" {
			b, err := gzip.NewReader(r.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			defer b.Close()
			r.Body = b
		}
		inner.ServeHTTP(w, r)
	})
}

// tracing add jaeger tracing middleware
func tracing(inner http.Handler) http.Handler {
	return nethttp.Middleware(
		opentracing.GlobalTracer(),
		inner,
		nethttp.OperationNameFunc(func(r *http.Request) string {
			return fmt.Sprintf("HTTP %s: %s", r.Method, r.URL.Path)
		}))
}

var RunnerChain = goresilience.RunnerChain

// BulkHead add bulk head pattern middleware based on https://github.com/slok/goresilience
// workers is the number of workers in the execution pool.
// maxWaitTime is the max time an incoming request will wait to execute before being dropped its execution and return 429 response.
func BulkHead(workers int, maxWaitTime time.Duration) func(inner http.Handler) http.Handler {
	runner := RunnerChain(
		bulkhead.NewMiddleware(bulkhead.Config{
			Workers:     workers,
			MaxWaitTime: maxWaitTime,
		}),
	)
	return func(inner http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			err := runner.Run(r.Context(), func(_ context.Context) error {
				inner.ServeHTTP(w, r)
				return nil
			})
			if err != nil {
				http.Error(w, "too many requests", http.StatusTooManyRequests)
			}
		})
	}
}

func BodyMaxBytes(n int64) func(inner http.Handler) http.Handler {
	return func(inner http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r2 := *r
			r2.Body = http.MaxBytesReader(w, r.Body, n)
			inner.ServeHTTP(w, &r2)
		})
	}
}
