package repository

import (
	"banking_system_golang/internal/models"
	"context"
	"database/sql"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func(r *UserRepository) Create(ctx context.Context, name, username, email, password string) (*models.User, error) {
	var user models.User
	q := `INSERT INTO users (name, username, email, password)
        VALUES ($1, $2, $3, $4)
        RETURNING id, name, username, email, created_at`
	err := r.db.QueryRowContext(ctx, q, name, username, email, password).Scan(&user.ID, &user.Name, &user.Username, &user.Email, &user.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func(r *UserRepository) Update(ctx context.Context, user *models.User) error {
	var user models.User
	q := `
		UPDATE users 
        SET name=$1, username=$2, email=$3, password=$4
        WHERE id=$5
	`
	_, err := r.db.ExecContext(ctx, q, user.Name, user.Username, user.Email, user.Password, user.ID)
	if err != nil {
		return nil, err
	}

	return err
}

func(r *UserRepository) Delete(ctx context.Context, id int64) error {
	var user models.User
	q := `DELETE FROM users WHERE id=$1`

	_, err := r.db.ExecContext(ctx, q, id)
	if err != nil {
        return nil, err
    }

	return err
}

func (r *UserRepository) FindByUsername(ctx context.Context, username string) (*models.User, error) {
    var user models.User
    query := `SELECT id, name, username, email, password, created_at FROM users WHERE username = $1`
    err := r.db.QueryRowContext(ctx, query, username).
        Scan(&user.ID, &user.Name, &user.Username, &user.Email, &user.Password, &user.CreatedAt)
    if err != nil {
        return nil, err
    }
    return &user, nil
}