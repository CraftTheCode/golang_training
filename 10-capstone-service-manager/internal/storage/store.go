package storage

import (
	"errors"
	"fmt"
	"golang_training/10-capstone-service-manager/internal/model"
	"sync"
	"time"
)

var (
	ErrNotFound = errors.New("service not found")
)

// Store defines the storage operations for monitored services.
type Store interface {
	GetAll() []model.MonitoredService
	GetByID(id int) (*model.MonitoredService, error)
	Create(name, url string) (*model.MonitoredService, error)
	UpdateHealth(id int, status model.HealthStatus, latencyMs int64) error
	Delete(id int) error
}

// MemoryStore provides a thread-safe in-memory store.
type MemoryStore struct {
	mu       sync.RWMutex
	services map[int]model.MonitoredService
	nextID   int
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		services: make(map[int]model.MonitoredService),
		nextID:   1,
	}
}

func (s *MemoryStore) GetAll() []model.MonitoredService {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]model.MonitoredService, 0, len(s.services))
	for _, svc := range s.services {
		result = append(result, svc)
	}
	return result
}

func (s *MemoryStore) GetByID(id int) (*model.MonitoredService, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	svc, exists := s.services[id]
	if !exists {
		return nil, ErrNotFound
	}
	return &svc, nil
}

func (s *MemoryStore) Create(name, url string) (*model.MonitoredService, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	svc := model.MonitoredService{
		ID:        s.nextID,
		Name:      name,
		URL:       url,
		Status:    model.StatusUnknown,
		LatencyMs: 0,
		CreatedAt: time.Now().UTC(),
	}

	s.services[svc.ID] = svc
	s.nextID++

	return &svc, nil
}

func (s *MemoryStore) UpdateHealth(id int, status model.HealthStatus, latencyMs int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	svc, exists := s.services[id]
	if !exists {
		return ErrNotFound
	}

	now := time.Now().UTC()
	svc.Status = status
	svc.LatencyMs = latencyMs
	svc.LastChecked = &now

	s.services[id] = svc
	return nil
}

func (s *MemoryStore) Delete(id int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.services[id]; !exists {
		return fmt.Errorf("cannot delete: %w", ErrNotFound)
	}

	delete(s.services, id)
	return nil
}
