package common

import (
	"fmt"
)

// AppError implements error
type AppError struct {
	Message string
	Errors  []*APIError
}

func (ae *AppError) Error() string {
	return ae.Message
}

// DataConflictError indicates conflicting data for a resource
type DataConflictError struct {
	Resource string
	Field    string
}

func (dce *DataConflictError) Error() string {
	return fmt.Sprintf("Conflicting data for %s field in %s resource", dce.Field, dce.Resource)
}
