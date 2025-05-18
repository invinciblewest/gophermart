package model

type AccrualResponse struct {
	Order   string      `json:"order"`
	Status  OrderStatus `json:"status"`
	Accrual Amount      `json:"accrual,omitempty"`
}
