package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type ServerLog struct {
	Timestamp time.Time `json:"timestamp"`
	Level     string    `json:"level"`
	Message   string    `json:"message"`
	LatencyMs int64     `json:"latency_ms"`
}

func main() {
	tmpDir := os.TempDir()
	filePath := filepath.Join(tmpDir, "demo_logs.json")
	defer os.Remove(filePath) // Cleanup when main returns

	fmt.Println("=== 1. Writing JSON to File ===")
	logs := []ServerLog{
		{Timestamp: time.Now().UTC(), Level: "INFO", Message: "Server boot complete", LatencyMs: 12},
		{Timestamp: time.Now().UTC(), Level: "WARN", Message: "Slow response from cache", LatencyMs: 450},
		{Timestamp: time.Now().UTC(), Level: "ERROR", Message: "Failed to persist transaction", LatencyMs: 980},
	}

	// Create and write with streaming encoder
	file, err := os.Create(filePath)
	if err != nil {
		fmt.Printf("Failed to create file: %v\n", err)
		return
	}

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(logs); err != nil {
		file.Close()
		fmt.Printf("Failed to encode JSON: %v\n", err)
		return
	}
	file.Close()
	fmt.Printf("Successfully wrote %d logs to %s\n", len(logs), filePath)

	fmt.Println("\n=== 2. Reading and Decoding JSON from File ===")
	readFile, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("Failed to open file: %v\n", err)
		return
	}
	defer readFile.Close()

	var loadedLogs []ServerLog
	if err := json.NewDecoder(readFile).Decode(&loadedLogs); err != nil {
		fmt.Printf("Failed to decode JSON: %v\n", err)
		return
	}

	for i, l := range loadedLogs {
		fmt.Printf("[%d] %s: %-5s -> %s (%dms)\n",
			i+1, l.Timestamp.Format("15:04:05"), l.Level, l.Message, l.LatencyMs)
	}
}
