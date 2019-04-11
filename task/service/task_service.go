package service

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/dheerajgopi/todo-api/models"
	"github.com/dheerajgopi/todo-api/task"
)

type taskService struct {
	taskRepo       task.Repository
	contextTimeout time.Duration
}

// New returns a new object implementing task.Service interface
func New(repo task.Repository, timeout time.Duration) task.Service {
	return &taskService{
		taskRepo:       repo,
		contextTimeout: timeout,
	}
}

func (service *taskService) Create(ctx context.Context, newTask *models.Task) error {
	timeoutContext, cancel := context.WithTimeout(ctx, service.contextTimeout)
	defer cancel()

	err := service.taskRepo.Create(timeoutContext, newTask)

	if err != nil {
		logrus.Error(err)
		return err
	}

	return nil
}
