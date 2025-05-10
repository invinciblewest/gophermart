package app

import (
	"context"
	"errors"
	"github.com/invinciblewest/gophermart/internal/client/accrual"
	"github.com/invinciblewest/gophermart/internal/logger"
	"github.com/invinciblewest/gophermart/internal/model"
	"github.com/invinciblewest/gophermart/internal/repository"
	"go.uber.org/zap"
	"sync"
	"time"
)

type AccrualProcessor struct {
	orderRepository repository.OrderRepository
	accrualClient   *accrual.Client
	workerCount     int
}

func NewAccrualProcessor(orderRepository repository.OrderRepository, accrualClient *accrual.Client) *AccrualProcessor {
	return &AccrualProcessor{
		orderRepository: orderRepository,
		accrualClient:   accrualClient,
	}
}

func (p *AccrualProcessor) Run(ctx context.Context, interval int, workerCount int) {
	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			p.processPendingOrders(ctx, workerCount)
		}
	}
}

func (p *AccrualProcessor) processPendingOrders(ctx context.Context, workerCount int) {
	orders, err := p.orderRepository.GetPending(ctx)
	if err != nil {
		if errors.Is(err, repository.ErrOrderNotFound) {
			logger.Log.Info("no pending orders found")
			return
		}
		logger.Log.Error("failed to get pending orders", zap.Error(err))
		return
	}

	jobs := make(chan model.Order)
	var wg sync.WaitGroup

	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for order := range jobs {
				p.processOrder(ctx, order)
			}
		}()
	}

outerLoop:
	for _, order := range orders {
		select {
		case jobs <- order:
		case <-ctx.Done():
			break outerLoop
		}
	}
	close(jobs)

	wg.Wait()
}

func (p *AccrualProcessor) processOrder(ctx context.Context, order model.Order) {
	reqCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	response, retryAfter, err := p.accrualClient.GetOrderInfo(reqCtx, order.Number)
	if err != nil {
		logger.Log.Info("failed to get order info", zap.String("order_number", order.Number), zap.Error(err))
		return
	}

	if retryAfter > 0 {
		logger.Log.Info("accrual service is busy, retrying later",
			zap.String("order_number", order.Number), zap.Int("retry_after", retryAfter))
		time.Sleep(time.Duration(retryAfter) * time.Second)
		return
	}

	if response == nil {
		logger.Log.Info("empty response from accrual service", zap.String("order_number", order.Number))
		return
	}

	var newAccrual *float64
	if response.Status == model.OrderStatusProcessed {
		newAccrual = &response.Accrual
	}

	if err = p.orderRepository.UpdateStatus(ctx, order.Number, response.Status, newAccrual); err != nil {
		logger.Log.Info("failed to update order accrual", zap.String("order_number", order.Number), zap.Error(err))
		return
	}
}
