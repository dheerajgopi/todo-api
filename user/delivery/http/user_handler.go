package http

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/dheerajgopi/todo-api/common"
	"github.com/dheerajgopi/todo-api/models"

	"github.com/dheerajgopi/todo-api/user"
	"github.com/gorilla/mux"
)

// UserHandler represents HTTP handler for users
type UserHandler struct {
	UserService user.Service
	App         *common.App
}

// New creates new HTTP handler for user
func New(router *mux.Router, service user.Service, app *common.App) {
	handler := &UserHandler{
		UserService: service,
		App:         app,
	}

	router.HandleFunc("/users", app.CreateHandler(handler.Create)).Methods("POST")
}

// Create will store new user
func (handler *UserHandler) Create(res http.ResponseWriter, req *http.Request, reqCtx *common.RequestContext) (int, interface{}, *common.AppError) {
	timeoutContext, cancel := context.WithTimeout(context.TODO(), handler.App.Config.RequestTimeout)
	defer cancel()
	defer req.Body.Close()

	decoder := json.NewDecoder(req.Body)
	var createUserReqBody CreateUserRequest
	err := decoder.Decode(&createUserReqBody)

	if err != nil {
		reqCtx.AddLogMessage("Invalid request body")
		requestBodyError := make([]*common.APIError, 0)
		requestBodyError = append(requestBodyError, &common.APIError{
			Message: "Invalid request body",
		})

		apiError := &common.AppError{
			Errors: requestBodyError,
		}

		return http.StatusBadRequest, nil, apiError
	}

	validationErrors := createUserReqBody.Validate()

	if len(validationErrors) > 0 {
		reqCtx.AddLogMessage("validation error")

		apiError := &common.AppError{
			Errors: validationErrors,
		}

		return http.StatusBadRequest, nil, apiError
	}

	now := time.Now()

	newUser := models.User{
		Name:      createUserReqBody.Name,
		Email:     createUserReqBody.Email,
		Passwd:    createUserReqBody.Password,
		IsActive:  true,
		CreatedAt: now,
		UpdatedAt: now,
	}

	switch err = handler.UserService.Create(timeoutContext, &newUser); err.(type) {
	case nil:
		break
	case *common.DataConflictError:
		dataConflictErr, _ := err.(*common.DataConflictError)
		conflictError := make([]*common.APIError, 0)
		conflictError = append(conflictError, &common.APIError{
			Message: "Conflicting data",
			Target:  dataConflictErr.Field,
		})

		apiError := &common.AppError{
			Message: dataConflictErr.Error(),
			Errors:  conflictError,
		}

		return http.StatusConflict, nil, apiError
	default:
		serverError := make([]*common.APIError, 0)
		serverError = append(serverError, &common.APIError{
			Message: "Internal server error",
		})

		apiError := &common.AppError{
			Message: err.Error(),
			Errors:  serverError,
		}

		return http.StatusInternalServerError, nil, apiError
	}

	userData := UserData{
		ID:        newUser.ID,
		Name:      newUser.Name,
		Email:     newUser.Email,
		IsActive:  newUser.IsActive,
		CreatedAt: newUser.CreatedAt,
		UpdatedAt: newUser.UpdatedAt,
	}

	return http.StatusCreated, userData, nil
}
