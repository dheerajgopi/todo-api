package repository

import (
	"context"
	"database/sql"

	"github.com/dheerajgopi/todo-api/models"
	"github.com/sirupsen/logrus"

	"github.com/dheerajgopi/todo-api/task"
)

type mySQLRepo struct {
	DB *sql.DB
}

// New will return new object which implements task.Repository
func New(db *sql.DB) task.Repository {
	return &mySQLRepo{
		DB: db,
	}
}

func (repo *mySQLRepo) getOne(ctx context.Context, query string, args ...interface{}) (*models.Task, error) {
	stmt, err := repo.DB.PrepareContext(ctx, query)

	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	row := stmt.QueryRowContext(ctx, args...)
	task := &models.Task{}
	userID := int64(0)

	err = row.Scan(
		&task.ID,
		&task.Title,
		&task.Description,
		&userID,
		&task.IsComplete,
		&task.CreatedAt,
		&task.UpdatedAt,
	)

	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	task.CreatedBy = models.User{
		ID: userID,
	}

	return task, nil
}

// GetByID will return task with the given id
func (repo *mySQLRepo) GetByID(ctx context.Context, id int64) (*models.Task, error) {
	query := `SELECT id, title, description, created_by, is_complete, created_at, updated_at
		FROM task WHERE id=?`
	return repo.getOne(ctx, query, id)
}

// Create will store new task entry
func (repo *mySQLRepo) Create(ctx context.Context, task *models.Task) error {
	query := `INSERT INTO task (title, description, created_by, is_complete, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?)`

	stmt, err := repo.DB.PrepareContext(ctx, query)

	defer func() {
		if err != nil {
			logrus.Error(err)
		}
	}()

	if err != nil {
		return err
	}

	res, err := stmt.ExecContext(
		ctx,
		task.Title,
		task.Description,
		task.CreatedBy.ID,
		task.IsComplete,
		task.CreatedAt,
		task.UpdatedAt,
	)

	if err != nil {
		return err
	}

	lastID, err := res.LastInsertId()

	if err != nil {
		return err
	}

	task.ID = lastID

	return nil
}
