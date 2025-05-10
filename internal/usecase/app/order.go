package app

import (
	"context"
	"errors"
	"github.com/invinciblewest/gophermart/internal/helper"
	"github.com/invinciblewest/gophermart/internal/model"
	"github.com/invinciblewest/gophermart/internal/repository"
	"github.com/invinciblewest/gophermart/internal/usecase"
)

type OrderUseCase struct {
	OrderRepository repository.OrderRepository
}

func NewOrderUseCase(orderRepository repository.OrderRepository) *OrderUseCase {
	return &OrderUseCase{
		OrderRepository: orderRepository,
	}
}

func (os *OrderUseCase) AddOrder(ctx context.Context, order *model.Order) error {
	if !helper.IsValidOrderNumber(order.Number) {
		return usecase.ErrInvalidOrderNumber
	}

	receivedOrder, err := os.OrderRepository.GetByNumber(ctx, order.Number)
	if err != nil && !errors.Is(err, repository.ErrOrderNotFound) {
		return err
	}

	if receivedOrder != nil {
		if order.UserID == receivedOrder.UserID {
			return usecase.ErrOrderAlreadyExists
		} else {
			return usecase.ErrOrderAlreadyExistsForAnotherUser
		}
	}

	order.Status = model.OrderStatusNew

	if err = os.OrderRepository.Add(ctx, order); err != nil {
		return err
	}

	return nil
}

func (os *OrderUseCase) GetByUser(ctx context.Context, userID int) ([]model.Order, error) {
	orders, err := os.OrderRepository.GetByUser(ctx, userID)
	if err != nil {
		if errors.Is(err, repository.ErrOrderNotFound) {
			return nil, usecase.ErrOrderNotFound
		}
		return nil, err
	}
	return orders, nil
}
