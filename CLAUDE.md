# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a Go logging library (`glog`) built as a wrapper around Uber's Zap logger. It provides:
- Structured logging with formatted output
- Context-based logging for tracking requests across function calls
- OpenTelemetry trace integration
- Named loggers for component separation

## Core Architecture

### Two-Layer Structure

1. **Top-level API** (`log.go`): Convenience functions that delegate to the zap implementation
   - `Infof()`, `Errorf()`, `Warnf()`, `Debugf()` - formatted logging
   - `Info()`, `Error()`, `Warn()`, `Debug()` - simple logging
   - `Log()` - returns default logger instance
   - `WithField()`, `WithError()` - chaining context
   - `SetDefaultLoggerConfig()` - configure the default logger

2. **Zap Implementation** (`zap/`): Direct Zap implementation
   - `Logger` wraps `zap.SugaredLogger`
   - Two singleton instances: `defaultLogger` (for external use) and `innerLogger` (with CallerSkip=1 for internal wrapper functions)
   - `NewLogger()` creates configured logger with custom encoders
   - Output format: `[timestamp] [level] [logger_name] file:line [] message {fields}`

3. **Context Logging** (`ctxlog.go`): Request-scoped logging via context.Context
   - `ToContext()` - inject logger into context
   - `AddField()` / `AddFields()` - accumulate fields throughout request lifecycle
   - `AddTopField()` - add priority fields (trace_id, user_id, pathname, client_ip)
   - `ExtractEntry()` - retrieve logger with accumulated fields and OTEL trace_id
   - `WithOTEL()` - extract OpenTelemetry trace context

### Key Patterns

**Context-based logging workflow:**
```go
// 1. Create root context with logger
ctx := ToContext(context.Background(), glog.DefaultLogger())
AddField(ctx, "app_name", "web")

// 2. In middleware: add request-specific fields
AddTraceID(ctx, requestID)

// 3. Throughout request: add more fields
AddField(ctx, "user_id", userID)

// 4. Log with all accumulated context
ExtractEntry(ctx).WithField("ip", "10.0.0.1").Info("request completed")
```

**Named loggers** for component separation:
```go
logger := glog.Log().Named("database")  // Output includes [database] prefix
```

**OpenTelemetry integration:**
- `ExtractEntry(ctx)` automatically includes `trace_id` if OTEL span is active
- Manual: `WithOTEL(ctx).Info("message")`

## Configuration

**Logger configuration** via `SetDefaultLoggerConfig()`:
```go
glog.SetDefaultLoggerConfig(common.Options{
    OutputPaths:      []string{common.PathStdout, "./app.log"},
    ErrorOutputPaths: []string{common.PathStderr},
    Encoding:         common.EncodeConsole,  // or common.EncodeJson
    Level:            common.DebugLevel,
})
```

**Options builders** (`common/const.go`):
- `WithStdoutOutputPath()`, `WithOutputPath(path)`
- `WithConsoleEncoding()`, `WithJsonEncoding()`
- `WithLevel(level)`, `WithCallerSkip(skip)`

## Testing

Run tests: `go test ./...`
Run specific test: `go test -run TestName`
Run benchmarks: `go test -bench=.`

Test files demonstrate usage patterns:
- `log_test.go`: basic logging, configuration, named loggers
- `ctxlog_test.go`: context-based logging patterns

## Key Constants

**Log levels** (`common/level.go`): DebugLevel < InfoLevel < WarnLevel < ErrorLevel < DPanicLevel < PanicLevel < FatalLevel

**Context keys** (`common/const.go`):
- `KeyTraceID` - request trace ID
- `KeyUserID` - user identifier
- `KeyPathname` - request path
- `KeyClientIP` - client IP address

## Module Information

- Module path: `github.com/gw123/glog`
- Go version: 1.17
- Main dependency: `go.uber.org/zap v1.24.0`
- OTEL support: `go.opentelemetry.io/otel/trace v1.1.0`
