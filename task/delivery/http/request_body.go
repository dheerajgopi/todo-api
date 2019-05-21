package http

import (
	"strings"
	"time"

	todoErr "github.com/dheerajgopi/todo-api/common/error"
)

// TaskData represents json structure for user
type TaskData struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	IsComplete  bool      `json:"isComplete"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// CreateTaskRequest represents request body for POST /tasks API
type CreateTaskRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

// ValidateAndBuild validates the request body for POST /tasks API
func (body *CreateTaskRequest) ValidateAndBuild() []*todoErr.APIErrorBody {
	trimmedTitle := strings.TrimSpace(body.Title)
	trimmedDescription := strings.TrimSpace(body.Description)

	validationErrors := make([]*todoErr.APIErrorBody, 0)

	if trimmedTitle == "" {
		validationErrors = append(validationErrors, &todoErr.APIErrorBody{
			Message: "Non-empty value is required",
			Target:  "title",
		})
	}

	body.Title = trimmedTitle
	body.Description = trimmedDescription

	return validationErrors
}
