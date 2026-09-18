package main

import (
	"context"
	"fmt"
	"time"
)

type contextKey string

const requestIDKey contextKey = "requestID"

// simulateLongTask simulates an operation that respects context cancellation
func simulateLongTask(ctx context.Context, taskName string, duration time.Duration) error {
	reqID := ctx.Value(requestIDKey)
	fmt.Printf("[%s] Task started (ReqID: %v, Duration: %v)...\n", taskName, reqID, duration)

	select {
	case <-time.After(duration):
		fmt.Printf("[%s] Task finished successfully!\n", taskName)
		return nil
	case <-ctx.Done():
		// Context was cancelled or timed out!
		fmt.Printf("[%s] Task aborted: %v\n", taskName, ctx.Err())
		return ctx.Err()
	}
}

func main() {
	fmt.Println("=== 1. Context Cancellation (WithCancel) ===")
	ctx, cancel := context.WithCancel(context.Background())

	// Attach request-scoped telemetry/ID
	ctxWithVal := context.WithValue(ctx, requestIDKey, "req-xyz-987")

	go func() {
		// Cancel after 100ms
		time.Sleep(100 * time.Millisecond)
		fmt.Println("[Main] Triggering manual cancel()...")
		cancel()
	}()

	err := simulateLongTask(ctxWithVal, "DatabaseBackup", 500*time.Millisecond)
	fmt.Printf("Outcome: %v\n", err)

	fmt.Println("\n=== 2. Context Timeout (WithTimeout) ===")
	// Timeout automatically cancels when deadline expires
	timeoutCtx, timeoutCancel := context.WithTimeout(context.Background(), 150*time.Millisecond)
	defer timeoutCancel() // Always call cancel to prevent resource leak

	errTimeout := simulateLongTask(timeoutCtx, "API_Call", 400*time.Millisecond)
	fmt.Printf("Outcome: %v\n", errTimeout)

	fmt.Println("\n=== 3. Context Success within Deadline ===")
	quickCtx, quickCancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer quickCancel()

	errQuick := simulateLongTask(quickCtx, "FastCacheLookup", 50*time.Millisecond)
	fmt.Printf("Outcome: %v\n", errQuick)
}
