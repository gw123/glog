package glog

import (
	"context"
	"testing"
)

func TestToContext(t *testing.T) {
	ctx := context.Background()
	logger := DefaultLogger()

	newCtx := ToContext(ctx, logger)
	if newCtx == nil {
		t.Error("ToContext() returned nil context")
	}

	if newCtx == ctx {
		t.Error("ToContext() should return a new context")
	}
}

func TestAddField_NilContext(t *testing.T) {
	ctx := context.Background()
	AddField(ctx, "key", "value")
}

func TestAddFields_NilContext(t *testing.T) {
	ctx := context.Background()
	fields := map[string]interface{}{"key": "value"}
	AddFields(ctx, fields)
}

func TestAddTopField(t *testing.T) {
	ctx := ToContext(context.Background(), DefaultLogger())
	AddTopField(ctx, "custom_key", "custom_value")

	logger := ExtractEntry(ctx)
	logger.Info("test with top field")
}

func TestAddTopField_NilContext(t *testing.T) {
	ctx := context.Background()
	AddTopField(ctx, "key", "value")
}

func TestAddTraceID(t *testing.T) {
	ctx := ToContext(context.Background(), DefaultLogger())
	testTraceID := "trace-123-456"

	AddTraceID(ctx, testTraceID)

	extractedID := ExtractTraceID(ctx)
	if extractedID != testTraceID {
		t.Errorf("ExtractTraceID() = %v, want %v", extractedID, testTraceID)
	}
}

func TestExtractTraceID_NoContext(t *testing.T) {
	ctx := context.Background()
	traceID := ExtractTraceID(ctx)
	if traceID != "" {
		t.Errorf("ExtractTraceID() from empty context should return empty string, got %v", traceID)
	}
}

func TestExtractTraceID_WithDebugPanic(t *testing.T) {
	oldIsDebug := IsDebug
	IsDebug = true
	defer func() {
		IsDebug = oldIsDebug
	}()

	ctx := context.Background()

	defer func() {
		if r := recover(); r == nil {
			t.Error("ExtractTraceID() should panic with IsDebug=true and no context logger")
		}
	}()

	ExtractTraceID(ctx)
}

func TestAddUserID(t *testing.T) {
	ctx := ToContext(context.Background(), DefaultLogger())
	testUserID := int64(12345)

	AddUserID(ctx, testUserID)

	extractedID := ExtractUserID(ctx)
	if extractedID != testUserID {
		t.Errorf("ExtractUserID() = %v, want %v", extractedID, testUserID)
	}
}

func TestExtractUserID_NoContext(t *testing.T) {
	ctx := context.Background()
	userID := ExtractUserID(ctx)
	if userID != 0 {
		t.Errorf("ExtractUserID() from empty context should return 0, got %v", userID)
	}
}

func TestExtractUserID_WithDebugPanic(t *testing.T) {
	oldIsDebug := IsDebug
	IsDebug = true
	defer func() {
		IsDebug = oldIsDebug
	}()

	ctx := context.Background()

	defer func() {
		if r := recover(); r == nil {
			t.Error("ExtractUserID() should panic with IsDebug=true and no context logger")
		}
	}()

	ExtractUserID(ctx)
}

func TestAddPathname(t *testing.T) {
	ctx := ToContext(context.Background(), DefaultLogger())
	testPathname := "/api/users/123"

	AddPathname(ctx, testPathname)

	extractedPath := ExtractPathname(ctx)
	if extractedPath != testPathname {
		t.Errorf("ExtractPathname() = %v, want %v", extractedPath, testPathname)
	}
}

func TestExtractPathname_NoContext(t *testing.T) {
	ctx := context.Background()
	pathname := ExtractPathname(ctx)
	if pathname != "" {
		t.Errorf("ExtractPathname() from empty context should return empty string, got %v", pathname)
	}
}

func TestExtractPathname_WithDebugPanic(t *testing.T) {
	oldIsDebug := IsDebug
	IsDebug = true
	defer func() {
		IsDebug = oldIsDebug
	}()

	ctx := context.Background()

	defer func() {
		if r := recover(); r == nil {
			t.Error("ExtractPathname() should panic with IsDebug=true and no context logger")
		}
	}()

	ExtractPathname(ctx)
}

func TestWithOTEL_WithLogger(t *testing.T) {
	ctx := ToContext(context.Background(), DefaultLogger())
	logger := WithOTEL(ctx)

	if logger == nil {
		t.Error("WithOTEL() returned nil")
	}

	logger.Info("test with OTEL")
}

