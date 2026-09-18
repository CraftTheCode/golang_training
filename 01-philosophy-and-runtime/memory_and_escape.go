package main

import (
	"bytes"
	"fmt"
)

// User demonstrates zero-value initialization for composite structs.
type User struct {
	ID        int
	Username  string
	IsActive  bool
	Roles     []string
	Profile   *Profile
}

type Profile struct {
	Bio string
}

// stackAllocated returns a value.
// The compiler keeps this on the stack frame of stackAllocated,
// copying the value back to the caller.
func stackAllocated() int {
	val := 100
	return val
}

// heapAllocated returns a pointer to a local variable.
// Escape analysis detects that '&val' escapes the function scope,
// so 'val' is allocated on the HEAP.
func heapAllocated() *int {
	val := 200
	return &val // Escapes to heap!
}

func main() {
	fmt.Println("=== 1. Zero Value Demonstration ===")
	var u User
	fmt.Printf("Default ID (int): %d\n", u.ID)
	fmt.Printf("Default Username (string): %q\n", u.Username)
	fmt.Printf("Default IsActive (bool): %t\n", u.IsActive)
	fmt.Printf("Default Roles (slice): %v (nil: %t)\n", u.Roles, u.Roles == nil)
	fmt.Printf("Default Profile (pointer): %v\n", u.Profile)

	fmt.Println("\n=== 2. The Zero-Value-Is-Ready-To-Use Idiom ===")
	// bytes.Buffer requires NO NewBuffer() constructor.
	// Its zero value is completely initialized and usable.
	var buf bytes.Buffer
	buf.WriteString("Go ")
	buf.WriteString("philosophy: ")
	buf.WriteString("zero values are safe and useful.")
	fmt.Println("Buffer output:", buf.String())

	fmt.Println("\n=== 3. Stack vs Heap & Escape Analysis ===")
	sVal := stackAllocated()
	hPtr := heapAllocated()

	fmt.Printf("Stack allocated value: %d\n", sVal)
	fmt.Printf("Heap allocated pointer value: %d at address: %p\n", *hPtr, hPtr)

	fmt.Println("\nTo inspect escape analysis decisions yourself, run:")
	fmt.Println("go build -gcflags=\"-m\" ./01-philosophy-and-runtime/memory_and_escape.go")
}
