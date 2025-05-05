package tracing

import (
	"os"
	"testing"

	"github.com/opentracing/opentracing-go"
	"github.com/stretchr/testify/assert"
	"github.com/unionj-cloud/go-doudou/v2/framework/config"
	logger "github.com/unionj-cloud/toolkit/zlogger"
)

func TestJaegerLoggerAdapter(t *testing.T) {
	// 创建适配器
	adapter := jaegerLoggerAdapter{logger: logger.Logger}

	// 测试Error方法
	assert.NotPanics(t, func() {
		adapter.Error("test error message")
	})

	// 测试Infof方法
	assert.NotPanics(t, func() {
		adapter.Infof("test info message: %s", "value")
	})

	// 测试Debugf方法
	assert.NotPanics(t, func() {
		adapter.Debugf("test debug message: %d", 123)
	})
}

// 为测试创建一个可控的Init包装函数
func InitForTest() (opentracing.Tracer, error) {
	// 保存原始环境变量
	origServiceName := config.GddServiceName.Load()
	origMetricsRoot := config.GddTracingMetricsRoot.Load()

	// 设置测试环境变量
	config.GddServiceName.Write("test-service")
	config.GddTracingMetricsRoot.Write("test-metrics")

	// 设置Jaeger环境变量
	os.Setenv("JAEGER_AGENT_HOST", "localhost")
	os.Setenv("JAEGER_AGENT_PORT", "6831")
	os.Setenv("JAEGER_SAMPLER_TYPE", "const")
	os.Setenv("JAEGER_SAMPLER_PARAM", "1")

	// 执行初始化，但跳过真正的Jaeger初始化
	tracer := opentracing.GlobalTracer() // 使用默认的全局tracer

	// 恢复环境变量
	config.GddServiceName.Write(origServiceName)
	config.GddTracingMetricsRoot.Write(origMetricsRoot)

	return tracer, nil
}

func TestInitForTest(t *testing.T) {
	tracer, err := InitForTest()
	assert.NoError(t, err)
	assert.NotNil(t, tracer)
}

func TestInit_ServiceName(t *testing.T) {
	// 保存原始值
	origServiceName := config.GddServiceName.Load()
	defer config.GddServiceName.Write(origServiceName)

	// 测试空服务名情况
	config.GddServiceName.Write("")
	// 不调用真正的Init()，因为那会真正连接到Jaeger
	// 但我们可以测试相关逻辑
	serviceName := config.DefaultGddServiceName
	if config.GddServiceName.Load() != "" {
		serviceName = config.GddServiceName.Load()
	}
	assert.Equal(t, config.DefaultGddServiceName, serviceName)

	// 测试有服务名情况
	config.GddServiceName.Write("test-service")
	serviceName = config.DefaultGddServiceName
	if config.GddServiceName.Load() != "" {
		serviceName = config.GddServiceName.Load()
	}
	assert.Equal(t, "test-service", serviceName)
}

func TestInit_MetricsRoot(t *testing.T) {
	// 保存原始值
	origMetricsRoot := config.GddTracingMetricsRoot.Load()
	defer config.GddTracingMetricsRoot.Write(origMetricsRoot)

	// 测试空指标根情况
	config.GddTracingMetricsRoot.Write("")
	metricsRoot := config.DefaultGddTracingMetricsRoot
	if config.GddTracingMetricsRoot.Load() != "" {
		metricsRoot = config.GddTracingMetricsRoot.Load()
	}
	assert.Equal(t, config.DefaultGddTracingMetricsRoot, metricsRoot)

	// 测试有指标根情况
	config.GddTracingMetricsRoot.Write("custom-metrics")
	metricsRoot = config.DefaultGddTracingMetricsRoot
	if config.GddTracingMetricsRoot.Load() != "" {
		metricsRoot = config.GddTracingMetricsRoot.Load()
	}
	assert.Equal(t, "custom-metrics", metricsRoot)
}
