package service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/dheerajgopi/todo-api/models"
	taskMock "github.com/dheerajgopi/todo-api/task/mock"
	"github.com/dheerajgopi/todo-api/task/service"
)

func TestCreate(t *testing.T) {
	now := time.Now()
	ctx := context.TODO()
	assert := assert.New(t)
	mockCtrl := gomock.NewController(t)

	defer mockCtrl.Finish()

	mockRepo := taskMock.NewRepository(mockCtrl)
	taskService := service.New(mockRepo)

	newTask := &models.Task{
		Title:       "testTitle",
		Description: "test description",
		CreatedBy: &models.User{
			ID: 1,
		},
		IsComplete: false,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	mockRepo.
		EXPECT().
		Create(ctx, newTask).
		Return(nil).
		Times(1)

	err := taskService.Create(ctx, newTask)

	assert.NoError(err)
	assert.Equal(newTask, newTask)
	assert.Equal("testTitle", newTask.Title)
	assert.Equal("test description", newTask.Description)
	assert.Equal(false, newTask.IsComplete)
	assert.Equal(now, newTask.CreatedAt)
	assert.Equal(now, newTask.UpdatedAt)
}

func TestCreateWithError(t *testing.T) {
	ctx := context.TODO()
	assert := assert.New(t)
	mockCtrl := gomock.NewController(t)

	defer mockCtrl.Finish()

	mockRepo := taskMock.NewRepository(mockCtrl)
	taskService := service.New(mockRepo)

	newTask := &models.Task{}

	mockRepo.
		EXPECT().
		Create(ctx, newTask).
		Return(errors.New("error")).
		Times(1)

	err := taskService.Create(ctx, newTask)

	assert.Error(err)
}
