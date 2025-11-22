package model

import "time"

type Book struct {
	Id          int    		`json:"id" gorm:"autoIncrement:true;primaryKey"`
	Title       string 		`json:"title" binding:"required"`
	Description string 		`json:"description" binding:"required"`
	Author      string 		`json:"author" binding:"required"`
	Price       int    		`json:"price" binding:"required"`
	Rating      int			`json:"rating" binding:"required"`
	CategoryID  int			`json:"category_id" gorm:"not null"`
	Category    Category	`json:"category" gorm:"foreignKey:CategoryID"`
	CreatedAt   time.Time	`json:"created_at"`
	UpdatedAt   time.Time	`json:"updated_at"`
}