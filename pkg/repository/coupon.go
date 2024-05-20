package repository

import (
	"readon/pkg/domain"
	interfaces "readon/pkg/repository/interface"

	"gorm.io/gorm"
)

type CouponDatabase struct {
	DB *gorm.DB
}

func NewCouponRepository(db *gorm.DB) interfaces.CouponRepository {
	return &CouponDatabase{
		DB: db,
	}
}

func (c CouponDatabase) CreateNewCoupon(coupon domain.Coupon) error {
	return nil
}
