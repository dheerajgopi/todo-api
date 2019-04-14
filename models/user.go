package models

import "time"

// User represents user table
type User struct {
	ID        int64
	Name      string
	Email     string
	Passwd    string
	IsActive  bool
	CreatedAt time.Time
	UpdatedAt time.Time
}
