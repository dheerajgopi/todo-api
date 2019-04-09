package repository_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/dheerajgopi/todo-api/user/repository"
)

func TestGetByID(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("Unexpected error while opening stub DB connection: %s", err)
	}

	defer db.Close()

	rows := sqlmock.
		NewRows([]string{"id", "name", "is_active", "created_at", "updated_at"}).
		AddRow(1, "test user", true, time.Now(), time.Now())

	userId := int64(1)
	query := "SELECT id, name, is_active, created_at, updated_at FROM user where id=\\?"

	prep := mock.ExpectPrepare(query)
	prep.ExpectQuery().WithArgs(userId).WillReturnRows(rows)

	repo := repository.Create(db)

	user, err := repo.GetByID(context.TODO(), userId)
	assert.NoError(t, err)
	assert.NotNil(t, user)
}
