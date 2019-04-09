package user

import (
	"context"

	"github.com/dheerajgopi/todo-api/models"
)

// Repository represents user's repository contract
type Repository interface {
	GetByID(ctx context.Context, id int64) (*models.User, error)
}
