package models

import (
	"time"
)

// Task represents task table
type Task struct {
	ID          int64     `json:"id"`
	Description string    `json:"description" validate:"required"`
	CreatedBy   User      `json:"user" validate:"required"`
	IsComplete  bool      `json:"isComplete"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}
