package interfaces

import (
	"readon/pkg/api/responses"
	"readon/pkg/models"
)

type CouponUsecase interface {
	CreateNewCoupon(newCoupon models.Coupon) responses.Response
	DeleteCoupon(couponID uint) responses.Response
	ListAllCoupon(pageDet models.Pagination) responses.Response
}
