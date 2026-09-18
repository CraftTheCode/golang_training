package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("=== 1. If with Short Statement ===")
	// Variable declared in if header is scoped strictly to the if-else block
	if length := len("Hello, Go!"); length > 5 {
		fmt.Printf("String is long (%d chars)\n", length)
	} else {
		fmt.Printf("String is short (%d chars)\n", length)
	}
	// 'length' is not accessible here

	fmt.Println("\n=== 2. For Loops (The Only Loop in Go) ===")
	// Standard 3-component loop
	for i := 1; i <= 3; i++ {
		fmt.Printf("Standard loop count: %d\n", i)
	}

	// While-style loop
	count := 3
	for count > 0 {
		fmt.Printf("While-style countdown: %d\n", count)
		count--
	}

	// For-range over slice
	fruits := []string{"Apple", "Banana", "Cherry"}
	for index, fruit := range fruits {
		fmt.Printf("Index %d: %s\n", index, fruit)
	}

	// For-range omitting index
	for _, fruit := range fruits {
		fmt.Printf("Fruit: %s\n", fruit)
	}

	fmt.Println("\n=== 3. Switch Statements ===")
	// In Go, switch does NOT fall through by default
	day := time.Wednesday
	switch day {
	case time.Saturday, time.Sunday:
		fmt.Println("It's the weekend!")
	case time.Monday:
		fmt.Println("Start of the work week.")
	case time.Wednesday:
		fmt.Println("Midweek Wednesday.")
	default:
		fmt.Println("Regular weekday.")
	}

	// Tagless switch (idiomatic replacement for long if/else if chains)
	score := 85
	switch {
	case score >= 90:
		fmt.Println("Grade: A")
	case score >= 80:
		fmt.Println("Grade: B")
	case score >= 70:
		fmt.Println("Grade: C")
	default:
		fmt.Println("Grade: F")
	}

	fmt.Println("\n=== 4. Defer Mechanics (LIFO & Argument Timing) ===")
	demonstrateDefer()
}

func demonstrateDefer() {
	// Deferred functions are stacked and executed in LIFO (Last-In, First-Out) order
	// when the enclosing function returns.
	fmt.Println("Starting demonstrateDefer...")

	defer fmt.Println("Deferred execution 1 (First deferred -> executed LAST)")
	defer fmt.Println("Deferred execution 2 (Second deferred -> executed SECOND)")
	defer fmt.Println("Deferred execution 3 (Last deferred -> executed FIRST)")

	// Arguments are evaluated IMMEDIATELY at defer definition time, NOT at execution time
	msg := "original message"
	defer func(captured string) {
		fmt.Printf("Deferred closure captured: %q, outer msg is now: %q\n", captured, msg)
	}(msg)

	msg = "mutated message"
	fmt.Println("Exiting demonstrateDefer function body...")
}
