package dto

import "time"

type PaymentCallbackRequest struct {
	OrderId           int    `json:"order_id"`
	TransactionId     string `json:"transaction_id"`
	TransactionStatus string `json:"transaction_status"`
	Transaction_time  time.Time `json:"transaction_time"`
}

type UpdatedShippingStatusRequest struct {
	OrderId int `json:"order_id"`
	Status string `json:"status"`
}

