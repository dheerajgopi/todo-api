package models

import "time"

// User represents user table
type User struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name" validate:"required"`
	IsActive  bool      `json:"isActive"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
