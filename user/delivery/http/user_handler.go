package http

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/dheerajgopi/todo-api/common"
	todoErr "github.com/dheerajgopi/todo-api/common/error"
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
	router.HandleFunc("/login", app.CreateHandler(handler.Login)).Methods("POST")
}

// Create will store new user
func (handler *UserHandler) Create(res http.ResponseWriter, req *http.Request, reqCtx *common.RequestContext) (int, interface{}, *todoErr.APIError) {
	timeoutInSec := time.Duration(handler.App.Config.Application.RequestTimeout) * time.Second
	timeoutContext, cancel := context.WithTimeout(context.TODO(), timeoutInSec)
	defer cancel()
	defer req.Body.Close()

	decoder := json.NewDecoder(req.Body)
	var createUserReqBody CreateUserRequest
	err := decoder.Decode(&createUserReqBody)

	if err != nil {
		reqCtx.AddLogMessage("Invalid request body")
		apiError := todoErr.NewAPIError("", &todoErr.APIErrorBody{
			Message: "Invalid request body",
		})

		return http.StatusBadRequest, nil, apiError
	}

	validationErrors := createUserReqBody.Validate()

	if len(validationErrors) > 0 {
		reqCtx.AddLogMessage("validation error")
		apiError := todoErr.NewAPIError("", validationErrors...)

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
	case *todoErr.DataConflictError:
		dataConflictErr, _ := err.(*todoErr.DataConflictError)

		apiError := todoErr.NewAPIError(dataConflictErr.Error(), &todoErr.APIErrorBody{
			Message: "Conflicting data",
			Target:  dataConflictErr.Field,
		})

		return http.StatusConflict, nil, apiError
	default:
		apiError := todoErr.NewAPIError(err.Error(), &todoErr.APIErrorBody{
			Message: "Internal server error",
		})

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

// Login will validate user credentials and return a token
func (handler *UserHandler) Login(res http.ResponseWriter, req *http.Request, reqCtx *common.RequestContext) (int, interface{}, *todoErr.APIError) {
	timeoutInSec := time.Duration(handler.App.Config.Application.RequestTimeout) * time.Second
	timeoutContext, cancel := context.WithTimeout(context.TODO(), timeoutInSec)
	defer cancel()
	defer req.Body.Close()

	decoder := json.NewDecoder(req.Body)
	var loginReqBody LoginRequest
	err := decoder.Decode(&loginReqBody)

	if err != nil {
		reqCtx.AddLogMessage("Invalid request body")
		apiError := todoErr.NewAPIError("", &todoErr.APIErrorBody{
			Message: "Invalid request body",
		})

		return http.StatusBadRequest, nil, apiError
	}

	validationErrors := loginReqBody.ValidateAndBuild()

	if len(validationErrors) > 0 {
		reqCtx.AddLogMessage("validation error")
		apiError := todoErr.NewAPIError("", validationErrors...)

		return http.StatusBadRequest, nil, apiError
	}

	token, err := handler.UserService.GenerateAuthToken(
		timeoutContext,
		loginReqBody.Email,
		loginReqBody.Passwd,
		handler.App.Config.Auth.Jwt.Secret,
	)

	switch err.(type) {
	case nil:
		break
	case *todoErr.ResourceNotFoundError:
		resourceNotFoundErr, _ := err.(*todoErr.ResourceNotFoundError)

		apiError := todoErr.NewAPIError(resourceNotFoundErr.Error(), &todoErr.APIErrorBody{
			Message: "Not found",
			Target:  "user",
		})

		return http.StatusNotFound, nil, apiError
	default:
		apiError := todoErr.NewAPIError(err.Error(), &todoErr.APIErrorBody{
			Message: "Internal server error",
		})

		return http.StatusInternalServerError, nil, apiError
	}

	loginResponse := LoginResponse{
		Token: token,
	}

	return http.StatusOK, loginResponse, nil
}
