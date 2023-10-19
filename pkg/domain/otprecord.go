package domain

import "gorm.io/gorm"

type Otp struct {
	gorm.Model
	Email string `form:"email"`
	Otp   string `form:"otp"`
}
