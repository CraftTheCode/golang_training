// Exercise 02: Log Line Counter
// Objective: Practice strings, slices, maps, loops, and conditional analysis.
//
// Task:
// Parse a batch of server log lines, extract log levels (INFO, WARN, ERROR, DEBUG),
// count the occurrences of each level, and calculate the overall error rate percentage.
package main

import (
	"fmt"
	"strings"
)

type LogStats struct {
	Total     int
	Counts    map[string]int
	ErrorRate float64
}

func AnalyzeLogs(logLines []string) LogStats {
	stats := LogStats{
		Counts: make(map[string]int),
	}

	for _, line := range logLines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}
		stats.Total++

		upper := strings.ToUpper(trimmed)
		switch {
		case strings.Contains(upper, "[ERROR]"):
			stats.Counts["ERROR"]++
		case strings.Contains(upper, "[WARN]"):
			stats.Counts["WARN"]++
		case strings.Contains(upper, "[INFO]"):
			stats.Counts["INFO"]++
		case strings.Contains(upper, "[DEBUG]"):
			stats.Counts["DEBUG"]++
		default:
			stats.Counts["UNKNOWN"]++
		}
	}

	if stats.Total > 0 {
		stats.ErrorRate = (float64(stats.Counts["ERROR"]) / float64(stats.Total)) * 100.0
	}

	return stats
}

func main() {
	sampleLogs := []string{
		"2026-09-18 10:00:00 [INFO] Server started on port 8080",
		"2026-09-18 10:00:05 [DEBUG] Loaded 45 routes",
		"2026-09-18 10:01:20 [WARN] High memory usage detected: 82%",
		"2026-09-18 10:02:15 [ERROR] Failed to connect to database: timeout",
		"2026-09-18 10:02:40 [INFO] Client 192.168.1.1 connected",
		"2026-09-18 10:03:00 [ERROR] Payment gateway returned 500",
		"2026-09-18 10:04:10 [INFO] Health check probe passed",
	}

	stats := AnalyzeLogs(sampleLogs)

	fmt.Println("=== Log Analysis Report ===")
	fmt.Printf("Total Log Lines: %d\n", stats.Total)
	for level, count := range stats.Counts {
		fmt.Printf("- %-8s: %d\n", level, count)
	}
	fmt.Printf("Error Rate: %.2f%%\n", stats.ErrorRate)
}
