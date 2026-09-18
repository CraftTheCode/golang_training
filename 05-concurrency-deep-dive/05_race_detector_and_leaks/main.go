package main

import (
	"context"
	"fmt"
	"runtime"
	"sync"
	"time"
)

// DemonstrateLeak shows how an abandoned unbuffered channel causes a permanent goroutine leak
func demonstrateGoroutineLeak() {
	fmt.Printf("Active goroutines BEFORE potential leak: %d\n", runtime.NumGoroutine())

	// Leaky function: spawned worker blocks forever writing to unbuffered channel
	// if the receiver abandons after a timeout
	leakyWorker := func() <-chan string {
		ch := make(chan string) // Unbuffered!
		go func() {
			time.Sleep(100 * time.Millisecond)
			ch <- "result" // If caller abandoned reading, this blocks FOREVER -> LEAK!
		}()
		return ch
	}

	// Caller abandons after 10ms
	ch := leakyWorker()
	select {
	case res := <-ch:
		fmt.Println("Received:", res)
	case <-time.After(10 * time.Millisecond):
		fmt.Println("Timed out! Worker goroutine is now orphaned and blocked.")
	}

	time.Sleep(150 * time.Millisecond)
	fmt.Printf("Active goroutines AFTER leak: %d (Leaked goroutine remains alive in memory!)\n", runtime.NumGoroutine())
}

// DemonstrateLeakPrevention shows how to fix leaks with a buffered channel or context
func demonstrateLeakPrevention() {
	fmt.Printf("\nActive goroutines BEFORE leak fix: %d\n", runtime.NumGoroutine())

	safeWorker := func(ctx context.Context) <-chan string {
		ch := make(chan string, 1) // Buffer capacity 1 allows goroutine to write and exit without blocking!
		go func() {
			time.Sleep(50 * time.Millisecond)
			select {
			case ch <- "safe-result":
			case <-ctx.Done():
				// Clean exit on context cancellation
			}
		}()
		return ch
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	ch := safeWorker(ctx)
	select {
	case res := <-ch:
		fmt.Println("Received:", res)
	case <-ctx.Done():
		fmt.Println("Timed out cleanly! Worker will not leak.")
	}

	time.Sleep(100 * time.Millisecond)
	fmt.Printf("Active goroutines AFTER leak fix: %d (No leaks!)\n", runtime.NumGoroutine())
}

// ThreadSafeMap demonstrates preventing data races
type SafeStats struct {
	mu     sync.Mutex
	counts map[string]int
}

func (s *SafeStats) Add(key string, delta int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.counts[key] += delta
}

func main() {
	fmt.Println("=== 1. Data Race Prevention ===")
	stats := SafeStats{counts: make(map[string]int)}
	var wg sync.WaitGroup

	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			stats.Add("page_views", 1)
		}()
	}
	wg.Wait()
	fmt.Println("Safely computed page_views:", stats.counts["page_views"])

	fmt.Println("\n=== 2. Goroutine Leaks and How to Avoid Them ===")
	demonstrateGoroutineLeak()
	demonstrateLeakPrevention()

	fmt.Println("\nTip: Always test concurrent code with the race detector:")
	fmt.Println("  go test -race ./...")
}
