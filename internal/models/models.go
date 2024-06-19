package models

import "time"

// NEW заказ загружен в систему, но не попал в обработку;
const NEW = "NEW"

// PROCESSING вознаграждение за заказ рассчитывается;
const PROCESSING = "PROCESSING"

// INVALID система расчёта вознаграждений отказала в расчёте;
const INVALID = "INVALID"

// PROCESSED данные по заказу проверены и информация о расчёте успешно
const PROCESSED = "PROCESSED"

type Order struct {
	OrderID    string    `json:"order_id"`
	Accrual    *float32  `json:"accrual"`
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
