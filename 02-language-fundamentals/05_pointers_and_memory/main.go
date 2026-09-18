package main

import (
	"fmt"
)

type Account struct {
	Owner   string
	Balance float64
}

// depositValue takes an Account by VALUE (copy).
// Modifications here do NOT affect the caller's struct!
func depositValue(acc Account, amount float64) {
	acc.Balance += amount
	fmt.Printf("[Inside depositValue] New Balance: %.2f (on copied struct)\n", acc.Balance)
}

// depositPointer takes an Account by POINTER.
// Modifications affect the original memory location!
func depositPointer(acc *Account, amount float64) {
	acc.Balance += amount
	fmt.Printf("[Inside depositPointer] New Balance: %.2f (on pointed struct)\n", acc.Balance)
}

func main() {
	fmt.Println("=== 1. Pointer Basics ===")
	val := 42
	ptr := &val // Pointer to val (*int)

	fmt.Printf("val value: %d\n", val)
	fmt.Printf("val address (&val): %p\n", ptr)
	fmt.Printf("Dereferenced pointer (*ptr): %d\n", *ptr)

	*ptr = 99 // Mutate via dereference
	fmt.Printf("val after *ptr = 99: %d\n", val)

	fmt.Println("\n=== 2. Value Semantics vs Pointer Semantics ===")
	myAccount := Account{
		Owner:   "Divya",
		Balance: 1000.00,
	}

	fmt.Printf("Initial Account: %+v\n", myAccount)

	// Attempt mutation with Value Semantics
	depositValue(myAccount, 250.00)
	fmt.Printf("After depositValue: %+v (Unchanged!)\n", myAccount)

	// Mutation with Pointer Semantics
	depositPointer(&myAccount, 250.00)
	fmt.Printf("After depositPointer: %+v (Successfully updated!)\n", myAccount)

	fmt.Println("\n=== 3. Pointer to Non-Existent Memory (nil) ===")
	var uninitialized *Account
	fmt.Printf("Uninitialized pointer: %v\n", uninitialized)
	// Dereferencing a nil pointer causes a panic:
	// *uninitialized = Account{} -> panic: runtime error: invalid memory address or nil pointer dereference
	if uninitialized == nil {
		fmt.Println("Safely checked: pointer is nil before dereferencing.")
	}
}
