package glog

import (
	"context"
	"errors"
	"github.com/gw123/glog/driver/zap_driver"
	"sync"
	"time"

	"github.com/gw123/glog/common"
	"go.opentelemetry.io/otel/trace"
)

type ctxLoggerMarker struct{}

type ctxLogger struct {
	logger    common.Logger
	fields    map[string]interface{}
	topFields map[string]interface{}
	mutex     sync.RWMutex
	lastTime  time.Time
}

var IsDebug bool = false

var (
	ctxLoggerKey = &ctxLoggerMarker{}
)

//添加日志字段到日志中间件(ctx_logrus)，添加的字段会在后面调用 info，debug，error 时候输出
func AddFields(ctx context.Context, fields map[string]interface{}) {
	l, ok := ctx.Value(ctxLoggerKey).(*ctxLogger)
	if !ok || l == nil {
		return
	}
	l.mutex.Lock()
	for k, v := range fields {
		l.fields[k] = v
	}
	l.mutex.Unlock()
}

func AddTopField(ctx context.Context, key string, val interface{}) {
	l, ok := ctx.Value(ctxLoggerKey).(*ctxLogger)
	if !ok || l == nil {
		return
	}

	l.mutex.Lock()
	l.topFields[key] = val
	l.mutex.Unlock()
}

//添加日志字段到日志中间件(ctx_logrus)，添加的字段会在后面调用 info，debug，error 时候输出
func AddField(ctx context.Context, key string, val interface{}) {
	l, ok := ctx.Value(ctxLoggerKey).(*ctxLogger)
	if !ok || l == nil {
		return
	}

	l.mutex.Lock()
	l.fields[key] = val
	l.mutex.Unlock()
}

// 添加一个追踪规矩id 用来聚合同一次请求, 注意要用返回的contxt 替换传入的ctx
func AddTraceID(ctx context.Context, traceID string) {
	AddTopField(ctx, common.KeyTraceID, traceID)
}

//export requestID
func ExtractTraceID(ctx context.Context) string {
	l, ok := ctx.Value(ctxLoggerKey).(*ctxLogger)
	if !ok || l == nil {
		if IsDebug {
			panic(errors.New("not set ctxLogger"))
		}
		return ""
	}

	l.mutex.RLock()
	val, ok := l.topFields[common.KeyTraceID].(string)
	l.mutex.RUnlock()
	if ok {
		return val
	}
	return ""
}

//export Otel traceID
func WithOTEL(ctx context.Context) common.Logger {
	l, ok := ctx.Value(ctxLoggerKey).(*ctxLogger)
	var logger common.Logger
	if !ok || l == nil {
		logger = DefaultLogger()
	} else {
		logger = l.logger
	}

	if span := trace.SpanContextFromContext(ctx); span.TraceID().IsValid() {
		return logger.WithField("trace_id", span.TraceID().String())
	}
	return logger
}

//add userID to ctx
func AddUserID(ctx context.Context, userID int64) {
	AddTopField(ctx, common.KeyUserID, userID)
}

//export userID
func ExtractUserID(ctx context.Context) int64 {
	l, ok := ctx.Value(ctxLoggerKey).(*ctxLogger)
	if !ok || l == nil {
		if IsDebug {
			panic(errors.New("not set ctxLogger"))
		}
		return 0
	}

	l.mutex.RLock()
	val, ok := l.topFields[common.KeyUserID].(int64)
	l.mutex.RUnlock()
	if ok {
		return val
	}
	return val
}

//add userID to ctx
func AddPathname(ctx context.Context, pathname string) {
	AddTopField(ctx, common.KeyPathname, pathname)
}

//export pathname
func ExtractPathname(ctx context.Context) string {
	l, ok := ctx.Value(ctxLoggerKey).(*ctxLogger)
	if !ok || l == nil {
		if IsDebug {
			panic(errors.New("not set ctxLogger"))
		}
		return ""
	}

	l.mutex.RLock()
	val, ok := l.topFields[common.KeyPathname].(string)
	l.mutex.RUnlock()
	if ok {
		return val
	}
	return val
}

// 添加logrus.Entry到context, 这个操作添加的logrus.Entry在后面AddFields和Extract都会使用到
func ToContext(ctx context.Context, entry common.Logger) context.Context {
	l := &ctxLogger{
		logger:    entry,
		fields:    map[string]interface{}{},
		topFields: map[string]interface{}{},
		mutex:     sync.RWMutex{},
		lastTime:  time.Time{},
	}
	return context.WithValue(ctx, ctxLoggerKey, l)
}

//extract ctx_logrus logrus_driver.Entry
func ExtractEntry(ctx context.Context) common.Logger {
	l, ok := ctx.Value(ctxLoggerKey).(*ctxLogger)
	if !ok || l == nil {
		return zap_driver.GetInnerLogger()
	}
	return l.logger.WithFields(l.topFields).WithFields(l.fields)
}
