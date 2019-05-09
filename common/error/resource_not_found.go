package error

import "fmt"

// ResourceNotFoundError indicates missing resource
type ResourceNotFoundError struct {
	Resource string
}

func (rnf *ResourceNotFoundError) Error() string {
	return fmt.Sprintf("%s not found", rnf.Resource)
}
