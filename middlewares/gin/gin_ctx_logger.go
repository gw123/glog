package gin

import (
	"github.com/gin-gonic/gin"
	"github.com/gw123/glog"
	glogCommon "github.com/gw123/glog/common"
	"github.com/sirupsen/logrus"
)

type ctxLogger struct {
	logger glogCommon.Logger
	fields logrus.Fields
}

type Record struct {
	TraceId   int64  `json:"trace_id"`
	CreatedAt int64  `json:"uint_64"`
	Point     string `json:"point"`
	Place     string `json:"place"`
}

type requestIdKey struct{}
type loggerKey struct{}

var (
	isDebug         = false
	ctxRequestIdKey = "&requestIdKey{}"
	ctxLoggerKey    = "&loggerKey{}"
)

func SetDebug(flag bool) {
	isDebug = flag
}

// 为了方便创建一个默认的Logger
func DefaultLogger() glogCommon.Logger {
	return glog.DefaultLogger()
}

// 为了方便创建一个默认的Logger
func DefaultJsonLogger() glogCommon.Logger {
	return DefaultLogger()
}

// 添加logrus.Entry到context, 这个操作添加的logrus.Entry在后面AddFields和Extract都会使用到
func ToContext(ctx *gin.Context, logger glogCommon.Logger) {
	if ctx == nil {
		return
	}

	l := &ctxLogger{
		logger: logger,
		fields: logrus.Fields{},
	}
	ctx.Set(ctxLoggerKey, l)
}

//添加日志字段到日志中间件(ctx_logrus)，添加的字段会在后面调用 info，debug，error 时候输出
func AddFields(ctx *gin.Context, fields logrus.Fields) {
	if ctx == nil {
		return
	}

	l, ok := ctx.Get(ctxLoggerKey)
	if !ok || l == nil {
		return
	}
	log, ok := l.(*ctxLogger)
	if !ok {
		return
	}
	for k, v := range fields {
		log.fields[k] = v
	}
}

//添加日志字段到日志中间件(ctx_logrus)，添加的字段会在后面调用 info，debug，error 时候输出
func AddField(ctx *gin.Context, key, val string) {
	if ctx == nil {
		return
	}

	l, ok := ctx.Get(ctxLoggerKey)
	if !ok || l == nil {
		return
	}
	log, ok := l.(*ctxLogger)
	if !ok {
		return
	}
	log.fields[key] = val
}

// 添加一个追踪规矩id 用来聚合同一次请求, 注意要用返回的contxt 替换传入的ctx
func AddRequestId(ctx *gin.Context, requestId string) {
	if ctx == nil {
		return
	}

	ctx.Set(ctxRequestIdKey, requestId)
}

//导出requestId
func ExtractRequestId(ctx *gin.Context) string {
	if ctx == nil {
		return ""
	}

	l, ok := ctx.Get(ctxRequestIdKey)
	if !ok {
		return ""
	}
	return l.(string)
}

//导出ctx_logrus日志库
func ExtractEntry(ctx *gin.Context) glogCommon.Logger {
	l, ok := ctx.Get(ctxLoggerKey)
	if !ok || l == nil {
		return glog.DefaultLogger()
	}
	log, ok := l.(*ctxLogger)
	if !ok {
		return glog.DefaultLogger()
	}
	if !ok || log == nil {
		return glog.DefaultLogger()
	}

	fields := make(map[string]interface{})
	for k, v := range log.fields {
		fields[k] = v
	}

	requestId := ExtractRequestId(ctx)
	if requestId != "" {
		fields[ctxRequestIdKey] = requestId
	}

	logger := log.logger
	return logger.WithFields(fields)
}
