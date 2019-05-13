package service_test

import (
	"context"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"

	"github.com/golang/mock/gomock"

	todoErr "github.com/dheerajgopi/todo-api/common/error"
	"github.com/dheerajgopi/todo-api/models"
	repoMock "github.com/dheerajgopi/todo-api/user/mock"
	"github.com/dheerajgopi/todo-api/user/service"
)

func TestCreate(t *testing.T) {
	now := time.Now()
	ctx := context.TODO()
	assert := assert.New(t)
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	userRepoMock := repoMock.NewRepository(mockCtrl)
	userService := service.New(userRepoMock)

	newUser := &models.User{
		Name:      "testName",
		Email:     "testName@email.com",
		Passwd:    "test",
		IsActive:  true,
		CreatedAt: now,
		UpdatedAt: now,
	}

	userRepoMock.
		EXPECT().
		GetByEmail(ctx, "testName@email.com").
		Return(nil, nil).
		Times(1)

	userRepoMock.
		EXPECT().
		Create(ctx, newUser).
		Return(nil).
		Times(1)

	err := userService.Create(ctx, newUser)

	assert.NoError(err)
	assert.Equal(newUser, newUser)
	assert.Equal("testName@email.com", newUser.Email)
	assert.Equal("testName", newUser.Name)
	assert.NotEqual("test", newUser.Passwd)
	assert.Equal(true, newUser.IsActive)
	assert.Equal(now, newUser.CreatedAt)
	assert.Equal(now, newUser.UpdatedAt)
}

func TestCreateForAlreadyExistingUser(t *testing.T) {
	now := time.Now()
	ctx := context.TODO()
	assert := assert.New(t)
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	userRepoMock := repoMock.NewRepository(mockCtrl)
	userService := service.New(userRepoMock)

	newUser := &models.User{
		Name:      "testName",
		Email:     "testName@email.com",
		Passwd:    "test",
		IsActive:  true,
		CreatedAt: now,
		UpdatedAt: now,
	}

	userRepoMock.
		EXPECT().
		GetByEmail(ctx, "testName@email.com").
		Return(newUser, nil).
		Times(1)

	dataConflictErr := &todoErr.DataConflictError{
		Resource: "user",
		Field:    "email",
	}

	err := userService.Create(ctx, newUser)

	assert.Error(err)
	assert.Equal(dataConflictErr, err)
}

func TestGenerateAuthToken(t *testing.T) {
	now := time.Now()
	ctx := context.TODO()
	assert := assert.New(t)
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	userRepoMock := repoMock.NewRepository(mockCtrl)
	userService := service.New(userRepoMock)

	email := "testName@email.com"
	passwd := "test"
	pswdHash, _ := bcrypt.GenerateFromPassword([]byte(passwd), bcrypt.DefaultCost)
	jwtSecret := "secret"

	newUser := &models.User{
		Name:      "testName",
		Email:     email,
		Passwd:    string(pswdHash),
		IsActive:  true,
		CreatedAt: now,
		UpdatedAt: now,
	}

	userRepoMock.
		EXPECT().
		GetByEmail(ctx, email).
		Return(newUser, nil).
		Times(1)

	token, _ := userService.GenerateAuthToken(ctx, newUser.Email, passwd, jwtSecret)

	parsedToken, _ := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})

	claims, _ := parsedToken.Claims.(jwt.MapClaims)

	assert.NoError(claims.Valid())
}

func TestGenerateAuthTokenForMissingUser(t *testing.T) {
	ctx := context.TODO()
	assert := assert.New(t)
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	userRepoMock := repoMock.NewRepository(mockCtrl)
	userService := service.New(userRepoMock)

	email := "testName@email.com"
	passwd := "test"
	jwtSecret := "secret"

	userRepoMock.
		EXPECT().
		GetByEmail(ctx, email).
		Return(nil, nil).
		Times(1)

	token, err := userService.GenerateAuthToken(ctx, email, passwd, jwtSecret)

	resourceNotFoundError := &todoErr.ResourceNotFoundError{
		Resource: "user",
	}

	assert.Empty(token)
	assert.Equal(resourceNotFoundError, err)
}

func TestGenerateAuthTokenForPasswordMismatch(t *testing.T) {
	now := time.Now()
	ctx := context.TODO()
	assert := assert.New(t)
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	userRepoMock := repoMock.NewRepository(mockCtrl)
	userService := service.New(userRepoMock)

	email := "testName@email.com"
	passwd := "test"
	pswdHash, _ := bcrypt.GenerateFromPassword([]byte(passwd), bcrypt.DefaultCost)
	jwtSecret := "secret"

	newUser := &models.User{
		Name:      "testName",
		Email:     email,
		Passwd:    string(pswdHash),
		IsActive:  true,
		CreatedAt: now,
		UpdatedAt: now,
	}

	userRepoMock.
		EXPECT().
		GetByEmail(ctx, email).
		Return(newUser, nil).
		Times(1)

	token, err := userService.GenerateAuthToken(ctx, newUser.Email, "invalid"+passwd, jwtSecret)

	expectedErr := &todoErr.PasswordMismatchError{}

	assert.Empty(token)
	assert.Equal(expectedErr, err)
}
