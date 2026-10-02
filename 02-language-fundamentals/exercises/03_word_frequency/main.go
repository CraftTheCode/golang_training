// Exercise 03: Word Frequency Counter
// Objective: Practice strings, runes, maps, structs, and custom sorting.
//
// Task:
// Tokenize an input passage, sanitize words (lowercase, strip non-alphanumeric runes),
// compute frequencies using a map, and display the Top N words sorted by count descending.
package main

import (
	"fmt"
	"sort"
	"strings"
	"unicode"
)

type WordCount struct {
	Word  string
	Count int
}

// sanitizeToken normalizes a word: lowercases and strips non-alphanumeric runes.
// Uses strings.Builder for efficient string construction — avoids creating
// intermediate strings on each rune append (unlike string concatenation with +).
func sanitizeToken(token string) string {
	var builder strings.Builder
	for _, r := range strings.ToLower(token) {
		// unicode.IsLetter/IsDigit work on runes (Unicode code points),
		// not bytes — so this correctly handles non-ASCII text.
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			builder.WriteRune(r)
		}
	}
	return builder.String()
}

// CountFrequencies tokenizes text by whitespace and counts occurrences of each cleaned word.
// strings.Fields splits on any whitespace (spaces, tabs, newlines) — more robust than strings.Split.
func CountFrequencies(text string) map[string]int {
	counts := make(map[string]int)
	tokens := strings.Fields(text)

	for _, token := range tokens {
		clean := sanitizeToken(token)
		if clean != "" {
			// map[key]++ works because Go zero-values missing keys to 0
			counts[clean]++
		}
	}
	return counts
}

// TopNWords returns the N most frequent words, sorted by count descending.
// Demonstrates: converting a map to a slice for sorting, and sort.Slice with a closure.
func TopNWords(counts map[string]int, n int) []WordCount {
	// Pre-allocate slice with capacity = map size to avoid repeated grow+copy.
	pairs := make([]WordCount, 0, len(counts))
	for w, c := range counts {
		pairs = append(pairs, WordCount{Word: w, Count: c})
	}

	// sort.Slice takes a "less" function — a closure that captures 'pairs'.
	// Sort descending by count; alphabetically as tie-breaker for deterministic output.
	sort.Slice(pairs, func(i, j int) bool {
		if pairs[i].Count == pairs[j].Count {
			return pairs[i].Word < pairs[j].Word
		}
		return pairs[i].Count > pairs[j].Count
	})

	// Guard against n > len(pairs) to prevent slice-out-of-bounds panic.
	if n > len(pairs) {
		n = len(pairs)
	}
	return pairs[:n]
}

func main() {
	passage := `
		Go is an open-source programming language that makes it easy to build simple,
		reliable, and efficient software. Go was designed at Google. Go has goroutines,
		channels, and a simple syntax. Simple software is reliable software.
	`

	freqs := CountFrequencies(passage)
	top5 := TopNWords(freqs, 5)

	fmt.Println("=== Top 5 Word Frequencies ===")
	for rank, item := range top5 {
		fmt.Printf("%d. %-12s : %d occurrences\n", rank+1, item.Word, item.Count)
	}
}
