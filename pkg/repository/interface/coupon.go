package interfaces

import (
	"readon/pkg/domain"
)

type CouponRepository interface {
	CreateNewCoupon(coupon domain.Coupon) error
}
