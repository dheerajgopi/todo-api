package repository_test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/dheerajgopi/todo-api/models"

	"github.com/stretchr/testify/assert"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/dheerajgopi/todo-api/user/repository"
)

func TestGetByID(t *testing.T) {
	assert := assert.New(t)
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("Unexpected error while opening stub DB connection: %s", err)
	}

	defer db.Close()

	rows := sqlmock.
		NewRows([]string{"id", "name", "email", "passwd", "is_active", "created_at", "updated_at"}).
		AddRow(1, "test user", "test@email.com", "passwd", true, time.Now(), time.Now())

	userID := int64(1)
	query := "SELECT id, name, email, passwd, is_active, created_at, updated_at FROM user WHERE id=\\?"

	prep := mock.ExpectPrepare(query)
	prep.ExpectQuery().WithArgs(userID).WillReturnRows(rows)

	repo := repository.New(db)

	user, err := repo.GetByID(context.TODO(), userID)
	assert.NoError(err)
	assert.NotNil(user)
}

func TestGetByIDWithNoRows(t *testing.T) {
	assert := assert.New(t)
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("Unexpected error while opening stub DB connection: %s", err)
	}

	defer db.Close()

	userID := int64(1)
	query := "SELECT id, name, email, passwd, is_active, created_at, updated_at FROM user WHERE id=\\?"

	prep := mock.ExpectPrepare(query)
	prep.ExpectQuery().WithArgs(userID).WillReturnError(sql.ErrNoRows)

	repo := repository.New(db)

	user, err := repo.GetByID(context.TODO(), userID)
	assert.NoError(err)
	assert.Nil(user)
}

func TestCreate(t *testing.T) {
	assert := assert.New(t)
	now := time.Now()
	user := &models.User{
		Name:      "name",
		Email:     "name@email.com",
		Passwd:    "secret",
		IsActive:  true,
		CreatedAt: now,
		UpdatedAt: now,
	}

	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("Unexpected error while opening stub DB connection: %s", err)
	}

	defer db.Close()

	query := "INSERT INTO user \\(name, email, passwd, is_active, created_at, updated_at\\) VALUES \\(\\?, \\?, \\?, \\?, \\?, \\?\\)"
	lastInsertID := int64(1)

	mock.ExpectBegin()
	mock.ExpectExec(query).WithArgs(
		user.Name,
		user.Email,
		user.Passwd,
		user.IsActive,
		user.CreatedAt,
		user.UpdatedAt,
	).WillReturnResult(sqlmock.NewResult(lastInsertID, 1))
	mock.ExpectCommit()

	repo := repository.New(db)

	err = repo.Create(context.TODO(), user)

	assert.NoError(err)
	assert.Equal(lastInsertID, user.ID)
}

func TestGetByEmail(t *testing.T) {
	assert := assert.New(t)
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("Unexpected error while opening stub DB connection: %s", err)
	}

	defer db.Close()

	rows := sqlmock.
		NewRows([]string{"id", "name", "email", "passwd", "is_active", "created_at", "updated_at"}).
		AddRow(1, "test user", "test@email.com", "passwd", true, time.Now(), time.Now())

	userEmail := "test@email.com"
	query := "SELECT id, name, email, passwd, is_active, created_at, updated_at FROM user WHERE email=\\?"

	prep := mock.ExpectPrepare(query)
	prep.ExpectQuery().WithArgs(userEmail).WillReturnRows(rows)

	repo := repository.New(db)

	user, err := repo.GetByEmail(context.TODO(), userEmail)
	assert.NoError(err)
	assert.NotNil(user)
}

func TestGetByEmailWithNoRows(t *testing.T) {
	assert := assert.New(t)
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("Unexpected error while opening stub DB connection: %s", err)
	}

	defer db.Close()

	userEmail := "test@email.com"
	query := "SELECT id, name, email, passwd, is_active, created_at, updated_at FROM user WHERE email=\\?"

	prep := mock.ExpectPrepare(query)
	prep.ExpectQuery().WithArgs(userEmail).WillReturnError(sql.ErrNoRows)

	repo := repository.New(db)

	user, err := repo.GetByEmail(context.TODO(), userEmail)
	assert.NoError(err)
	assert.Nil(user)
}
