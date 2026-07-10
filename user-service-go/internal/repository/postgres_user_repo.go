package repository

import (
	"context"
	"database/sql"
	"errors"
	"strings"
)

var ErrUserNotFound = errors.New("user not found")

type User struct {
	UserID      string
	Email       string
	DisplayName string
}

type PostgresUserRepo struct {
	db *sql.DB
}

func NewPostgresUserRepo(db *sql.DB) *PostgresUserRepo {
	return &PostgresUserRepo{db: db}
}

func (r *PostgresUserRepo) Init(ctx context.Context) error {
	if ctx == nil {
		return errors.New("context is required")
	}
	if r == nil || r.db == nil {
		return errors.New("user repository is not initialized")
	}
	_, err := r.db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS users (
			user_id TEXT PRIMARY KEY,
			email TEXT NOT NULL,
			display_name TEXT NOT NULL
		)
	`)
	if err != nil {
		return err
	}

	var count int
	if err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM users`).Scan(&count); err != nil {
		return err
	}

	if count == 0 {
		_, err = r.db.ExecContext(ctx, `
			INSERT INTO users (user_id, email, display_name)
			VALUES
				('u-100', 'ava@example.com', 'Ava Reader'),
				('u-101', 'liam@example.com', 'Liam Bookworm')
		`)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *PostgresUserRepo) GetByID(ctx context.Context, userID string) (*User, error) {
	if ctx == nil {
		return nil, errors.New("context is required")
	}
	userID = strings.TrimSpace(userID)
	if userID == "" {
		return nil, errors.New("user_id is required")
	}
	if len(userID) > 128 {
		return nil, errors.New("user_id must not exceed 128 characters")
	}
	row := r.db.QueryRowContext(ctx, `
		SELECT user_id, email, display_name
		FROM users
		WHERE user_id = ?
	`, userID)

	var user User
	err := row.Scan(&user.UserID, &user.Email, &user.DisplayName)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}
