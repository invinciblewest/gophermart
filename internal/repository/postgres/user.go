package postgres

import (
	"context"
	"database/sql"
	"errors"
	"github.com/invinciblewest/gophermart/internal/helper"
	"github.com/invinciblewest/gophermart/internal/model"
	"github.com/jackc/pgerrcode"
	"github.com/lib/pq"
)

type PGUserRepository struct {
	db *sql.DB
}

func NewPGUserRepository(db *sql.DB) *PGUserRepository {
	return &PGUserRepository{
		db: db,
	}
}

func (r *PGUserRepository) Create(ctx context.Context, user *model.User) error {
	if user.Login == "" || user.Password == "" {
		return helper.ErrEmptyLoginOrPassword
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
			return helper.ErrUserAlreadyExists
		}
		return err
	}

	return nil
}

func (r *PGUserRepository) GetByLogin(ctx context.Context, login string) (*model.User, error) {
	if login == "" {
		return nil, helper.ErrEmptyLoginOrPassword
	}

	var user model.User

	err := r.db.QueryRowContext(ctx,
		`SELECT id, login, password, created_at FROM users WHERE login = $1`,
		login).Scan(&user.ID, &user.Login, &user.Password, &user.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, helper.ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}
