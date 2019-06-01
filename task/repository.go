package task

import (
	"context"

	"github.com/dheerajgopi/todo-api/common"
	"github.com/dheerajgopi/todo-api/models"
)

// Repository represents task's repository contract
type Repository interface {
	GetAllByUserID(ctx context.Context, userID int64, page *common.Page) ([]*models.Task, error)
	GetByID(ctx context.Context, id int64) (*models.Task, error)
	Create(ctx context.Context, task *models.Task) error
}
