package service

import (
	"context"

	"github.com/dheerajgopi/todo-api/common"
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

func (service *userService) Create(ctx context.Context, newUser *models.User) error {
	email := newUser.Email

	existingUser, err := service.userRepo.GetByEmail(ctx, email)

	if err != nil {
		return err
	}

	if existingUser != nil {
		err := common.DataConflictError{
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
