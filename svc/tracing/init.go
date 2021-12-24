package tracing

import (
	"fmt"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/uber/jaeger-client-go/config"
	"github.com/uber/jaeger-client-go/rpcmetrics"
	"github.com/uber/jaeger-lib/metrics"
	"io"
)

// Init returns an instance of Jaeger Tracer.
func Init(service string, metricsFactory metrics.Factory) (opentracing.Tracer, io.Closer) {
	cfg := &config.Configuration{
		Sampler:  &config.SamplerConfig{},
		Reporter: &config.ReporterConfig{},
	}
	cfg.ServiceName = service
	cfg.Sampler.Type = "const"
	cfg.Sampler.Param = 1
	cfg.Reporter.LogSpans = true

	_, err := cfg.FromEnv()
	if err != nil {
		logrus.Panic(errors.Wrap(err, "cannot parse Jaeger env vars"))
	}
	jaegerLogger := jaegerLoggerAdapter{logger: logrus.StandardLogger()}
	metricsFactory = metricsFactory.Namespace(metrics.NSOptions{Name: service, Tags: nil})
	tracer, closer, err := cfg.NewTracer(
		config.Logger(jaegerLogger),
		config.Metrics(metricsFactory),
		config.Observer(rpcmetrics.NewObserver(metricsFactory, rpcmetrics.DefaultNameNormalizer)),
	)
	if err != nil {
		logrus.Panic(errors.Wrap(err, "cannot initialize Jaeger Tracer"))
	}
	return tracer, closer
}

type jaegerLoggerAdapter struct {
	logger *logrus.Logger
}

func (l jaegerLoggerAdapter) Error(msg string) {
	l.logger.Error(msg)
}

func (l jaegerLoggerAdapter) Infof(msg string, args ...interface{}) {
	l.logger.Info(fmt.Sprintf(msg, args...))
}
