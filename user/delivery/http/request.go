package http

import (
	"strings"
	"time"

	todoErr "github.com/dheerajgopi/todo-api/common/error"
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
func (body *CreateUserRequest) Validate() []*todoErr.APIErrorBody {
	name := body.Name
	email := body.Email
	pswd := body.Password

	validator := validator.New()
	validationErrors := make([]*todoErr.APIErrorBody, 0)

	if strings.TrimSpace(name) == "" {
		validationErrors = append(validationErrors, &todoErr.APIErrorBody{
			Message: "Non-empty value is required",
			Target:  "name",
		})
	}

	if strings.TrimSpace(email) == "" {
		validationErrors = append(validationErrors, &todoErr.APIErrorBody{
			Message: "Non-empty value is required",
			Target:  "email",
		})
	} else {
		emailErr := validator.Var(email, "email")

		if emailErr != nil {
			validationErrors = append(validationErrors, &todoErr.APIErrorBody{
				Message: "Invalid value",
				Target:  "email",
			})
		}
	}

	trimmedPswd := strings.TrimSpace(pswd)

	if trimmedPswd == "" {
		validationErrors = append(validationErrors, &todoErr.APIErrorBody{
			Message: "Non-empty value is required",
			Target:  "password",
		})
	} else if len(trimmedPswd) < 6 {
		validationErrors = append(validationErrors, &todoErr.APIErrorBody{
			Message: "Length should be 6 or more",
			Target:  "password",
		})
	}

	return validationErrors
}

// LoginRequest represents request body for POST /login API
type LoginRequest struct {
	Email  string `json:"email"`
	Passwd string `json:"password"`
}

// ValidateAndBuild validates the request body for POST /login API
func (body *LoginRequest) ValidateAndBuild() []*todoErr.APIErrorBody {
	email := body.Email
	passwd := body.Passwd

	validator := validator.New()
	validationErrors := make([]*todoErr.APIErrorBody, 0)

	trimmedEmail := strings.TrimSpace(email)
	trimmedPassword := strings.TrimSpace(passwd)

	if trimmedEmail == "" {
		validationErrors = append(validationErrors, &todoErr.APIErrorBody{
			Message: "Non-empty value is required",
			Target:  "email",
		})
	} else {
		emailErr := validator.Var(trimmedEmail, "email")

		if emailErr != nil {
			validationErrors = append(validationErrors, &todoErr.APIErrorBody{
				Message: "Invalid value",
				Target:  "email",
			})
		}
	}

	if trimmedPassword == "" {
		validationErrors = append(validationErrors, &todoErr.APIErrorBody{
			Message: "Non-empty value is required",
			Target:  "password",
		})
	}

	body.Email = trimmedEmail
	body.Passwd = trimmedPassword

	return validationErrors
}
