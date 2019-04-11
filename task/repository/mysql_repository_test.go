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
		NewRows([]string{"id", "description", "created_by", "is_complete", "created_at", "updated_at"}).
		AddRow(1, "description", 1, false, time.Now(), time.Now())

	taskId := int64(1)
	query := "SELECT id, description, created_by, is_complete, created_at, updated_at FROM task WHERE id=\\?"

	prep := mock.ExpectPrepare(query)
	prep.ExpectQuery().WithArgs(taskId).WillReturnRows(rows)

	repo := repository.New(db)

	task, err := repo.GetByID(context.TODO(), taskId)

	assert.NoError(t, err)
	assert.NotNil(t, task)
}

func TestCreate(t *testing.T) {
	now := time.Now()
	task := &models.Task{
		Description: "description",
		CreatedBy: models.User{
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

	query := "INSERT INTO task \\(description, created_by, is_complete, created_at, updated_at\\) VALUES \\(\\?, \\?, \\?, \\?, \\?\\)"

	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(
		task.Description,
		task.CreatedBy.ID,
		task.IsComplete,
		task.CreatedAt,
		task.UpdatedAt,
	).WillReturnResult(sqlmock.NewResult(2, 1))

	repo := repository.New(db)

	err = repo.Create(context.TODO(), task)

	assert.NoError(t, err)
	assert.Equal(t, int64(2), task.ID)
}
