package calc

import (
	"math"
	"testing"
)

func TestAdd(t *testing.T) {
	// Canonical Table-Driven Test Structure
	tests := []struct {
		name     string
		a, b     int
		expected int
	}{
		{"Positive numbers", 2, 3, 5},
		{"Negative numbers", -5, -3, -8},
		{"Zero identity", 42, 0, 42},
		{"Cancelling terms", 10, -10, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Add(tt.a, tt.b)
			if result != tt.expected {
				t.Errorf("Add(%d, %d) = %d; want %d", tt.a, tt.b, result, tt.expected)
			}
		})
	}
}

func TestDivide(t *testing.T) {
	tests := []struct {
		name      string
		a, b      float64
		expected  float64
		expectErr bool
	}{
		{"Exact division", 10, 2, 5.0, false},
		{"Fractional division", 1, 3, 1.0 / 3.0, false},
		{"Division by zero error", 10, 0, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Divide(tt.a, tt.b)

			if (err != nil) != tt.expectErr {
				t.Fatalf("Divide(%.2f, %.2f) error = %v; expectErr = %v", tt.a, tt.b, err, tt.expectErr)
			}

			if !tt.expectErr {
				if math.Abs(result-tt.expected) > 1e-9 {
					t.Errorf("Divide(%.2f, %.2f) = %f; want %f", tt.a, tt.b, result, tt.expected)
				}
			}
		})
	}
}
