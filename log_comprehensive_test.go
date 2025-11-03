package glog

import (
	"errors"
	"testing"

	"github.com/gw123/glog/common"
)

func TestDefaultLogger_NotNil(t *testing.T) {
	logger := DefaultLogger()
	if logger == nil {
		t.Error("DefaultLogger() should never return nil")
	}
}

func TestLog_NotNil(t *testing.T) {
	logger := Log()
	if logger == nil {
		t.Error("Log() should never return nil")
	}
}

func TestLog_SameAsDefaultLogger(t *testing.T) {
	logger1 := DefaultLogger()
	logger2 := Log()

	logger1.Info("test from DefaultLogger")
	logger2.Info("test from Log")
}

func TestSetDefaultLoggerConfig_Success(t *testing.T) {
	err := SetDefaultLoggerConfig(
		common.Options{},
		common.WithConsoleEncoding(),
		common.WithLevel(common.DebugLevel),
		common.WithStdoutOutputPath(),
	)
	if err != nil {
		t.Errorf("SetDefaultLoggerConfig() error = %v", err)
	}

	DefaultLogger().Debug("test after config change")
}

func TestSetDefaultLoggerConfig_WithFile(t *testing.T) {
	err := SetDefaultLoggerConfig(
		common.Options{},
		common.WithOutputPath("./test_log_api.log"),
		common.WithConsoleEncoding(),
		common.WithLevel(common.InfoLevel),
	)
	if err != nil {
		t.Errorf("SetDefaultLoggerConfig() with file error = %v", err)
	}

	Info("test log to file")
}

func TestError(t *testing.T) {
	Error("test error message")
}

func TestErrorf(t *testing.T) {
	Errorf("test error message with format: %s, %d", "param", 123)
}

func TestWarn(t *testing.T) {
	Warn("test warn message")
}

func TestWarnf(t *testing.T) {
	Warnf("test warn message with format: %s, %d", "param", 456)
}

func TestInfo(t *testing.T) {
	Info("test info message")
}

func TestInfof(t *testing.T) {
	Infof("test info message with format: %s, %d", "param", 789)
}

func TestDebug(t *testing.T) {
	SetDefaultLoggerConfig(
		common.Options{},
		common.WithLevel(common.DebugLevel),
	)
	Debug("test debug message")
}

func TestDebugf(t *testing.T) {
	SetDefaultLoggerConfig(
		common.Options{},
		common.WithLevel(common.DebugLevel),
	)
	Debugf("test debug message with format: %s, %d", "param", 111)
}

func TestWithField(t *testing.T) {
	logger := WithField("key", "value")
	if logger == nil {
		t.Error("WithField() returned nil")
	}
	logger.Info("test with field")
}

func TestWithField_Multiple(t *testing.T) {
	logger := WithField("key1", "value1").
		WithField("key2", 123).
		WithField("key3", true)

	if logger == nil {
		t.Error("Chained WithField() returned nil")
	}

	logger.Info("test with multiple fields")
}

func TestWithErr(t *testing.T) {
	testErr := errors.New("test error")
	logger := WithErr(testErr)

	if logger == nil {
		t.Error("WithErr() returned nil")
	}

	logger.Error("test with error")
}

func TestWithErr_Nil(t *testing.T) {
	logger := WithErr(nil)
	if logger == nil {
		t.Error("WithErr(nil) returned nil")
	}
	logger.Info("test with nil error")
}

func TestWithError(t *testing.T) {
	testErr := errors.New("test error via WithError")
	logger := WithError(testErr)

	if logger == nil {
		t.Error("WithError() returned nil")
	}

	logger.Error("test with error")
}

func TestWithError_Nil(t *testing.T) {
	logger := WithError(nil)
	if logger == nil {
		t.Error("WithError(nil) returned nil")
	}
	logger.Info("test with nil error")
}

func TestAllLogLevels(t *testing.T) {
	SetDefaultLoggerConfig(
		common.Options{},
		common.WithLevel(common.DebugLevel),
		common.WithConsoleEncoding(),
	)

	Debug("debug level")
	Debugf("debug level %d", 1)
	Info("info level")
	Infof("info level %d", 2)
	Warn("warn level")
	Warnf("warn level %d", 3)
	Error("error level")
	Errorf("error level %d", 4)
}

func TestCombinedFieldsAndError(t *testing.T) {
	testErr := errors.New("combined test error")
	logger := WithField("request_id", "req-123").
		WithField("user_id", 456).
		WithError(testErr)

	logger.Error("test combined fields and error")
}

func TestLoggerChaining(t *testing.T) {
	logger := DefaultLogger().
		WithField("service", "test-service").
		WithField("version", "1.0.0")

	logger.Info("test logger chaining")
	logger.Named("component").Info("test named after chaining")
}

func TestConcurrentLogging(t *testing.T) {
	done := make(chan bool)

	for i := 0; i < 100; i++ {
		go func(id int) {
			Infof("concurrent message %d", id)
			WithField("id", id).Info("concurrent with field")
			done <- true
		}(i)
	}

	for i := 0; i < 100; i++ {
		<-done
	}
}

func TestConcurrentConfigChange(t *testing.T) {
	done := make(chan bool)

	for i := 0; i < 10; i++ {
		go func(id int) {
			SetDefaultLoggerConfig(
				common.Options{},
				common.WithLevel(common.InfoLevel),
			)
			Infof("message during config change %d", id)
			done <- true
		}(i)
	}

	for i := 0; i < 10; i++ {
		<-done
	}
}

func BenchmarkInfo(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Info("benchmark message")
	}
}

func BenchmarkInfof(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Infof("benchmark message %d", i)
	}
}

func BenchmarkWithField(b *testing.B) {
	for i := 0; i < b.N; i++ {
		WithField("key", i).Info("benchmark with field")
	}
}

func BenchmarkWithError(b *testing.B) {
	testErr := errors.New("benchmark error")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		WithError(testErr).Error("benchmark with error")
	}
}

func BenchmarkDefaultLogger(b *testing.B) {
	for i := 0; i < b.N; i++ {
		DefaultLogger()
	}
}

func BenchmarkConcurrentLogging(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			Infof("concurrent benchmark %d", i)
			i++
		}
	})
}
