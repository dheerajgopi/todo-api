package error

import "fmt"

// DataConflictError indicates conflicting data for a resource
type DataConflictError struct {
	Resource string
	Field    string
}

func (dce *DataConflictError) Error() string {
	return fmt.Sprintf("Conflicting data for %s field in %s resource", dce.Field, dce.Resource)
}
