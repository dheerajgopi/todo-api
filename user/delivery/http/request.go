package http

import (
	"strings"
	"time"

	"github.com/dheerajgopi/todo-api/common"
	"gopkg.in/go-playground/validator.v9"
)

// UserData represents json structure for user
type UserData struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	IsActive  bool      `json:"isActive"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// CreateUserRequest represents request body for POST /users API
type CreateUserRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Validate validates the request body for POST /users API
func (body *CreateUserRequest) Validate() []*common.APIError {
	name := body.Name
	email := body.Email
	pswd := body.Password

	validator := validator.New()
	validationErrors := make([]*common.APIError, 0)

	if strings.TrimSpace(name) == "" {
		validationErrors = append(validationErrors, &common.APIError{
			Message: "Non-empty value is required",
			Target:  "name",
		})
	}

	if strings.TrimSpace(email) == "" {
		validationErrors = append(validationErrors, &common.APIError{
			Message: "Non-empty value is required",
			Target:  "email",
		})
	} else {
		emailErr := validator.Var(email, "email")

		if emailErr != nil {
			validationErrors = append(validationErrors, &common.APIError{
				Message: "Invalid value",
				Target:  "email",
			})
		}
	}

	trimmedPswd := strings.TrimSpace(pswd)

	if trimmedPswd == "" {
		validationErrors = append(validationErrors, &common.APIError{
			Message: "Non-empty value is required",
			Target:  "password",
		})
	} else if len(trimmedPswd) < 6 {
		validationErrors = append(validationErrors, &common.APIError{
			Message: "Length should be 6 or more",
			Target:  "password",
		})
	}

	return validationErrors
}
