package worker

import (
	"context"
	"golang_training/10-capstone-service-manager/internal/config"
	"golang_training/10-capstone-service-manager/internal/model"
	"golang_training/10-capstone-service-manager/internal/storage"
	"log/slog"
	"net/http"
	"sync"
	"time"
)

type HealthMonitor struct {
	store      storage.Store
	cfg        *config.Config
	httpClient *http.Client
	logger     *slog.Logger
}

func NewHealthMonitor(store storage.Store, cfg *config.Config, logger *slog.Logger) *HealthMonitor {
	return &HealthMonitor{
		store: store,
		cfg:   cfg,
		httpClient: &http.Client{
			Timeout: cfg.ProbeTimeout,
		},
		logger: logger,
	}
}

// Start launches the periodic background health checker until ctx is cancelled.
func (m *HealthMonitor) Start(ctx context.Context) {
	ticker := time.NewTicker(m.cfg.CheckInterval)
	defer ticker.Stop()

	m.logger.Info("Background health monitor started", "interval", m.cfg.CheckInterval.String())

	// Run initial check immediately
	m.runBatchCheck(ctx)

	for {
		select {
		case <-ctx.Done():
			m.logger.Info("Health monitor stopping gracefully...")
			return
		case <-ticker.C:
			m.runBatchCheck(ctx)
		}
	}
}

// CheckSingleService performs an immediate probe of a single service and updates storage.
func (m *HealthMonitor) CheckSingleService(ctx context.Context, svc model.MonitoredService) model.HealthCheckResult {
	start := time.Now()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, svc.URL, nil)
	if err != nil {
		_ = m.store.UpdateHealth(svc.ID, model.StatusUnhealthy, 0)
		return model.HealthCheckResult{ServiceID: svc.ID, Status: model.StatusUnhealthy, Error: err}
	}

	resp, err := m.httpClient.Do(req)
	latency := time.Since(start).Milliseconds()

	if err != nil {
		_ = m.store.UpdateHealth(svc.ID, model.StatusUnhealthy, latency)
		return model.HealthCheckResult{ServiceID: svc.ID, Status: model.StatusUnhealthy, LatencyMs: latency, Error: err}
	}
	defer resp.Body.Close()

	var status model.HealthStatus
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		if latency > m.cfg.DegradedThreshold {
			status = model.StatusDegraded
		} else {
			status = model.StatusHealthy
		}
	} else {
		status = model.StatusUnhealthy
	}

	_ = m.store.UpdateHealth(svc.ID, status, latency)

	return model.HealthCheckResult{
		ServiceID: svc.ID,
		Status:    status,
		LatencyMs: latency,
	}
}

func (m *HealthMonitor) runBatchCheck(ctx context.Context) {
	services := m.store.GetAll()
	if len(services) == 0 {
		return
	}

	jobs := make(chan model.MonitoredService, len(services))
	var wg sync.WaitGroup

	// Bounded worker pool
	for i := 0; i < m.cfg.WorkerPoolSize; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for svc := range jobs {
				select {
				case <-ctx.Done():
					return
				default:
					res := m.CheckSingleService(ctx, svc)
					m.logger.Debug("Health probe completed",
						"service", svc.Name,
						"status", res.Status,
						"latency_ms", res.LatencyMs)
				}
			}
		}()
	}

	for _, svc := range services {
		jobs <- svc
	}
	close(jobs)

	wg.Wait()
}
