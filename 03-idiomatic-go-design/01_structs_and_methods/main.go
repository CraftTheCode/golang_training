package main

import (
	"fmt"
	"math"
)

// Shape represents a 2D geometric entity.
type Circle struct {
	Radius float64
}

// Area uses a VALUE receiver because it only reads fields without mutation.
func (c Circle) Area() float64 {
	return math.Pi * c.Radius * c.Radius
}

// Scale uses a POINTER receiver because it must mutate the struct in-place.
func (c *Circle) Scale(factor float64) {
	c.Radius *= factor
}

// -------------------------------------------------------------
// Composition (Embedding) Demonstration
// -------------------------------------------------------------

type Address struct {
	City    string
	Country string
}

func (a Address) LocationString() string {
	return fmt.Sprintf("%s, %s", a.City, a.Country)
}

// User embeds Address directly.
// All fields and methods of Address are PROMOTED to User!
type User struct {
	ID   int
	Name string
	Address
}

func main() {
	fmt.Println("=== 1. Value vs Pointer Receivers ===")
	c := Circle{Radius: 5.0}
	fmt.Printf("Initial Circle: Radius=%.2f, Area=%.2f\n", c.Radius, c.Area())

	// Scale requires pointer semantics
	c.Scale(2.0)
	fmt.Printf("After Scale(2.0): Radius=%.2f, Area=%.2f\n", c.Radius, c.Area())

	fmt.Println("\n=== 2. Composition (Struct Embedding) ===")
	u := User{
		ID:   101,
		Name: "Divya",
		Address: Address{
			City:    "Bengaluru",
			Country: "India",
		},
	}

	// Direct field access via promotion
	fmt.Printf("User: %s from %s (direct access)\n", u.Name, u.City)

	// Direct method call via promotion
	fmt.Printf("Promoted Method: %s\n", u.LocationString())

	// Explicit inner struct access still works
	fmt.Printf("Explicit: %s\n", u.Address.LocationString())
}
