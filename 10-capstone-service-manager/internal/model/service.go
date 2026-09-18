package model

import "time"

type HealthStatus string

const (
	StatusHealthy   HealthStatus = "HEALTHY"
	StatusDegraded  HealthStatus = "DEGRADED"
	StatusUnhealthy HealthStatus = "UNHEALTHY"
	StatusUnknown   HealthStatus = "UNKNOWN"
)

// MonitoredService represents a service registered with GSM for monitoring.
type MonitoredService struct {
	ID          int          `json:"id"`
	Name        string       `json:"name"`
	URL         string       `json:"url"`
	Status      HealthStatus `json:"status"`
	LatencyMs   int64        `json:"latency_ms"`
	LastChecked *time.Time   `json:"last_checked,omitempty"`
	CreatedAt   time.Time    `json:"created_at"`
}

// CreateServiceRequest defines the payload for registering a new service.
type CreateServiceRequest struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

// HealthCheckResult represents the outcome of probing a single service.
type HealthCheckResult struct {
	ServiceID int
	Status    HealthStatus
	LatencyMs int64
	Error     error
}
