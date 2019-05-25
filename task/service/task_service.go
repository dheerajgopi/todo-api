package service

import (
	"context"

	"github.com/dheerajgopi/todo-api/models"
	"github.com/dheerajgopi/todo-api/task"
)

type taskService struct {
	taskRepo task.Repository
}

// New returns a new object implementing task.Service interface
func New(repo task.Repository) task.Service {
	return &taskService{
		taskRepo: repo,
	}
}

// Create creates a new task
func (service *taskService) Create(ctx context.Context, newTask *models.Task) error {
	return service.taskRepo.Create(ctx, newTask)
}

// List returns tasks created by an user
func (service *taskService) List(ctx context.Context, userID int64) ([]*models.Task, error) {
	return service.taskRepo.GetAllByUserID(ctx, userID)
}
