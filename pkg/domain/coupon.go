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
	ApplicableOn       string
	ApplicableCategory string
	ApplicableProduct  string
	ValidFrom          time.Time
	ValidTill          time.Time
	MaxQuantity        int
	IsBound            bool
	Expired            bool
}

// currently the coupon user has belongsto relationship to coupon
// but shold also have a has many relation with user ? ? ?
type UserCoupon struct {
	gorm.Model
	UserID     uint
	CouponID   uint
	CouponCode string
	Redeemed   bool
	RedeemedOn uint
	Coupon     Coupon
}
