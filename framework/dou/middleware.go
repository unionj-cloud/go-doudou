package dou

import (
	"bytes"
	"context"
	"crypto/subtle"
	"encoding/json"
	"fmt"
	"github.com/apolloconfig/agollo/v4/storage"
	"github.com/ascarter/requestid"
	"github.com/felixge/httpsnoop"
	"github.com/opentracing-contrib/go-stdlib/nethttp"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/slok/goresilience"
	"github.com/slok/goresilience/bulkhead"
	"github.com/uber/jaeger-client-go"
	"github.com/unionj-cloud/go-doudou/framework/configmgr"
	"github.com/unionj-cloud/go-doudou/framework/internal/config"
	"github.com/unionj-cloud/go-doudou/framework/logger"
	"github.com/unionj-cloud/go-doudou/toolkit/stringutils"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"runtime/debug"
	"strings"
	"time"
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
				DataId:   "__" + dataId + "__" + "ddhttp",
				OnChange: CallbackOnChange(listener),
			})
		}
	case config.ApolloConfigType:
		configmgr.ApolloClient.AddChangeListener(listener)
	default:
		logrus.Warnf("[go-doudou] from ddhttp pkg: unknown config type: %s\n", configType)
	}
}

func init() {
	InitialiseRemoteConfigListener()
}

// metrics logs some metrics for http request
func metrics(inner http.Handler) http.Handler {
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

// borrowed from httputil unexported function drainBody
func copyReqBody(b io.ReadCloser) (r1, r2 io.ReadCloser, err error) {
	if b == nil || b == http.NoBody {
		// No copying needed. Preserve the magic sentinel meaning of NoBody.
		return http.NoBody, http.NoBody, nil
	}
	var buf bytes.Buffer
	if _, err = buf.ReadFrom(b); err != nil {
		return nil, b, err
	}
	if err = b.Close(); err != nil {
		return nil, b, err
	}
	return ioutil.NopCloser(&buf), ioutil.NopCloser(bytes.NewReader(buf.Bytes())), nil
}

func copyRespBody(b *bytes.Buffer) (b1, b2 *bytes.Buffer, err error) {
	if b == nil {
		return
	}
	var buf bytes.Buffer
	if _, err = buf.ReadFrom(b); err != nil {
		return nil, b, err
	}
	return &buf, bytes.NewBuffer(buf.Bytes()), nil
}

func jsonMarshalIndent(data interface{}, prefix, indent string, disableHTMLEscape bool) (string, error) {
	b := &bytes.Buffer{}
	encoder := json.NewEncoder(b)
	encoder.SetEscapeHTML(!disableHTMLEscape)
	encoder.SetIndent(prefix, indent)
	if err := encoder.Encode(data); err != nil {
		return "", errors.Errorf("failed to marshal data to JSON, %s", err)
	}
	return b.String(), nil
}

func getReqBody(cp io.ReadCloser, r *http.Request) string {
	var contentType string
	if len(r.Header["Content-Type"]) > 0 {
		contentType = r.Header["Content-Type"][0]
	}
	var reqBody string
	if cp != nil {
		if strings.Contains(contentType, "multipart/form-data") {
			r.Body = cp
			if err := r.ParseMultipartForm(32 << 20); err == nil {
				reqBody = r.Form.Encode()
				if unescape, err := url.QueryUnescape(reqBody); err == nil {
					reqBody = unescape
				}
			} else {
				logger.Debug("call r.ParseMultipartForm(32 << 20) error: ", err)
			}
		} else if strings.Contains(contentType, "application/json") {
			data := make(map[string]interface{})
			if err := json.NewDecoder(cp).Decode(&data); err == nil {
				b, _ := json.MarshalIndent(data, "", "    ")
				reqBody = string(b)
			} else {
				logger.Debug("call json.NewDecoder(reqBodyCopy).Decode(&data) error: ", err)
			}
		} else {
			var buf bytes.Buffer
			if _, err := buf.ReadFrom(cp); err == nil {
				data := []rune(buf.String())
				end := len(data)
				if end > 1000 {
					end = 1000
				}
				reqBody = string(data[:end])
				if strings.Contains(contentType, "application/x-www-form-urlencoded") {
					if unescape, err := url.QueryUnescape(reqBody); err == nil {
						reqBody = unescape
					}
				}
			} else {
				logger.Debug("call buf.ReadFrom(reqBodyCopy) error: ", err)
			}
		}
	}
	return reqBody
}

func getRespBody(rec *httptest.ResponseRecorder) string {
	var (
		respBody string
		err      error
	)
	if strings.Contains(rec.Result().Header.Get("Content-Type"), "application/json") {
		var respBodyCopy *bytes.Buffer
		if respBodyCopy, rec.Body, err = copyRespBody(rec.Body); err == nil {
			data := make(map[string]interface{})
			if err := json.NewDecoder(rec.Body).Decode(&data); err == nil {
				b, _ := json.MarshalIndent(data, "", "    ")
				respBody = string(b)
			} else {
				logger.Debug("call json.NewDecoder(rec.Body).Decode(&data) error: ", err)
			}
		} else {
			logger.Debug("call respBodyCopy.ReadFrom(rec.Body) error: ", err)
		}
		rec.Body = respBodyCopy
	} else {
		data := []rune(rec.Body.String())
		end := len(data)
		if end > 1000 {
			end = 1000
		}
		respBody = string(data[:end])
	}
	return respBody
}

// log logs http request body and response body for debugging
func log(inner http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			reqBodyCopy io.ReadCloser
			err         error
			traceId     string
		)
		if reqBodyCopy, r.Body, err = copyReqBody(r.Body); err != nil {
			logger.Debug("call copyReqBody(r.Body) error: ", err)
		}

		rec := httptest.NewRecorder()
		inner.ServeHTTP(rec, r)

		reqBody := getReqBody(reqBodyCopy, r)
		start := time.Now()
		rid, _ := requestid.FromContext(r.Context())
		span := opentracing.SpanFromContext(r.Context())
		if jspan, ok := span.(*jaeger.Span); ok {
			traceId = jspan.SpanContext().TraceID().String()
		}
		respBody := getRespBody(rec)
		reqQuery := r.URL.RawQuery
		if unescape, err := url.QueryUnescape(reqQuery); err == nil {
			reqQuery = unescape
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
			"reqQuery":          reqQuery,
			"reqBody":           reqBody,
			"respBody":          respBody,
			"statusCode":        rec.Result().StatusCode,
			"respHeader":        rec.Result().Header,
			"respContentLength": rec.Body.Len(),
			"elapsedTime":       time.Since(start).String(),
			"elapsed":           time.Since(start).Milliseconds(),
			"span":              fmt.Sprint(span),
			"traceId":           traceId,
		}
		var reqLog string
		if reqLog, err = jsonMarshalIndent(fields, "", "    ", true); err != nil {
			reqLog = fmt.Sprintf("call jsonMarshalIndent(fields, \"\", \"    \", true) error: %s", err)
		}
		logger.WithFields(fields).Infoln(reqLog)

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
