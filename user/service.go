package user

import (
	"context"

	"github.com/dheerajgopi/todo-api/models"
)

// Service represents user's service contract
type Service interface {
	Create(ctx context.Context, newUser *models.User) error
}
