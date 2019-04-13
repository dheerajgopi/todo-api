package http

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/dheerajgopi/todo-api/models"

	"github.com/dheerajgopi/todo-api/task"

	"github.com/gorilla/mux"
)

// TaskHandler represents HTTP handler for tasks
type TaskHandler struct {
	TaskService    task.Service
	contextTimeout time.Duration
}

// New creates new HTTP handler for task
func New(router *mux.Router, service task.Service, timeout time.Duration) {
	handler := &TaskHandler{
		TaskService:    service,
		contextTimeout: timeout,
	}

	router.HandleFunc("/tasks", handler.Create).Methods("POST")
}

// Create will store new task
func (handler *TaskHandler) Create(res http.ResponseWriter, req *http.Request) {
	timeoutContext, cancel := context.WithTimeout(context.TODO(), handler.contextTimeout)
	defer cancel()

	newTask, err := validateAndBuildCreateTask(req)
	defer req.Body.Close()

	if err != nil {
		payload := make(map[string]interface{})
		statusCode := http.StatusBadRequest
		errorData := make(map[string]string)
		errorData["message"] = "invalid request body"

		payload["status"] = statusCode
		payload["error"] = errorData
		payload["data"] = make(map[string]string)

		respondWithJSON(res, statusCode, payload)

		return
	}

	var user models.User
	user.ID = int64(1)
	newTask.CreatedBy = user

	err = handler.TaskService.Create(timeoutContext, newTask)

	if err != nil {
		payload := make(map[string]interface{})
		statusCode := http.StatusBadRequest
		errorData := make(map[string]string)
		errorData["message"] = "request cannot be processed"

		payload["status"] = statusCode
		payload["error"] = errorData
		payload["data"] = make(map[string]string)

		respondWithJSON(res, statusCode, payload)

		return
	}

	payload := make(map[string]interface{})
	statusCode := http.StatusCreated

	errorData := make(map[string]string)

	data := make(map[string]interface{})
	data["task"] = newTask

	payload["status"] = statusCode
	payload["error"] = errorData
	payload["data"] = data

	respondWithJSON(res, statusCode, payload)
}

func validateAndBuildCreateTask(req *http.Request) (*models.Task, error) {
	var task models.Task
	decoder := json.NewDecoder(req.Body)

	err := decoder.Decode(&task)

	if err != nil {
		return nil, err
	}

	now := time.Now()
	task.CreatedAt = now
	task.UpdatedAt = now

	return &task, nil
}

func respondWithJSON(res http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(code)
	res.Write(response)
}
