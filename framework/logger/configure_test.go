package logger

import (
	"bytes"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"github.com/unionj-cloud/go-doudou/v2/framework/config"
)

func TestLoggerOptions(t *testing.T) {
	// 测试 WithWritter
	logger := logrus.New()
	buffer := &bytes.Buffer{}
	WithWritter(buffer)(logger)
	logger.Info("test")
	assert.Contains(t, buffer.String(), "test")

	// 测试 WithFormatter
	logger = logrus.New()
	formatter := &logrus.JSONFormatter{}
	WithFormatter(formatter)(logger)
	assert.Equal(t, formatter, logger.Formatter)

	// 测试 WithReportCaller
	logger = logrus.New()
	WithReportCaller(true)(logger)
	assert.True(t, logger.ReportCaller)
}

func TestDefaultFormatter(t *testing.T) {
	// 保存原始值
	originalFormat := config.GddLogFormat.Load()
	defer config.GddLogFormat.Write(originalFormat)

	// 测试 JSON 格式
	config.GddLogFormat.Write("json")
	formatter := defaultFormatter()
	assert.IsType(t, &logrus.JSONFormatter{}, formatter)

	// 测试 Text 格式
	config.GddLogFormat.Write("text")
	formatter = defaultFormatter()
	assert.IsType(t, &logrus.TextFormatter{}, formatter)

	// 测试默认格式
	config.GddLogFormat.Write("")
	formatter = defaultFormatter()
	assert.Nil(t, formatter)
}

func TestLogLevel_Decode(t *testing.T) {
	tests := []struct {
		name  string
		value string
		want  LogLevel
	}{
		{"panic", "panic", LogLevel(logrus.PanicLevel)},
		{"fatal", "fatal", LogLevel(logrus.FatalLevel)},
		{"error", "error", LogLevel(logrus.ErrorLevel)},
		{"warn", "warn", LogLevel(logrus.WarnLevel)},
		{"debug", "debug", LogLevel(logrus.DebugLevel)},
		{"trace", "trace", LogLevel(logrus.TraceLevel)},
		{"default", "unknown", LogLevel(logrus.InfoLevel)},
		{"empty", "", LogLevel(logrus.InfoLevel)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var ll LogLevel
			err := ll.Decode(tt.value)
			assert.NoError(t, err)
			assert.Equal(t, tt.want, ll)
		})
	}
}

func TestInit(t *testing.T) {
	// 保存原始设置
	originalLevel := config.GddLogLevel.Load()
	defer config.GddLogLevel.Write(originalLevel)

	// 测试初始化
	config.GddLogLevel.Write("debug")
	Init()
	assert.Equal(t, logrus.DebugLevel, logrus.StandardLogger().GetLevel())

	// 测试带选项的初始化
	buffer := &bytes.Buffer{}
	Init(WithWritter(buffer))
	logrus.Info("test init")
	assert.Contains(t, buffer.String(), "test init")
}
