package model

import "time"

type OrderItem struct {
	ID        int `json:"id" gorm:"primaryKey;autoIncrement"`
	OrderID   int `json:"order_id"`
	BookID    int `json:"book_id"`
	Price     int `json:"price"`
	Qty       int `json:"qty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}