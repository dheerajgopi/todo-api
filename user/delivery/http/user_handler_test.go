package http_test

import (
	"encoding/json"
	"errors"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/dheerajgopi/todo-api/common"
	_errors "github.com/dheerajgopi/todo-api/common/error"
	"github.com/dheerajgopi/todo-api/config"
	"github.com/dheerajgopi/todo-api/user"
	_userHandler "github.com/dheerajgopi/todo-api/user/delivery/http"
	mock "github.com/dheerajgopi/todo-api/user/mock"
	"github.com/golang/mock/gomock"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestCreateWithNoRequestBody(t *testing.T) {
	assert := assert.New(t)
	mockCtrl := gomock.NewController(t)
	mockService := mock.NewService(mockCtrl)
	handler := setupHandler(mockService)
	reqCtx := setupRequestContext(handler.App)
	req := httptest.NewRequest("POST", "/users", strings.NewReader(""))

	status, data, err := handler.Create(httptest.NewRecorder(), req, reqCtx)

	assert.Equal(400, status)
	assert.Nil(data)
	assert.Error(err)
	assert.Equal(1, len(err.Body))
	assert.Equal("Invalid request body", err.Body[0].Message)
}

func TestCreateWithBlankName(t *testing.T) {
	payload, _ := json.Marshal(&_userHandler.CreateUserRequest{
		Name:     "",
		Email:    "testuser@mail.com",
		Password: "secret",
	})

	assert := assert.New(t)
	mockCtrl := gomock.NewController(t)
	mockService := mock.NewService(mockCtrl)
	handler := setupHandler(mockService)
	reqCtx := setupRequestContext(handler.App)
	req := httptest.NewRequest("POST", "/users", strings.NewReader(string(payload)))

	status, data, err := handler.Create(httptest.NewRecorder(), req, reqCtx)

	assert.Equal(400, status)
	assert.Nil(data)
	assert.Error(err)
	assert.Equal(1, len(err.Body))
	assert.Equal("Non-empty value is required", err.Body[0].Message)
	assert.Equal("name", err.Body[0].Target)
}

func TestCreateWithBlankEmail(t *testing.T) {
	payload, _ := json.Marshal(&_userHandler.CreateUserRequest{
		Name:     "testuser",
		Email:    "",
		Password: "secret",
	})

	assert := assert.New(t)
	mockCtrl := gomock.NewController(t)
	mockService := mock.NewService(mockCtrl)
	handler := setupHandler(mockService)
	reqCtx := setupRequestContext(handler.App)
	req := httptest.NewRequest("POST", "/users", strings.NewReader(string(payload)))

	status, data, err := handler.Create(httptest.NewRecorder(), req, reqCtx)

	assert.Equal(400, status)
	assert.Nil(data)
	assert.Error(err)
	assert.Equal(1, len(err.Body))
	assert.Equal("Non-empty value is required", err.Body[0].Message)
	assert.Equal("email", err.Body[0].Target)
}

func TestCreateWithInvalidEmail(t *testing.T) {
	payload, _ := json.Marshal(&_userHandler.CreateUserRequest{
		Name:     "testuser",
		Email:    "testuser@",
		Password: "secret",
	})

	assert := assert.New(t)
	mockCtrl := gomock.NewController(t)
	mockService := mock.NewService(mockCtrl)
	handler := setupHandler(mockService)
	reqCtx := setupRequestContext(handler.App)
	req := httptest.NewRequest("POST", "/users", strings.NewReader(string(payload)))

	status, data, err := handler.Create(httptest.NewRecorder(), req, reqCtx)

	assert.Equal(400, status)
	assert.Nil(data)
	assert.Error(err)
	assert.Equal(1, len(err.Body))
	assert.Equal("Invalid value", err.Body[0].Message)
	assert.Equal("email", err.Body[0].Target)
}

func TestCreateWithBlankPassword(t *testing.T) {
	payload, _ := json.Marshal(&_userHandler.CreateUserRequest{
		Name:     "testuser",
		Email:    "testuser@mail.com",
		Password: "",
	})

	assert := assert.New(t)
	mockCtrl := gomock.NewController(t)
	mockService := mock.NewService(mockCtrl)
	handler := setupHandler(mockService)
	reqCtx := setupRequestContext(handler.App)
	req := httptest.NewRequest("POST", "/users", strings.NewReader(string(payload)))

	status, data, err := handler.Create(httptest.NewRecorder(), req, reqCtx)

	assert.Equal(400, status)
	assert.Nil(data)
	assert.Error(err)
	assert.Equal(1, len(err.Body))
	assert.Equal("Non-empty value is required", err.Body[0].Message)
	assert.Equal("password", err.Body[0].Target)
}

func TestCreateWithInvalidPasswordLength(t *testing.T) {
	payload, _ := json.Marshal(&_userHandler.CreateUserRequest{
		Name:     "testuser",
		Email:    "testuser@mail.com",
		Password: "pass",
	})

	assert := assert.New(t)
	mockCtrl := gomock.NewController(t)
	mockService := mock.NewService(mockCtrl)
	handler := setupHandler(mockService)
	reqCtx := setupRequestContext(handler.App)
	req := httptest.NewRequest("POST", "/users", strings.NewReader(string(payload)))

	status, data, err := handler.Create(httptest.NewRecorder(), req, reqCtx)

	assert.Equal(400, status)
	assert.Nil(data)
	assert.Error(err)
	assert.Equal(1, len(err.Body))
	assert.Equal("Length should be 6 or more", err.Body[0].Message)
	assert.Equal("password", err.Body[0].Target)
}

func TestCreateWithDataConflict(t *testing.T) {
	reqBody := &_userHandler.CreateUserRequest{
		Name:     "testuser",
		Email:    "testuser@mail.com",
		Password: "secret",
	}

	payload, _ := json.Marshal(reqBody)

	assert := assert.New(t)
	mockCtrl := gomock.NewController(t)
	mockService := mock.NewService(mockCtrl)
	handler := setupHandler(mockService)
	handler.UserService = mockService
	reqCtx := setupRequestContext(handler.App)
	req := httptest.NewRequest("POST", "/users", strings.NewReader(string(payload)))

	conflictErr := &_errors.DataConflictError{
		Resource: "user",
		Field:    "email",
	}

	mockService.
		EXPECT().
		Create(gomock.Any(), gomock.Any()).
		Return(conflictErr).
		Times(1)

	status, data, err := handler.Create(httptest.NewRecorder(), req, reqCtx)

	assert.Equal(409, status)
	assert.Nil(data)
	assert.NotNil(err)
	assert.Equal(1, len(err.Body))
	assert.Equal("Conflicting data", err.Body[0].Message)
	assert.Equal("email", err.Body[0].Target)
}

func TestCreateWithServerError(t *testing.T) {
	reqBody := &_userHandler.CreateUserRequest{
		Name:     "testuser",
		Email:    "testuser@mail.com",
		Password: "secret",
	}

	payload, _ := json.Marshal(reqBody)

	assert := assert.New(t)
	mockCtrl := gomock.NewController(t)
	mockService := mock.NewService(mockCtrl)
	handler := setupHandler(mockService)
	handler.UserService = mockService
	reqCtx := setupRequestContext(handler.App)
	req := httptest.NewRequest("POST", "/users", strings.NewReader(string(payload)))

	mockService.
		EXPECT().
		Create(gomock.Any(), gomock.Any()).
		Return(errors.New("server error")).
		Times(1)

	status, data, err := handler.Create(httptest.NewRecorder(), req, reqCtx)

	assert.Equal(500, status)
	assert.Nil(data)
	assert.NotNil(err)
	assert.Equal(1, len(err.Body))
	assert.Equal("Internal server error", err.Body[0].Message)
}

func TestCreate(t *testing.T) {
	reqBody := &_userHandler.CreateUserRequest{
		Name:     "testuser",
		Email:    "testuser@mail.com",
		Password: "secret",
	}

	payload, _ := json.Marshal(reqBody)

	assert := assert.New(t)
	mockCtrl := gomock.NewController(t)
	mockService := mock.NewService(mockCtrl)
	handler := setupHandler(mockService)
	handler.UserService = mockService
	reqCtx := setupRequestContext(handler.App)
	req := httptest.NewRequest("POST", "/users", strings.NewReader(string(payload)))

	mockService.
		EXPECT().
		Create(gomock.Any(), gomock.Any()).
		Return(nil).
		Times(1)

	status, data, err := handler.Create(httptest.NewRecorder(), req, reqCtx)

	responseData := data.(*_userHandler.CreateUserResponse)

	assert.Equal(201, status)
	assert.NotNil(data)
	assert.Nil(err)
	assert.Equal(reqBody.Name, responseData.User.Name)
	assert.Equal(reqBody.Email, responseData.User.Email)
	assert.Equal(true, responseData.User.IsActive)
}

func TestLoginWithNoRequestBody(t *testing.T) {
	assert := assert.New(t)
	mockCtrl := gomock.NewController(t)
	mockService := mock.NewService(mockCtrl)
	handler := setupHandler(mockService)
	reqCtx := setupRequestContext(handler.App)
	req := httptest.NewRequest("POST", "/login", strings.NewReader(""))

	status, data, err := handler.Login(httptest.NewRecorder(), req, reqCtx)

	assert.Equal(400, status)
	assert.Nil(data)
	assert.Error(err)
	assert.Equal(1, len(err.Body))
	assert.Equal("Invalid request body", err.Body[0].Message)
}

func TestLoginWithBlankEmail(t *testing.T) {
	payload, _ := json.Marshal(&_userHandler.LoginRequest{
		Email:  "",
		Passwd: "secret",
	})

	assert := assert.New(t)
	mockCtrl := gomock.NewController(t)
	mockService := mock.NewService(mockCtrl)
	handler := setupHandler(mockService)
	reqCtx := setupRequestContext(handler.App)
	req := httptest.NewRequest("POST", "/login", strings.NewReader(string(payload)))

	status, data, err := handler.Login(httptest.NewRecorder(), req, reqCtx)

	assert.Equal(400, status)
	assert.Nil(data)
	assert.Error(err)
	assert.Equal(1, len(err.Body))
	assert.Equal("Non-empty value is required", err.Body[0].Message)
	assert.Equal("email", err.Body[0].Target)
}

func TestLoginWithBlankPassword(t *testing.T) {
	payload, _ := json.Marshal(&_userHandler.LoginRequest{
		Email:  "testuser@mail.com",
		Passwd: "",
	})

	assert := assert.New(t)
	mockCtrl := gomock.NewController(t)
	mockService := mock.NewService(mockCtrl)
	handler := setupHandler(mockService)
	reqCtx := setupRequestContext(handler.App)
	req := httptest.NewRequest("POST", "/login", strings.NewReader(string(payload)))

	status, data, err := handler.Login(httptest.NewRecorder(), req, reqCtx)

	assert.Equal(400, status)
	assert.Nil(data)
	assert.Error(err)
	assert.Equal(1, len(err.Body))
	assert.Equal("Non-empty value is required", err.Body[0].Message)
	assert.Equal("password", err.Body[0].Target)
}

func TestLoginWithInvalidEmail(t *testing.T) {
	payload, _ := json.Marshal(&_userHandler.LoginRequest{
		Email:  "testuser@",
		Passwd: "secret",
	})

	assert := assert.New(t)
	mockCtrl := gomock.NewController(t)
	mockService := mock.NewService(mockCtrl)
	handler := setupHandler(mockService)
	reqCtx := setupRequestContext(handler.App)
	req := httptest.NewRequest("POST", "/login", strings.NewReader(string(payload)))

	status, data, err := handler.Login(httptest.NewRecorder(), req, reqCtx)

	assert.Equal(400, status)
	assert.Nil(data)
	assert.Error(err)
	assert.Equal(1, len(err.Body))
	assert.Equal("Invalid value", err.Body[0].Message)
	assert.Equal("email", err.Body[0].Target)
}

func TestLoginWithDataConflict(t *testing.T) {
	reqBody := &_userHandler.LoginRequest{
		Email:  "testuser@mail.com",
		Passwd: "secret",
	}

	payload, _ := json.Marshal(reqBody)

	assert := assert.New(t)
	mockCtrl := gomock.NewController(t)
	mockService := mock.NewService(mockCtrl)
	handler := setupHandler(mockService)
	handler.UserService = mockService
	reqCtx := setupRequestContext(handler.App)
	req := httptest.NewRequest("POST", "/login", strings.NewReader(string(payload)))

	notFoundErr := &_errors.ResourceNotFoundError{
		Resource: "user",
	}

	mockService.
		EXPECT().
		GenerateAuthToken(gomock.Any(), reqBody.Email, reqBody.Passwd, handler.App.Config.Auth.Jwt.Secret).
		Return("", notFoundErr).
		Times(1)

	status, data, err := handler.Login(httptest.NewRecorder(), req, reqCtx)

	assert.Equal(404, status)
	assert.Nil(data)
	assert.NotNil(err)
	assert.Equal(1, len(err.Body))
	assert.Equal("Not found", err.Body[0].Message)
	assert.Equal("user", err.Body[0].Target)
}

func TestLoginWithMismatchingPassword(t *testing.T) {
	reqBody := &_userHandler.LoginRequest{
		Email:  "testuser@mail.com",
		Passwd: "secret",
	}

	payload, _ := json.Marshal(reqBody)

	assert := assert.New(t)
	mockCtrl := gomock.NewController(t)
	mockService := mock.NewService(mockCtrl)
	handler := setupHandler(mockService)
	handler.UserService = mockService
	reqCtx := setupRequestContext(handler.App)
	req := httptest.NewRequest("POST", "/login", strings.NewReader(string(payload)))

	passwordMismatchErr := &_errors.PasswordMismatchError{}

	mockService.
		EXPECT().
		GenerateAuthToken(gomock.Any(), reqBody.Email, reqBody.Passwd, handler.App.Config.Auth.Jwt.Secret).
		Return("", passwordMismatchErr).
		Times(1)

	status, data, err := handler.Login(httptest.NewRecorder(), req, reqCtx)

	assert.Equal(403, status)
	assert.Nil(data)
	assert.NotNil(err)
	assert.Equal(1, len(err.Body))
	assert.Equal("Invalid email/password", err.Body[0].Message)
}

func TestLoginWithServerError(t *testing.T) {
	reqBody := &_userHandler.LoginRequest{
		Email:  "testuser@mail.com",
		Passwd: "secret",
	}

	payload, _ := json.Marshal(reqBody)

	assert := assert.New(t)
	mockCtrl := gomock.NewController(t)
	mockService := mock.NewService(mockCtrl)
	handler := setupHandler(mockService)
	handler.UserService = mockService
	reqCtx := setupRequestContext(handler.App)
	req := httptest.NewRequest("POST", "/login", strings.NewReader(string(payload)))

	serverErr := errors.New("server error")

	mockService.
		EXPECT().
		GenerateAuthToken(gomock.Any(), reqBody.Email, reqBody.Passwd, handler.App.Config.Auth.Jwt.Secret).
		Return("", serverErr).
		Times(1)

	status, data, err := handler.Login(httptest.NewRecorder(), req, reqCtx)

	assert.Equal(500, status)
	assert.Nil(data)
	assert.NotNil(err)
	assert.Equal(1, len(err.Body))
	assert.Equal("Internal server error", err.Body[0].Message)
}

func TestLogin(t *testing.T) {
	reqBody := &_userHandler.LoginRequest{
		Email:  "testuser@mail.com",
		Passwd: "secret",
	}

	payload, _ := json.Marshal(reqBody)

	assert := assert.New(t)
	mockCtrl := gomock.NewController(t)
	mockService := mock.NewService(mockCtrl)
	handler := setupHandler(mockService)
	handler.UserService = mockService
	reqCtx := setupRequestContext(handler.App)
	req := httptest.NewRequest("POST", "/login", strings.NewReader(string(payload)))

	mockService.
		EXPECT().
		GenerateAuthToken(gomock.Any(), reqBody.Email, reqBody.Passwd, handler.App.Config.Auth.Jwt.Secret).
		Return("token", nil).
		Times(1)

	status, data, err := handler.Login(httptest.NewRecorder(), req, reqCtx)

	loginResponse := data.(_userHandler.LoginResponse)

	assert.Equal(200, status)
	assert.Nil(err)
	assert.NotNil(data)
	assert.Equal("token", loginResponse.Token)
}

func setupHandler(mockService user.Service) *_userHandler.UserHandler {
	app := &common.App{
		Logger: logrus.New(),
		Config: &config.Config{},
	}

	appSettings := &config.ApplicationSetting{
		RequestTimeout: 5,
	}

	authSettings := &config.AuthSetting{
		Jwt: &config.JwtSetting{
			Secret: "secret",
		},
	}

	app.Config.Application = appSettings
	app.Config.Auth = authSettings

	handler := &_userHandler.UserHandler{
		UserService: mockService,
		App:         app,
	}

	return handler
}

func setupRequestContext(app *common.App) *common.RequestContext {
	reqCtx := &common.RequestContext{
		RequestID: "dummyRequestID",
		LogEntry: app.Logger.WithFields(
			logrus.Fields{},
		),
	}

	return reqCtx
}
