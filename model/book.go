package model

import "time"

type Book struct {
	Id          int    		`json:"id" binding:"required"`
	Title       string 		`json:"title" binding:"required"`
	Description string 		`json:"description" binding:"required"`
	Author      string 		`json:"author" binding:"required"`
	Price       int    		`json:"price" binding:"required"`
	Rating      int			`json:"rating" binding:"required"`
	CreatedAt   time.Time	`json:"created_at"`
	UpdatedAt   time.Time	`json:"updated_at"`
}