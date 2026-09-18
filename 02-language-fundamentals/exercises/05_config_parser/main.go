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

type Config struct {
	values map[string]string
}

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

func (c *Config) GetString(key, fallback string) string {
	if val, ok := c.values[key]; ok {
		return val
	}
	return fallback
}

func (c *Config) GetInt(key string, fallback int) int {
	if val, ok := c.values[key]; ok {
		if parsed, err := strconv.Atoi(val); err == nil {
			return parsed
		}
	}
	return fallback
}

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
