package postgres

import (
	"context"
	"database/sql"
	"errors"
	"github.com/invinciblewest/gophermart/internal/model"
	"github.com/jackc/pgerrcode"
	"github.com/lib/pq"
)

func (r *PGRepository) CreateUser(ctx context.Context, user *model.User) error {
	if user.Login == "" || user.Password == "" {
		return model.ErrEmptyLoginOrPassword
	}

	err := r.db.QueryRowContext(
		ctx,
		"INSERT INTO users (login, password) VALUES ($1, $2) RETURNING id, created_at",
		user.Login,
		user.Password,
	).Scan(&user.ID, &user.CreatedAt)

	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == pgerrcode.UniqueViolation {
			return model.ErrUserAlreadyExists
		}
		return err
	}

	return nil
}

func (r *PGRepository) GetUserByLogin(ctx context.Context, login string) (*model.User, error) {
	if login == "" {
		return nil, model.ErrEmptyLoginOrPassword
	}

	var user model.User

	err := r.db.QueryRowContext(ctx,
		`SELECT id, login, password, created_at FROM users WHERE login = $1`,
		login).Scan(&user.ID, &user.Login, &user.Password, &user.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, model.ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}
