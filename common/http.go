package common

// APIError represents json structure of error in API response
type APIError struct {
	Message string `json:"message,omitempty"`
	Target  string `json:"target,omitempty"`
}

// APIResponse represents json structure of API response
type APIResponse struct {
	Status int         `json:"status"`
	Errors []*APIError `json:"errors"`
	Data   interface{} `json:"data"`
}
