// Exercise 05: Configuration Parser
// Objective: Practice strings, parsing, type casting, and structured errors.
//
// Task:
// Parse an INI/env-style multi-line string into a structured configuration map.
// Ignore comments ('#' and '//') and blank lines. Provide helper methods to fetch
// values as String, Int, and Bool with default fallbacks.
package main

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"
)

// Config stores parsed key-value pairs from a configuration string.
// The values map is unexported (lowercase 'v') — external code must
// use the typed accessor methods (GetString, GetInt, GetBool) instead
// of reading raw strings directly. This is encapsulation in Go.
type Config struct {
	values map[string]string
}

// ParseConfig reads an INI/env-style config string and returns a populated Config.
// Returns (*Config, error) — the pointer avoids copying the map on return.
// Uses bufio.Scanner for line-by-line parsing, which handles large inputs efficiently
// without loading the entire string into memory as a slice of lines.
func ParseConfig(raw string) (*Config, error) {
	cfg := &Config{
		values: make(map[string]string),
	}

	scanner := bufio.NewScanner(strings.NewReader(raw))
	lineNum := 0
	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines and comment lines
		if line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, "//") {
			continue
		}

		// SplitN limits to 2 parts so values containing '=' (like URLs) aren't broken.
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("line %d: malformed config entry (missing '='): %q", lineNum, line)
		}

		key := strings.TrimSpace(parts[0])
		val := strings.TrimSpace(parts[1])

		// Strip surrounding quotes if present
		val = strings.Trim(val, `"'`)
		cfg.values[key] = val
	}

	return cfg, scanner.Err()
}

// GetString retrieves a config value as a string, with a fallback default.
// The fallback pattern avoids nil/error checks at every call site — callers always
// get a usable value. This is a common Go API design pattern.
func (c *Config) GetString(key, fallback string) string {
	if val, ok := c.values[key]; ok {
		return val
	}
	return fallback
}

// GetInt retrieves a config value as an int, using strconv.Atoi for parsing.
// If the key is missing OR the value isn't a valid integer, returns the fallback.
// This double-guard (comma-ok + parse error) makes the API safe against bad config.
func (c *Config) GetInt(key string, fallback int) int {
	if val, ok := c.values[key]; ok {
		if parsed, err := strconv.Atoi(val); err == nil {
			return parsed
		}
	}
	return fallback
}

// GetBool retrieves a config value as a bool using strconv.ParseBool.
// strconv.ParseBool recognizes: "1", "t", "true", "TRUE", "0", "f", "false", "FALSE".
func (c *Config) GetBool(key string, fallback bool) bool {
	if val, ok := c.values[key]; ok {
		if parsed, err := strconv.ParseBool(val); err == nil {
			return parsed
		}
	}
	return fallback
}

func main() {
	rawConfig := `
		# Application Configuration
		APP_NAME = "GoLearningApp"
		PORT = 8080
		DEBUG_MODE = true

		// Worker pool settings
		MAX_WORKERS = 16
		TIMEOUT_SECONDS = 30
	`

	cfg, err := ParseConfig(rawConfig)
	if err != nil {
		fmt.Println("Error parsing config:", err)
		return
	}

	fmt.Println("=== Parsed Configuration ===")
	fmt.Printf("App Name: %s\n", cfg.GetString("APP_NAME", "DefaultApp"))
	fmt.Printf("Port: %d\n", cfg.GetInt("PORT", 3000))
	fmt.Printf("Debug Mode: %t\n", cfg.GetBool("DEBUG_MODE", false))
	fmt.Printf("Max Workers: %d\n", cfg.GetInt("MAX_WORKERS", 4))
	fmt.Printf("Database URL (Fallback): %s\n", cfg.GetString("DATABASE_URL", "localhost:5432"))
}
