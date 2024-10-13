package rest

// borrow code from https://github.com/TannerGabriel/learning-go/tree/master/advanced-programs/PrometheusHTTPServer
// Many thanks to TannerGabriel https://github.com/TannerGabriel
// Post link https://gabrieltanner.org/blog/collecting-prometheus-metrics-in-golang written by TannerGabriel
import (
	"net/http"
	"runtime"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/unionj-cloud/go-doudou/v2/framework/buildinfo"
	"github.com/unionj-cloud/toolkit/constants"
	"github.com/unionj-cloud/toolkit/stringutils"
)

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

// NewResponseWriter creates new responseWriter
func NewResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{w, http.StatusOK}
}

// WriteHeader set header to code
func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

var countRequests = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "go_doudou_http_request_count",
		Help: "Number of http requests.",
	},
	[]string{"path", "method", "status"},
)

var httpDuration = prometheus.NewHistogramVec(prometheus.HistogramOpts{
	Name: "go_doudou_http_response_time_seconds",
	Help: "Duration of HTTP requests.",
}, []string{"path", "method"})

// PrometheusMiddleware returns http HandlerFunc for prometheus matrix
func PrometheusMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		method := r.Method
		timer := prometheus.NewTimer(httpDuration.WithLabelValues(path, method))

		rw := NewResponseWriter(w)
		next.ServeHTTP(rw, r)

		statusCode := rw.statusCode

		countRequests.WithLabelValues(path, method, strconv.Itoa(statusCode)).Inc()

		timer.ObserveDuration()
	})
}

func init() {
	prometheus.Register(countRequests)
	prometheus.Register(httpDuration)
	buildTime := buildinfo.BuildTime
	if stringutils.IsNotEmpty(buildinfo.BuildTime) {
		if t, err := time.Parse(constants.FORMAT15, buildinfo.BuildTime); err == nil {
			buildTime = t.Local().Format(constants.FORMAT8)
		}
	}
	appInfo := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "go_doudou_app_info",
		Help: "Information about the go-doudou app",
		ConstLabels: prometheus.Labels{
			"goVer":     runtime.Version(),
			"gddVer":    buildinfo.GddVer,
			"buildUser": buildinfo.BuildUser,
			"buildTime": buildTime,
		},
	})
	appInfo.Set(1)
	prometheus.Register(appInfo)

	prometheus.Register(prometheus.NewGaugeFunc(prometheus.GaugeOpts{
		Name:        "go_doudou_gomaxprocs",
		Help:        "The value of gomaxprocs",
		ConstLabels: nil,
	}, func() float64 {
		return float64(runtime.GOMAXPROCS(0))
	}))

	prometheus.Register(prometheus.NewGaugeFunc(prometheus.GaugeOpts{
		Name:        "go_doudou_numcpu",
		Help:        "The value of numcpu",
		ConstLabels: nil,
	}, func() float64 {
		return float64(runtime.NumCPU())
	}))
}
