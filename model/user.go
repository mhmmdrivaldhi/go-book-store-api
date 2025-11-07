package model

import "time"

type User struct {
	ID        int    `json:"id" gorm:"primaryKey;autoIncrement:true"`
	Name      string `json:"name" gorm:"varchar(100)"`
	Email     string `json:"email" gorm:"varchar(100);unique"`
	Password  string `json:"password" gorm:"varchar(255)"`
	Role      string `json:"role" gorm:"varchar(50);default:'user'"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
