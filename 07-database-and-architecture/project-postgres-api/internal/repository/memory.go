package repository

import (
	"context"
	"golang_training/07-database-and-architecture/project-postgres-api/internal/domain"
	"sync"
	"time"
)

// MemoryRepository implements domain.TaskRepository entirely in-memory for testing.
type MemoryRepository struct {
	mu     sync.RWMutex
	tasks  map[int]domain.Task
	nextID int
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		tasks:  make(map[int]domain.Task),
		nextID: 1,
	}
}

func (m *MemoryRepository) GetAll(ctx context.Context) ([]domain.Task, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	list := make([]domain.Task, 0, len(m.tasks))
	for _, t := range m.tasks {
		list = append(list, t)
	}
	return list, nil
}

func (m *MemoryRepository) GetByID(ctx context.Context, id int) (*domain.Task, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	task, ok := m.tasks[id]
	if !ok {
		return nil, domain.ErrNotFound
	}
	return &task, nil
}

func (m *MemoryRepository) Create(ctx context.Context, title string) (*domain.Task, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	task := domain.Task{
		ID:        m.nextID,
		Title:     title,
		Completed: false,
		CreatedAt: time.Now().UTC(),
	}
	m.tasks[task.ID] = task
	m.nextID++
	return &task, nil
}

func (m *MemoryRepository) UpdateStatus(ctx context.Context, id int, completed bool) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	task, ok := m.tasks[id]
	if !ok {
		return domain.ErrNotFound
	}
	task.Completed = completed
	m.tasks[id] = task
	return nil
}

func (m *MemoryRepository) Delete(ctx context.Context, id int) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.tasks[id]; !ok {
		return domain.ErrNotFound
	}
	delete(m.tasks, id)
	return nil
}
