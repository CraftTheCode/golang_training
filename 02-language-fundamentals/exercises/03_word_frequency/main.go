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

func sanitizeToken(token string) string {
	var builder strings.Builder
	for _, r := range strings.ToLower(token) {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			builder.WriteRune(r)
		}
	}
	return builder.String()
}

func CountFrequencies(text string) map[string]int {
	counts := make(map[string]int)
	tokens := strings.Fields(text)

	for _, token := range tokens {
		clean := sanitizeToken(token)
		if clean != "" {
			counts[clean]++
		}
	}
	return counts
}

func TopNWords(counts map[string]int, n int) []WordCount {
	pairs := make([]WordCount, 0, len(counts))
	for w, c := range counts {
		pairs = append(pairs, WordCount{Word: w, Count: c})
	}

	// Sort descending by count, tie-breaker alphabetically
	sort.Slice(pairs, func(i, j int) bool {
		if pairs[i].Count == pairs[j].Count {
			return pairs[i].Word < pairs[j].Word
		}
		return pairs[i].Count > pairs[j].Count
	})

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
