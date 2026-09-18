// Package config provides application configuration loading and validation.
package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// Config holds the application configuration parameters.
type Config struct {
	AppName         string
	Port            int
	Environment     string // "development", "staging", "production"
	ShutdownTimeout time.Duration
	EnableMetrics   bool
}

// ConfigError represents a structured configuration validation error.
type ConfigError struct {
	Field   string
	Reason  string
}

func (e *ConfigError) Error() string {
	return fmt.Sprintf("invalid config for field %q: %s", e.Field, e.Reason)
}

// EnvLookup is an abstraction for fetching environment variables, allowing easy testing.
type EnvLookup func(key string) (string, bool)

// Load reads configuration using the system environment variables and applies defaults.
func Load() (*Config, error) {
	return LoadWithLookup(os.LookupEnv)
}

// LoadWithLookup reads configuration using the provided lookup function.
func LoadWithLookup(lookup EnvLookup) (*Config, error) {
	cfg := &Config{
		AppName:         "GoApp",
		Port:            8080,
		Environment:     "development",
		ShutdownTimeout: 10 * time.Second,
		EnableMetrics:   true,
	}

	if val, ok := lookup("APP_NAME"); ok && strings.TrimSpace(val) != "" {
		cfg.AppName = strings.TrimSpace(val)
	}

	if val, ok := lookup("PORT"); ok {
		port, err := strconv.Atoi(val)
		if err != nil {
			return nil, &ConfigError{Field: "PORT", Reason: "must be an integer"}
		}
		cfg.Port = port
	}

	if val, ok := lookup("ENVIRONMENT"); ok {
		cfg.Environment = strings.ToLower(strings.TrimSpace(val))
	}

	if val, ok := lookup("SHUTDOWN_TIMEOUT_SECONDS"); ok {
		secs, err := strconv.Atoi(val)
		if err != nil {
			return nil, &ConfigError{Field: "SHUTDOWN_TIMEOUT_SECONDS", Reason: "must be an integer"}
		}
		cfg.ShutdownTimeout = time.Duration(secs) * time.Second
	}

	if val, ok := lookup("ENABLE_METRICS"); ok {
		parsed, err := strconv.ParseBool(val)
		if err != nil {
			return nil, &ConfigError{Field: "ENABLE_METRICS", Reason: "must be a boolean (true/false)"}
		}
		cfg.EnableMetrics = parsed
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

// Validate ensures that all configuration values adhere to business constraints.
func (c *Config) Validate() error {
	if c.Port < 1 || c.Port > 65535 {
		return &ConfigError{Field: "PORT", Reason: "port must be between 1 and 65535"}
	}

	switch c.Environment {
	case "development", "staging", "production":
		// valid
	default:
		return &ConfigError{
			Field:  "ENVIRONMENT",
			Reason: fmt.Sprintf("environment %q must be one of: development, staging, production", c.Environment),
		}
	}

	if c.ShutdownTimeout <= 0 {
		return &ConfigError{Field: "SHUTDOWN_TIMEOUT_SECONDS", Reason: "timeout must be positive"}
	}

	return nil
}
