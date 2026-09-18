package main

import (
	"bytes"
	"fmt"
	"io"
	"strings"
)

// Greeter is a small, single-method interface.
type Greeter interface {
	Greet() string
}

type EnglishSpeaker struct{ Name string }

// Implicit implementation: EnglishSpeaker implements Greeter without any "implements" keyword
func (e EnglishSpeaker) Greet() string {
	return fmt.Sprintf("Hello, my name is %s!", e.Name)
}

type SpanishSpeaker struct{ Name string }

func (s SpanishSpeaker) Greet() string {
	return fmt.Sprintf("¡Hola, me llamo %s!", s.Name)
}

// Broadcaster accepts any Greeter interface (Dependency Inversion)
func BroadcastGreeting(g Greeter) {
	fmt.Println("[Broadcast]", g.Greet())
}

// InspectType demonstrates Type Switches
func InspectType(val any) {
	switch v := val.(type) {
	case string:
		fmt.Printf("It's a string of length %d: %q\n", len(v), v)
	case int:
		fmt.Printf("It's an int multiplied by 2: %d\n", v*2)
	case Greeter:
		fmt.Printf("It's a Greeter! Greet output: %s\n", v.Greet())
	default:
		fmt.Printf("Unknown type: %T\n", v)
	}
}

func main() {
	fmt.Println("=== 1. Implicit Interfaces ===")
	e := EnglishSpeaker{Name: "Alice"}
	s := SpanishSpeaker{Name: "Carlos"}

	BroadcastGreeting(e)
	BroadcastGreeting(s)

	fmt.Println("\n=== 2. Standard Library io.Reader / io.Writer ===")
	// Standard interfaces are universal. Any reader can be copied to any writer.
	src := strings.NewReader("Streaming data from string reader...")
	var dst bytes.Buffer

	bytesWritten, _ := io.Copy(&dst, src)
	fmt.Printf("Copied %d bytes to buffer: %s\n", bytesWritten, dst.String())

	fmt.Println("\n=== 3. Type Assertions and Type Switches ===")
	InspectType("Go is awesome")
	InspectType(42)
	InspectType(e)

	// Direct Type Assertion
	var anyVal any = "Sample text"
	strVal, ok := anyVal.(string)
	if ok {
		fmt.Printf("Successful type assertion: %s\n", strVal)
	}

	fmt.Println("\n=== 4. The Nil Interface Trap ===")
	// An interface is ONLY nil when both its type AND value are nil.
	var customErr *CustomError = nil
	var iErr error = customErr

	// iErr is NOT nil because it holds type metadata (*CustomError)!
	fmt.Printf("Is iErr == nil? %t (Gotcha: type is %T, value is %v)\n", iErr == nil, iErr, iErr)
}

type CustomError struct{}

func (c *CustomError) Error() string { return "custom error" }
