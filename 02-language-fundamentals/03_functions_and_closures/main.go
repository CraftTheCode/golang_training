package main

import (
	"errors"
	"fmt"
)

// divide demonstrates idiomatic multiple return values (result, error).
func divide(a, b float64) (float64, error) {
	if b == 0 {
		return 0, errors.New("cannot divide by zero")
	}
	return a / b, nil
}

// calculateStats demonstrates named return values.
// Named return values document the intent of returns in godoc.
func calculateStats(values []int) (min int, max int, sum int) {
	if len(values) == 0 {
		return 0, 0, 0
	}
	min, max = values[0], values[0]
	for _, v := range values {
		if v < min {
			min = v
		}
		if v > max {
			max = v
		}
		sum += v
	}
	// Explicit return is preferred over naked return for readability
	return min, max, sum
}

// sumAll demonstrates variadic functions.
func sumAll(numbers ...int) int {
	total := 0
	for _, n := range numbers {
		total += n
	}
	return total
}

// makeCounter demonstrates a closure capturing lexical scope state.
func makeCounter(initial int) func() int {
	count := initial
	return func() int {
		count++
		return count
	}
}

func main() {
	fmt.Println("=== 1. Multiple Return Values ===")
	res, err := divide(10, 2)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Printf("10 / 2 = %.2f\n", res)
	}

	_, errZero := divide(10, 0)
	if errZero != nil {
		fmt.Println("Handled expected error:", errZero)
	}

	fmt.Println("\n=== 2. Named Returns ===")
	nums := []int{15, 3, 42, 8, 23}
	min, max, sum := calculateStats(nums)
	fmt.Printf("Nums: %v -> Min: %d, Max: %d, Sum: %d\n", nums, min, max, sum)

	fmt.Println("\n=== 3. Variadic Functions ===")
	fmt.Printf("sumAll(1, 2, 3): %d\n", sumAll(1, 2, 3))
	fmt.Printf("sumAll(nums...): %d\n", sumAll(nums...))

	fmt.Println("\n=== 4. Closures and State Encapsulation ===")
	counterA := makeCounter(0)
	counterB := makeCounter(100)

	fmt.Println("Counter A:", counterA()) // 1
	fmt.Println("Counter A:", counterA()) // 2
	fmt.Println("Counter B:", counterB()) // 101
	fmt.Println("Counter A:", counterA()) // 3
}
