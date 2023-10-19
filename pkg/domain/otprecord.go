package domain

import "gorm.io/gorm"

type Otp struct {
	gorm.Model
	Email string `form:"email" json:"email"`
	Otp   string `form:"otp" json:"otp"`
}
