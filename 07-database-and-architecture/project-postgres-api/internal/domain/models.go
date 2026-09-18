package domain

import (
	"context"
	"errors"
	"time"
)

var (
	ErrNotFound = errors.New("entity not found")
)

// Task represents our core domain entity.
type Task struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Completed bool      `json:"completed"`
	CreatedAt time.Time `json:"created_at"`
}

// TaskRepository defines the repository interface.
// Notice that the domain knows NOTHING about PostgreSQL or database/sql!
type TaskRepository interface {
	GetAll(ctx context.Context) ([]Task, error)
	GetByID(ctx context.Context, id int) (*Task, error)
	Create(ctx context.Context, title string) (*Task, error)
	UpdateStatus(ctx context.Context, id int, completed bool) error
	Delete(ctx context.Context, id int) error
}
