package models

import "time"

type Order struct {
	OrderID    string    `json:"order_id"`
	Accrual    float32   `json:"accrual"`
	Status     string    `json:"status"`
	UploadedAt time.Time `json:"uploaded_at"`
	UserID     int       `json:"user_id"`
}

type Withdraw struct {
	OrderNum    string    `json:"order_num"`
	ProcessedAt time.Time `json:"processed_at"`
	Sum         float32   `json:"sum"`
	UserID      int       `json:"user_id"`
}