func TestWithOTEL_NoLogger(t *testing.T) {
	ctx := context.Background()
	logger := WithOTEL(ctx)

	if logger == nil {
		t.Error("WithOTEL() should return default logger when no logger in context")
	}
}

func TestExtractEntry_WithAllFields(t *testing.T) {
	ctx := ToContext(context.Background(), DefaultLogger())

	AddTraceID(ctx, "trace-123")
	AddUserID(ctx, 999)
	AddPathname(ctx, "/test")
	AddField(ctx, "request_id", "req-456")
	AddFields(ctx, map[string]interface{}{
		"method": "GET",
		"status": 200,
	})

	logger := ExtractEntry(ctx)
	if logger == nil {
		t.Error("ExtractEntry() returned nil")
	}

	logger.Info("test with all fields")
}

func TestExtractEntry_EmptyContext(t *testing.T) {
	ctx := context.Background()
	logger := ExtractEntry(ctx)

	if logger == nil {
		t.Error("ExtractEntry() should return default logger for empty context")
	}
}

func TestAddField_Overwrite(t *testing.T) {
	ctx := ToContext(context.Background(), DefaultLogger())

	AddField(ctx, "key", "value1")
	AddField(ctx, "key", "value2")

	logger := ExtractEntry(ctx)
	logger.Info("test overwrite field")
}

func TestAddFields_Empty(t *testing.T) {
	ctx := ToContext(context.Background(), DefaultLogger())

	AddFields(ctx, map[string]interface{}{})

	logger := ExtractEntry(ctx)
	logger.Info("test empty fields")
}

func TestConcurrentFieldOperations(t *testing.T) {
	ctx := ToContext(context.Background(), DefaultLogger())
	done := make(chan bool)

	for i := 0; i < 50; i++ {
		go func(id int) {
			AddField(ctx, "field", id)
			AddTopField(ctx, "top", id)
			AddFields(ctx, map[string]interface{}{
				"multi1": id,
				"multi2": id * 2,
			})
			done <- true
		}(i)
	}

	for i := 0; i < 50; i++ {
		<-done
	}

	logger := ExtractEntry(ctx)
	logger.Info("test concurrent operations")
}

func TestConcurrentExtractEntry(t *testing.T) {
	ctx := ToContext(context.Background(), DefaultLogger())
	AddField(ctx, "test", "value")

	done := make(chan bool)

	for i := 0; i < 100; i++ {
		go func() {
			logger := ExtractEntry(ctx)
			logger.Info("concurrent extract")
			done <- true
		}()
	}

	for i := 0; i < 100; i++ {
		<-done
	}
}

func TestMultipleToContext(t *testing.T) {
	ctx1 := ToContext(context.Background(), DefaultLogger())
	AddField(ctx1, "ctx", "ctx1")

	logger2 := DefaultLogger().Named("logger2")
	ctx2 := ToContext(context.Background(), logger2)
	AddField(ctx2, "ctx", "ctx2")

	logger1 := ExtractEntry(ctx1)
	logger1.Info("from ctx1")

	logger2Result := ExtractEntry(ctx2)
	logger2Result.Info("from ctx2")
}

func TestNestedContext(t *testing.T) {
	ctx := ToContext(context.Background(), DefaultLogger())
	AddField(ctx, "level1", "value1")

	childCtx := context.WithValue(ctx, "child", "data")
	AddField(childCtx, "level2", "value2")

	logger := ExtractEntry(childCtx)
	logger.Info("nested context test")
}

func BenchmarkToContext(b *testing.B) {
	logger := DefaultLogger()
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ToContext(ctx, logger)
	}
}

func BenchmarkAddField(b *testing.B) {
	ctx := ToContext(context.Background(), DefaultLogger())

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		AddField(ctx, "key", i)
	}
}

func BenchmarkAddFields(b *testing.B) {
	ctx := ToContext(context.Background(), DefaultLogger())
	fields := map[string]interface{}{
		"key1": "value1",
		"key2": 123,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		AddFields(ctx, fields)
	}
}

func BenchmarkExtractEntry_WithFields(b *testing.B) {
	ctx := ToContext(context.Background(), DefaultLogger())
	AddField(ctx, "key1", "value1")
	AddField(ctx, "key2", "value2")
	AddTopField(ctx, "trace_id", "123")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ExtractEntry(ctx)
	}
}

func BenchmarkConcurrentAddField(b *testing.B) {
	ctx := ToContext(context.Background(), DefaultLogger())

	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			AddField(ctx, "key", i)
			i++
		}
	})
}
