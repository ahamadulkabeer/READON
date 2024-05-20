package interfaces

import (
	"readon/pkg/api/responses"
	"readon/pkg/models"
)

type CouponUsecase interface {
	AddNewCoupon(newCoupon models.Coupon) responses.Response
}
