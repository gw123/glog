package glog

import (
	"context"
	"encoding/json"
	"errors"
	"sync"

	"github.com/sirupsen/logrus"
)

type ctxLoggerMarker struct{}

type ctxLogger struct {
	logger    *logrus.Entry
	fields    map[string]interface{}
	topFields map[string]interface{}
	mutex     sync.RWMutex
}

var IsDebug bool = false

const KeyRequestID = "request_id"
const KeyUserID = "user_id"

const TimeFormat = "2006-01-02 15:04:05"

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
	l.fields[key] = val
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
func AddRequestID(ctx context.Context, requestID string) {
	AddField(ctx, KeyRequestID, requestID)
}

//export requestID
func ExtractRequestID(ctx context.Context) string {
	l, ok := ctx.Value(ctxLoggerKey).(*ctxLogger)
	if !ok || l == nil {
		if IsDebug {
			panic(errors.New("not set ctxLogger"))
		}
		return ""
	}

	l.mutex.RLock()
	val, ok := l.fields[KeyRequestID].(string)
	l.mutex.RUnlock()
	if ok {
		return val
	}
	return ""
}

//add userID to ctx
func AddUserID(ctx context.Context, userID int64) {
	AddField(ctx, KeyUserID, userID)
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
	val, ok := l.fields[KeyUserID].(int64)
	l.mutex.RUnlock()
	if ok {
		return val
	}
	return val
}

// 添加logrus.Entry到context, 这个操作添加的logrus.Entry在后面AddFields和Extract都会使用到
func ToContext(ctx context.Context, entry *logrus.Entry) context.Context {
	l := &ctxLogger{
		logger:    entry,
		fields:    map[string]interface{}{},
		topFields: map[string]interface{}{},
		mutex:     sync.RWMutex{},
	}
	return context.WithValue(ctx, ctxLoggerKey, l)
}

//extract ctx_logrus logrus.Entry
func ExtractEntry(ctx context.Context) *logrus.Entry {
	l, ok := ctx.Value(ctxLoggerKey).(*ctxLogger)
	if !ok || l == nil {
		return logrus.NewEntry(logrus.New())
	}

	requestID := ExtractRequestID(ctx)
	userID := ExtractUserID(ctx)

	l.mutex.Lock()
	l.topFields[KeyRequestID] = requestID
	l.topFields[KeyUserID] = userID

	levelTwo, err := json.Marshal(l.fields)
	if err == nil {
		l.topFields["extra"] = string(levelTwo)
	} else {
		l.topFields["extra"] = err.Error()
	}

	l.mutex.Unlock()
	return l.logger.WithFields(l.topFields)
}
