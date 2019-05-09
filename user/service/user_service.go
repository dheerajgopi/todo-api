package service

import (
	"context"
	"time"

	"github.com/dgrijalva/jwt-go"
	_err "github.com/dheerajgopi/todo-api/common/error"
	"github.com/dheerajgopi/todo-api/models"
	"github.com/dheerajgopi/todo-api/user"
	"golang.org/x/crypto/bcrypt"
)

type userService struct {
	userRepo user.Repository
}

// New returns a new object implementing user.Service interface
func New(repo user.Repository) user.Service {
	return &userService{
		userRepo: repo,
	}
}

// Create creates a new user
func (service *userService) Create(ctx context.Context, newUser *models.User) error {
	email := newUser.Email

	existingUser, err := service.userRepo.GetByEmail(ctx, email)

	if err != nil {
		return err
	}

	if existingUser != nil {
		err := _err.DataConflictError{
			Resource: "user",
			Field:    "email",
		}
		return &err
	}

	pswd := newUser.Passwd
	pswdHash, err := bcrypt.GenerateFromPassword([]byte(pswd), bcrypt.DefaultCost)

	if err != nil {
		return err
	}

	newUser.Passwd = string(pswdHash)

	err = service.userRepo.Create(ctx, newUser)

	if err != nil {
		return err
	}

	return nil
}

// GenerateAuthToken validates the password and generates JWT
func (service *userService) GenerateAuthToken(ctx context.Context, email string, pswd string, secret string) (string, error) {
	user, err := service.userRepo.GetByEmail(ctx, email)

	if err != nil {
		return "", err
	}

	if user == nil {
		resourceNotFoundError := _err.ResourceNotFoundError{
			Resource: "user",
		}

		return "", &resourceNotFoundError
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Passwd), []byte(pswd))

	if err != nil {
		return "", err
	}

	now := time.Now()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss":    "todo-api",
		"sub":    "user",
		"userId": user.ID,
		"name":   user.Name,
		"email":  user.Email,
		"iat":    now,
		"exp":    now.Add(1 * time.Hour),
	})

	signedToken, err := token.SignedString([]byte(secret))

	if err != nil {
		return "", err
	}

	return signedToken, nil
}
