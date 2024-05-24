package domain

import (
	"time"

	"gorm.io/gorm"
)

type Coupon struct {
	gorm.Model
	Name           string `gorm:"unique; not null"`
	Description    string
	DiscountType   string
	DiscountAmount int
	ValidFrom      time.Time
	ValidTill      time.Time
	MaxQuantity    int
	IsBound        bool
	Expired        bool
}

type UserCoupon struct {
	gorm.Model
	UserID     uint
	CouponID   uint
	CouponCode string
	Redeemed   bool
	Coupon     Coupon
}
