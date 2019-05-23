package http

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/dheerajgopi/todo-api/common"
	todoErr "github.com/dheerajgopi/todo-api/common/error"
	"github.com/dheerajgopi/todo-api/common/middlewares"
	"github.com/dheerajgopi/todo-api/models"
	"github.com/dheerajgopi/todo-api/task"
	"github.com/gorilla/mux"
)

// TaskHandler represents HTTP handler for tasks
type TaskHandler struct {
	TaskService task.Service
	App         *common.App
}

// New creates new HTTP handler for task
func New(router *mux.Router, service task.Service, app *common.App) {
	handler := &TaskHandler{
		TaskService: service,
		App:         app,
	}

	jwtMiddleware := middlewares.JwtValidator(app.Config.Auth.Jwt.Secret)

	router.HandleFunc("/tasks", app.CreateHandler(jwtMiddleware(handler.Create))).Methods("POST")
}

// Create will store new task
func (handler *TaskHandler) Create(res http.ResponseWriter, req *http.Request, reqCtx *common.RequestContext) (int, interface{}, *todoErr.APIError) {
	defer req.Body.Close()

	decoder := json.NewDecoder(req.Body)
	var createTaskReqBody CreateTaskRequest
	err := decoder.Decode(&createTaskReqBody)

	if err != nil {
		apiError := todoErr.NewAPIError("", &todoErr.APIErrorBody{
			Message: "Invalid request body",
		})

		return http.StatusBadRequest, nil, apiError
	}

	validationErrors := createTaskReqBody.ValidateAndBuild()

	if len(validationErrors) > 0 {
		apiError := todoErr.NewAPIError("", validationErrors...)

		return http.StatusBadRequest, nil, apiError
	}

	now := time.Now()

	newTask := &models.Task{
		Title:       createTaskReqBody.Title,
		Description: createTaskReqBody.Description,
		CreatedBy: &models.User{
			ID: reqCtx.UserID,
		},
		IsComplete: false,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	switch err = handler.TaskService.Create(context.TODO(), newTask); err.(type) {
	case nil:
		break
	default:
		apiError := todoErr.NewAPIError(err.Error(), &todoErr.APIErrorBody{
			Message: "Internal server error",
		})

		return http.StatusInternalServerError, nil, apiError
	}

	taskData := &TaskData{
		ID:          newTask.ID,
		Title:       newTask.Title,
		Description: newTask.Description,
		IsComplete:  newTask.IsComplete,
		CreatedAt:   newTask.CreatedAt,
		UpdatedAt:   newTask.UpdatedAt,
	}

	responseData := &CreateTaskResponse{
		Task: taskData,
	}

	return http.StatusCreated, responseData, nil
}
