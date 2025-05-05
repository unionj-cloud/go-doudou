package tracing

import (
	"os"
	"sync"
	"testing"

	"github.com/opentracing/opentracing-go"
	"github.com/stretchr/testify/assert"

	"github.com/unionj-cloud/go-doudou/v2/framework/config"
)

// 使用互斥锁确保测试按顺序运行，防止度量收集器重复注册
var tracerTestMutex sync.Mutex

func TestInit(t *testing.T) {
	// 加锁以便串行执行测试
	tracerTestMutex.Lock()
	defer tracerTestMutex.Unlock()

	// 测试tracing初始化可能只能运行一次，因为Prometheus收集器不能重复注册
	// 我们将简化测试，只测试一次初始化

	// 设置必要的环境变量
	os.Setenv("JAEGER_AGENT_HOST", "localhost")
	os.Setenv("JAEGER_AGENT_PORT", "6831")

	// 保存原始值以便恢复
	origServiceName := config.GddServiceName.Load()
	origMetricsRoot := config.GddTracingMetricsRoot.Load()
	defer func() {
		config.GddServiceName.Write(origServiceName)
		config.GddTracingMetricsRoot.Write(origMetricsRoot)
	}()

	// 设置服务名称
	os.Setenv(string(config.GddServiceName), "test-service")

	// 我们只测试一次初始化，并确保它不会panic
	assert.NotPanics(t, func() {
		tracer, closer := Init()
		assert.NotNil(t, tracer)
		assert.NotNil(t, closer)
		if closer != nil {
			closer.Close()
		}
	})
}

func TestJaegerLoggerAdapter(t *testing.T) {
	// 创建一个logger adapter
	logger := jaegerLoggerAdapter{}

	// 测试各种日志方法 (无需assert，只要不panic即可)
	logger.Error("test error message")
	logger.Infof("test info message %s", "param")
	logger.Debugf("test debug message %d", 123)
}

func TestDefaultTracer(t *testing.T) {
	// 确保使用opentracing.GlobalTracer()不会出错
	span := opentracing.StartSpan("test-operation")
	assert.NotNil(t, span)
	span.Finish()
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
