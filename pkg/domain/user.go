package domain

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name       string `json:"name" gorm:"not null" validate:"required,name"`
	Email      string `json:"email" gorm:"unique;not null" validate:"required,email"`
	Password   string `json:"password" gorm:"not null" validate:"required,password"`
	Wallet     float64
	Premium    bool
	Permission bool `gorm:"not null;default:true"`
}
