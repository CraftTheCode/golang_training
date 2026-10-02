// 02_json_and_files demonstrates file I/O and JSON streaming in Go:
//   - JSON struct tags (`json:"..."`) mapping Go struct fields to JSON keys
//   - Streaming with json.NewEncoder / json.NewDecoder vs in-memory json.Marshal / json.Unmarshal
//   - Resource lifecycle management using `defer` for guaranteed file closing and cleanup
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// ServerLog demonstrates struct tags for JSON serialization.
// In Go, struct fields MUST be exported (capitalized: Timestamp, Level...) to be visible
// to the external `encoding/json` package. Struct tags tell the encoder how to rename them in JSON.
type ServerLog struct {
	Timestamp time.Time `json:"timestamp"`  // Serialized as "timestamp" in ISO-8601 format
	Level     string    `json:"level"`      // Serialized as "level"
	Message   string    `json:"message"`    // Serialized as "message"
	LatencyMs int64     `json:"latency_ms"` // Serialized as "latency_ms"
}

func main() {
	tmpDir := os.TempDir()
	filePath := filepath.Join(tmpDir, "demo_logs.json")

	// defer schedules a function call to run immediately before the surrounding function returns.
	// This guarantees cleanup of our temporary file even if a runtime error occurs below!
	defer os.Remove(filePath)

	fmt.Println("=== 1. Writing JSON to File (Streaming Encoder) ===")
	logs := []ServerLog{
		{Timestamp: time.Now().UTC(), Level: "INFO", Message: "Server boot complete", LatencyMs: 12},
		{Timestamp: time.Now().UTC(), Level: "WARN", Message: "Slow response from cache", LatencyMs: 450},
		{Timestamp: time.Now().UTC(), Level: "ERROR", Message: "Failed to persist transaction", LatencyMs: 980},
	}

	// ── json.NewEncoder vs json.Marshal ──
	// json.Marshal(logs) -> serializes everything into an in-memory []byte slice.
	// json.NewEncoder(file).Encode(logs) -> streams JSON directly into the io.Writer (the file)
	//   without allocating a large byte slice in memory!
	// Rule: Use Encoder/Decoder for files, network streams, and HTTP bodies.
	//       Use Marshal/Unmarshal when you already have or need a []byte in memory.
	file, err := os.Create(filePath)
	if err != nil {
		fmt.Printf("Failed to create file: %v\n", err)
		return
	}

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ") // Pretty-print with 2-space indentation
	if err := encoder.Encode(logs); err != nil {
		file.Close()
		fmt.Printf("Failed to encode JSON: %v\n", err)
		return
	}
	file.Close()
	fmt.Printf("Successfully wrote %d logs to %s\n", len(logs), filePath)

	fmt.Println("\n=== 2. Reading and Decoding JSON from File (Streaming Decoder) ===")
	readFile, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("Failed to open file: %v\n", err)
		return
	}
	// Idiomatic pattern: always defer Close() right after checking err != nil
	defer readFile.Close()

	// json.NewDecoder reads directly from the file stream and deserializes into our struct slice.
	var loadedLogs []ServerLog
	if err := json.NewDecoder(readFile).Decode(&loadedLogs); err != nil {
		fmt.Printf("Failed to decode JSON: %v\n", err)
		return
	}

	for i, l := range loadedLogs {
		fmt.Printf("[%d] %s: %-5s -> %s (%dms)\n",
			i+1, l.Timestamp.Format("15:04:05"), l.Level, l.Message, l.LatencyMs)
	}

	fmt.Println("\n=== Key Takeaways ===")
	fmt.Println("• Capitalized fields are exported and visible to encoding/json")
	fmt.Println("• Struct tags `json:\"key\"` control the JSON attribute name")
	fmt.Println("• json.NewEncoder/Decoder stream directly to/from io.Reader/Writer (memory efficient)")
	fmt.Println("• defer guarantees cleanup even if functions return early on error")
}

