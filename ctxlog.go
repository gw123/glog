package glog

import (
	"context"
	"github.com/sirupsen/logrus"
	"sync"
)

type ctxLoggerMarker struct{}
type ctxRequestIDMarker struct{}
type ctxUserIDMarker struct{}
type ctxComIDMarker struct{}

type ctxLogger struct {
	logger *logrus.Entry
	fields logrus.Fields
	mutex  sync.Mutex
}

type Record struct {
	TraceID   int64  `json:"trace_id"`
	CreatedAt int64  `json:"uint_64"`
	Point     string `json:"point"`
	Place     string `json:"place"`
}

const RequestID = "request_id"
const UserID = "user_id"
const ComID = "com_id"
const TimeFormat = "2006-01-02 15:04:05"

var (
	ctxLoggerKey    = &ctxLoggerMarker{}
	ctxRequestIDKey = &ctxRequestIDMarker{}
	ctxUserIDKey    = &ctxUserIDMarker{}
	ctxComIDKey     = &ctxComIDMarker{}

	defaultLogger *logrus.Logger
	isDebug       = false
)

func SetDebug(flag bool) {
	isDebug = flag
}

// 为了方便创建一个默认的Logger
func DefaultLogger() *logrus.Logger {
	if defaultLogger == nil {
		defaultLogger = logrus.New()
	}
	defaultLogger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat:  TimeFormat,
		DisableTimestamp: false,
		DataKey:          "",
		FieldMap:         nil,
		CallerPrettyfier: nil,
		PrettyPrint:      isDebug,
	})
	return defaultLogger
}

func NewDefaultEntry() *logrus.Entry {
	return logrus.NewEntry(DefaultLogger())
}

//添加日志字段到日志中间件(ctx_logrus)，添加的字段会在后面调用 info，debug，error 时候输出
func AddFields(ctx context.Context, fields logrus.Fields) {
	l, ok := ctx.Value(ctxLoggerKey).(*ctxLogger)
	if !ok || l == nil {
		return
	}
	for k, v := range fields {
		l.fields[k] = v
	}
}

//添加日志字段到日志中间件(ctx_logrus)，添加的字段会在后面调用 info，debug，error 时候输出
func AddField(ctx context.Context, key, val string) {
	l, ok := ctx.Value(ctxLoggerKey).(*ctxLogger)
	if !ok || l == nil {
		return
	}
	l.fields[key] = val
}

// 添加一个追踪规矩id 用来聚合同一次请求, 注意要用返回的contxt 替换传入的ctx
func AddRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, ctxRequestIDKey, requestID)
}

//导出requestID
func ExtractRequestID(ctx context.Context) string {
	l, ok := ctx.Value(ctxRequestIDKey).(string)
	if !ok {
		return ""
	}
	return l
}

//add userID to ctx
func AddUserID(ctx context.Context, userID int64) context.Context {
	return context.WithValue(ctx, ctxUserIDKey, userID)
}

//export userID
func ExtractUserID(ctx context.Context) int64 {
	l, ok := ctx.Value(ctxUserIDKey).(int64)
	if !ok {
		return 0
	}
	return l
}

//add comID to ctx
func AddComID(ctx context.Context, comID int64) context.Context {
	return context.WithValue(ctx, ctxComIDKey, comID)
}

//export comID
func ExtractComID(ctx context.Context) int64 {
	l, ok := ctx.Value(ctxComIDKey).(int64)
	if !ok {
		return 0
	}
	return l
}

// 添加logrus.Entry到context, 这个操作添加的logrus.Entry在后面AddFields和Extract都会使用到
func ToContext(ctx context.Context, entry *logrus.Entry) context.Context {
	l := &ctxLogger{
		logger: entry,
		fields: logrus.Fields{},
		mutex:  sync.Mutex{},
	}

	return context.WithValue(ctx, ctxLoggerKey, l)
}

//extract ctx_logrus logrus.Entry
func ExtractEntry(ctx context.Context) *logrus.Entry {
	l, ok := ctx.Value(ctxLoggerKey).(*ctxLogger)
	if !ok || l == nil {
		return logrus.NewEntry(logrus.New())
	}

	l.mutex.Lock()
	defer l.mutex.Unlock()

	fields := logrus.Fields{}
	for k, v := range l.fields {
		fields[k] = v
	}

	requestID := ExtractRequestID(ctx)
	fields[RequestID] = requestID

	userID := ExtractUserID(ctx)
	fields[UserID] = userID
	return l.logger.WithFields(fields)
}
