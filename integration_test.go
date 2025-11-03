package glog

import (
	"context"
	"os"
	"testing"

	"github.com/gw123/glog/common"
)

func TestIntegration_FullWorkflow(t *testing.T) {
	err := SetDefaultLoggerConfig(
		common.Options{},
		common.WithConsoleEncoding(),
		common.WithLevel(common.DebugLevel),
		common.WithOutputPath("./test_integration.log"),
	)
	if err != nil {
		t.Fatalf("SetDefaultLoggerConfig() error = %v", err)
	}
	//defer os.Remove("./test_integration.log")

	ctx := context.Background()
	ctx = ToContext(ctx, DefaultLogger())

	AddTraceID(ctx, "integration-trace-001")
	AddUserID(ctx, 999)
	AddPathname(ctx, "/api/integration/test")
	AddField(ctx, "request_id", "req-integration-001")
	AddFields(ctx, map[string]interface{}{
		"method":      "POST",
		"status_code": 200,
		"duration_ms": 150,
	})

	logger := ExtractEntry(ctx)
	logger.Info("Integration test: full workflow")
	logger.WithField("extra", "data").Debug("Debug message with extra field")

	traceID := ExtractTraceID(ctx)
	if traceID != "integration-trace-001" {
		t.Errorf("TraceID mismatch: got %v, want integration-trace-001", traceID)
	}
}

func TestIntegration_MultipleContexts(t *testing.T) {
	ctx1 := ToContext(context.Background(), DefaultLogger().Named("service1"))
	AddTraceID(ctx1, "trace-ctx1")
	AddField(ctx1, "service", "service1")

	ctx2 := ToContext(context.Background(), DefaultLogger().Named("service2"))
	AddTraceID(ctx2, "trace-ctx2")
	AddField(ctx2, "service", "service2")

	logger1 := ExtractEntry(ctx1)
	logger1.Info("Message from context 1")

	logger2 := ExtractEntry(ctx2)
	logger2.Info("Message from context 2")

	if ExtractTraceID(ctx1) == ExtractTraceID(ctx2) {
		t.Error("Different contexts should have different trace IDs")
	}
}

func TestIntegration_ErrorHandlingWithContext(t *testing.T) {
	ctx := ToContext(context.Background(), DefaultLogger())
	AddTraceID(ctx, "error-test-trace")

	simulateError := func(ctx context.Context) error {
		logger := ExtractEntry(ctx)
		err := os.ErrNotExist
		logger.WithError(err).Error("Simulated error occurred")
		return err
	}

	err := simulateError(ctx)
	if err == nil {
		t.Error("Expected error but got nil")
	}
}

func TestIntegration_NestedFunctionCalls(t *testing.T) {
	ctx := ToContext(context.Background(), DefaultLogger())
	AddTraceID(ctx, "nested-trace")

	level2 := func(ctx context.Context) {
		AddField(ctx, "level", 2)
		ExtractEntry(ctx).Info("Level 2 function")

		level3 := func(ctx context.Context) {
			AddField(ctx, "level", 3)
			ExtractEntry(ctx).Info("Level 3 function")
		}
		level3(ctx)
	}

	level1 := func(ctx context.Context) {
		AddField(ctx, "level", 1)
		ExtractEntry(ctx).Info("Level 1 function")
		level2(ctx)
	}

	level1(ctx)
}

func TestIntegration_ConcurrentRequests(t *testing.T) {
	done := make(chan bool)

	for i := 0; i < 50; i++ {
		go func(requestID int) {
			ctx := ToContext(context.Background(), DefaultLogger())
			AddTraceID(ctx, "concurrent-trace-"+string(rune(requestID)))
			AddUserID(ctx, int64(requestID))
			AddField(ctx, "request_num", requestID)

			logger := ExtractEntry(ctx)
			logger.Infof("Processing request %d", requestID)

			AddField(ctx, "status", "completed")
			logger = ExtractEntry(ctx)
			logger.Infof("Completed request %d", requestID)

			done <- true
		}(i)
	}

	for i := 0; i < 50; i++ {
		<-done
	}
}

func TestIntegration_DynamicLogLevelChange(t *testing.T) {
	SetDefaultLoggerConfig(
		common.Options{},
		common.WithLevel(common.InfoLevel),
	)

	Debug("This should not appear")
	Info("This should appear")

	SetDefaultLoggerConfig(
		common.Options{},
		common.WithLevel(common.DebugLevel),
	)

	Debug("This should now appear")
	Info("This should still appear")
}

