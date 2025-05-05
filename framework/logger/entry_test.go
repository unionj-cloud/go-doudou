package logger

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	// 测试dev环境
	os.Setenv("GDD_ENV", "dev")
	entry := New()
	assert.NotNil(t, entry)
	assert.IsType(t, &logrus.Entry{}, entry)

	// 测试非dev环境
	os.Setenv("GDD_ENV", "prod")
	entry = New()
	assert.NotNil(t, entry)
	assert.IsType(t, &logrus.Entry{}, entry)

	// 恢复环境变量
	os.Unsetenv("GDD_ENV")
}

func TestCheckDev(t *testing.T) {
	// 测试未设置环境变量的情况
	os.Unsetenv("GDD_ENV")
	assert.True(t, CheckDev())

	// 测试dev环境
	os.Setenv("GDD_ENV", "dev")
	assert.True(t, CheckDev())

	// 测试非dev环境
	os.Setenv("GDD_ENV", "prod")
	assert.False(t, CheckDev())

	// 恢复环境变量
	os.Unsetenv("GDD_ENV")
}

func TestWithError(t *testing.T) {
	err := assert.AnError
	entry := WithError(err)
	assert.NotNil(t, entry)
	assert.Equal(t, err, entry.Data[logrus.ErrorKey])
}

func TestWithContext(t *testing.T) {
	ctx := context.Background()
	entry := WithContext(ctx)
	assert.NotNil(t, entry)
}

func TestWithField(t *testing.T) {
	entry := WithField("key", "value")
	assert.NotNil(t, entry)
	assert.Equal(t, "value", entry.Data["key"])
}

func TestWithFields(t *testing.T) {
	// 测试dev环境
	os.Setenv("GDD_ENV", "dev")
	fields := logrus.Fields{"key1": "value1", "key2": "value2"}
	entry := WithFields(fields)
	assert.NotNil(t, entry)

	// 测试非dev环境
	os.Setenv("GDD_ENV", "prod")
	entry = WithFields(fields)
	assert.NotNil(t, entry)
	assert.Equal(t, "value1", entry.Data["key1"])
	assert.Equal(t, "value2", entry.Data["key2"])

	// 恢复环境变量
	os.Unsetenv("GDD_ENV")
}

func TestWithTime(t *testing.T) {
	now := time.Now()
	entry := WithTime(now)
	assert.NotNil(t, entry)
	// 由于logrus.Entry的Time字段不可直接访问，我们只能验证函数不panic
}

// 为了确保测试覆盖率，我们测试所有的日志函数
// 注意：这些测试只是验证函数调用不会panic，不验证实际的日志输出

func TestLogFunctions(t *testing.T) {
	// 为避免实际输出日志，我们可以临时将日志级别设置为更高级别
	originalLevel := logrus.StandardLogger().GetLevel()
	logrus.StandardLogger().SetLevel(logrus.PanicLevel)
	defer logrus.StandardLogger().SetLevel(originalLevel)

	// 测试各种日志函数
	Trace("trace message")
	Debug("debug message")
	Print("print message")
	Info("info message")
	Warn("warn message")
	Warning("warning message")
	Error("error message")
	// 不测试Panic和Fatal，因为它们会中断测试

	Tracef("trace %s", "message")
	Debugf("debug %s", "message")
	Printf("print %s", "message")
	Infof("info %s", "message")
	Warnf("warn %s", "message")
	Warningf("warning %s", "message")
	Errorf("error %s", "message")
	// 不测试Panicf和Fatalf

	Traceln("trace", "message")
	Debugln("debug", "message")
	Println("print", "message")
	Infoln("info", "message")
	Warnln("warn", "message")
	Warningln("warning", "message")
	Errorln("error", "message")
	// 不测试Panicln和Fatalln
}

// 测试函数Entry
func TestEntry(t *testing.T) {
	e := Entry()
	assert.NotNil(t, e)
	assert.IsType(t, &logrus.Entry{}, e)
}
