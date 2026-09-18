package api

import (
	"encoding/json"
	"errors"
	"golang_training/10-capstone-service-manager/internal/model"
	"golang_training/10-capstone-service-manager/internal/storage"
	"golang_training/10-capstone-service-manager/internal/worker"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Server struct {
	store   storage.Store
	monitor *worker.HealthMonitor
	mux     *http.ServeMux
}

func NewServer(store storage.Store, monitor *worker.HealthMonitor) *Server {
	s := &Server{
		store:   store,
		monitor: monitor,
		mux:     http.NewServeMux(),
	}
	s.routes()
	return s
}

func (s *Server) Handler() http.Handler {
	return s.mux
}

func (s *Server) routes() {
	s.mux.HandleFunc("GET /health", s.handleHealth)
	s.mux.HandleFunc("GET /api/v1/services", s.handleListServices)
	s.mux.HandleFunc("POST /api/v1/services", s.handleCreateService)
	s.mux.HandleFunc("GET /api/v1/services/{id}", s.handleGetService)
	s.mux.HandleFunc("POST /api/v1/services/{id}/check", s.handleCheckServiceNow)
	s.mux.HandleFunc("DELETE /api/v1/services/{id}", s.handleDeleteService)
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{
		"status":    "UP",
		"timestamp": time.Now().UTC(),
	})
}

func (s *Server) handleListServices(w http.ResponseWriter, r *http.Request) {
	services := s.store.GetAll()
	writeJSON(w, http.StatusOK, services)
}

func (s *Server) handleCreateService(w http.ResponseWriter, r *http.Request) {
	var req model.CreateServiceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	name := strings.TrimSpace(req.Name)
	url := strings.TrimSpace(req.URL)

	if name == "" || url == "" {
		writeError(w, http.StatusUnprocessableEntity, "name and url are required")
		return
	}

	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		writeError(w, http.StatusUnprocessableEntity, "url must begin with http:// or https://")
		return
	}

	svc, err := s.store.Create(name, url)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to create service")
		return
	}

	writeJSON(w, http.StatusCreated, svc)
}

func (s *Server) handleGetService(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid service id")
		return
	}

	svc, err := s.store.GetByID(id)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			writeError(w, http.StatusNotFound, "service not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}

	writeJSON(w, http.StatusOK, svc)
}

func (s *Server) handleCheckServiceNow(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid service id")
		return
	}

	svc, err := s.store.GetByID(id)
	if err != nil {
		writeError(w, http.StatusNotFound, "service not found")
		return
	}

	result := s.monitor.CheckSingleService(r.Context(), *svc)
	writeJSON(w, http.StatusOK, result)
}

func (s *Server) handleDeleteService(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid service id")
		return
	}

	if err := s.store.Delete(id); err != nil {
		writeError(w, http.StatusNotFound, "service not found")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
