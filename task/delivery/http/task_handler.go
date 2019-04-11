package http

import (
	"fmt"
	"net/http"

	"github.com/dheerajgopi/todo-api/task"

	"github.com/gorilla/mux"
)

// TaskHandler represents HTTP handler for tasks
type TaskHandler struct {
	TaskService task.Service
}

func dummyTask(res http.ResponseWriter, req *http.Request) {
	fmt.Println("tasks")
}

// New creates new HTTP handler for task
func New(router *mux.Router, service task.Service) {
	handler := &TaskHandler{
		TaskService: service,
	}

	router.HandleFunc("/tasks", dummyTask).Methods("POST")
}
