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
	"github.com/unionj-cloud/go-doudou/svc/logger"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"runtime/debug"
	"strings"
	"time"
)

// Metrics logs some metrics for http request
func Metrics(inner http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		m := httpsnoop.CaptureMetrics(inner, w, r)
		logger.WithFields(logrus.Fields{
			"remoteAddr": r.RemoteAddr,
			"httpMethod": r.Method,
			"requestUri": r.URL.RequestURI(),
			"requestUrl": r.URL.String(),
			"statusCode": m.Code,
			"written":    m.Written,
			"duration":   m.Duration.String(),
		}).Info(fmt.Sprintf("%s\t%s\t%s\t%d\t%d\t%s\n",
			r.RemoteAddr,
			r.Method,
			r.URL,
			m.Code,
			m.Written,
			m.Duration.String()))
	})
}

// Logger logs http request body and response body for debugging
func Logger(inner http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.RequestURI(), "/go-doudou/") || os.Getenv("GDD_LOG_LEVEL") != "debug" {
			inner.ServeHTTP(w, r)
			return
		}
		var reqBodyCopy io.ReadCloser
		if r.Body != nil {
			var buf bytes.Buffer
			if _, err := buf.ReadFrom(r.Body); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			if err := r.Body.Close(); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			r.Body, reqBodyCopy = io.NopCloser(&buf), io.NopCloser(bytes.NewReader(buf.Bytes()))
		}

		rec := httptest.NewRecorder()
		inner.ServeHTTP(rec, r)

		var reqBody string
		var reqContentType string
		if len(r.Header["Content-Type"]) > 0 {
			reqContentType = r.Header["Content-Type"][0]
		}
		if stringutils.IsNotEmpty(reqContentType) {
			if strings.Contains(reqContentType, "multipart/form-data") && reqBodyCopy != nil {
				r.Body = reqBodyCopy
				if err := r.ParseMultipartForm(32 << 20); err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}
				reqBody = r.Form.Encode()
			} else if strings.Contains(reqContentType, "application/json") {
				data := make(map[string]interface{})
				if err := json.NewDecoder(reqBodyCopy).Decode(&data); err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}
				b, _ := json.MarshalIndent(data, "", "    ")
				reqBody = string(b)
			} else {
				var buf bytes.Buffer
				if _, err := buf.ReadFrom(reqBodyCopy); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				reqBody = buf.String()
			}
		}
		start := time.Now()
		rid, _ := requestid.FromContext(r.Context())
		span := opentracing.SpanFromContext(r.Context())
		// Example:
		//  POST /usersvc/pageusers HTTP/1.1
		//  Host: localhost:6060
		//  Content-Length: 80
		//  Content-Type: application/json
		//  User-Agent: go-resty/2.6.0 (https://github.com/go-resty/resty)
		//  X-Request-Id: d1e4dc83-18be-493e-be5b-2e0faaca90ec
		//
		//  {"filter":{"dept":99,"name":"Jack"},"page":{"orders":null,"pageNo":2,"size":10}}
		contentType := rec.Result().Header.Get("Content-Type")
		var respBody string
		if strings.Contains(contentType, "text/html") || strings.Contains(contentType, "text/plain") {
			respBody = rec.Body.String()
		} else if strings.Contains(contentType, "application/json") {
			data := make(map[string]interface{})
			if err := json.NewDecoder(rec.Body).Decode(&data); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			b, _ := json.MarshalIndent(data, "", "    ")
			respBody = string(b)
		}
		fields := logrus.Fields{
			"remoteAddr":        r.RemoteAddr,
			"httpMethod":        r.Method,
			"requestUri":        r.URL.RequestURI(),
			"requestUrl":        r.URL.String(),
			"proto":             r.Proto,
			"host":              r.Host,
			"reqContentLength":  r.ContentLength,
			"reqHeader":         r.Header,
			"requestId":         rid,
			"reqBody":           reqBody,
			"respBody":          respBody,
			"statusCode":        rec.Result().StatusCode,
			"respHeader":        rec.Result().Header,
			"respContentLength": rec.Body.Len(),
			"elapsedTime":       time.Since(start).String(),
			"elapsed":           time.Since(start).Milliseconds(),
			"span":              fmt.Sprint(span),
		}
		log, _ := json.MarshalIndent(fields, "", "    ")
		logger.WithFields(fields).Debugln(string(log))

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
				logger.Errorf("panic: %+v\n\nstacktrace from panic: %s\n", e, string(debug.Stack()))
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
