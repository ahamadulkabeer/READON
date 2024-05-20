package usecase

import (
	"readon/pkg/api/responses"
	"readon/pkg/models"
	interfaces "readon/pkg/repository/interface"
	services "readon/pkg/usecase/interface"
)

type CouponUsecase struct {
	CouponRepo interfaces.CouponRepository
}

func NewCouponUseCase(couponRepo interfaces.CouponRepository) services.CouponUsecase {
	return &CouponUsecase{
		CouponRepo: couponRepo,
	}
}

func (c CouponUsecase) AddNewCoupon(newCoupon models.Coupon) responses.Response {
	return responses.Response{}
}
