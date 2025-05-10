package usecase

import (
	"context"
	"errors"
	"github.com/invinciblewest/gophermart/internal/model"
)

var (
	ErrOrderNotFound                    = errors.New("order not found")
	ErrInvalidOrderNumber               = errors.New("invalid order number")
	ErrOrderAlreadyExists               = errors.New("order already exists")
	ErrOrderAlreadyExistsForAnotherUser = errors.New("order already exists for another user")
	ErrInvalidWithdrawSum               = errors.New("invalid withdraw sum")
	ErrWithdrawalNotFound               = errors.New("withdrawal not found")
)

type AuthUseCase interface {
	GenerateToken(userID int) (string, error)
	ParseToken(tokenStr string) (int, error)
	HashPassword(password string) string
	VerifyPassword(user *model.User, password string) bool
}

type OrderUseCase interface {
	AddOrder(ctx context.Context, order *model.Order) error
	GetByUser(ctx context.Context, userID int) ([]model.Order, error)
}

type UserUseCase interface {
	RegisterAndLogin(ctx context.Context, user *model.User) (string, error)
	Login(ctx context.Context, user model.User) (string, error)
}

type BalanceUseCase interface {
	GetUserBalance(ctx context.Context, userID int) (*model.Balance, error)
	WithdrawBalance(ctx context.Context, userID int, request model.WithdrawRequest) error
	GetWithdrawals(ctx context.Context, userID int) ([]model.Withdrawal, error)
}
