# glog

[![Go Version](https://img.shields.io/badge/go-%3E%3D1.17-blue.svg)](https://golang.org/)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)

基于 [Uber Zap](https://github.com/uber-go/zap) 封装的高性能结构化日志库，提供简洁的 API 和强大的上下文追踪功能。

## 特性

- **高性能**: 基于 Uber Zap，提供出色的日志性能
- **结构化日志**: 支持字段化的结构化日志输出
- **上下文追踪**: 通过 `context.Context` 实现请求级别的日志追踪
- **OpenTelemetry 集成**: 自动提取 trace_id，支持分布式追踪
- **命名日志器**: 支持为不同组件创建独立的命名日志器
- **灵活配置**: 支持多种输出格式（Console/JSON）和输出目标
- **并发安全**: 完全支持高并发场景下的日志记录

## 安装

```bash
go get github.com/gw123/glog
```

## 快速开始

### 基础使用

```go
package main

import "github.com/gw123/glog"

func main() {
    // 简单日志输出
    glog.Info("Application started")
    glog.Infof("Server listening on port %d", 8080)
    
    // 带字段的结构化日志
    glog.WithField("user_id", 12345).Info("User logged in")
    
    // 链式调用添加多个字段
    glog.Log().
        WithField("method", "GET").
        WithField("path", "/api/users").
        WithField("status", 200).
        Info("HTTP request processed")
    
    // 错误日志
    err := someOperation()
    glog.WithError(err).Error("Operation failed")
}
```

**输出示例:**
```
[2025-11-04 10:30:15.891] [info] [] main.go:6 []  Application started
[2025-11-04 10:30:15.892] [info] [] main.go:7 []  Server listening on port 8080 {"user_id": 12345}
[2025-11-04 10:30:15.893] [info] [] main.go:10 []  User logged in {"method": "GET", "path": "/api/users", "status": 200}
```

### 上下文日志追踪

`glog` 的核心功能之一是通过 `context.Context` 实现请求级别的日志追踪，在整个请求链路中自动传递和累积日志字段。

```go
package main

import (
    "context"
    "github.com/gw123/glog"
)

func main() {
    // 1. 创建根 context 并注入 logger
    ctx := context.Background()
    ctx = glog.ToContext(ctx, glog.DefaultLogger())
    
    // 2. 添加应用全局字段
    glog.AddField(ctx, "app_name", "web-service")
    glog.AddField(ctx, "version", "1.0.0")
    
    // 3. 在中间件中添加请求级别的追踪信息
    handleRequest(ctx)
}

func handleRequest(ctx context.Context) {
    // 添加请求追踪 ID 和用户信息
    glog.AddTraceID(ctx, "trace-abc-123")
    glog.AddUserID(ctx, 12345)
    glog.AddPathname(ctx, "/api/orders")
    glog.AddField(ctx, "client_ip", "192.168.1.100")
    
    // 提取 logger，自动包含之前添加的所有字段
    logger := glog.ExtractEntry(ctx)
    logger.Info("Request received")
    
    // 在业务逻辑中继续添加字段
    processOrder(ctx)
    
    logger.Info("Request completed")
}

func processOrder(ctx context.Context) {
    // 继续添加业务相关字段
    glog.AddField(ctx, "order_id", "ORD-12345")
    glog.AddField(ctx, "amount", 299.99)
    
    // 提取的 logger 包含从 root context 到此处的所有字段
    glog.ExtractEntry(ctx).Info("Processing order")
    
    // 调用其他服务
    validatePayment(ctx)
}

func validatePayment(ctx context.Context) {
    glog.AddField(ctx, "payment_method", "credit_card")
    
    // 所有日志都包含完整的上下文链路信息
    glog.ExtractEntry(ctx).Info("Payment validated")
}
```

**输出示例:**
```json
{
  "level": "info",
  "time": "2025-11-04 10:30:15",
  "msg": "Request received",
  "trace_id": "trace-abc-123",
  "user_id": 12345,
  "pathname": "/api/orders",
  "client_ip": "192.168.1.100",
  "app_name": "web-service",
  "version": "1.0.0"
}
{
  "level": "info",
  "time": "2025-11-04 10:30:15",
  "msg": "Processing order",
  "trace_id": "trace-abc-123",
  "user_id": 12345,
  "pathname": "/api/orders",
  "client_ip": "192.168.1.100",
  "app_name": "web-service",
  "version": "1.0.0",
  "order_id": "ORD-12345",
  "amount": 299.99
}
{
  "level": "info",
  "time": "2025-11-04 10:30:15",
  "msg": "Payment validated",
  "trace_id": "trace-abc-123",
  "user_id": 12345,
  "pathname": "/api/orders",
  "client_ip": "192.168.1.100",
  "app_name": "web-service",
  "version": "1.0.0",
  "order_id": "ORD-12345",
  "amount": 299.99,
  "payment_method": "credit_card"
}
```

### 命名日志器

为不同组件创建独立的命名日志器，便于日志过滤和分析：

```go
// 数据库日志器
dbLogger := glog.Log().Named("database")
dbLogger.Info("Database connection established")
dbLogger.WithField("pool_size", 10).Info("Connection pool initialized")

// 缓存日志器
cacheLogger := glog.Log().Named("cache")
cacheLogger.Info("Redis connected")

// API 日志器
apiLogger := glog.Log().Named("api")
apiLogger.WithField("port", 8080).Info("API server started")
```

**输出示例:**
```
[2025-11-04 10:30:15.891] [info] [database] db.go:15 []  Database connection established
[2025-11-04 10:30:15.892] [info] [cache] cache.go:22 []  Redis connected
[2025-11-04 10:30:15.893] [info] [api] server.go:45 []  API server started {"port": 8080}
```

## 自定义配置

### 基础配置

```go
import (
    "github.com/gw123/glog"
    "github.com/gw123/glog/common"
)

func main() {
    // 配置日志器
    err := glog.SetDefaultLoggerConfig(
        common.Options{},
        common.WithLevel(common.DebugLevel),        // 设置日志级别
        common.WithConsoleEncoding(),               // 使用 Console 格式
        common.WithOutputPath("./logs/app.log"),    // 输出到文件
        common.WithStdoutOutputPath(),              // 同时输出到标准输出
    )
    if err != nil {
        panic(err)
    }
    
    // 现在可以看到 debug 日志
    glog.Debug("This is a debug message")
    glog.Info("Application configured")
}
```

### JSON 格式输出

```go
err := glog.SetDefaultLoggerConfig(
    common.Options{},
    common.WithLevel(common.InfoLevel),
    common.WithJsonEncoding(),                  // 使用 JSON 格式
    common.WithOutputPath("./logs/app.json"),
)
```

### 可用的配置选项

**日志级别:**
- `common.DebugLevel` - 调试级别
- `common.InfoLevel` - 信息级别（默认）
- `common.WarnLevel` - 警告级别
- `common.ErrorLevel` - 错误级别
- `common.DPanicLevel` - 开发环境 panic
- `common.PanicLevel` - Panic 级别
- `common.FatalLevel` - Fatal 级别

**编码格式:**
- `common.WithConsoleEncoding()` - 人类可读的 Console 格式
- `common.WithJsonEncoding()` - 机器可解析的 JSON 格式

**输出目标:**
- `common.WithStdoutOutputPath()` - 标准输出
- `common.WithOutputPath(path)` - 自定义文件路径
- `common.WithErrorOutputPath(path)` - 错误日志输出路径

## OpenTelemetry 集成

`glog` 自动集成 OpenTelemetry，在调用 `ExtractEntry(ctx)` 时自动提取 trace_id：

```go
import (
    "context"
    "go.opentelemetry.io/otel"
    "github.com/gw123/glog"
)

func handleRequest(ctx context.Context) {
    // 假设已经通过 OTEL 中间件注入了 span
    // ExtractEntry 会自动提取 trace_id
    logger := glog.ExtractEntry(ctx)
    logger.Info("Processing with OTEL trace")
    
    // 或者显式使用 WithOTEL
    glog.WithOTEL(ctx).Info("Explicit OTEL integration")
}
```

## API 参考

### 顶层日志函数

```go
// 简单日志
glog.Debug(args ...interface{})
glog.Info(args ...interface{})
glog.Warn(args ...interface{})
glog.Error(args ...interface{})

// 格式化日志
glog.Debugf(format string, args ...interface{})
glog.Infof(format string, args ...interface{})
glog.Warnf(format string, args ...interface{})
glog.Errorf(format string, args ...interface{})

// 获取 logger 实例
logger := glog.Log()
defaultLogger := glog.DefaultLogger()
```

### 字段操作

```go
// 添加单个字段
glog.WithField(key string, value interface{})

// 添加错误字段
glog.WithError(err error)

// 创建命名 logger
logger := glog.Log().Named("component-name")
```

### 上下文日志函数

```go
// 注入 logger 到 context
ctx = glog.ToContext(ctx, logger)

// 添加字段到 context
glog.AddField(ctx, key string, value interface{})
glog.AddFields(ctx, fields map[string]interface{})

// 添加顶级字段（优先级更高）
glog.AddTopField(ctx, key string, value interface{})

// 添加预定义字段
glog.AddTraceID(ctx, traceID string)
glog.AddUserID(ctx, userID int64)
glog.AddPathname(ctx, pathname string)
glog.AddClientIP(ctx, clientIP string)

// 从 context 提取 logger（自动包含所有字段和 OTEL trace_id）
logger := glog.ExtractEntry(ctx)

// 从 context 提取字段值
traceID := glog.ExtractTraceID(ctx)
userID := glog.ExtractUserID(ctx)

// 显式 OTEL 集成
logger := glog.WithOTEL(ctx)
```

## 示例代码

在 `demo/` 目录下提供了完整的示例代码：

### 基础示例 (`demo/basic/`)
展示所有基础功能，包括：
- 基础日志输出
- 字段化日志
- 上下文日志追踪
- 错误日志
- 命名日志器
- 自定义配置
- 并发日志
- 不同日志级别

运行示例：
```bash
go run demo/basic/main.go
```

### Web 应用示例 (`demo/webapp/`)
模拟真实 Web 应用场景，展示：
- 应用启动流程
- HTTP 请求处理
- 数据库和缓存操作
- 请求链路追踪
- 优雅关闭

运行示例：
```bash
go run demo/webapp/main.go
```

### 性能测试 (`demo/performance/`)
性能基准测试和压力测试。

运行示例：
```bash
go run demo/performance/main.go
```

## 最佳实践

### 1. Web 应用中的使用模式

```go
// 在应用启动时配置 logger
func init() {
    glog.SetDefaultLoggerConfig(
        common.Options{},
        common.WithLevel(common.InfoLevel),
        common.WithJsonEncoding(),
        common.WithOutputPath("./logs/app.log"),
        common.WithStdoutOutputPath(),
    )
}

// 在 HTTP 中间件中初始化请求上下文
func LoggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // 创建请求 context
        ctx := glog.ToContext(r.Context(), glog.Log().Named("http"))
        
        // 添加请求追踪信息
        requestID := generateRequestID()
        glog.AddTraceID(ctx, requestID)
        glog.AddPathname(ctx, r.URL.Path)
        glog.AddClientIP(ctx, r.RemoteAddr)
        glog.AddField(ctx, "method", r.Method)
        
        // 记录请求开始
        glog.ExtractEntry(ctx).Info("Request started")
        
        // 处理请求
        next.ServeHTTP(w, r.WithContext(ctx))
        
        // 记录请求结束
        glog.ExtractEntry(ctx).Info("Request completed")
    })
}

// 在业务逻辑中使用
func CreateOrder(ctx context.Context, order Order) error {
    glog.AddField(ctx, "order_id", order.ID)
    glog.AddField(ctx, "amount", order.Amount)
    
    logger := glog.ExtractEntry(ctx)
    logger.Info("Creating order")
    
    // 业务逻辑...
    
    logger.Info("Order created successfully")
    return nil
}
```

### 2. 微服务中的追踪

在微服务架构中，建议结合 OpenTelemetry 使用：

```go
import (
    "go.opentelemetry.io/otel"
    "github.com/gw123/glog"
)

func ServiceA(ctx context.Context) {
    // OTEL 会自动注入 span 到 context
    tracer := otel.Tracer("service-a")
    ctx, span := tracer.Start(ctx, "ServiceA.Process")
    defer span.End()
    
    // ExtractEntry 会自动提取 OTEL trace_id
    logger := glog.ExtractEntry(ctx)
    logger.Info("Service A processing")
    
    // 调用其他服务，trace_id 会自动传递
    ServiceB(ctx)
}
```

### 3. 错误处理

```go
func ProcessData(ctx context.Context, data []byte) error {
    logger := glog.ExtractEntry(ctx)
    
    result, err := parseData(data)
    if err != nil {
        logger.WithError(err).
            WithField("data_size", len(data)).
            Error("Failed to parse data")
        return err
    }
    
    logger.WithField("records_count", len(result)).
        Info("Data processed successfully")
    return nil
}
```

## 性能

`glog` 基于 Uber Zap 构建，提供出色的性能表现：

- **零内存分配**: 在生产环境中几乎无内存分配
- **高吞吐量**: 支持每秒百万级日志记录
- **低延迟**: 微秒级的日志记录延迟
- **并发安全**: 原生支持高并发场景

运行性能测试：
```bash
go test -bench=. -benchmem
```

## 许可证

MIT License

## 贡献

欢迎提交 Issue 和 Pull Request！

## 链接

- [Uber Zap](https://github.com/uber-go/zap)
- [OpenTelemetry Go](https://github.com/open-telemetry/opentelemetry-go)
