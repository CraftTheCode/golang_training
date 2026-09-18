package service

import (
	"errors"
	"fmt"
	"golang_training/04-cli-applications/project-gomanager/internal/model"
	"golang_training/04-cli-applications/project-gomanager/internal/repository"
	"strings"
	"time"
)

// TaskService orchestrates task-related use cases and domain validation.
type TaskService struct {
	repo       repository.TaskRepository
	configRepo repository.ConfigRepository
}

func NewTaskService(repo repository.TaskRepository, configRepo repository.ConfigRepository) *TaskService {
	return &TaskService{
		repo:       repo,
		configRepo: configRepo,
	}
}

func (s *TaskService) CreateTask(title string) (*model.Task, error) {
	cleanTitle := strings.TrimSpace(title)
	if cleanTitle == "" {
		return nil, errors.New("task title cannot be empty")
	}
	if len(cleanTitle) < 3 {
		return nil, errors.New("task title must be at least 3 characters long")
	}

	task := model.Task{
		Title:     cleanTitle,
		Status:    "pending",
		CreatedAt: time.Now().UTC(),
	}

	return s.repo.Create(task)
}

func (s *TaskService) ListTasks() ([]model.Task, error) {
	return s.repo.GetAll()
}

func (s *TaskService) DeleteTask(id int) error {
	if id <= 0 {
		return fmt.Errorf("invalid task id: %d (must be positive)", id)
	}
	return s.repo.Delete(id)
}

func (s *TaskService) Initialize() error {
	defaultCfg := &model.Config{
		StoragePath: ".gomanager",
		Version:     "1.0.0",
	}
	return s.configRepo.SaveConfig(defaultCfg)
}

func (s *TaskService) GetConfig() (*model.Config, error) {
	return s.configRepo.GetConfig()
}
