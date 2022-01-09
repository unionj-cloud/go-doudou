package ddhttp

import (
	"bytes"
	"context"
	"crypto/subtle"
	"encoding/json"
	"fmt"
	"github.com/ascarter/requestid"
	"github.com/felixge/httpsnoop"
	"github.com/opentracing-contrib/go-stdlib/nethttp"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/slok/goresilience"
	"github.com/slok/goresilience/bulkhead"
	"github.com/unionj-cloud/go-doudou/stringutils"
	"github.com/unionj-cloud/go-doudou/svc/config"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"
	"time"
)

// Metrics logs some metrics for http request
func Metrics(inner http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		m := httpsnoop.CaptureMetrics(inner, w, r)
		logrus.Printf(
			"%s\t%s\t%s\t%d\t%d\t%s\n",
			r.RemoteAddr,
			r.Method,
			r.URL,
			m.Code,
			m.Written,
			m.Duration,
		)
	})
}

// Logger logs http request body and response body for debugging
func Logger(inner http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.RequestURI(), "/go-doudou/") || os.Getenv("GDD_LOG_LEVEL") != "debug" {
			inner.ServeHTTP(w, r)
			return
		}
		x, err := httputil.DumpRequest(r, true)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		rec := httptest.NewRecorder()
		inner.ServeHTTP(rec, r)

		rawReq := string(x)
		if len(r.Header["Content-Type"]) > 0 && strings.Contains(r.Header["Content-Type"][0], "multipart/form-data") {
			r.Body = ioutil.NopCloser(bytes.NewReader(x))
			if err := r.ParseMultipartForm(32 << 20); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			rawReq = r.Form.Encode()
		}
		start := time.Now()
		rid, _ := requestid.FromContext(r.Context())
		span := opentracing.SpanFromContext(r.Context())
		hlog := HttpLog{
			ClientIp:          r.RemoteAddr,
			HttpMethod:        r.Method,
			Uri:               r.URL.RequestURI(),
			Proto:             r.Proto,
			Host:              r.Host,
			ReqContentLength:  r.ContentLength,
			ReqHeader:         r.Header,
			RequestId:         rid,
			RawReq:            rawReq,
			RespBody:          rec.Body.String(),
			StatusCode:        rec.Result().StatusCode,
			RespHeader:        rec.Result().Header,
			RespContentLength: rec.Body.Len(),
			ElapsedTime:       time.Since(start).String(),
			Elapsed:           time.Since(start).Milliseconds(),
			Span:              fmt.Sprint(span),
		}
		log, _ := json.MarshalIndent(hlog, "", "    ")
		logrus.Debugln(string(log))

		header := rec.Result().Header
		for k, v := range header {
			w.Header()[k] = v
		}
		w.WriteHeader(rec.Result().StatusCode)
		rec.Body.WriteTo(w)
	})
}

// Rest set Content-Type to application/json
func Rest(inner http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if stringutils.IsEmpty(w.Header().Get("Content-Type")) {
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		}
		inner.ServeHTTP(w, r)
	})
}

// BasicAuth adds http basic auth validation
func BasicAuth(inner http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username := config.GddManageUser.Load()
		password := config.GddManagePass.Load()
		if stringutils.IsNotEmpty(username) || stringutils.IsNotEmpty(password) {
			user, pass, ok := r.BasicAuth()

			if !ok || subtle.ConstantTimeCompare([]byte(user), []byte(username)) != 1 || subtle.ConstantTimeCompare([]byte(pass), []byte(password)) != 1 {
				w.Header().Set("WWW-Authenticate", `Basic realm="Provide user name and password"`)
				w.WriteHeader(401)
				w.Write([]byte("Unauthorised.\n"))
				return
			}
		}
		inner.ServeHTTP(w, r)
	})
}

// Recover handles panic from processing incoming http request
func Recover(inner http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if e := recover(); e != nil {
				statusCode := http.StatusInternalServerError
				if err, ok := e.(error); ok {
					if errors.Is(err, context.Canceled) {
						statusCode = http.StatusBadRequest
					}
				}
				logrus.Errorf("panic: %+v\n\nstacktrace from panic: %s\n", e, string(debug.Stack()))
				http.Error(w, fmt.Sprintf("%v", e), statusCode)
			}
		}()
		inner.ServeHTTP(w, r)
	})
}

// Tracing add jaeger tracing middleware
func Tracing(inner http.Handler) http.Handler {
	return nethttp.Middleware(
		opentracing.GlobalTracer(),
		inner,
		nethttp.OperationNameFunc(func(r *http.Request) string {
			return "HTTP " + r.Method + " " + r.URL.Path
		}))
}

// BulkHead add bulk head pattern middleware based on https://github.com/slok/goresilience
// workers is the number of workers in the execution pool.
// maxWaitTime is the max time an incoming request will wait to execute before being dropped its execution and return 429 response.
func BulkHead(workers int, maxWaitTime time.Duration) func(inner http.Handler) http.Handler {
	runner := goresilience.RunnerChain(
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
