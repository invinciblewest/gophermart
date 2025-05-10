package repository

import (
	"context"
	"errors"
	"github.com/invinciblewest/gophermart/internal/model"
)

var (
	ErrEmptyLoginOrPassword = errors.New("login or password is empty")
	ErrUserAlreadyExists    = errors.New("user already exists")
	ErrUserNotFound         = errors.New("user not found")
	ErrInvalidPassword      = errors.New("invalid password")
	ErrOrderNotFound        = errors.New("order not found")
	ErrWithdrawalNotFound   = errors.New("withdrawal not found")
)

type UserRepository interface {
	Create(ctx context.Context, user *model.User) error
	GetByLogin(ctx context.Context, login string) (*model.User, error)
}

type OrderRepository interface {
	Add(ctx context.Context, order *model.Order) error
	GetByUser(ctx context.Context, userID int) ([]model.Order, error)
	GetByNumber(ctx context.Context, number string) (*model.Order, error)
	UpdateStatus(ctx context.Context, number string, status model.OrderStatus, accrual *float64) error
	GetSumAccrualByUser(ctx context.Context, userID int) (float64, error)
	GetPending(ctx context.Context) ([]model.Order, error)
}

type WithdrawalRepository interface {
	GetSumWithdrawnByUser(ctx context.Context, userID int) (float64, error)
	Create(ctx context.Context, withdrawal *model.Withdrawal) error
	GetByUser(ctx context.Context, userID int) ([]model.Withdrawal, error)
}
