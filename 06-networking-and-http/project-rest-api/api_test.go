package main

import (
	"bytes"
	"encoding/json"
	"golang_training/06-networking-and-http/project-rest-api/handler"
	"net/http"
	"net/http/httptest"
	"testing"
)

func setupTestMux() (*http.ServeMux, *handler.TaskHandler) {
	h := handler.NewTaskHandler()
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", h.Health)
	mux.HandleFunc("GET /api/v1/tasks", h.ListTasks)
	mux.HandleFunc("POST /api/v1/tasks", h.CreateTask)
	mux.HandleFunc("GET /api/v1/tasks/{id}", h.GetTask)
	mux.HandleFunc("DELETE /api/v1/tasks/{id}", h.DeleteTask)
	return mux, h
}

func TestHealthEndpoint(t *testing.T) {
	mux, _ := setupTestMux()

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}
}

func TestTaskLifecycle(t *testing.T) {
	mux, _ := setupTestMux()

	// 1. Create Task
	payload := []byte(`{"title":"Write comprehensive unit tests"}`)
	createReq := httptest.NewRequest(http.MethodPost, "/api/v1/tasks", bytes.NewBuffer(payload))
	createReq.Header.Set("Content-Type", "application/json")
	wCreate := httptest.NewRecorder()

	mux.ServeHTTP(wCreate, createReq)

	if wCreate.Code != http.StatusCreated {
		t.Fatalf("expected status 201 Created, got %d", wCreate.Code)
	}

	var createdTask handler.Task
	if err := json.NewDecoder(wCreate.Body).Decode(&createdTask); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if createdTask.ID != 1 || createdTask.Title != "Write comprehensive unit tests" {
		t.Errorf("unexpected task payload: %+v", createdTask)
	}

	// 2. Fetch Task by ID
	getReq := httptest.NewRequest(http.MethodGet, "/api/v1/tasks/1", nil)
	wGet := httptest.NewRecorder()
	mux.ServeHTTP(wGet, getReq)

	if wGet.Code != http.StatusOK {
		t.Fatalf("expected status 200 OK, got %d", wGet.Code)
	}

	// 3. Delete Task
	delReq := httptest.NewRequest(http.MethodDelete, "/api/v1/tasks/1", nil)
	wDel := httptest.NewRecorder()
	mux.ServeHTTP(wDel, delReq)

	if wDel.Code != http.StatusNoContent {
		t.Fatalf("expected status 204 NoContent, got %d", wDel.Code)
	}

	// 4. Verify 404 after deletion
	get404Req := httptest.NewRequest(http.MethodGet, "/api/v1/tasks/1", nil)
	w404 := httptest.NewRecorder()
	mux.ServeHTTP(w404, get404Req)

	if w404.Code != http.StatusNotFound {
		t.Errorf("expected status 404 NotFound after deletion, got %d", w404.Code)
	}
}
