package main

import (
	"fmt"
	"sync"
	"time"
)

// -------------------------------------------------------------
// Pattern 1: Worker Pool
// -------------------------------------------------------------

type Job struct {
	ID    int
	Value int
}

type Result struct {
	Job    Job
	Output int
}

func worker(id int, jobs <-chan Job, results chan<- Result, wg *sync.WaitGroup) {
	defer wg.Done()
	for job := range jobs {
		// Simulate processing time
		time.Sleep(10 * time.Millisecond)
		results <- Result{
			Job:    job,
			Output: job.Value * job.Value,
		}
	}
}

// -------------------------------------------------------------
// Pattern 2: Pipeline (Generator -> Square -> Consumer)
// -------------------------------------------------------------

func generateNumbers(nums ...int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for _, n := range nums {
			out <- n
		}
	}()
	return out
}

func squareStage(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for n := range in {
			out <- n * n
		}
	}()
	return out
}

// -------------------------------------------------------------
// Pattern 3: Rate Limiting
// -------------------------------------------------------------

func demonstrateRateLimiting() {
	requests := make(chan int, 5)
	for i := 1; i <= 5; i++ {
		requests <- i
	}
	close(requests)

	// Rate limiter: 1 event every 50ms
	limiter := time.Tick(50 * time.Millisecond)

	for req := range requests {
		<-limiter // Block until next tick
		fmt.Printf("[RateLimiter] Processed request #%d at %s\n",
			req, time.Now().Format("15:04:05.000"))
	}
}

func main() {
	fmt.Println("=== 1. Worker Pool Pattern ===")
	numJobs := 6
	numWorkers := 3

	jobs := make(chan Job, numJobs)
	results := make(chan Result, numJobs)
	var wg sync.WaitGroup

	// Spin up fixed pool of 3 workers
	for w := 1; w <= numWorkers; w++ {
		wg.Add(1)
		go worker(w, jobs, results, &wg)
	}

	// Submit jobs
	for j := 1; j <= numJobs; j++ {
		jobs <- Job{ID: j, Value: j * 10}
	}
	close(jobs) // Signal workers that no more jobs will be sent

	// Wait for workers in background then close results channel
	go func() {
		wg.Wait()
		close(results)
	}()

	for res := range results {
		fmt.Printf("Job #%d: input=%d -> output=%d\n", res.Job.ID, res.Job.Value, res.Output)
	}

	fmt.Println("\n=== 2. Pipeline Pattern ===")
	// Numbers flow through pipeline stages cleanly
	source := generateNumbers(2, 3, 4, 5)
	squared := squareStage(source)

	for val := range squared {
		fmt.Printf("Pipeline output: %d\n", val)
	}

	fmt.Println("\n=== 3. Rate Limiting Pattern ===")
	demonstrateRateLimiting()
}
