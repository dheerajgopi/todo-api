package task

import (
	"context"

	"github.com/dheerajgopi/todo-api/common"
	"github.com/dheerajgopi/todo-api/models"
)

// Service represents task service contract
type Service interface {
	Create(ctx context.Context, newTask *models.Task) error
	List(ctx context.Context, userID int64, page *common.Page) ([]*models.Task, error)
}
