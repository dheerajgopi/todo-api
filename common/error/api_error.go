package error

// NewAPIError creates a new application error object with the provided message and errors
func NewAPIError(message string, apiErrorBodyList ...*APIErrorBody) *APIError {
	return &APIError{
		Message: message,
		Body:    apiErrorBodyList,
	}
}

// APIError implements error
type APIError struct {
	Message string
	Body    []*APIErrorBody
}

func (ae *APIError) Error() string {
	return ae.Message
}

// APIErrorBody represents json structure of error in API response
type APIErrorBody struct {
	Message string `json:"message,omitempty"`
	Target  string `json:"target,omitempty"`
}

// NewAPIErrorBody returns a new instance of APIError with the
// message and target.
func NewAPIErrorBody(message string, target string) *APIErrorBody {
	return &APIErrorBody{
		Message: message,
		Target:  target,
	}
}
