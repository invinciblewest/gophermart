package postgres

import (
	"context"
	"database/sql"
	"errors"
	"github.com/invinciblewest/gophermart/internal/helper"
	"github.com/invinciblewest/gophermart/internal/logger"
	"github.com/invinciblewest/gophermart/internal/model"
	"github.com/jackc/pgerrcode"
	"github.com/lib/pq"
	"go.uber.org/zap"
)

type PGRepository struct {
	db *sql.DB
}

func NewPGRepository(db *sql.DB) *PGRepository {
	return &PGRepository{
		db: db,
	}
}

func (r *PGRepository) CreateUser(ctx context.Context, user *model.User) error {
	if user.Login == "" || user.Password == "" {
		return helper.ErrEmptyLoginOrPassword
	}

	err := r.db.QueryRowContext(
		ctx,
		"INSERT INTO users (login, password) VALUES ($1, $2) RETURNING id, created_at",
		user.Login,
		user.Password,
	).Scan(&user.ID, &user.CreatedAt)

	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == pgerrcode.UniqueViolation {
			return helper.ErrUserAlreadyExists
		}
		return err
	}

	return nil
}

func (r *PGRepository) GetUserByLogin(ctx context.Context, login string) (*model.User, error) {
	if login == "" {
		return nil, helper.ErrEmptyLoginOrPassword
	}

	var user model.User

	err := r.db.QueryRowContext(ctx,
		`SELECT id, login, password, created_at FROM users WHERE login = $1`,
		login).Scan(&user.ID, &user.Login, &user.Password, &user.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, helper.ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}

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
		return nil, helper.ErrOrderNotFound
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
			return nil, helper.ErrOrderNotFound
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
		return nil, helper.ErrOrderNotFound
	}

	return orders, nil
}

func (r *PGRepository) CreateWithdrawal(ctx context.Context, withdrawal *model.Withdrawal) error {
	query := `INSERT INTO withdrawals (user_id, order_number, amount)
		VALUES ($1, $2, $3) RETURNING id, processed_at`
	err := r.db.QueryRowContext(ctx, query,
		withdrawal.UserID, withdrawal.OrderNumber, withdrawal.Amount,
	).Scan(&withdrawal.ID, &withdrawal.ProcessedAt)
	if err != nil {
		return err
	}
	return nil
}

func (r *PGRepository) GetWithdrawalByUser(ctx context.Context, userID int) ([]model.Withdrawal, error) {
	query := `SELECT id, user_id, order_number, amount, processed_at FROM withdrawals WHERE user_id = $1`
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		if err = rows.Close(); err != nil {
			logger.Log.Info("failed to close rows", zap.Error(err))
		}
	}(rows)

	var withdrawals []model.Withdrawal
	for rows.Next() {
		var withdrawal model.Withdrawal
		if err = rows.Scan(&withdrawal.ID, &withdrawal.UserID, &withdrawal.OrderNumber,
			&withdrawal.Amount, &withdrawal.ProcessedAt); err != nil {
			return nil, err
		}
		withdrawals = append(withdrawals, withdrawal)
	}

	if len(withdrawals) == 0 {
		return nil, helper.ErrWithdrawalNotFound
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return withdrawals, nil
}

func (r *PGRepository) GetBalanceByUser(ctx context.Context, userID int) (*model.Balance, error) {
	var balance model.Balance

	query := `SELECT
	  COALESCE(accrual_sum, 0) - COALESCE(withdrawn_sum, 0) AS current,
	  COALESCE(withdrawn_sum, 0) AS withdrawn
	FROM
	  (SELECT SUM(accrual) AS accrual_sum FROM orders WHERE user_id = $1 AND status = $2) o,
	  (SELECT SUM(amount) AS withdrawn_sum FROM withdrawals WHERE user_id = $1) w`

	err := r.db.QueryRowContext(ctx, query, userID, model.OrderStatusProcessed).Scan(&balance.Current, &balance.Withdrawn)
	if err != nil {
		return nil, err
	}

	return &balance, nil
}
