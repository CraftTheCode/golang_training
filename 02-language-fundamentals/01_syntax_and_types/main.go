package main

import (
	"fmt"
	"math"
)

func main() {
	fmt.Println("=== Variable Declarations ===")

	// 1. Explicit declaration with type
	var age int = 30

	// 2. Type-inferred declaration (package or function level)
	var name = "Gopher"

	// 3. Short variable declaration syntax := (function-scoped only)
	score := 99.5

	// 4. Multiple variable declaration
	var (
		city    string = "Bengaluru"
		country string = "India"
		active  bool   = true
	)

	fmt.Printf("User: %s, Age: %d, Score: %.1f, Location: %s, %s, Active: %t\n",
		name, age, score, city, country, active)

	fmt.Println("\n=== Numeric Types and Zero Values ===")
	var defaultInt int
	var defaultFloat float64
	var defaultBool bool
	var defaultString string

	fmt.Printf("int zero value: %d\n", defaultInt)
	fmt.Printf("float64 zero value: %f\n", defaultFloat)
	fmt.Printf("bool zero value: %t\n", defaultBool)
	fmt.Printf("string zero value: %q\n", defaultString)

	fmt.Println("\n=== Explicit Type Conversions ===")
	// Go never coerces types implicitly. Conversions must be explicit.
	var integerVal int = 42
	var floatVal float64 = float64(integerVal)
	var unsignedVal uint = uint(integerVal)

	fmt.Printf("Converted: int(%d) -> float64(%.2f) -> uint(%d)\n",
		integerVal, floatVal, unsignedVal)

	// Math functions require float64
	hypot := math.Sqrt(float64(3*3 + 4*4))
	fmt.Printf("Hypotenuse of 3 and 4: %.2f\n", hypot)

	fmt.Println("\n=== Constants and Iota ===")
	// Constants are evaluated at compile-time
	const Pi = 3.14159
	const (
		StatusPending = iota // 0
		StatusActive         // 1
		StatusSuspended      // 2
		StatusTerminated     // 3
	)

	fmt.Printf("Pi: %f\n", Pi)
	fmt.Printf("Statuses: Pending=%d, Active=%d, Suspended=%d, Terminated=%d\n",
		StatusPending, StatusActive, StatusSuspended, StatusTerminated)
}
