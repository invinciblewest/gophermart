package postgres

import (
	"context"
	"database/sql"
	"github.com/invinciblewest/gophermart/internal/logger"
	"github.com/invinciblewest/gophermart/internal/model"
	"go.uber.org/zap"
)

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
		return nil, model.ErrWithdrawalNotFound
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return withdrawals, nil
}
