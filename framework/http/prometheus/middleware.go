package prometheus

// borrow code from https://github.com/TannerGabriel/learning-go/tree/master/advanced-programs/PrometheusHTTPServer
// Many thanks to TannerGabriel https://github.com/TannerGabriel
// Post link https://gabrieltanner.org/blog/collecting-prometheus-metrics-in-golang written by TannerGabriel
import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"net/http"
	"strconv"
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

var httpDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
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
}
