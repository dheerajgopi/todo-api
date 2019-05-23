package http

// CreateTaskResponse represents response for POST /task API
type CreateTaskResponse struct {
	Task *TaskData `json:"task"`
}
