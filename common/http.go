package common

import todoErr "github.com/dheerajgopi/todo-api/common/error"

// APIResponse represents json structure of API response
type APIResponse struct {
	Status int                     `json:"status"`
	Errors []*todoErr.APIErrorBody `json:"errors"`
	Data   interface{}             `json:"data"`
}
