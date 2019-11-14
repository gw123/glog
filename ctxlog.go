package glog

import (
	"context"
	"github.com/sirupsen/logrus"
)

type ctxLoggerMarker struct{}

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

var (
	ctxLoggerKey = &ctxLoggerMarker{}
)

// 添加logrus.Entry到context, 这个操作添加的logrus.Entry在后面AddFields和Extract都会使用到
func ToContext(ctx context.Context, entry *logrus.Entry) context.Context {
	l := &ctxLogger{
		logger: entry,
		fields: logrus.Fields{},
	}
	return context.WithValue(ctx, ctxLoggerKey, l)
}

//添加日志字段到日志中间件(ctx_logrus)
func AddFields(ctx context.Context, fields logrus.Fields) {
	l, ok := ctx.Value(ctxLoggerKey).(*ctxLogger)
	if !ok || l == nil {
		return
	}
	for k, v := range fields {
		l.fields[k] = v
	}
}

//导出ctx_logrus日志库
func Extract(ctx context.Context) *logrus.Entry {
	l, ok := ctx.Value(ctxLoggerKey).(*ctxLogger)
	if !ok || l == nil {
		return logrus.NewEntry(logrus.New())
	}

	fields := logrus.Fields{}

	// Add logrus fields added until now.
	for k, v := range l.fields {
		fields[k] = v
	}

	return l.logger.WithFields(fields)
}
