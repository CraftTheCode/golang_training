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

const (
	AbsoluteZeroCelsius    = -273.15
	AbsoluteZeroFahrenheit = -459.67
	AbsoluteZeroKelvin     = 0.0
)

func CelsiusToFahrenheit(c float64) (float64, error) {
	if c < AbsoluteZeroCelsius {
		return 0, fmt.Errorf("invalid temperature: %.2f C is below absolute zero", c)
	}
	return (c * 9 / 5) + 32, nil
}

func FahrenheitToCelsius(f float64) (float64, error) {
	if f < AbsoluteZeroFahrenheit {
		return 0, fmt.Errorf("invalid temperature: %.2f F is below absolute zero", f)
	}
	return (f - 32) * 5 / 9, nil
}

func CelsiusToKelvin(c float64) (float64, error) {
	if c < AbsoluteZeroCelsius {
		return 0, errors.New("temperature cannot be below absolute zero")
	}
	return c + 273.15, nil
}

func main() {
	testTemps := []float64{0.0, 100.0, 37.0, -300.0}

	fmt.Println("=== Temperature Converter Exercise ===")
	for _, c := range testTemps {
		f, err := CelsiusToFahrenheit(c)
		if err != nil {
			fmt.Printf("Celsius: %6.2f C -> Error: %v\n", c, err)
			continue
		}
		k, _ := CelsiusToKelvin(c)
		fmt.Printf("Celsius: %6.2f C -> %6.2f F | %6.2f K\n", c, f, k)
	}
}
