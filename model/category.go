package model

import "time"

type Category struct {
	ID        int    `json:"id" gorm:"primaryKey;autoIncrement:true"`
	Name      string `json:"name" gorm:"varchar(100)"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Books     []Book `json:"books" gorm:"foreignKey:CategoryID"`
}
