package glog

import (
	"context"
	"errors"
	"sync"

	"github.com/gw123/glog/common"
	"go.opentelemetry.io/otel/trace"
)

type ctxLoggerMarker struct{}

type ctxLogger struct {
	logger    common.Logger
	fields    map[string]interface{}
	topFields map[string]interface{}
	mutex     sync.RWMutex
}

var IsDebug bool = false

var (
	ctxLoggerKey = &ctxLoggerMarker{}
)

// AddFields adds multiple log fields to the context logger
func AddFields(ctx context.Context, fields map[string]interface{}) {
	l, ok := ctx.Value(ctxLoggerKey).(*ctxLogger)
	if !ok || l == nil {
		return
	}
	l.mutex.Lock()
	newFields := make(map[string]interface{}, len(l.fields)+len(fields))
	for k, v := range l.fields {
		newFields[k] = v
	}
	for k, v := range fields {
		newFields[k] = v
	}
	l.fields = newFields
	l.mutex.Unlock()
}

func AddTopField(ctx context.Context, key string, val interface{}) {
	l, ok := ctx.Value(ctxLoggerKey).(*ctxLogger)
	if !ok || l == nil {
		return
	}

	l.mutex.Lock()
	newTopFields := make(map[string]interface{}, len(l.topFields)+1)
	for k, v := range l.topFields {
		newTopFields[k] = v
	}
	newTopFields[key] = val
	l.topFields = newTopFields
	l.mutex.Unlock()
}

// AddField adds a single log field to the context logger
func AddField(ctx context.Context, key string, val interface{}) {
	l, ok := ctx.Value(ctxLoggerKey).(*ctxLogger)
	if !ok || l == nil {
		return
	}

	l.mutex.Lock()
	newFields := make(map[string]interface{}, len(l.fields)+1)
	for k, v := range l.fields {
		newFields[k] = v
	}
	newFields[key] = val
	l.fields = newFields
	l.mutex.Unlock()
}

// AddTraceID adds a trace ID for aggregating logs from the same request
func AddTraceID(ctx context.Context, traceID string) {
	AddTopField(ctx, common.KeyTraceID, traceID)
}

// ExtractTraceID extracts the trace ID from the context
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

	if span := trace.SpanContextFromContext(ctx); span.TraceID().IsValid() {
		return span.TraceID().String()
	}
	return ""
}

// WithOTEL extracts OpenTelemetry trace ID from context and adds it to the logger
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

// AddUserID add userID to ctx
func AddUserID(ctx context.Context, userID int64) {
	AddTopField(ctx, common.KeyUserID, userID)
}

// ExtractUserID userID
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
	return 0
}

// AddPathname add userID to ctx
func AddPathname(ctx context.Context, pathname string) {
	AddTopField(ctx, common.KeyPathname, pathname)
}

// ExtractPathname pathname
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
	return ""
}

// ToContext adds a logger to the context for use in subsequent log operations
func ToContext(ctx context.Context, entry common.Logger) context.Context {
	l := &ctxLogger{
		logger:    entry,
		fields:    map[string]interface{}{},
		topFields: map[string]interface{}{},
		mutex:     sync.RWMutex{},
	}
	return context.WithValue(ctx, ctxLoggerKey, l)
}

// ExtractEntry extracts the logger from context with all accumulated fields
func ExtractEntry(ctx context.Context) common.Logger {
	var logger common.Logger
	l, ok := ctx.Value(ctxLoggerKey).(*ctxLogger)
	if ok && l != nil {
		l.mutex.RLock()
		topFields := l.topFields
		fields := l.fields
		l.mutex.RUnlock()
		logger = l.logger.WithFields(topFields).WithFields(fields)
	} else {
		logger = DefaultLogger()
	}
	if tID := ExtractTraceID(ctx); tID != "" {
		return logger.WithField("trace_id", tID)
	}
	return logger
}
