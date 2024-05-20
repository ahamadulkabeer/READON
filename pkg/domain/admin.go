package domain

import "gorm.io/gorm"

type Admin struct {
	gorm.Model
	Name       string `json:"name" gorm:"not null"`
	Email      string `json:"email" gorm:"unique;not null"`
	Password   string `json:"password" gorm:"not null"`
	Permission bool
}
