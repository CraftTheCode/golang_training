package main

import (
	"fmt"
)

// Person demonstrates a struct with fields and tags
type Person struct {
	Name string
	Age  int
	Tags []string
}

func main() {
	fmt.Println("=== 1. Arrays vs Slices ===")
	// Array: Fixed length, part of type definition [3]int != [4]int
	var fixedArr [3]int = [3]int{10, 20, 30}
	arrCopy := fixedArr // Entire array is copied by value!
	arrCopy[0] = 999
	fmt.Printf("Original Array: %v, Copied Array: %v\n", fixedArr, arrCopy)

	// Slice: Dynamic window onto an underlying array
	// make([]T, length, capacity)
	slice := make([]int, 2, 4)
	slice[0] = 100
	slice[1] = 200
	fmt.Printf("Slice: %v | len: %d | cap: %d\n", slice, len(slice), cap(slice))

	// Appending within capacity (no reallocation)
	slice = append(slice, 300)
	fmt.Printf("After 1 append: %v | len: %d | cap: %d\n", slice, len(slice), cap(slice))

	// Appending beyond capacity (triggers reallocation and capacity doubling)
	slice = append(slice, 400, 500)
	fmt.Printf("After appending beyond cap: %v | len: %d | cap: %d\n", slice, len(slice), cap(slice))

	fmt.Println("\n=== 2. Sub-slicing and Shared Underlying Array ===")
	original := []string{"A", "B", "C", "D", "E"}
	subSlice := original[1:3] // Views "B", "C"
	fmt.Println("Original:", original)
	fmt.Println("SubSlice:", subSlice)

	// Modifying subSlice modifies original because they share the underlying array!
	subSlice[0] = "MUTATED_B"
	fmt.Println("After mutating subSlice[0]:")
	fmt.Println("Original:", original)

	fmt.Println("\n=== 3. Maps and the Comma-Ok Idiom ===")
	// Maps must be initialized before writing, or will panic on nil map assignment
	scores := make(map[string]int)
	scores["Alice"] = 95
	scores["Bob"] = 80

	// Safe retrieval: Comma-ok idiom distinguishes between zero-value and non-existent key
	score, ok := scores["Charlie"]
	if !ok {
		fmt.Println("Key 'Charlie' does not exist in map! (Default zero value would be:", score, ")")
	}

	scoreAlice, ok := scores["Alice"]
	if ok {
		fmt.Printf("Alice's score: %d\n", scoreAlice)
	}

	// Deleting a key
	delete(scores, "Bob")
	fmt.Println("Scores after deleting Bob:", scores)

	fmt.Println("\n=== 4. Structs ===")
	p := Person{
		Name: "Gopher",
		Age:  15,
		Tags: []string{"developer", "concurrent", "compiled"},
	}
	fmt.Printf("Person: %+v\n", p)
}
