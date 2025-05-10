package app

import (
	"context"
	"errors"
	"github.com/invinciblewest/gophermart/internal/helper"
	"github.com/invinciblewest/gophermart/internal/model"
	"github.com/invinciblewest/gophermart/internal/repository"
	"github.com/invinciblewest/gophermart/internal/usecase"
)

type BalanceUseCase struct {
	orderRepository      repository.OrderRepository
	withdrawalRepository repository.WithdrawalRepository
}

func NewBalanceUseCase(orderRepository repository.OrderRepository, withdrawalRepository repository.WithdrawalRepository) *BalanceUseCase {
	return &BalanceUseCase{
		orderRepository:      orderRepository,
		withdrawalRepository: withdrawalRepository,
	}
}

func (b *BalanceUseCase) GetUserBalance(ctx context.Context, userID int) (*model.Balance, error) {
	var balance model.Balance
	totalAccrual, err := b.orderRepository.GetSumAccrualByUser(ctx, userID)
	if err != nil {
		return &balance, err
	}

	totalWithdrawn, err := b.withdrawalRepository.GetSumWithdrawnByUser(ctx, userID)
	if err != nil {
		return &balance, err
	}

	balance.Current = totalAccrual - totalWithdrawn
	balance.Withdrawn = totalWithdrawn

	return &balance, nil
}

func (b *BalanceUseCase) WithdrawBalance(ctx context.Context, userID int, withdrawRequest model.WithdrawRequest) error {
	if !helper.IsValidOrderNumber(withdrawRequest.Order) {
		return usecase.ErrInvalidOrderNumber
	}

	balance, err := b.GetUserBalance(ctx, userID)
	if err != nil {
		return err
	}

	if balance.Current < withdrawRequest.Sum {
		return usecase.ErrInvalidWithdrawSum
	}

	withdrawal := &model.Withdrawal{
		UserID:      userID,
		OrderNumber: withdrawRequest.Order,
		Amount:      withdrawRequest.Sum,
	}

	if err = b.withdrawalRepository.Create(ctx, withdrawal); err != nil {
		return err
	}

	return nil
}

func (b *BalanceUseCase) GetWithdrawals(ctx context.Context, userID int) ([]model.Withdrawal, error) {
	withdrawals, err := b.withdrawalRepository.GetByUser(ctx, userID)
	if err != nil {
		if errors.Is(err, repository.ErrWithdrawalNotFound) {
			return nil, usecase.ErrWithdrawalNotFound
		}
		return nil, err
	}

	return withdrawals, nil
}
