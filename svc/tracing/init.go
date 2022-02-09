package tracing

import (
	"fmt"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/uber/jaeger-client-go/config"
	"github.com/uber/jaeger-client-go/rpcmetrics"
	"github.com/uber/jaeger-lib/metrics"
	jprom "github.com/uber/jaeger-lib/metrics/prometheus"
	"github.com/unionj-cloud/go-doudou/stringutils"
	ddconfig "github.com/unionj-cloud/go-doudou/svc/config"
	"github.com/unionj-cloud/go-doudou/svc/logger"
	"github.com/unionj-cloud/go-doudou/svc/registry"
	"io"
)

// Init returns an instance of Jaeger Tracer.
func Init() (opentracing.Tracer, io.Closer) {
	service := ddconfig.GddServiceName.Load()
	cfg := &config.Configuration{
		Sampler:  &config.SamplerConfig{},
		Reporter: &config.ReporterConfig{},
	}
	if registry.LocalNode() != nil {
		cfg.ServiceName = fmt.Sprintf("%s:%s", service, registry.LocalNode().Name)
	} else {
		cfg.ServiceName = service
	}
	cfg.Sampler.Type = "const"
	cfg.Sampler.Param = 1
	cfg.Reporter.LogSpans = true
	_, err := cfg.FromEnv()
	if err != nil {
		logger.Panic(errors.Wrap(err, "cannot parse Jaeger env vars"))
	}
	jaegerLogger := jaegerLoggerAdapter{logger: logger.Entry()}
	metricsRoot := ddconfig.DefaultGddTracingMetricsRoot
	if stringutils.IsNotEmpty(ddconfig.GddTracingMetricsRoot.Load()) {
		metricsRoot = ddconfig.GddTracingMetricsRoot.Load()
	}
	metricsFactory := jprom.New().Namespace(metrics.NSOptions{Name: metricsRoot, Tags: nil}).
		Namespace(metrics.NSOptions{Name: service, Tags: nil})
	if registry.LocalNode() != nil {
		metricsFactory = metricsFactory.Namespace(metrics.NSOptions{Name: registry.LocalNode().Name, Tags: nil})
	}
	tracer, closer, err := cfg.NewTracer(
		config.Logger(jaegerLogger),
		config.Metrics(metricsFactory),
		config.Observer(rpcmetrics.NewObserver(metricsFactory, rpcmetrics.DefaultNameNormalizer)),
	)
	if err != nil {
		logger.Panic(errors.Wrap(err, "cannot initialize Jaeger Tracer"))
	}
	return tracer, closer
}

type jaegerLoggerAdapter struct {
	logger *logrus.Entry
}

func (l jaegerLoggerAdapter) Error(msg string) {
	l.logger.Error(msg)
}

func (l jaegerLoggerAdapter) Infof(msg string, args ...interface{}) {
	l.logger.Info(fmt.Sprintf(msg, args...))
}

func (l jaegerLoggerAdapter) Debugf(msg string, args ...interface{}) {
	l.logger.Debug(fmt.Sprintf(msg, args...))
}
