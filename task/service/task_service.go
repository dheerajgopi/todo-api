package service

import (
	"context"

	"github.com/sirupsen/logrus"

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

func (service *taskService) Create(ctx context.Context, newTask *models.Task) error {
	err := service.taskRepo.Create(ctx, newTask)

	if err != nil {
		logrus.Error(err)
		return err
	}

	return nil
}
