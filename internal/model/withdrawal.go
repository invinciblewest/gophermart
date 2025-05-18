package model

import "time"

type Withdrawal struct {
	ID          int       `json:"-"`
	UserID      int       `json:"-"`
	OrderNumber string    `json:"order"`
	Amount      Amount    `json:"sum"`
	ProcessedAt time.Time `json:"processed_at"`
}

type WithdrawRequest struct {
	Order string `json:"order"`
	Sum   Amount `json:"sum"`
}
