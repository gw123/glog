package glog

import (
	"context"
	"github.com/sirupsen/logrus"
)

type ctxLoggerMarker struct{}
type ctxRequestIdMarker struct{}

type ctxLogger struct {
	logger *logrus.Entry
	fields logrus.Fields
}

type Record struct {
	TraceId   int64  `json:"trace_id"`
	CreatedAt int64  `json:"uint_64"`
	Point     string `json:"point"`
	Place     string `json:"place"`
}

const RequestId  = "RequestId"
const TimeFormat = "2006-01-02 15:04:05"

var (
	ctxLoggerKey    = &ctxLoggerMarker{}
	ctxRequestIdKey = &ctxRequestIdMarker{}

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

// 添加logrus.Entry到context, 这个操作添加的logrus.Entry在后面AddFields和Extract都会使用到
func ToContext(ctx context.Context, entry *logrus.Entry) context.Context {
	l := &ctxLogger{
		logger: entry,
		fields: logrus.Fields{},
	}

	return context.WithValue(ctx, ctxLoggerKey, l)
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
func AddRequestId(ctx context.Context, requestId string) context.Context {
	return context.WithValue(ctx, ctxRequestIdKey, requestId)
}

//导出requestId
func ExtractRequestId(ctx context.Context) string {
	l, ok := ctx.Value(ctxRequestIdKey).(string)
	if !ok {
		return ""
	}
	return l
}

//导出ctx_logrus日志库
func ExtractEntry(ctx context.Context) *logrus.Entry {
	l, ok := ctx.Value(ctxLoggerKey).(*ctxLogger)
	if !ok || l == nil {
		return logrus.NewEntry(logrus.New())
	}

	fields := logrus.Fields{}
	for k, v := range l.fields {
		fields[k] = v
	}

	requestId := ExtractRequestId(ctx)
	if requestId != "" {
		fields[RequestId] = requestId
	}
	return l.logger.WithFields(fields)
}
