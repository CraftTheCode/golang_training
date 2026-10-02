// 05_generics_intro demonstrates Go generics (type parameters), introduced in Go 1.18.
//
// Before generics, you needed separate functions for each type (MinInt, MinFloat, etc.)
// or used interface{}/any with runtime type assertions (unsafe and verbose).
// Generics let you write type-safe, reusable code with compile-time checks.
//
// Run:
//
//	go run ./03-idiomatic-go-design/05_generics_intro/
package main

import (
	"cmp"
	"fmt"
)

// ---------------------------------------------------------------------------
// 1. Generic Function: Min
// ---------------------------------------------------------------------------

// Min returns the smaller of two values.
// [T cmp.Ordered] means T can be any type that supports < comparison:
// all integers, floats, and strings.
//
// Without generics, you'd need: MinInt, MinFloat64, MinString, etc.
func Min[T cmp.Ordered](a, b T) T {
	if a < b {
		return a
	}
	return b
}

// ---------------------------------------------------------------------------
// 2. Generic Function: Contains
// ---------------------------------------------------------------------------

// Contains checks if a slice contains a target value.
// [T comparable] means T must support == (ints, strings, structs without slices/maps, etc.)
//
// This single function replaces: ContainsInt, ContainsString, ContainsBool, etc.
func Contains[T comparable](slice []T, target T) bool {
	for _, v := range slice {
		if v == target {
			return true
		}
	}
	return false
}

// ---------------------------------------------------------------------------
// 3. Generic Function: Filter
// ---------------------------------------------------------------------------

// Filter returns a new slice containing only elements that satisfy the predicate.
// The predicate is a function — combining generics with first-class functions.
func Filter[T any](slice []T, predicate func(T) bool) []T {
	result := make([]T, 0)
	for _, v := range slice {
		if predicate(v) {
			result = append(result, v)
		}
	}
	return result
}

// ---------------------------------------------------------------------------
// 4. Generic Type: Stack (generic data structure)
// ---------------------------------------------------------------------------

// Stack is a generic LIFO data structure.
// [T any] means it can hold any type — Stack[int], Stack[string], etc.
// The concrete type is specified when creating an instance.
type Stack[T any] struct {
	items []T
}

// Push adds an element to the top of the stack.
// The method receiver uses [T any] to match the type parameter of Stack.
func (s *Stack[T]) Push(item T) {
	s.items = append(s.items, item)
}

// Pop removes and returns the top element.
// Returns the value and a bool indicating success (false if stack is empty).
// This follows Go's comma-ok idiom.
func (s *Stack[T]) Pop() (T, bool) {
	if len(s.items) == 0 {
		var zero T // zero value of any type T
		return zero, false
	}
	top := s.items[len(s.items)-1]
	s.items = s.items[:len(s.items)-1]
	return top, true
}

// Size returns the number of elements in the stack.
func (s *Stack[T]) Size() int {
	return len(s.items)
}

// ---------------------------------------------------------------------------
// main
// ---------------------------------------------------------------------------

func main() {
	fmt.Println("=== 1. Generic Min Function ===")
	// The compiler infers T from the arguments — no need to write Min[int](3, 7)
	fmt.Printf("Min(3, 7) = %d\n", Min(3, 7))
	fmt.Printf("Min(3.14, 2.71) = %.2f\n", Min(3.14, 2.71))
	fmt.Printf("Min(\"banana\", \"apple\") = %s\n", Min("banana", "apple"))

	fmt.Println("\n=== 2. Generic Contains Function ===")
	names := []string{"Alice", "Bob", "Charlie"}
	fmt.Printf("Contains(%v, \"Bob\") = %t\n", names, Contains(names, "Bob"))
	fmt.Printf("Contains(%v, \"Eve\") = %t\n", names, Contains(names, "Eve"))

	numbers := []int{10, 20, 30, 40}
	fmt.Printf("Contains(%v, 30) = %t\n", numbers, Contains(numbers, 30))

	fmt.Println("\n=== 3. Generic Filter Function ===")
	scores := []int{45, 82, 90, 33, 76, 95, 60}
	// Filter with an inline predicate function
	passing := Filter(scores, func(s int) bool { return s >= 60 })
	fmt.Printf("Scores: %v\n", scores)
	fmt.Printf("Passing (>=60): %v\n", passing)

	words := []string{"Go", "is", "a", "systems", "language"}
	long := Filter(words, func(w string) bool { return len(w) > 2 })
	fmt.Printf("Words: %v\n", words)
	fmt.Printf("Long words (>2 chars): %v\n", long)

	fmt.Println("\n=== 4. Generic Stack Data Structure ===")
	// Stack[string] — a stack that holds strings
	var strStack Stack[string]
	strStack.Push("first")
	strStack.Push("second")
	strStack.Push("third")
	fmt.Printf("Stack size: %d\n", strStack.Size())

	if val, ok := strStack.Pop(); ok {
		fmt.Printf("Popped: %s\n", val)
	}
	fmt.Printf("Stack size after pop: %d\n", strStack.Size())

	// Stack[int] — same Stack type, different concrete type
	var intStack Stack[int]
	intStack.Push(100)
	intStack.Push(200)
	if val, ok := intStack.Pop(); ok {
		fmt.Printf("Popped from int stack: %d\n", val)
	}

	fmt.Println("\n=== Key Takeaways ===")
	fmt.Println("• [T constraint] declares a type parameter with a constraint")
	fmt.Println("• cmp.Ordered: types that support < > <= >= (numbers, strings)")
	fmt.Println("• comparable: types that support == != ")
	fmt.Println("• any: alias for interface{} — no restrictions")
	fmt.Println("• The compiler infers T from arguments — explicit types rarely needed")
}
