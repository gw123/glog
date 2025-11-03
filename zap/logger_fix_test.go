package zap

import (
	"testing"

	"github.com/gw123/glog/common"
)

func TestLogger_WithFields(t *testing.T) {
	logger, err := NewLogger(common.Options{}, common.WithConsoleEncoding(), common.WithLevel(common.DebugLevel))
	if err != nil {
		t.Fatalf("NewLogger() error = %v", err)
	}

	fields := map[string]interface{}{
		"key1": "value1",
		"key2": 123,
		"key3": true,
	}

	newLogger := logger.WithFields(fields)
	if newLogger == nil {
		t.Error("WithFields() returned nil")
	}

	newLogger.Info("test message with multiple fields")
}

func TestSetDefaultLoggerConfig_ThreadSafety(t *testing.T) {
	done := make(chan bool)

	for i := 0; i < 10; i++ {
		go func() {
			err := SetDefaultLoggerConfig(common.Options{}, common.WithConsoleEncoding(), common.WithLevel(common.InfoLevel))
			if err != nil {
				t.Errorf("SetDefaultLoggerConfig() error = %v", err)
			}
			done <- true
		}()
	}

	for i := 0; i < 10; i++ {
		<-done
	}

	for i := 0; i < 100; i++ {
		go func() {
			logger := DefaultLogger()
			if logger == nil {
				t.Error("DefaultLogger() returned nil")
			}
			done <- true
		}()
	}

	for i := 0; i < 100; i++ {
		<-done
	}
}

func TestNewLogger_DirectoryCreation(t *testing.T) {
	testPath := "./test_logs/app.log"
	logger, err := NewLogger(
		common.Options{},
		common.WithOutputPath(testPath),
		common.WithConsoleEncoding(),
		common.WithLevel(common.DebugLevel),
	)

	if err != nil {
		t.Fatalf("NewLogger() error = %v, should create directory automatically", err)
	}

	if logger == nil {
		t.Error("NewLogger() returned nil")
	}

	logger.Info("test log to file")
}
