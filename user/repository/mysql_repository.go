package repository

import (
	"context"
	"database/sql"

	"github.com/sirupsen/logrus"

	"github.com/dheerajgopi/todo-api/models"
	"github.com/dheerajgopi/todo-api/user"
)

type mySQLUserRepo struct {
	DB *sql.DB
}

// New will return new object which implements user.Repository
func New(db *sql.DB) user.Repository {
	return &mySQLUserRepo{
		DB: db,
	}
}

func (repo *mySQLUserRepo) getOne(ctx context.Context, query string, args ...interface{}) (*models.User, error) {
	stmt, err := repo.DB.PrepareContext(ctx, query)

	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	row := stmt.QueryRowContext(ctx, args...)
	user := &models.User{}

	err = row.Scan(
		&user.ID,
		&user.Name,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	return user, nil
}

// GetByID will return user with the given id
func (repo *mySQLUserRepo) GetByID(ctx context.Context, id int64) (*models.User, error) {
	query := `SELECT id, name, is_active, created_at, updated_at FROM user WHERE id=?`
	return repo.getOne(ctx, query, id)
}
