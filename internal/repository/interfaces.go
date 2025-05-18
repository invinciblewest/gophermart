package repository

import (
	"context"
	"github.com/invinciblewest/gophermart/internal/model"
)

type UserRepository interface {
	Create(ctx context.Context, user *model.User) error
	GetByLogin(ctx context.Context, login string) (*model.User, error)
}

type OrderRepository interface {
	Add(ctx context.Context, order *model.Order) error
	GetByUser(ctx context.Context, userID int) ([]model.Order, error)
	GetByNumber(ctx context.Context, number string) (*model.Order, error)
	UpdateStatus(ctx context.Context, number string, status model.OrderStatus, accrual *model.Amount) error
	GetSumAccrualByUser(ctx context.Context, userID int) (model.Amount, error)
	GetPending(ctx context.Context) ([]model.Order, error)
}

type WithdrawalRepository interface {
	GetSumWithdrawnByUser(ctx context.Context, userID int) (model.Amount, error)
	Create(ctx context.Context, withdrawal *model.Withdrawal) error
	GetByUser(ctx context.Context, userID int) ([]model.Withdrawal, error)
}
