package prometheus

// borrow code from https://github.com/TannerGabriel/learning-go/tree/master/advanced-programs/PrometheusHTTPServer
// Many thanks to TannerGabriel https://github.com/TannerGabriel
// Post link https://gabrieltanner.org/blog/collecting-prometheus-metrics-in-golang written by TannerGabriel
import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/unionj-cloud/go-doudou/framework/buildinfo"
	logger "github.com/unionj-cloud/go-doudou/toolkit/zlogger"
	"github.com/unionj-cloud/go-doudou/toolkit/load"
	"github.com/unionj-cloud/go-doudou/toolkit/process"
	"net/http"
	"runtime"
	"strconv"
	"sync"
	"time"
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

var processPool sync.Pool

func init() {
	processPool = sync.Pool{
		New: func() interface{} {
			return process.NewCurrentProcess()
		},
	}
	prometheus.Register(countRequests)
	prometheus.Register(httpDuration)
	appInfo := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "go_doudou_app_info",
		Help: "Information about the go-doudou app",
		ConstLabels: prometheus.Labels{
			"goVer":     runtime.Version(),
			"gddVer":    buildinfo.GddVer,
			"buildUser": buildinfo.BuildUser,
			"buildTime": buildinfo.BuildTime,
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

	prometheus.Register(prometheus.NewGaugeFunc(prometheus.GaugeOpts{
		Name:        "go_doudou_system_cpu_usage",
		Help:        "The \"recent cpu usage\" for the whole system",
		ConstLabels: nil,
	}, func() float64 {
		p := processPool.Get().(*process.Process)
		defer processPool.Put(p)
		if p.Err != nil {
			return 0
		}
		ts, err := p.Times1()
		if err != nil {
			logger.Error().Err(err).Msg("")
			return 0
		}
		return ts.System
	}))
	prometheus.Register(prometheus.NewGaugeFunc(prometheus.GaugeOpts{
		Name:        "go_doudou_process_cpu_usage",
		Help:        "The \"recent cpu usage\" for the go-doudou process",
		ConstLabels: nil,
	}, func() float64 {
		p := processPool.Get().(*process.Process)
		defer processPool.Put(p)
		if p.Err != nil {
			return 0
		}
		ts, err := p.Times1()
		if err != nil {
			logger.Error().Err(err).Msg("")
			return 0
		}
		return ts.User
	}))
	prometheus.Register(prometheus.NewGaugeFunc(prometheus.GaugeOpts{
		Name:        "go_doudou_process_start_time_millis",
		Help:        "Start time of the process since unix epoch.",
		ConstLabels: nil,
	}, func() float64 {
		p := processPool.Get().(*process.Process)
		defer processPool.Put(p)
		if p.Err != nil {
			return float64(time.Unix(0, 0).UnixMilli())
		}
		start, err := p.CreateTime()
		if err != nil {
			logger.Error().Err(err).Msg("")
			return float64(time.Unix(0, 0).UnixMilli())
		}
		return float64(start)
	}))
	prometheus.Register(prometheus.NewGaugeFunc(prometheus.GaugeOpts{
		Name:        "go_doudou_process_uptime_millis",
		Help:        "The uptime of the go-doudou process.",
		ConstLabels: nil,
	}, func() float64 {
		p := processPool.Get().(*process.Process)
		defer processPool.Put(p)
		if p.Err != nil {
			return float64(time.Since(time.Unix(0, 0).Local()).Milliseconds())
		}
		start, err := p.CreateTime()
		if err != nil {
			logger.Error().Err(err).Msg("")
			return float64(time.Since(time.Unix(0, 0).Local()).Milliseconds())
		}
		return float64(time.Since(time.UnixMilli(start).Local()).Milliseconds())
	}))

	prometheus.Register(prometheus.NewGaugeFunc(prometheus.GaugeOpts{
		Name:        "go_doudou_system_load_average_1m",
		Help:        "The sum of the number of runnable entities queued to available processors and the number of runnable entities running on the available processors averaged over a period of time",
		ConstLabels: nil,
	}, func() float64 {
		as, err := load.Avg1()
		if err != nil {
			logger.Error().Err(err).Msg("")
			return 0
		}
		return as.Load1
	}))
}
