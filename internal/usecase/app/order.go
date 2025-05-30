package app

import (
	"context"
	"errors"
	"github.com/invinciblewest/gophermart/internal/helper"
	"github.com/invinciblewest/gophermart/internal/model"
	"github.com/invinciblewest/gophermart/internal/repository"
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
		return model.ErrInvalidOrderNumber
	}

	receivedOrder, err := os.OrderRepository.GetOrderByNumber(ctx, order.Number)
	if err != nil && !errors.Is(err, model.ErrOrderNotFound) {
		return err
	}

	if receivedOrder != nil {
		if order.UserID == receivedOrder.UserID {
			return model.ErrOrderAlreadyExists
		} else {
			return model.ErrOrderAlreadyExistsForAnotherUser
		}
	}

	order.Status = model.OrderStatusNew

	if err = os.OrderRepository.AddOrder(ctx, order); err != nil {
		return err
	}

	return nil
}

func (os *OrderUseCase) GetByUser(ctx context.Context, userID int) ([]model.Order, error) {
	orders, err := os.OrderRepository.GetOrderByUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	return orders, nil
}
