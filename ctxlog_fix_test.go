package glog

import (
	"context"
	"testing"
)

func TestAddField_ConcurrentSafety(t *testing.T) {
	ctx := ToContext(context.Background(), DefaultLogger())
	done := make(chan bool)

	for i := 0; i < 100; i++ {
		go func(id int) {
			AddField(ctx, "key", id)
			done <- true
		}(i)
	}

	for i := 0; i < 100; i++ {
		<-done
	}
}

func TestExtractEntry_WithOTEL(t *testing.T) {
	ctx := ToContext(context.Background(), DefaultLogger())
	AddTraceID(ctx, "test-trace-123")
	AddUserID(ctx, 12345)
	AddPathname(ctx, "/api/test")

	logger := ExtractEntry(ctx)
	if logger == nil {
		t.Error("ExtractEntry() returned nil")
	}

	traceID := ExtractTraceID(ctx)
	if traceID != "test-trace-123" {
		t.Errorf("ExtractTraceID() = %v, want %v", traceID, "test-trace-123")
	}

	userID := ExtractUserID(ctx)
	if userID != 12345 {
		t.Errorf("ExtractUserID() = %v, want %v", userID, 12345)
	}

	pathname := ExtractPathname(ctx)
	if pathname != "/api/test" {
		t.Errorf("ExtractPathname() = %v, want %v", pathname, "/api/test")
	}
}

func TestAddFields_Multiple(t *testing.T) {
	ctx := ToContext(context.Background(), DefaultLogger())

	fields := map[string]interface{}{
		"request_id": "req-123",
		"user_agent": "test-agent",
		"ip":         "127.0.0.1",
	}

	AddFields(ctx, fields)

	logger := ExtractEntry(ctx)
	logger.Info("test message with multiple fields")
}

func TestExtractEntry_NoContext(t *testing.T) {
	ctx := context.Background()
	logger := ExtractEntry(ctx)
	if logger == nil {
		t.Error("ExtractEntry() should return default logger when context has no logger")
	}
	logger.Info("should use default logger")
}

func TestIsDebug_Panic(t *testing.T) {
	oldIsDebug := IsDebug
	IsDebug = true
	defer func() {
		IsDebug = oldIsDebug
	}()

	ctx := context.Background()

	defer func() {
		if r := recover(); r == nil {
			t.Error("ExtractTraceID() should panic when IsDebug=true and no logger in context")
		}
	}()

	ExtractTraceID(ctx)
}

func TestWithOTEL_NoContext(t *testing.T) {
	ctx := context.Background()
	logger := WithOTEL(ctx)
	if logger == nil {
		t.Error("WithOTEL() should return default logger when context has no logger")
	}
}

func BenchmarkExtractEntry(b *testing.B) {
	ctx := ToContext(context.Background(), DefaultLogger())
	AddField(ctx, "test_key", "test_value")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger := ExtractEntry(ctx)
		_ = logger
	}
}

func BenchmarkAddField_Concurrent(b *testing.B) {
	ctx := ToContext(context.Background(), DefaultLogger())

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			AddField(ctx, "key", "value")
		}
	})
}

func TestALl(t *testing.T) {
	//1. 一般是在应用的入口创建一个根context
	ctx := context.Background()
	//2. 新建立entry 使用ToContext将entry传入context
	entry := DefaultLogger()
	//3. 加入一些应用全局的描述
	entry.WithField("app", "xytschol")
	//4. 然后将这个根context传到应用框架中

	//5. 在中间件里调用AddRequestID(ctx) 记录一次请求的同一个 requestID
	ctx = ToContext(ctx, entry)
	AddTraceID(ctx, "10000001")
	AddPathname(ctx, "/home/index")

	//6. 在action 或者service 等地方记录日志记录日志
	entry = ExtractEntry(ctx).WithField("ip", "10.0.0.1")
	entry = ExtractEntry(ctx).WithField("ip2", "10.0.0.1")
	for i := 0; i < 20; i++ {
		entry.WithField("key", i).Infof("TestContent abc %d", i)
	}
	//输出结果
	//{"RequestID":"10000001","app_name":"web","ip":"10.0.0.1","level":"info","msg":"TestContent","time":"2020-03-17 20:34:14"}
}
