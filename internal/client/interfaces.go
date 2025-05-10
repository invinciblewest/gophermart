package client

import (
	"context"
	"github.com/invinciblewest/gophermart/internal/model"
)

type AccrualClient interface {
	GetOrderInfo(ctx context.Context, orderNumber string) (*model.AccrualResponse, int, error)
}
