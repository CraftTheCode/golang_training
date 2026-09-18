package bench

import (
	"fmt"
	"testing"
)

var sampleStrings = []string{
	"Go", "is", "a", "statically", "typed,", "compiled", "programming", "language",
	"designed", "at", "Google", "by", "Robert", "Griesemer,", "Rob", "Pike,",
	"and", "Ken", "Thompson.", "Go", "is", "syntactically", "similar", "to", "C.",
}

func BenchmarkConcatPlus(b *testing.B) {
	b.ReportAllocs() // Reports B/op and allocs/op
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = ConcatPlus(sampleStrings)
	}
}

func BenchmarkConcatBuilder(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = ConcatBuilder(sampleStrings)
	}
}

func Example() {
	// Demonstrates godoc runnable example
	out := ConcatBuilder([]string{"Hello, ", "World!"})
	fmt.Println(out)
	// Output: Hello, World!
}
