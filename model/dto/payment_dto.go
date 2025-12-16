package dto

type UpdatePaymentStatusRequest struct {
	PaymentStatus string `json:"payment_status" binding:"required,oneof=Paid Unpaid"`
}

type UpdateOrderStatusRequest struct {
	OrderStatus string `json:"order_status" binding:"required,oneof=Processing Shipping Delivered"`
}


