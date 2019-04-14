package http

// CreateUserResponse represents response for POST /users API
type CreateUserResponse struct {
	User UserData `json:"user"`
}
