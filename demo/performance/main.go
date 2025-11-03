package main

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/gw123/glog"
	"github.com/gw123/glog/common"
)

func main() {
	println("=== Performance & Stress Test Demo ===\n")

	// Configure logger for performance testing
	setupHighPerformanceLogger()

	// 1. High-frequency logging
	highFrequencyLogging()

	// 2. Concurrent logging stress test
	concurrentStressTest()

	// 3. Large field logging
	largeFieldLogging()

	// 4. Context accumulation test
	contextAccumulationTest()

	println("\n=== Performance Demo Complete ===")
}

func setupHighPerformanceLogger() {
	err := glog.SetDefaultLoggerConfig(
		common.Options{},
		common.WithLevel(common.InfoLevel),
		common.WithConsoleEncoding(),
		common.WithOutputPath("./demo/performance.log"),
	)
	if err != nil {
		fmt.Printf("Failed to configure logger: %v\n", err)
	}
}

func highFrequencyLogging() {
	println("\n--- 1. High-Frequency Logging (10,000 messages) ---")

	start := time.Now()
	count := 10000

	for i := 0; i < count; i++ {
		glog.Infof("High-frequency log message #%d", i)
	}

	duration := time.Since(start)
	fmt.Printf("Logged %d messages in %v (%.2f msgs/sec)\n",
		count, duration, float64(count)/duration.Seconds())
}

func concurrentStressTest() {
	println("\n--- 2. Concurrent Logging Stress Test (100 goroutines) ---")

	workers := 100
	messagesPerWorker := 100
	done := make(chan bool, workers)

	start := time.Now()

	for i := 0; i < workers; i++ {
		go func(workerID int) {
			ctx := context.Background()
			ctx = glog.ToContext(ctx, glog.DefaultLogger())

			glog.AddTraceID(ctx, fmt.Sprintf("worker-%d", workerID))
			glog.AddField(ctx, "worker_id", workerID)

			logger := glog.ExtractEntry(ctx)

			for j := 0; j < messagesPerWorker; j++ {
				logger.Infof("Message %d from worker %d", j, workerID)
			}

			done <- true
		}(i)
	}

	// Wait for all workers
	for i := 0; i < workers; i++ {
		<-done
	}

	duration := time.Since(start)
	totalMessages := workers * messagesPerWorker
	fmt.Printf("Logged %d messages concurrently in %v (%.2f msgs/sec)\n",
		totalMessages, duration, float64(totalMessages)/duration.Seconds())
}

func largeFieldLogging() {
	println("\n--- 3. Large Field Logging ---")

	// Create large data structures
	largeArray := make([]int, 100)
	for i := range largeArray {
		largeArray[i] = rand.Intn(1000)
	}

	largeMap := make(map[string]interface{})
	for i := 0; i < 50; i++ {
		largeMap[fmt.Sprintf("key_%d", i)] = fmt.Sprintf("value_%d", i)
	}

	largeString := string(make([]byte, 1000))

	start := time.Now()

	glog.Log().
		WithField("large_array", largeArray).
		WithField("large_map", largeMap).
		WithField("large_string_len", len(largeString)).
		Info("Logging with large fields")

	duration := time.Since(start)
	fmt.Printf("Logged large fields in %v\n", duration)
}

func contextAccumulationTest() {
	println("\n--- 4. Context Field Accumulation Test ---")

	ctx := context.Background()
	ctx = glog.ToContext(ctx, glog.DefaultLogger())

	start := time.Now()

	// Accumulate many fields
	for i := 0; i < 100; i++ {
		glog.AddField(ctx, fmt.Sprintf("field_%d", i), i)
	}

	// Log with accumulated context
	glog.ExtractEntry(ctx).Info("Message with 100 accumulated fields")

	duration := time.Since(start)
	fmt.Printf("Accumulated 100 fields and logged in %v\n", duration)
}

// Additional performance scenarios
func performanceTestSuite() {
	println("\n=== Additional Performance Tests ===")

	// Test 1: WithField chaining performance
	chainedFieldsTest()

	// Test 2: Named logger performance
	namedLoggerTest()

	// Test 3: Error logging performance
	errorLoggingTest()
}

func chainedFieldsTest() {
	println("\n--- Chained Fields Performance ---")

	count := 1000
	start := time.Now()

	for i := 0; i < count; i++ {
		glog.Log().
			WithField("field1", i).
			WithField("field2", i*2).
			WithField("field3", i*3).
			WithField("field4", i*4).
			WithField("field5", i*5).
			Info("Chained fields message")
	}

	duration := time.Since(start)
	fmt.Printf("Logged %d messages with chained fields in %v\n", count, duration)
}

func namedLoggerTest() {
	println("\n--- Named Logger Performance ---")

	logger := glog.Log().Named("performance").Named("test").Named("nested")

	count := 1000
	start := time.Now()

	for i := 0; i < count; i++ {
		logger.Infof("Named logger message #%d", i)
	}

	duration := time.Since(start)
	fmt.Printf("Logged %d messages with nested named logger in %v\n", count, duration)
}

func errorLoggingTest() {
	println("\n--- Error Logging Performance ---")

	testErr := errors.New("test error for performance")

	count := 1000
	start := time.Now()

	for i := 0; i < count; i++ {
		glog.WithError(testErr).Errorf("Error message #%d", i)
	}

	duration := time.Since(start)
	fmt.Printf("Logged %d error messages in %v\n", count, duration)
}
