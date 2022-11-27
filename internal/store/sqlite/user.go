package sqlite

import (
	"context"
	"database/sql"

	"github.com/iryzzh/gophkeeper/internal/models"
	"github.com/iryzzh/gophkeeper/internal/store"
	"github.com/pkg/errors"
)

type UserRepository struct {
	db *sql.DB
}

func (r *UserRepository) FindByID(ctx context.Context, userID string) (*models.User, error) {
	u := &models.User{}

	err := r.db.QueryRowContext(ctx,
		"SELECT user_id, login, password FROM users WHERE user_id = $1",
		userID).Scan(&u.ID, &u.Login, &u.PasswordHash)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, store.ErrUserNotFound
	}

	return u, err
}

func (r *UserRepository) Create(ctx context.Context, user *models.User) (*models.User, error) {
	if u, _ := r.FindByLogin(ctx, user.Login); u != nil {
		return u, store.ErrUserAlreadyExists
	}

	_, err := r.db.ExecContext(ctx,
		"INSERT INTO users (user_id, login, password) VALUES ($1, $2, $3)",
		user.ID, user.Login, user.PasswordHash)
	if err != nil {
		return nil, store.ErrUserCreateFailed
	}

	return user, nil
}

func (r *UserRepository) FindByLogin(ctx context.Context, login string) (*models.User, error) {
	u := &models.User{}

	err := r.db.QueryRowContext(ctx,
		"SELECT user_id, login, password FROM users WHERE login = $1",
		login).Scan(&u.ID, &u.Login, &u.PasswordHash)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, store.ErrUserNotFound
	}

	return u, err
}
