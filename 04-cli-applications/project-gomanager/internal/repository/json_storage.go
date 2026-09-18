package repository

import (
	"encoding/json"
	"errors"
	"fmt"
	"golang_training/04-cli-applications/project-gomanager/internal/model"
	"os"
	"path/filepath"
	"sync"
)

var (
	ErrTaskNotFound = errors.New("task not found")
)

// TaskRepository defines the persistence contract for tasks.
type TaskRepository interface {
	GetAll() ([]model.Task, error)
	GetByID(id int) (*model.Task, error)
	Create(task model.Task) (*model.Task, error)
	Delete(id int) error
}

// ConfigRepository defines the persistence contract for app configuration.
type ConfigRepository interface {
	GetConfig() (*model.Config, error)
	SaveConfig(cfg *model.Config) error
}

// FileStorage implements TaskRepository and ConfigRepository via local JSON files.
type FileStorage struct {
	baseDir    string
	tasksFile  string
	configFile string
	mu         sync.RWMutex
}

// NewFileStorage creates or opens storage in the target base directory.
func NewFileStorage(baseDir string) (*FileStorage, error) {
	if err := os.MkdirAll(baseDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create storage directory: %w", err)
	}

	return &FileStorage{
		baseDir:    baseDir,
		tasksFile:  filepath.Join(baseDir, "tasks.json"),
		configFile: filepath.Join(baseDir, "config.json"),
	}, nil
}

func (s *FileStorage) GetAll() ([]model.Task, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if _, err := os.Stat(s.tasksFile); os.IsNotExist(err) {
		return []model.Task{}, nil
	}

	data, err := os.ReadFile(s.tasksFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read tasks file: %w", err)
	}

	if len(data) == 0 {
		return []model.Task{}, nil
	}

	var tasks []model.Task
	if err := json.Unmarshal(data, &tasks); err != nil {
		return nil, fmt.Errorf("failed to unmarshal tasks: %w", err)
	}
	return tasks, nil
}

func (s *FileStorage) GetByID(id int) (*model.Task, error) {
	tasks, err := s.GetAll()
	if err != nil {
		return nil, err
	}
	for _, t := range tasks {
		if t.ID == id {
			return &t, nil
		}
	}
	return nil, ErrTaskNotFound
}

func (s *FileStorage) Create(task model.Task) (*model.Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var tasks []model.Task
	if _, err := os.Stat(s.tasksFile); err == nil {
		data, err := os.ReadFile(s.tasksFile)
		if err == nil && len(data) > 0 {
			_ = json.Unmarshal(data, &tasks)
		}
	}

	// Determine next ID
	maxID := 0
	for _, t := range tasks {
		if t.ID > maxID {
			maxID = t.ID
		}
	}
	task.ID = maxID + 1

	tasks = append(tasks, task)

	encoded, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to encode tasks: %w", err)
	}

	if err := os.WriteFile(s.tasksFile, encoded, 0644); err != nil {
		return nil, fmt.Errorf("failed to persist tasks: %w", err)
	}

	return &task, nil
}

func (s *FileStorage) Delete(id int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	var tasks []model.Task
	if _, err := os.Stat(s.tasksFile); os.IsNotExist(err) {
		return ErrTaskNotFound
	}

	data, err := os.ReadFile(s.tasksFile)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(data, &tasks); err != nil {
		return err
	}

	foundIndex := -1
	for i, t := range tasks {
		if t.ID == id {
			foundIndex = i
			break
		}
	}

	if foundIndex == -1 {
		return ErrTaskNotFound
	}

	// Remove from slice
	tasks = append(tasks[:foundIndex], tasks[foundIndex+1:]...)

	encoded, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(s.tasksFile, encoded, 0644)
}

func (s *FileStorage) GetConfig() (*model.Config, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if _, err := os.Stat(s.configFile); os.IsNotExist(err) {
		return &model.Config{
			StoragePath: s.baseDir,
			Version:     "1.0.0",
		}, nil
	}

	data, err := os.ReadFile(s.configFile)
	if err != nil {
		return nil, err
	}

	var cfg model.Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func (s *FileStorage) SaveConfig(cfg *model.Config) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.configFile, data, 0644)
}
