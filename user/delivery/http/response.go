package http

// CreateUserResponse represents response for POST /users API
type CreateUserResponse struct {
	User UserData `json:"user"`
}

// LoginResponse represents response for POST /login API
type LoginResponse struct {
	Token string `json:"token"`
}
