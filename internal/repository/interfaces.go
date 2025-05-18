package repository

import (
	"context"
	"github.com/invinciblewest/gophermart/internal/model"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *model.User) error
	GetUserByLogin(ctx context.Context, login string) (*model.User, error)
}

type OrderRepository interface {
	AddOrder(ctx context.Context, order *model.Order) error
	GetOrderByUser(ctx context.Context, userID int) ([]model.Order, error)
	GetOrderByNumber(ctx context.Context, number string) (*model.Order, error)
	UpdateOrderStatus(ctx context.Context, number string, status model.OrderStatus, accrual *model.Amount) error
	GetPendingOrders(ctx context.Context) ([]model.Order, error)
}

type WithdrawalRepository interface {
	CreateWithdrawal(ctx context.Context, withdrawal *model.Withdrawal) error
	GetWithdrawalByUser(ctx context.Context, userID int) ([]model.Withdrawal, error)
}

type BalanceRepository interface {
	GetBalanceByUser(ctx context.Context, userID int) (*model.Balance, error)
}
