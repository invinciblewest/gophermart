package postgres

import (
	"context"
	"database/sql"
	"errors"
	"github.com/invinciblewest/gophermart/internal/logger"
	"github.com/invinciblewest/gophermart/internal/model"
	"go.uber.org/zap"
)

func (r *PGRepository) AddOrder(ctx context.Context, order *model.Order) error {
	err := r.db.QueryRowContext(ctx,
		"INSERT INTO orders (number, user_id, status, accrual) VALUES ($1, $2, $3, $4) RETURNING id, uploaded_at",
		order.Number, order.UserID, order.Status, order.Accrual).Scan(&order.ID, &order.UploadedAt)
	if err != nil {
		return err
	}
	return nil
}

func (r *PGRepository) GetOrderByUser(ctx context.Context, userID int) ([]model.Order, error) {
	rows, err := r.db.QueryContext(ctx,
		"SELECT id, number, user_id, status, accrual, uploaded_at FROM orders WHERE user_id = $1", userID)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		if err = rows.Close(); err != nil {
			logger.Log.Info("failed to close rows", zap.Error(err))
		}
	}(rows)

	var orders []model.Order
	for rows.Next() {
		var order model.Order
		if err = rows.Scan(&order.ID, &order.Number, &order.UserID, &order.Status, &order.Accrual, &order.UploadedAt); err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	if len(orders) == 0 {
		return nil, model.ErrOrderNotFound
	}

	return orders, nil
}

func (r *PGRepository) GetOrderByNumber(ctx context.Context, number string) (*model.Order, error) {
	var order model.Order
	err := r.db.QueryRowContext(ctx,
		"SELECT id, number, user_id, status, accrual, uploaded_at FROM orders WHERE number = $1",
		number).Scan(&order.ID, &order.Number, &order.UserID, &order.Status, &order.Accrual, &order.UploadedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, model.ErrOrderNotFound
		}
		return nil, err
	}
	return &order, nil
}

func (r *PGRepository) UpdateOrderStatus(ctx context.Context, number string, status model.OrderStatus, accrual *model.Amount) error {
	_, err := r.db.ExecContext(ctx,
		"UPDATE orders SET status = $1, accrual = $2 WHERE number = $3",
		status, accrual, number)
	if err != nil {
		return err
	}
	return nil
}

func (r *PGRepository) GetPendingOrders(ctx context.Context) ([]model.Order, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, number, user_id, status, accrual, uploaded_at FROM orders WHERE status IN ($1, $2)`,
		model.OrderStatusNew, model.OrderStatusProcessing)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		if err = rows.Close(); err != nil {
			logger.Log.Info("failed to close rows", zap.Error(err))
		}
	}(rows)

	var orders []model.Order
	for rows.Next() {
		var order model.Order
		if err = rows.Scan(&order.ID, &order.Number, &order.UserID, &order.Status, &order.Accrual, &order.UploadedAt); err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	if len(orders) == 0 {
		return nil, model.ErrOrderNotFound
	}

	return orders, nil
}
