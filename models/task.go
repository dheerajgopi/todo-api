package models

import (
	"time"
)

// Task represents task table
type Task struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title" validate:"required"`
	Description string    `json:"description"`
	CreatedBy   User      `json:"createdBy" validate:"required"`
	IsComplete  bool      `json:"isComplete"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}
