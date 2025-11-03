package zap

import (
	"os"
	"strings"
	"testing"

	"github.com/gw123/glog/common"
)

func TestLogger_WithField(t *testing.T) {
	logger, err := NewLogger(common.Options{}, common.WithConsoleEncoding(), common.WithLevel(common.DebugLevel))
	if err != nil {
		t.Fatalf("NewLogger() error = %v", err)
	}

	newLogger := logger.WithField("key", "value")
	if newLogger == nil {
		t.Error("WithField() returned nil")
	}

	newLogger.Info("test with field")
}

func TestLogger_WithError_Nil(t *testing.T) {
	logger, err := NewLogger(common.Options{}, common.WithConsoleEncoding())
	if err != nil {
		t.Fatalf("NewLogger() error = %v", err)
	}

	newLogger := logger.WithError(nil)
	if newLogger == nil {
		t.Error("WithError(nil) should not return nil")
	}

	newLogger.Info("test with nil error")
}

func TestLogger_Warningf(t *testing.T) {
	logger, err := NewLogger(common.Options{}, common.WithConsoleEncoding(), common.WithLevel(common.WarnLevel))
	if err != nil {
		t.Fatalf("NewLogger() error = %v", err)
	}

	logger.Warningf("test warning %s", "message")
	logger.Warning("test warning")
}

func TestLogger_Named(t *testing.T) {
	logger, err := NewLogger(common.Options{}, common.WithConsoleEncoding())
	if err != nil {
		t.Fatalf("NewLogger() error = %v", err)
	}

	named := logger.Named("component")
	if named == nil {
		t.Error("Named() returned nil")
	}

	named.Info("test named logger")
}

func TestNewLogger_DefaultOptions(t *testing.T) {
	logger, err := NewLogger(common.Options{})
	if err != nil {
		t.Fatalf("NewLogger() with default options error = %v", err)
	}

	if logger == nil {
		t.Error("NewLogger() returned nil")
	}

	logger.Info("test with default options")
}

func TestNewLogger_WithAllOptions(t *testing.T) {
	logger, err := NewLogger(
		common.Options{},
		common.WithConsoleEncoding(),
		common.WithLevel(common.DebugLevel),
		common.WithStdoutOutputPath(),
		common.WithStderrErrorOutputPath(),
		common.WithCallerSkip(0),
	)
	if err != nil {
		t.Fatalf("NewLogger() error = %v", err)
	}

	if logger == nil {
		t.Error("NewLogger() returned nil")
	}

	logger.Debug("test debug")
	logger.Info("test info")
	logger.Warn("test warn")
	logger.Error("test error")
}

func TestNewLogger_JsonEncoding(t *testing.T) {
	logger, err := NewLogger(
		common.Options{},
		common.WithJsonEncoding(),
		common.WithLevel(common.DebugLevel),
	)
	if err != nil {
		t.Fatalf("NewLogger() error = %v", err)
	}

	logger.WithField("key1", "value1").WithField("key2", 123).Info("test json encoding")
}

func TestNewLogger_MultipleOutputPaths(t *testing.T) {
	tmpFile := "./test_output_multi.log"
	defer os.Remove(tmpFile)

	logger, err := NewLogger(
		common.Options{},
		common.WithOutputPath(tmpFile),
		common.WithStdoutOutputPath(),
		common.WithConsoleEncoding(),
	)
	if err != nil {
		t.Fatalf("NewLogger() error = %v", err)
	}

	testMessage := "test multiple outputs"
	logger.Info(testMessage)

	content, err := os.ReadFile(tmpFile)
	if err != nil {
		t.Fatalf("Failed to read log file: %v", err)
	}

	if !strings.Contains(string(content), testMessage) {
		t.Errorf("Log file should contain %q, got %q", testMessage, string(content))
	}
}

func TestNewLogger_InvalidDirectory(t *testing.T) {
	invalidPath := "/invalid/path/that/cannot/be/created/\x00/test.log"
	_, err := NewLogger(
		common.Options{},
		common.WithOutputPath(invalidPath),
	)
	if err == nil {
		t.Error("NewLogger() with invalid path should return error")
	}
}

func TestSetDefaultLoggerConfig_Error(t *testing.T) {
	invalidPath := "/invalid/path/\x00/test.log"
	err := SetDefaultLoggerConfig(
		common.Options{},
		common.WithOutputPath(invalidPath),
	)
	if err == nil {
		t.Error("SetDefaultLoggerConfig() with invalid path should return error")
	}
}

func TestDefaultLogger_NotNil(t *testing.T) {
	logger := DefaultLogger()
	if logger == nil {
		t.Error("DefaultLogger() should never return nil")
	}
}

func TestGetInnerLogger_NotNil(t *testing.T) {
	logger := GetInnerLogger()
	if logger == nil {
		t.Error("GetInnerLogger() should never return nil")
	}
}

func TestLogger_AllLevels(t *testing.T) {
	logger, err := NewLogger(
		common.Options{},
		common.WithConsoleEncoding(),
		common.WithLevel(common.DebugLevel),
	)
	if err != nil {
		t.Fatalf("NewLogger() error = %v", err)
	}

	logger.Debug("debug message")
	logger.Debugf("debug message %d", 1)
	logger.Info("info message")
	logger.Infof("info message %d", 2)
	logger.Warn("warn message")
	logger.Warnf("warn message %d", 3)
	logger.Warning("warning message")
	logger.Warningf("warning message %d", 4)
	logger.Error("error message")
	logger.Errorf("error message %d", 5)
}

func TestLogger_WithFieldsChaining(t *testing.T) {
	logger, err := NewLogger(common.Options{}, common.WithConsoleEncoding())
	if err != nil {
		t.Fatalf("NewLogger() error = %v", err)
	}

	fields1 := map[string]interface{}{"key1": "value1"}
	fields2 := map[string]interface{}{"key2": "value2"}

	chainedLogger := logger.WithFields(fields1).WithFields(fields2)
	if chainedLogger == nil {
		t.Error("Chained WithFields() returned nil")
	}

	chainedLogger.Info("test chained fields")
}

func TestLogger_WithFieldAndError(t *testing.T) {
	logger, err := NewLogger(common.Options{}, common.WithConsoleEncoding())
	if err != nil {
		t.Fatalf("NewLogger() error = %v", err)
	}

	testErr := os.ErrNotExist
	combinedLogger := logger.WithField("request_id", "123").WithError(testErr)
	if combinedLogger == nil {
		t.Error("Combined WithField and WithError returned nil")
	}

	combinedLogger.Error("test combined field and error")
}

func BenchmarkLogger_Info(b *testing.B) {
	logger := DefaultLogger()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info("benchmark message")
	}
}

func BenchmarkLogger_Infof(b *testing.B) {
	logger := DefaultLogger()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Infof("benchmark message %d", i)
	}
}

func BenchmarkLogger_WithField(b *testing.B) {
	logger := DefaultLogger()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.WithField("key", i).Info("benchmark with field")
	}
}

func BenchmarkLogger_WithFields(b *testing.B) {
	logger := DefaultLogger()
	fields := map[string]interface{}{
		"key1": "value1",
		"key2": 123,
		"key3": true,
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.WithFields(fields).Info("benchmark with fields")
	}
}

func BenchmarkSetDefaultLoggerConfig(b *testing.B) {
	for i := 0; i < b.N; i++ {
		SetDefaultLoggerConfig(
			common.Options{},
			common.WithConsoleEncoding(),
			common.WithLevel(common.InfoLevel),
		)
	}
}
