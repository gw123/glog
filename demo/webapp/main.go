package main

import (
	"context"
	"fmt"
	"time"

	"github.com/gw123/glog"
	"github.com/gw123/glog/common"
)

func main() {
	println("=== Advanced Demo: Web Application Simulation ===\n")

	// Configure logger
	setupLogger()

	// Simulate web server lifecycle
	simulateWebServer()

	println("\n=== Advanced Demo Complete ===")
}

func setupLogger() {
	err := glog.SetDefaultLoggerConfig(
		common.Options{},
		common.WithLevel(common.DebugLevel),
		common.WithConsoleEncoding(),
		common.WithOutputPath("./demo/web_app.log"),
		common.WithStdoutOutputPath(),
	)
	if err != nil {
		fmt.Printf("Failed to configure logger: %v\n", err)
	}
}

func simulateWebServer() {
	// 1. Application startup
	startupLogger := glog.Log().Named("startup")
	startupLogger.Info("Starting web application")
	startupLogger.WithField("version", "1.0.0").Info("Application version")
	startupLogger.WithField("environment", "development").Info("Environment")

	// 2. Database initialization
	initDatabase()

	// 3. Cache initialization
	initCache()

	// 4. API server start
	startAPIServer()

	// 5. Handle incoming requests
	println("\nSimulating incoming HTTP requests...")
	handleHTTPRequests()

	// 6. Graceful shutdown
	gracefulShutdown()
}

func initDatabase() {
	dbLogger := glog.Log().Named("database")
	dbLogger.Info("Initializing database connection")

	// Simulate connection
	time.Sleep(time.Millisecond * 50)

	dbLogger.WithField("host", "localhost").
		WithField("port", 5432).
		WithField("database", "myapp").
		Info("Database connected successfully")

	dbLogger.WithField("max_connections", 20).
		WithField("idle_connections", 5).
		Debug("Connection pool configured")
}

func initCache() {
	cacheLogger := glog.Log().Named("cache")
	cacheLogger.Info("Initializing Redis cache")

	// Simulate connection
	time.Sleep(time.Millisecond * 30)

	cacheLogger.WithField("host", "localhost").
		WithField("port", 6379).
		Info("Cache connected successfully")

	cacheLogger.WithField("max_retries", 3).
		WithField("timeout_ms", 5000).
		Debug("Cache configuration loaded")
}

func startAPIServer() {
	apiLogger := glog.Log().Named("api")
	apiLogger.Info("Starting API server")

	apiLogger.WithField("host", "0.0.0.0").
		WithField("port", 8080).
		Info("API server listening")

	apiLogger.WithField("endpoints", []string{
		"/api/users",
		"/api/orders",
		"/api/products",
	}).Debug("Registered API endpoints")
}

func handleHTTPRequests() {
	// Simulate multiple concurrent HTTP requests
	requests := []struct {
		method string
		path   string
		userID int64
	}{
		{"GET", "/api/users/123", 123},
		{"POST", "/api/orders", 456},
		{"GET", "/api/products", 789},
		{"PUT", "/api/users/123", 123},
		{"DELETE", "/api/orders/999", 456},
	}

	for i, req := range requests {
		// Each request gets its own context and trace ID
		requestID := fmt.Sprintf("req-%d-%d", time.Now().Unix(), i)
		handleRequest(requestID, req.method, req.path, req.userID)
		time.Sleep(time.Millisecond * 100) // Simulate time between requests
	}
}

func handleRequest(requestID, method, path string, userID int64) {
	// Create request context
	ctx := context.Background()
	ctx = glog.ToContext(ctx, glog.Log().Named("http"))

	// Add request metadata
	glog.AddTraceID(ctx, requestID)
	glog.AddUserID(ctx, userID)
	glog.AddPathname(ctx, path)
	glog.AddField(ctx, "method", method)
	glog.AddField(ctx, "client_ip", "192.168.1.100")

	logger := glog.ExtractEntry(ctx)

	// Request started
	startTime := time.Now()
	logger.Info("Request received")

	// Simulate request processing
	processRequest(ctx, method, path)

	// Request completed
	duration := time.Since(startTime)
	logger.WithField("duration_ms", duration.Milliseconds()).
		WithField("status", 200).
		Info("Request completed")
}

func processRequest(ctx context.Context, method, path string) {
	// Simulate authentication
	authenticateUser(ctx)

	// Simulate business logic based on path
	switch path {
	case "/api/users/123":
		getUserProfile(ctx)
	case "/api/orders":
		createOrder(ctx)
	case "/api/products":
		listProducts(ctx)
	default:
		glog.ExtractEntry(ctx).Debug("Processing generic request")
	}

	// Simulate database query
	executeQuery(ctx)

	// Simulate cache operation
	cacheOperation(ctx)
}

func authenticateUser(ctx context.Context) {
	logger := glog.ExtractEntry(ctx).Named("auth")
	logger.Debug("Authenticating user")

	time.Sleep(time.Millisecond * 10)

	userID := glog.ExtractUserID(ctx)
	logger.WithField("authenticated_user_id", userID).Debug("User authenticated")
}

func getUserProfile(ctx context.Context) {
	logger := glog.ExtractEntry(ctx).Named("service")
	logger.Info("Fetching user profile")

	glog.AddField(ctx, "operation", "get_user_profile")
	time.Sleep(time.Millisecond * 20)

	logger.WithField("profile_fields", []string{"name", "email", "avatar"}).
		Debug("User profile retrieved")
}

func createOrder(ctx context.Context) {
	logger := glog.ExtractEntry(ctx).Named("service")
	logger.Info("Creating new order")

	glog.AddField(ctx, "operation", "create_order")
	glog.AddField(ctx, "order_items", 3)
	glog.AddField(ctx, "total_amount", 299.99)

	time.Sleep(time.Millisecond * 30)

	logger.WithField("order_id", "ORD-12345").Info("Order created successfully")
}

func listProducts(ctx context.Context) {
	logger := glog.ExtractEntry(ctx).Named("service")
	logger.Info("Listing products")

	glog.AddField(ctx, "operation", "list_products")
	glog.AddField(ctx, "page", 1)
	glog.AddField(ctx, "page_size", 20)

	time.Sleep(time.Millisecond * 15)

	logger.WithField("total_products", 150).
		WithField("returned_products", 20).
		Debug("Products retrieved")
}

func executeQuery(ctx context.Context) {
	logger := glog.ExtractEntry(ctx).Named("database")
	logger.Debug("Executing database query")

	time.Sleep(time.Millisecond * 25)

	logger.WithField("query_type", "SELECT").
		WithField("rows_affected", 1).
		WithField("execution_time_ms", 25).
		Debug("Query executed successfully")
}

func cacheOperation(ctx context.Context) {
	logger := glog.ExtractEntry(ctx).Named("cache")
	logger.Debug("Cache operation")

	time.Sleep(time.Millisecond * 5)

	logger.WithField("cache_hit", true).
		WithField("key", "user:profile:123").
		Debug("Cache retrieved")
}

func gracefulShutdown() {
	shutdownLogger := glog.Log().Named("shutdown")
	shutdownLogger.Info("Starting graceful shutdown")

	shutdownLogger.Debug("Draining active connections")
	time.Sleep(time.Millisecond * 50)

	shutdownLogger.Debug("Closing database connections")
	time.Sleep(time.Millisecond * 30)

	shutdownLogger.Debug("Closing cache connections")
	time.Sleep(time.Millisecond * 20)

	shutdownLogger.Info("Application shutdown complete")
}
