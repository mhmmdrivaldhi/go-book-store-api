package dto

import "time"

type CreateOrderRequest struct {
	UserId int `json:"user_id"`
}

type OrderResponse struct {
	ID            int                 `json:"id"`
	UserId        int                 `json:"user_id"`
	PaymentStatus string              `json:"payment_status"`
	OrderStatus   string              `json:"order_status"`
	ExpiredAt     time.Time           `json:"expired_at"`
	TotalQty      int                 `json:"total_qty"`
	TotalAmount    int                 `json:"total_amount"`
	Items         []OrderItemResponse `json:"items"`
}

type OrderItemResponse struct {
	BookID int `json:"book_id"`
	Qty    int `json:"qty"`
	Price  int `json:"price"`
}