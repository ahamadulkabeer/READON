package interfaces

import (
	"readon/pkg/domain"
	"readon/pkg/models"
)

type CouponRepository interface {
	CreateNewCoupon(coupon domain.Coupon) (domain.Coupon, error)
	DeleteCoupon(couponID uint) error
	ListAllCoupon(pageDet models.Pagination) ([]domain.Coupon, error)
	ListCoupons(userID int, pageDet models.Pagination) ([]domain.UserCoupon, error)
	GetCouponByID(couponID uint) (domain.Coupon, error)
	IssueCoupon(userCoupon domain.UserCoupon) error
	ListCouponsbyUser(userID uint) ([]domain.UserCoupon, error)
	UserHasCoupon(userID uint, couponCode string) (bool, domain.UserCoupon, error)
	MarkCouponAsRedemed(couponCode string, orderID uint) error
	DeleteCouponUser(couponCode string) error
	MarkCouponAsNotRedeemed(orderID uint) error
}
