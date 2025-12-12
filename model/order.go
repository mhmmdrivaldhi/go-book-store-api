package model

import "time"

type Order struct {
	ID           int    `json:"id" gorm:"primaryKey;autoIncrement:true"`
	UserID       int    `json:"user_id" gorm:"not null"`
	User         User   `json:"user" gorm:"foreignKey:user_id"`
	TotalAmount  int    `json:"total_amount"`
	TotalQty     int    `json:"total_qty"`
	PaymentStatus  string `json:"payment_status"`
	OrderStatus       string `json:"order_status"`
	ExpiredAt time.Time `json:"expired_at"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	OrderItems   []OrderItem `json:"order_items" gorm:"foreignKey:order_id"`
}

