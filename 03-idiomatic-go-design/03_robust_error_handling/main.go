package main

import (
	"errors"
	"fmt"
	"os"
)

// 1. Sentinel Errors: Predefined error values
var (
	ErrNotFound      = errors.New("resource not found")
	ErrPermission    = errors.New("permission denied")
	ErrInvalidInput  = errors.New("invalid input provided")
)

// 2. Custom Error Type carrying structured context
type ValidationError struct {
	Field   string
	Message string
}

func (v *ValidationError) Error() string {
	return fmt.Sprintf("validation failed on field '%s': %s", v.Field, v.Message)
}

// simulateDatabaseQuery simulates an operation that might fail
func simulateDatabaseQuery(id int) (string, error) {
	if id <= 0 {
		return "", &ValidationError{
			Field:   "ID",
			Message: "must be a positive integer",
		}
	}
	if id == 404 {
		// Wrap the sentinel error using %w
		return "", fmt.Errorf("user lookup failed: %w", ErrNotFound)
	}
	return fmt.Sprintf("User #%d", id), nil
}

func main() {
	fmt.Println("=== 1. Error Wrapping and errors.Is ===")
	_, err := simulateDatabaseQuery(404)
	if err != nil {
		fmt.Printf("Raw error: %v\n", err)

		// errors.Is checks the entire wrapped error chain!
		if errors.Is(err, ErrNotFound) {
			fmt.Println("Matched ErrNotFound using errors.Is (even though wrapped!)")
		}
	}

	fmt.Println("\n=== 2. Custom Errors and errors.As ===")
	_, errVal := simulateDatabaseQuery(-5)
	if errVal != nil {
		fmt.Printf("Raw error: %v\n", errVal)

		// errors.As extracts the custom error type from the chain
		var valErr *ValidationError
		if errors.As(errVal, &valErr) {
			fmt.Printf("Successfully extracted ValidationError: Field=%s, Message=%s\n",
				valErr.Field, valErr.Message)
		}
	}

	fmt.Println("\n=== 3. Standard Library Error Inspection ===")
	_, errFile := os.Open("non-existent-file.txt")
	if errors.Is(errFile, os.ErrNotExist) {
		fmt.Println("Cleanly verified file does not exist using os.ErrNotExist")
	}
}
