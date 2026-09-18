package model

import "time"

// Task represents an operational task item.
type Task struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Status    string    `json:"status"` // "pending", "in_progress", "completed"
	CreatedAt time.Time `json:"created_at"`
}

// Config represents persistent application settings.
type Config struct {
	StoragePath string `json:"storage_path"`
	Version     string `json:"version"`
}