func TestEdgeCase_EmptyStrings(t *testing.T) {
	ctx := ToContext(context.Background(), DefaultLogger())
	AddTraceID(ctx, "")
	AddPathname(ctx, "")
	AddField(ctx, "", "")
	AddField(ctx, "key", "")

	logger := ExtractEntry(ctx)
	logger.Info("Test with empty strings")
}

func TestEdgeCase_NilValues(t *testing.T) {
	ctx := ToContext(context.Background(), DefaultLogger())
	AddField(ctx, "nil_value", nil)
	AddFields(ctx, map[string]interface{}{
		"key1": nil,
		"key2": "value",
	})

	logger := ExtractEntry(ctx)
	logger.Info("Test with nil values")
}

func TestEdgeCase_LargeFieldValues(t *testing.T) {
	ctx := ToContext(context.Background(), DefaultLogger())

	largeString := string(make([]byte, 10000))
	AddField(ctx, "large_field", largeString)

	largeMap := make(map[string]interface{})
	for i := 0; i < 100; i++ {
		largeMap["key"+string(rune(i))] = i
	}
	AddFields(ctx, largeMap)

	logger := ExtractEntry(ctx)
	logger.Info("Test with large field values")
}

func TestEdgeCase_SpecialCharacters(t *testing.T) {
	ctx := ToContext(context.Background(), DefaultLogger())
	AddField(ctx, "special", "!@#$%^&*(){}[]|\\:;\"'<>,.?/~`")
	AddTraceID(ctx, "trace-with-特殊字符-🎉")

	logger := ExtractEntry(ctx)
	logger.Info("Test with special characters: 测试中文 🚀")
}

func TestEdgeCase_MaxUserID(t *testing.T) {
	ctx := ToContext(context.Background(), DefaultLogger())
	AddUserID(ctx, 9223372036854775807)

	extractedID := ExtractUserID(ctx)
	if extractedID != 9223372036854775807 {
		t.Errorf("Max int64 UserID mismatch: got %v", extractedID)
	}
}

func TestEdgeCase_NegativeUserID(t *testing.T) {
	ctx := ToContext(context.Background(), DefaultLogger())
	AddUserID(ctx, -1)

	extractedID := ExtractUserID(ctx)
	if extractedID != -1 {
		t.Errorf("Negative UserID mismatch: got %v", extractedID)
	}
}

func TestEdgeCase_VeryLongTraceID(t *testing.T) {
	ctx := ToContext(context.Background(), DefaultLogger())
	longTraceID := string(make([]byte, 1000))
	for i := range longTraceID {
		longTraceID = string(append([]byte(longTraceID[:i]), 'a'))
	}

	AddTraceID(ctx, longTraceID)
	extractedID := ExtractTraceID(ctx)

	if len(extractedID) != len(longTraceID) {
		t.Errorf("Long TraceID length mismatch: got %d, want %d", len(extractedID), len(longTraceID))
	}
}

func TestEdgeCase_RapidFieldUpdates(t *testing.T) {
	ctx := ToContext(context.Background(), DefaultLogger())

	for i := 0; i < 1000; i++ {
		AddField(ctx, "counter", i)
	}

	logger := ExtractEntry(ctx)
	logger.Info("Test with rapid field updates")
}

func TestEdgeCase_ContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	ctx = ToContext(ctx, DefaultLogger())
	AddField(ctx, "test", "value")

	cancel()

	logger := ExtractEntry(ctx)
	logger.Info("Test with cancelled context")
}

func BenchmarkIntegration_FullWorkflow(b *testing.B) {
	ctx := ToContext(context.Background(), DefaultLogger())

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		AddTraceID(ctx, "benchmark-trace")
		AddUserID(ctx, int64(i))
		AddField(ctx, "iteration", i)
		logger := ExtractEntry(ctx)
		logger.Info("Benchmark iteration")
	}
}

func BenchmarkIntegration_ConcurrentWorkflow(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			ctx := ToContext(context.Background(), DefaultLogger())
			AddTraceID(ctx, "benchmark-concurrent")
			AddField(ctx, "id", i)
			ExtractEntry(ctx).Info("Concurrent benchmark")
			i++
		}
	})
}
