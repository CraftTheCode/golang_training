package api

import (
	"bytes"
	"encoding/json"
	"golang_training/10-capstone-service-manager/internal/config"
	"golang_training/10-capstone-service-manager/internal/model"
	"golang_training/10-capstone-service-manager/internal/storage"
	"golang_training/10-capstone-service-manager/internal/worker"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
)

func setupTestServer() *Server {
	store := storage.NewMemoryStore()
	cfg := config.Load()
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	monitor := worker.NewHealthMonitor(store, cfg, logger)
	return NewServer(store, monitor)
}

func TestServer_ServiceLifecycle(t *testing.T) {
	server := setupTestServer()
	handler := server.Handler()

	// 1. Health check probe
	reqHealth := httptest.NewRequest(http.MethodGet, "/health", nil)
	wHealth := httptest.NewRecorder()
	handler.ServeHTTP(wHealth, reqHealth)

	if wHealth.Code != http.StatusOK {
		t.Fatalf("expected health status 200, got %d", wHealth.Code)
	}

	// 2. Register Service
	payload := []byte(`{"name":"Local Test Service","url":"http://127.0.0.1:9999"}`)
	reqCreate := httptest.NewRequest(http.MethodPost, "/api/v1/services", bytes.NewBuffer(payload))
	reqCreate.Header.Set("Content-Type", "application/json")
	wCreate := httptest.NewRecorder()
	handler.ServeHTTP(wCreate, reqCreate)

	if wCreate.Code != http.StatusCreated {
		t.Fatalf("expected create status 201, got %d", wCreate.Code)
	}

	var created model.MonitoredService
	if err := json.NewDecoder(wCreate.Body).Decode(&created); err != nil {
		t.Fatalf("failed to decode created service: %v", err)
	}

	if created.ID != 1 || created.Name != "Local Test Service" {
		t.Errorf("unexpected created service: %+v", created)
	}

	// 3. List Services
	reqList := httptest.NewRequest(http.MethodGet, "/api/v1/services", nil)
	wList := httptest.NewRecorder()
	handler.ServeHTTP(wList, reqList)

	var list []model.MonitoredService
	_ = json.NewDecoder(wList.Body).Decode(&list)
	if len(list) != 1 {
		t.Errorf("expected 1 service in list, got %d", len(list))
	}

	// 4. Delete Service
	reqDelete := httptest.NewRequest(http.MethodDelete, "/api/v1/services/1", nil)
	wDelete := httptest.NewRecorder()
	handler.ServeHTTP(wDelete, reqDelete)

	if wDelete.Code != http.StatusNoContent {
		t.Fatalf("expected delete status 204, got %d", wDelete.Code)
	}

	// 5. Verify 404 after deletion
	reqGetDeleted := httptest.NewRequest(http.MethodGet, "/api/v1/services/1", nil)
	wGetDeleted := httptest.NewRecorder()
	handler.ServeHTTP(wGetDeleted, reqGetDeleted)

	if wGetDeleted.Code != http.StatusNotFound {
		t.Errorf("expected status 404 after deletion, got %d", wGetDeleted.Code)
	}
}
