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

func NewResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{w, http.StatusOK}
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

var TotalRequests = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Number of get requests.",
	},
	[]string{"path", "method"},
)

var responseStatus = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "response_status",
		Help: "Status of HTTP response",
	},
	[]string{"status"},
)

var httpDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
	Name: "http_response_time_seconds",
	Help: "Duration of HTTP requests.",
}, []string{"path", "method"})

func PrometheusMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		method := r.Method
		timer := prometheus.NewTimer(httpDuration.WithLabelValues(path, method))

		rw := NewResponseWriter(w)
		next.ServeHTTP(rw, r)

		statusCode := rw.statusCode

		responseStatus.WithLabelValues(strconv.Itoa(statusCode)).Inc()
		TotalRequests.WithLabelValues(path, method).Inc()

		timer.ObserveDuration()
	})
}

func init() {
	prometheus.Register(TotalRequests)
	prometheus.Register(responseStatus)
	prometheus.Register(httpDuration)
}
