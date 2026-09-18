package bench

import "strings"

// ConcatPlus repeatedly concatenates using the '+' operator (triggers O(N^2) memory allocations)
func ConcatPlus(parts []string) string {
	res := ""
	for _, p := range parts {
		res += p
	}
	return res
}

// ConcatBuilder uses strings.Builder (reallocates logarithmically or preallocates)
func ConcatBuilder(parts []string) string {
	var builder strings.Builder
	for _, p := range parts {
		builder.WriteString(p)
	}
	return builder.String()
}
