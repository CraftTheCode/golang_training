package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	Port              int
	CheckInterval     time.Duration
	ProbeTimeout      time.Duration
	WorkerPoolSize    int
	DegradedThreshold int64 // Latency threshold in ms above which service is marked degraded
}

func Load() *Config {
	cfg := &Config{
		Port:              8080,
		CheckInterval:     15 * time.Second,
		ProbeTimeout:      3 * time.Second,
		WorkerPoolSize:    5,
		DegradedThreshold: 800, // 800ms
	}

	if val := os.Getenv("PORT"); val != "" {
		if p, err := strconv.Atoi(val); err == nil {
			cfg.Port = p
		}
	}

	if val := os.Getenv("CHECK_INTERVAL_SECONDS"); val != "" {
		if sec, err := strconv.Atoi(val); err == nil {
			cfg.CheckInterval = time.Duration(sec) * time.Second
		}
	}

	if val := os.Getenv("WORKER_POOL_SIZE"); val != "" {
		if size, err := strconv.Atoi(val); err == nil {
			cfg.WorkerPoolSize = size
		}
	}

	return cfg
}
