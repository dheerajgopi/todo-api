package repository

import (
	"context"
	"database/sql"

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
		return nil, err
	}

	row := stmt.QueryRowContext(ctx, args...)
	user := &models.User{}

	err = row.Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Passwd,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	switch err {
	case nil:
	case sql.ErrNoRows:
		return nil, nil
	default:
		return nil, err
	}

	return user, nil
}

// GetByID will return user with the given id
func (repo *mySQLUserRepo) GetByID(ctx context.Context, id int64) (*models.User, error) {
	query := `SELECT id, name, email, passwd, is_active, created_at, updated_at FROM user WHERE id=?`
	return repo.getOne(ctx, query, id)
}

// Create will store new user entry
func (repo *mySQLUserRepo) Create(ctx context.Context, user *models.User) error {
	query := `INSERT INTO user (name, email, passwd, is_active, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?)`

	tx, err := repo.DB.BeginTx(ctx, nil)

	if err != nil {
		return err
	}

	res, err := tx.Exec(
		query,
		&user.Name,
		&user.Email,
		&user.Passwd,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		tx.Rollback()
		return err
	}

	lastID, err := res.LastInsertId()

	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()

	if err != nil {
		return err
	}

	user.ID = lastID

	return nil
}

// GetByEmail will return user with the given email
func (repo *mySQLUserRepo) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `SELECT id, name, email, passwd, is_active, created_at, updated_at FROM user WHERE email=?`
	return repo.getOne(ctx, query, email)
}
