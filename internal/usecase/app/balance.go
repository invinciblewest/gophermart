package app

import (
	"context"
	"github.com/invinciblewest/gophermart/internal/helper"
	"github.com/invinciblewest/gophermart/internal/model"
	"github.com/invinciblewest/gophermart/internal/repository"
)

type BalanceUseCase struct {
	balanceRepository    repository.BalanceRepository
	withdrawalRepository repository.WithdrawalRepository
}

func NewBalanceUseCase(balanceRepository repository.BalanceRepository, withdrawalRepository repository.WithdrawalRepository) *BalanceUseCase {
	return &BalanceUseCase{
		balanceRepository:    balanceRepository,
		withdrawalRepository: withdrawalRepository,
	}
}

func (b *BalanceUseCase) GetUserBalance(ctx context.Context, userID int) (*model.Balance, error) {
	balance, err := b.balanceRepository.GetBalanceByUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	return balance, nil
}

func (b *BalanceUseCase) WithdrawBalance(ctx context.Context, userID int, withdrawRequest model.WithdrawRequest) error {
	if !helper.IsValidOrderNumber(withdrawRequest.Order) {
		return model.ErrInvalidOrderNumber
	}

	balance, err := b.GetUserBalance(ctx, userID)
	if err != nil {
		return err
	}

	if balance.Current < withdrawRequest.Sum {
		return model.ErrInvalidWithdrawSum
	}

	withdrawal := &model.Withdrawal{
		UserID:      userID,
		OrderNumber: withdrawRequest.Order,
		Amount:      withdrawRequest.Sum,
	}

	if err = b.withdrawalRepository.CreateWithdrawal(ctx, withdrawal); err != nil {
		return err
	}

	return nil
}

func (b *BalanceUseCase) GetWithdrawals(ctx context.Context, userID int) ([]model.Withdrawal, error) {
	withdrawals, err := b.withdrawalRepository.GetWithdrawalByUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	return withdrawals, nil
}
