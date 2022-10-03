package tracing

import (
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/uber/jaeger-client-go/config"
	"github.com/uber/jaeger-client-go/rpcmetrics"
	"github.com/uber/jaeger-lib/metrics"
	jprom "github.com/uber/jaeger-lib/metrics/prometheus"
	ddconfig "github.com/unionj-cloud/go-doudou/framework/internal/config"
	logger "github.com/unionj-cloud/go-doudou/toolkit/zlogger"
	"github.com/unionj-cloud/go-doudou/toolkit/stringutils"
	"io"
)

// Init returns an instance of Jaeger Tracer.
func Init() (opentracing.Tracer, io.Closer) {
	cfg := &config.Configuration{
		Sampler:  &config.SamplerConfig{},
		Reporter: &config.ReporterConfig{},
	}
	service := ddconfig.DefaultGddServiceName
	if stringutils.IsNotEmpty(ddconfig.GddServiceName.Load()) {
		service = ddconfig.GddServiceName.Load()
	}
	cfg.ServiceName = service
	cfg.Sampler.Type = "const"
	cfg.Sampler.Param = 1
	cfg.Reporter.LogSpans = true
	_, err := cfg.FromEnv()
	if err != nil {
		logger.Panic().Err(errors.Wrap(err, "[go-doudou] cannot parse Jaeger env vars")).Msg("")
	}
	jaegerLogger := jaegerLoggerAdapter{logger: logger.Logger}
	metricsRoot := ddconfig.DefaultGddTracingMetricsRoot
	if stringutils.IsNotEmpty(ddconfig.GddTracingMetricsRoot.Load()) {
		metricsRoot = ddconfig.GddTracingMetricsRoot.Load()
	}
	metricsFactory := jprom.New().Namespace(metrics.NSOptions{Name: metricsRoot, Tags: nil})
	tracer, closer, err := cfg.NewTracer(
		config.Logger(jaegerLogger),
		config.Metrics(metricsFactory),
		config.Observer(rpcmetrics.NewObserver(metricsFactory, rpcmetrics.DefaultNameNormalizer)),
	)
	if err != nil {
		logger.Panic().Err(errors.Wrap(err, "[go-doudou] cannot initialize Jaeger Tracer")).Msg("")
	}
	return tracer, closer
}

type jaegerLoggerAdapter struct {
	logger zerolog.Logger
}

func (l jaegerLoggerAdapter) Error(msg string) {
	l.logger.Error().Msg(msg)
}

func (l jaegerLoggerAdapter) Infof(msg string, args ...interface{}) {
	l.logger.Info().Msgf(msg, args...)
}

func (l jaegerLoggerAdapter) Debugf(msg string, args ...interface{}) {
	l.logger.Debug().Msgf(msg, args...)
}
