package service

import (
	"errors"
	"golang_training/04-cli-applications/project-gomanager/internal/model"
	"testing"
)

// MockTaskRepository implements repository.TaskRepository in-memory for testing.
type MockTaskRepository struct {
	tasks  map[int]model.Task
	nextID int
}

func NewMockTaskRepository() *MockTaskRepository {
	return &MockTaskRepository{
		tasks:  make(map[int]model.Task),
		nextID: 1,
	}
}

func (m *MockTaskRepository) GetAll() ([]model.Task, error) {
	var list []model.Task
	for _, t := range m.tasks {
		list = append(list, t)
	}
	return list, nil
}

func (m *MockTaskRepository) GetByID(id int) (*model.Task, error) {
	t, ok := m.tasks[id]
	if !ok {
		return nil, errors.New("not found")
	}
	return &t, nil
}

func (m *MockTaskRepository) Create(task model.Task) (*model.Task, error) {
	task.ID = m.nextID
	m.nextID++
	m.tasks[task.ID] = task
	return &task, nil
}

func (m *MockTaskRepository) Delete(id int) error {
	if _, ok := m.tasks[id]; !ok {
		return errors.New("not found")
	}
	delete(m.tasks, id)
	return nil
}

// MockConfigRepository implements repository.ConfigRepository in-memory.
type MockConfigRepository struct {
	cfg *model.Config
}

func (m *MockConfigRepository) GetConfig() (*model.Config, error) {
	return m.cfg, nil
}

func (m *MockConfigRepository) SaveConfig(cfg *model.Config) error {
	m.cfg = cfg
	return nil
}

func TestTaskService_CreateTask(t *testing.T) {
	mockRepo := NewMockTaskRepository()
	mockCfg := &MockConfigRepository{}
	svc := NewTaskService(mockRepo, mockCfg)

	tests := []struct {
		name    string
		title   string
		wantErr bool
	}{
		{
			name:    "Valid task",
			title:   "Build Dockerfile",
			wantErr: false,
		},
		{
			name:    "Empty task title",
			title:   "   ",
			wantErr: true,
		},
		{
			name:    "Too short title",
			title:   "ab",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task, err := svc.CreateTask(tt.title)
			if (err != nil) != tt.wantErr {
				t.Fatalf("CreateTask() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && task.Title != tt.title {
				t.Errorf("expected title %q, got %q", tt.title, task.Title)
			}
		})
	}
}

func TestTaskService_DeleteTask(t *testing.T) {
	mockRepo := NewMockTaskRepository()
	mockCfg := &MockConfigRepository{}
	svc := NewTaskService(mockRepo, mockCfg)

	task, _ := svc.CreateTask("Test Task")

	// Delete existing
	if err := svc.DeleteTask(task.ID); err != nil {
		t.Fatalf("unexpected delete error: %v", err)
	}

	// Delete non-existent
	if err := svc.DeleteTask(999); err == nil {
		t.Errorf("expected error when deleting non-existent task, got nil")
	}

	// Delete invalid ID
	if err := svc.DeleteTask(-1); err == nil {
		t.Errorf("expected error for negative task ID, got nil")
	}
}
