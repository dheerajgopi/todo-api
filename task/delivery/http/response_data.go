package http

// CreateTaskResponse represents response for POST /tasks API
type CreateTaskResponse struct {
	Task *TaskData `json:"task"`
}

// ListTaskResponse represents response for GET /tasks API
type ListTaskResponse struct {
	Tasks    []*TaskData `json:"tasks"`
	PageInfo Cursor      `json:"pageInfo"`
}

// Cursor is the pointer to next page (keyset pagination)
type Cursor struct {
	Next    string `json:"next"`
	HasNext bool   `json:"hasNext"`
}
