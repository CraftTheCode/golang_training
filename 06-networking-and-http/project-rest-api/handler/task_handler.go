package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Task struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Completed bool      `json:"completed"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateTaskRequest struct {
	Title string `json:"title"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type TaskHandler struct {
	mu     sync.RWMutex
	tasks  map[int]Task
	nextID int
}

func NewTaskHandler() *TaskHandler {
	return &TaskHandler{
		tasks:  make(map[int]Task),
		nextID: 1,
	}
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, ErrorResponse{Error: message})
}

func (h *TaskHandler) Health(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{
		"status": "HEALTHY",
		"time":   time.Now().UTC().Format(time.RFC3339),
	})
}

func (h *TaskHandler) ListTasks(w http.ResponseWriter, r *http.Request) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	list := make([]Task, 0, len(h.tasks))
	for _, t := range h.tasks {
		list = append(list, t)
	}
	writeJSON(w, http.StatusOK, list)
}

func (h *TaskHandler) GetTask(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid task id: must be integer")
		return
	}

	h.mu.RLock()
	defer h.mu.RUnlock()

	task, ok := h.tasks[id]
	if !ok {
		writeError(w, http.StatusNotFound, "task not found")
		return
	}

	writeJSON(w, http.StatusOK, task)
}

func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	var req CreateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	title := strings.TrimSpace(req.Title)
	if title == "" {
		writeError(w, http.StatusUnprocessableEntity, "title cannot be empty")
		return
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	task := Task{
		ID:        h.nextID,
		Title:     title,
		Completed: false,
		CreatedAt: time.Now().UTC(),
	}
	h.tasks[task.ID] = task
	h.nextID++

	writeJSON(w, http.StatusCreated, task)
}

func (h *TaskHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid task id")
		return
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	if _, ok := h.tasks[id]; !ok {
		writeError(w, http.StatusNotFound, "task not found")
		return
	}

	delete(h.tasks, id)
	w.WriteHeader(http.StatusNoContent)
}
