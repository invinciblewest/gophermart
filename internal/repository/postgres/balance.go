package postgres

import (
	"context"
	"github.com/invinciblewest/gophermart/internal/model"
)

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
