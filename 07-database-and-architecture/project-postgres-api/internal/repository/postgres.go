package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"golang_training/07-database-and-architecture/project-postgres-api/internal/domain"
	"time"
)

// PostgresRepository implements domain.TaskRepository against a PostgreSQL database.
type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) GetAll(ctx context.Context) ([]domain.Task, error) {
	query := `SELECT id, title, completed, created_at FROM tasks ORDER BY id ASC`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query tasks: %w", err)
	}
	defer rows.Close() // Mandatory to prevent connection pool leaks!

	var tasks []domain.Task
	for rows.Next() {
		var t domain.Task
		if err := rows.Scan(&t.ID, &t.Title, &t.Completed, &t.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan task row: %w", err)
		}
		tasks = append(tasks, t)
	}

	// Always check rows.Err() after iteration completes!
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error during row iteration: %w", err)
	}

	return tasks, nil
}

func (r *PostgresRepository) GetByID(ctx context.Context, id int) (*domain.Task, error) {
	query := `SELECT id, title, completed, created_at FROM tasks WHERE id = $1`

	var t domain.Task
	err := r.db.QueryRowContext(ctx, query, id).Scan(&t.ID, &t.Title, &t.Completed, &t.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("failed to query task by id: %w", err)
	}

	return &t, nil
}

func (r *PostgresRepository) Create(ctx context.Context, title string) (*domain.Task, error) {
	query := `
		INSERT INTO tasks (title, completed, created_at)
		VALUES ($1, false, $2)
		RETURNING id, created_at
	`

	now := time.Now().UTC()
	var id int
	err := r.db.QueryRowContext(ctx, query, title, now).Scan(&id, &now)
	if err != nil {
		return nil, fmt.Errorf("failed to insert task: %w", err)
	}

	return &domain.Task{
		ID:        id,
		Title:     title,
		Completed: false,
		CreatedAt: now,
	}, nil
}

func (r *PostgresRepository) UpdateStatus(ctx context.Context, id int, completed bool) error {
	query := `UPDATE tasks SET completed = $1 WHERE id = $2`

	res, err := r.db.ExecContext(ctx, query, completed, id)
	if err != nil {
		return fmt.Errorf("failed to update task: %w", err)
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *PostgresRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM tasks WHERE id = $1`

	res, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete task: %w", err)
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return domain.ErrNotFound
	}
	return nil
}
