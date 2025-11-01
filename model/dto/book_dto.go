package dto

import "time"

type CreateBookRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description" binding:"required"`
	Author      string `json:"author" binding:"required"`
	Price       int    `json:"price" binding:"required"`
	Rating      int    `json:"rating" binding:"required"`
}

type UpdateBookRequest struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
	Author      *string `json:"author"`
	Price       *int    `json:"price"`
	Rating      *int    `json:"rating"`
}

type BookResponse struct {
	Id          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Author      string `json:"author"`
	Price       int    `json:"price"`
	Rating      int    `json:"rating"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type GeneralResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}