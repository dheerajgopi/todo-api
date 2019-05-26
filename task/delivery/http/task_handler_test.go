package http_test

import (
	"encoding/json"
	"errors"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/dheerajgopi/todo-api/common"
	"github.com/dheerajgopi/todo-api/config"
	"github.com/dheerajgopi/todo-api/models"
	"github.com/dheerajgopi/todo-api/task"
	_taskHandler "github.com/dheerajgopi/todo-api/task/delivery/http"
	mock "github.com/dheerajgopi/todo-api/task/mock"
	"github.com/golang/mock/gomock"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestCreateWithNoRequestBody(t *testing.T) {
	assert := assert.New(t)
	mockCtrl := gomock.NewController(t)
	mockService := mock.NewService(mockCtrl)
	handler := setupHandler(mockService)
	reqCtx := setupRequestContext(handler.App)
	req := httptest.NewRequest("POST", "/tasks", strings.NewReader(""))

	status, data, err := handler.Create(httptest.NewRecorder(), req, reqCtx)

	assert.Equal(400, status)
	assert.Nil(data)
	assert.Error(err)
	assert.Equal(1, len(err.Body))
	assert.Equal("Invalid request body", err.Body[0].Message)
}

func TestCreateWithBlankTitle(t *testing.T) {
	payload, _ := json.Marshal(&_taskHandler.CreateTaskRequest{
		Title:       " ",
		Description: "test description",
	})

	assert := assert.New(t)
	mockCtrl := gomock.NewController(t)
	mockService := mock.NewService(mockCtrl)
	handler := setupHandler(mockService)
	reqCtx := setupRequestContext(handler.App)
	req := httptest.NewRequest("POST", "/tasks", strings.NewReader(string(payload)))

	status, data, err := handler.Create(httptest.NewRecorder(), req, reqCtx)

	assert.Equal(400, status)
	assert.Nil(data)
	assert.Error(err)
	assert.Equal(1, len(err.Body))
	assert.Equal("Non-empty value is required", err.Body[0].Message)
	assert.Equal("title", err.Body[0].Target)
}

func TestCreateWithServerError(t *testing.T) {
	reqBody := &_taskHandler.CreateTaskRequest{
		Title:       "test title",
		Description: "test description",
	}

	payload, _ := json.Marshal(reqBody)

	assert := assert.New(t)
	mockCtrl := gomock.NewController(t)
	mockService := mock.NewService(mockCtrl)
	handler := setupHandler(mockService)
	reqCtx := setupRequestContext(handler.App)
	req := httptest.NewRequest("POST", "/tasks", strings.NewReader(string(payload)))

	mockService.
		EXPECT().
		Create(gomock.Any(), gomock.Any()).
		Return(errors.New("server error")).
		Times(1)

	status, data, err := handler.Create(httptest.NewRecorder(), req, reqCtx)

	assert.Equal(500, status)
	assert.Nil(data)
	assert.NotNil(err)
	assert.Equal(1, len(err.Body))
	assert.Equal("Internal server error", err.Body[0].Message)
}

func TestCreate(t *testing.T) {
	reqBody := &_taskHandler.CreateTaskRequest{
		Title:       "test title",
		Description: "test description",
	}

	payload, _ := json.Marshal(reqBody)

	assert := assert.New(t)
	mockCtrl := gomock.NewController(t)
	mockService := mock.NewService(mockCtrl)
	handler := setupHandler(mockService)
	reqCtx := setupRequestContext(handler.App)
	req := httptest.NewRequest("POST", "/tasks", strings.NewReader(string(payload)))

	mockService.
		EXPECT().
		Create(gomock.Any(), gomock.Any()).
		Return(nil).
		Times(1)

	status, data, err := handler.Create(httptest.NewRecorder(), req, reqCtx)

	responseData := data.(*_taskHandler.CreateTaskResponse)

	assert.Equal(201, status)
	assert.NotNil(data)
	assert.Nil(err)
	assert.Equal(reqBody.Title, responseData.Task.Title)
	assert.Equal(reqBody.Description, responseData.Task.Description)
	assert.Equal(false, responseData.Task.IsComplete)
}

func TestListWithServerError(t *testing.T) {
	assert := assert.New(t)
	mockCtrl := gomock.NewController(t)
	mockService := mock.NewService(mockCtrl)
	handler := setupHandler(mockService)
	reqCtx := setupRequestContext(handler.App)

	req := httptest.NewRequest("GET", "/tasks", nil)

	mockService.
		EXPECT().
		List(gomock.Any(), gomock.Any()).
		Return(nil, errors.New("server error")).
		Times(1)

	status, data, err := handler.List(httptest.NewRecorder(), req, reqCtx)

	assert.Equal(500, status)
	assert.Nil(data)
	assert.NotNil(err)
	assert.Equal(1, len(err.Body))
	assert.Equal("Internal server error", err.Body[0].Message)
}

func TestList(t *testing.T) {
	now := time.Now()
	assert := assert.New(t)
	mockCtrl := gomock.NewController(t)
	mockService := mock.NewService(mockCtrl)
	handler := setupHandler(mockService)
	reqCtx := setupRequestContext(handler.App)

	req := httptest.NewRequest("GET", "/tasks", nil)

	expectedData := make([]*models.Task, 0)
	expectedData = append(expectedData, &models.Task{
		ID:          1,
		Title:       "test title",
		Description: "test description",
		CreatedBy: &models.User{
			ID: reqCtx.UserID,
		},
		IsComplete: false,
		CreatedAt:  now,
		UpdatedAt:  now,
	})

	mockService.
		EXPECT().
		List(gomock.Any(), gomock.Any()).
		Return(expectedData, nil).
		Times(1)

	status, data, err := handler.List(httptest.NewRecorder(), req, reqCtx)

	actualData := data.(*_taskHandler.ListTaskResponse)

	assert.Equal(200, status)
	assert.Nil(err)
	assert.NotNil(data)
	assert.Equal(1, len(actualData.Tasks))
	assert.Equal("test title", actualData.Tasks[0].Title)
	assert.Equal("test description", actualData.Tasks[0].Description)
	assert.Equal(false, actualData.Tasks[0].IsComplete)
	assert.Equal(now, actualData.Tasks[0].CreatedAt)
	assert.Equal(now, actualData.Tasks[0].UpdatedAt)
}

func setupHandler(mockService task.Service) *_taskHandler.TaskHandler {
	app := &common.App{
		Logger: logrus.New(),
		Config: &config.Config{},
	}

	handler := &_taskHandler.TaskHandler{
		TaskService: mockService,
		App:         app,
	}

	return handler
}

func setupRequestContext(app *common.App) *common.RequestContext {
	reqCtx := &common.RequestContext{
		RequestID: "dummyRequestID",
		LogEntry: app.Logger.WithFields(
			logrus.Fields{},
		),
	}

	return reqCtx
}
