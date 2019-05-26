package repository_test

import (
	"context"
	"testing"
	"time"

	"github.com/dheerajgopi/todo-api/models"

	"github.com/stretchr/testify/assert"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/dheerajgopi/todo-api/task/repository"
)

func TestGetByID(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("Unexpected error while opening stub DB connection: %s", err)
	}

	defer db.Close()

	rows := sqlmock.
		NewRows([]string{"id", "title", "description", "created_by", "is_complete", "created_at", "updated_at"}).
		AddRow(1, "title", "description", 1, false, time.Now(), time.Now())

	taskID := int64(1)
	query := "SELECT id, title, description, created_by, is_complete, created_at, updated_at FROM task WHERE id=\\?"

	prep := mock.ExpectPrepare(query)
	prep.ExpectQuery().WithArgs(taskID).WillReturnRows(rows)

	repo := repository.New(db)

	task, err := repo.GetByID(context.TODO(), taskID)

	assert.NoError(t, err)
	assert.NotNil(t, task)
}

func TestCreate(t *testing.T) {
	now := time.Now()
	task := &models.Task{
		Title:       "title",
		Description: "description",
		CreatedBy: &models.User{
			ID: int64(1),
		},
		IsComplete: false,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("Unexpected error while opening stub DB connection: %s", err)
	}

	defer db.Close()

	query := "INSERT INTO task \\(title, description, created_by, is_complete, created_at, updated_at\\) VALUES \\(\\?, \\?, \\?, \\?, \\?, \\?\\)"

	mock.ExpectBegin()
	mock.ExpectExec(query).WithArgs(
		task.Title,
		task.Description,
		task.CreatedBy.ID,
		task.IsComplete,
		task.CreatedAt,
		task.UpdatedAt,
	).WillReturnResult(sqlmock.NewResult(2, 1))
	mock.ExpectCommit()

	repo := repository.New(db)

	err = repo.Create(context.TODO(), task)

	assert.NoError(t, err)
	assert.Equal(t, int64(2), task.ID)
}

func TestGetAllByUserID(t *testing.T) {
	assert := assert.New(t)
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("Unexpected error while opening stub DB connection: %s", err)
	}

	defer db.Close()

	rows := sqlmock.
		NewRows([]string{"id", "title", "description", "created_by", "is_complete", "created_at", "updated_at"}).
		AddRow(1, "title", "description", 1, false, time.Now(), time.Now())

	userID := int64(1)
	query := "SELECT id, title, description, created_by, is_complete, created_at, updated_at FROM task WHERE created_by=\\?"

	prep := mock.ExpectPrepare(query)
	prep.ExpectQuery().WithArgs(userID).WillReturnRows(rows)

	repo := repository.New(db)

	tasks, err := repo.GetAllByUserID(context.TODO(), userID)

	assert.NoError(err)
	assert.NotNil(tasks)
	assert.Equal(1, len(tasks))
}
