package dto

import "time"

type OrderResponse struct {
	ID            int                 `json:"id"`
	UserId        int                 `json:"user_id"`
	User		  UserResponse        `json:"user"`
	TotalQty      int                 `json:"total_qty"`
	TotalAmount    int                 `json:"total_amount"`
	OrderStatus   string              `json:"order_status"`
	PaymentStatus string              `json:"payment_status"`
	ShippingStatus string              `json:"shipping_status"`
	ExpiredAt     time.Time           `json:"expired_at"`
	Items         []OrderItemResponse `json:"items"`
}

type OrderItemResponse struct {
	BookID int `json:"book_id"`
	Qty    int `json:"qty"`
	Price  int `json:"price"`
}