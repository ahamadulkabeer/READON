package domain

import (
	"time"

	"gorm.io/gorm"
)

type Coupon struct {
	gorm.Model
	Name               string `gorm:"unique; not null"`
	Description        string
	Prefix             string
	DiscountType       string
	DiscountAmount     int
	MaxDiscount        float64
	ApplicableOn       string
	ApplicableCategory string
	ApplicableProduct  string
	ValidFrom          time.Time
	ValidTill          time.Time
	Limited            bool // the number of coupons that can be issued
	MaxQuantity        int  //only relevent if Limited is true
	IsBound            bool
	Expired            bool
}

type UserCoupon struct {
	gorm.Model
	UserID     uint
	CouponID   uint
	CouponCode string
	Redeemed   bool
	RedeemedOn uint
	Coupon     Coupon
}
