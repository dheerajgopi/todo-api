package http

// CreateTaskResponse represents response for POST /tasks API
type CreateTaskResponse struct {
	Task *TaskData `json:"task"`
}

// ListTaskResponse represents response for GET /tasks API
type ListTaskResponse struct {
	Tasks []*TaskData `json:"tasks"`
}
