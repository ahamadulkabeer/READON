package interfaces

import (
	"readon/pkg/domain"
	"readon/pkg/models"
)

type CouponRepository interface {
	CreateNewCoupon(coupon domain.Coupon) (domain.Coupon, error)
	DeleteCoupon(couponID uint) error
	ListAllCoupon(pageDet models.Pagination) ([]domain.Coupon, error)
}
