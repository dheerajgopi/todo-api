package user

import (
	"context"

	"github.com/dheerajgopi/todo-api/models"
)

// Repository represents user's repository contract
type Repository interface {
	GetByID(ctx context.Context, id int64) (*models.User, error)
	Create(ctx context.Context, user *models.User) error
	GetByEmail(ctx context.Context, email string) (*models.User, error)
}
