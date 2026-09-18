package main

import (
	"fmt"
	"time"
)

// sendData demonstrates send-only channel direction chan<-
func sendData(ch chan<- int, count int) {
	for i := 1; i <= count; i++ {
		ch <- i
		time.Sleep(20 * time.Millisecond)
	}
	close(ch) // Only the sender should close a channel!
}

func main() {
	fmt.Println("=== 1. Unbuffered Channel (Synchronous Rendezvous) ===")
	unbuffered := make(chan string)

	go func() {
		fmt.Println("[Goroutine] Ready to send message...")
		unbuffered <- "Hello from concurrent goroutine!"
		fmt.Println("[Goroutine] Message received by receiver!")
	}()

	time.Sleep(50 * time.Millisecond)
	msg := <-unbuffered
	fmt.Printf("[Main] Received: %q\n", msg)

	fmt.Println("\n=== 2. Buffered Channel & for-range over Channel ===")
	buffered := make(chan int, 3) // Buffer size 3 (can hold 3 items without blocking)

	go sendData(buffered, 5)

	// for-range automatically terminates when the channel is closed!
	for val := range buffered {
		fmt.Printf("[Main] Received from stream: %d\n", val)
	}

	fmt.Println("\n=== 3. Multiplexing with Select and Timeouts ===")
	ch1 := make(chan string)
	ch2 := make(chan string)

	go func() {
		time.Sleep(100 * time.Millisecond)
		ch1 <- "Result from service 1"
	}()

	go func() {
		time.Sleep(200 * time.Millisecond)
		ch2 <- "Result from service 2"
	}()

	// Select waits until one of the cases is ready
	for i := 0; i < 2; i++ {
		select {
		case res1 := <-ch1:
			fmt.Printf("[Select] Got: %s\n", res1)
		case res2 := <-ch2:
			fmt.Printf("[Select] Got: %s\n", res2)
		case <-time.After(300 * time.Millisecond):
			fmt.Println("[Select] Timed out waiting for response!")
		}
	}

	fmt.Println("\n=== 4. Non-blocking Channel Operations via default ===")
	chNonBlock := make(chan int)
	select {
	case val := <-chNonBlock:
		fmt.Printf("Received: %d\n", val)
	default:
		fmt.Println("No message available immediately (non-blocking fallback executed).")
	}
}
