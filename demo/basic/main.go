package main

import (
	"context"
	"errors"
	"time"

	"github.com/gw123/glog"
	"github.com/gw123/glog/common"
)

func main() {
	println("=== glog Demo Examples ===\n")

	// 1. Basic logging
	basicLogging()

	// 2. Logger with fields
	loggerWithFields()

	// 3. Context-based logging
	contextLogging()

	// 4. Error logging
	errorLogging()

	// 5. Named loggers
	namedLoggers()

	// 6. Custom configuration
	customConfiguration()

	// 7. Concurrent logging
	concurrentLogging()

	// 8. Different log levels
	logLevels()

	println("\n=== Demo Complete ===")
}

// 1. Basic logging examples
func basicLogging() {
	println("\n--- 1. Basic Logging ---")

	glog.Info("Simple info message")
	glog.Infof("Formatted info: %s = %d", "count", 42)
	glog.Warn("Warning message")
	glog.Warnf("Formatted warning: %s", "something went wrong")
	glog.Error("Error message")
	glog.Errorf("Formatted error: %v", errors.New("sample error"))
}

// 2. Logger with fields
func loggerWithFields() {
	println("\n--- 2. Logger with Fields ---")

	// Single field
	glog.WithField("user_id", 12345).Info("User logged in")

	// Multiple fields chained
	glog.Log().
		WithField("request_id", "req-001").
		WithField("method", "GET").
		WithField("path", "/api/users").
		WithField("status", 200).
		Info("HTTP request processed")

	// With error
	err := errors.New("database connection failed")
	glog.WithError(err).Error("Failed to connect to database")
}

// 3. Context-based logging
func contextLogging() {
	println("\n--- 3. Context-based Logging ---")

	// Create context with logger
	ctx := context.Background()
	ctx = glog.ToContext(ctx, glog.DefaultLogger())

	// Add trace ID (simulating request tracking)
	glog.AddTraceID(ctx, "trace-abc-123")
	glog.AddUserID(ctx, 999)
	glog.AddPathname(ctx, "/api/orders")

	// Add custom fields
	glog.AddField(ctx, "client_ip", "192.168.1.100")
	glog.AddField(ctx, "user_agent", "Mozilla/5.0")

	// Extract logger with all context
	logger := glog.ExtractEntry(ctx)
	logger.Info("Request started")

	// Simulate processing
	processOrder(ctx)

	logger.Info("Request completed")
}

func processOrder(ctx context.Context) {
	glog.AddField(ctx, "order_id", "ORD-12345")
	glog.ExtractEntry(ctx).Info("Processing order")

	// Simulate order validation
	validateOrder(ctx)

	// Simulate payment
	processPayment(ctx)
}

func validateOrder(ctx context.Context) {
	glog.AddField(ctx, "validation_step", "inventory_check")
	glog.ExtractEntry(ctx).Info("Validating order")
}

func processPayment(ctx context.Context) {
	glog.AddField(ctx, "payment_method", "credit_card")
	glog.AddField(ctx, "amount", 99.99)
	glog.ExtractEntry(ctx).Info("Processing payment")
}

// 4. Error logging examples
func errorLogging() {
	println("\n--- 4. Error Logging ---")

	// Simple error
	err1 := errors.New("file not found")
	glog.WithError(err1).Error("Failed to read configuration")

	// Error with context
	ctx := context.Background()
	ctx = glog.ToContext(ctx, glog.DefaultLogger())
	glog.AddTraceID(ctx, "err-trace-001")

	err2 := errors.New("connection timeout")
	glog.ExtractEntry(ctx).
		WithError(err2).
		WithField("timeout_seconds", 30).
		Error("Database connection failed")
}

// 5. Named loggers for different components
func namedLoggers() {
	println("\n--- 5. Named Loggers ---")

	// Database logger
	dbLogger := glog.Log().Named("database")
	dbLogger.Info("Database connection established")
	dbLogger.WithField("pool_size", 10).Info("Connection pool initialized")

	// Cache logger
	cacheLogger := glog.Log().Named("cache")
	cacheLogger.Info("Redis connection established")
	cacheLogger.WithField("ttl_seconds", 3600).Info("Cache configuration loaded")

	// API logger
	apiLogger := glog.Log().Named("api")
	apiLogger.Info("API server started")
	apiLogger.WithField("port", 8080).Info("Listening on port")
}

// 6. Custom configuration
func customConfiguration() {
	println("\n--- 6. Custom Configuration ---")

	// Save original config
	println("Configuring logger with custom settings...")

	// Configure with debug level and file output
	err := glog.SetDefaultLoggerConfig(
		common.Options{},
		common.WithLevel(common.DebugLevel),
		common.WithConsoleEncoding(),
		common.WithOutputPath("./demo/demo.log"),
		common.WithStdoutOutputPath(),
	)

	if err != nil {
		println("Error configuring logger:", err.Error())
		return
	}

	// Now debug logs will appear
	glog.Debug("This is a debug message")
	glog.Debugf("Debug with format: %s", "now visible")
	glog.Info("Info message with custom config")

	println("Logs are being written to demo/demo.log")
}

// 7. Concurrent logging simulation
func concurrentLogging() {
	println("\n--- 7. Concurrent Logging ---")

	done := make(chan bool)

	// Simulate multiple goroutines logging
	for i := 0; i < 5; i++ {
		go func(workerID int) {
			ctx := context.Background()
			ctx = glog.ToContext(ctx, glog.DefaultLogger())

			glog.AddTraceID(ctx, "worker-trace")
			glog.AddField(ctx, "worker_id", workerID)

			logger := glog.ExtractEntry(ctx)

			logger.Infof("Worker %d started", workerID)
			time.Sleep(time.Millisecond * 10) // Simulate work
			logger.Infof("Worker %d processing", workerID)
			time.Sleep(time.Millisecond * 10) // Simulate more work
			logger.Infof("Worker %d completed", workerID)

			done <- true
		}(i)
	}

	// Wait for all workers
	for i := 0; i < 5; i++ {
		<-done
	}

	println("All workers completed")
}

// 8. Different log levels demonstration
func logLevels() {
	println("\n--- 8. Log Levels ---")

	// Ensure debug level is enabled
	glog.SetDefaultLoggerConfig(
		common.Options{},
		common.WithLevel(common.DebugLevel),
		common.WithConsoleEncoding(),
	)

	println("Demonstrating all log levels:")

	glog.Debug("DEBUG: Detailed information for debugging")
	glog.Info("INFO: General information about application flow")
	glog.Warn("WARN: Warning about potential issues")
	glog.Error("ERROR: Error that needs attention")

	// With different configurations
	println("\nChanging to INFO level (debug won't show):")
	glog.SetDefaultLoggerConfig(
		common.Options{},
		common.WithLevel(common.InfoLevel),
		common.WithConsoleEncoding(),
	)

	glog.Debug("This debug message won't appear")
	glog.Info("But this info message will")
}
