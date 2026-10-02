// Exercise 01: Temperature Converter
// Objective: Practice variables, constants, type conversions, functions, and control flow.
//
// Task:
// Implement a temperature conversion utility that converts between Celsius, Fahrenheit, and Kelvin.
// Ensure values below Absolute Zero (0 Kelvin / -273.15 C / -459.67 F) return a descriptive error.
package main

import (
	"errors"
	"fmt"
)

// Constants for absolute zero in each scale.
// Using constants instead of magic numbers makes validation readable and maintainable.
// Untyped constants (no explicit type) allow flexible use with any float type.
const (
	AbsoluteZeroCelsius    = -273.15
	AbsoluteZeroFahrenheit = -459.67
	AbsoluteZeroKelvin     = 0.0
)

// CelsiusToFahrenheit converts Celsius to Fahrenheit.
// Returns (result, error) — the idiomatic Go pattern for operations that can fail.
// Uses fmt.Errorf to create a descriptive error with the invalid value embedded.
func CelsiusToFahrenheit(c float64) (float64, error) {
	if c < AbsoluteZeroCelsius {
		return 0, fmt.Errorf("invalid temperature: %.2f C is below absolute zero", c)
	}
	return (c * 9 / 5) + 32, nil
}

// FahrenheitToCelsius converts Fahrenheit to Celsius.
// Uses the same (result, error) pattern for consistency across the API.
func FahrenheitToCelsius(f float64) (float64, error) {
	if f < AbsoluteZeroFahrenheit {
		return 0, fmt.Errorf("invalid temperature: %.2f F is below absolute zero", f)
	}
	return (f - 32) * 5 / 9, nil
}

// CelsiusToKelvin converts Celsius to Kelvin.
// Uses errors.New instead of fmt.Errorf — appropriate when no dynamic values are needed in the message.
func CelsiusToKelvin(c float64) (float64, error) {
	if c < AbsoluteZeroCelsius {
		return 0, errors.New("temperature cannot be below absolute zero")
	}
	return c + 273.15, nil
}

func main() {
	// Test with a mix of valid temps and one invalid (-300°C is below absolute zero)
	testTemps := []float64{0.0, 100.0, 37.0, -300.0}

	fmt.Println("=== Temperature Converter Exercise ===")
	for _, c := range testTemps {
		f, err := CelsiusToFahrenheit(c)
		if err != nil {
			fmt.Printf("Celsius: %6.2f C -> Error: %v\n", c, err)
			continue // Skip to next temperature — don't attempt further conversions for invalid input
		}
		k, _ := CelsiusToKelvin(c)
		fmt.Printf("Celsius: %6.2f C -> %6.2f F | %6.2f K\n", c, f, k)
	}
}
