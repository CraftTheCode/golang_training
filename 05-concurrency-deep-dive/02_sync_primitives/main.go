package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

// SafeCounter protects an integer counter using sync.Mutex
type SafeCounter struct {
	mu    sync.Mutex
	count int
}

func (c *SafeCounter) Inc() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.count++
}

func (c *SafeCounter) Value() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.count
}

// ThreadSafeCache uses sync.RWMutex to allow multiple readers simultaneously
type ThreadSafeCache struct {
	rwMu  sync.RWMutex
	items map[string]string
}

func (c *ThreadSafeCache) Get(key string) (string, bool) {
	c.rwMu.RLock() // Multiple goroutines can hold RLock concurrently!
	defer c.rwMu.RUnlock()
	val, ok := c.items[key]
	return val, ok
}

func (c *ThreadSafeCache) Set(key, value string) {
	c.rwMu.Lock() // Exclusive lock (blocks all readers and writers)
	defer c.rwMu.Unlock()
	c.items[key] = value
}

func main() {
	fmt.Println("=== 1. sync.WaitGroup & sync.Mutex ===")
	var wg sync.WaitGroup
	counter := SafeCounter{}

	numWorkers := 10
	iterationsPerWorker := 100

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < iterationsPerWorker; j++ {
				counter.Inc()
			}
		}()
	}

	wg.Wait()
	fmt.Printf("Expected count: %d, Actual count: %d\n", numWorkers*iterationsPerWorker, counter.Value())

	fmt.Println("\n=== 2. Lock-free Atomic Counters (sync/atomic) ===")
	var atomicCounter atomic.Int64

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < iterationsPerWorker; j++ {
				atomicCounter.Add(1) // Lock-free hardware atomic instruction
			}
		}()
	}

	wg.Wait()
	fmt.Printf("Atomic Counter Value: %d\n", atomicCounter.Load())

	fmt.Println("\n=== 3. Thread-Safe Lazy Initialization (sync.Once) ===")
	var once sync.Once
	initialize := func() {
		fmt.Println("[Init] Expensive initialization executed exactly ONCE!")
	}

	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			once.Do(initialize) // Executed only by the first goroutine; others block until done
		}(i)
	}

	wg.Wait()

	fmt.Println("\n=== 4. sync.RWMutex Cache Demonstration ===")
	cache := ThreadSafeCache{items: make(map[string]string)}
	cache.Set("cluster-status", "HEALTHY")

	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func(readerID int) {
			defer wg.Done()
			val, _ := cache.Get("cluster-status")
			fmt.Printf("Reader #%d read: %s\n", readerID, val)
		}(i)
	}

	wg.Wait()
	time.Sleep(10 * time.Millisecond)
}
