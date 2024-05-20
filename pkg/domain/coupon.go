package domain

import (
	"time"

	"gorm.io/gorm"
)

type Coupon struct {
	gorm.Model
	CouponName      string `gorm:"unique; not null"`
	Description     string
	DeductionType   string
	DeductionAmount int
	ValidFrom       time.Time
	ValidTill       time.Time
}

type UserCoupon struct {
	gorm.Model
	UserID     uint
	CouponID   uint
	CouponCode string
	Redeemed   bool
	Coupon     Coupon
}
