package usecase

import (
	"fmt"
	"log"
	"net/http"
	"readon/pkg/api/responses"
	"readon/pkg/domain"
	"readon/pkg/models"
	interfaces "readon/pkg/repository/interface"
	services "readon/pkg/usecase/interface"
	"time"

	"github.com/jinzhu/copier"
)

type CouponUsecase struct {
	CouponRepo interfaces.CouponRepository
}

func NewCouponUseCase(couponRepo interfaces.CouponRepository) services.CouponUsecase {
	return &CouponUsecase{
		CouponRepo: couponRepo,
	}
}

func (c CouponUsecase) CreateNewCoupon(newCoupon models.Coupon) responses.Response {
	var coupon domain.Coupon
	// validating coupon details // futher validation needed
	if len(newCoupon.Name) < 4 {
		return responses.ClientReponse(http.StatusBadRequest, "coupon name should be atleast 4 letters", nil, nil)
	}
	err := copier.Copy(&coupon, &newCoupon)
	if err != nil {
		return responses.ClientReponse(http.StatusBadRequest, "error while binding data", err, nil)
	}
	coupon.ValidFrom, err = time.Parse("2006-01-02", newCoupon.ValidFrom)
	if err != nil {
		fmt.Println("err :", err)
		return responses.ClientReponse(http.StatusBadRequest, "error while binding 'validfrom'", err, nil)
	}
	coupon.ValidTill, err = time.Parse("2006-01-02", newCoupon.ValidTill)
	if err != nil {
		fmt.Println("err :", err)
		return responses.ClientReponse(http.StatusBadRequest, "error while binding 'validtill'", err, nil)
	}
	fmt.Println("copied :", coupon)
	coupon, err = c.CouponRepo.CreateNewCoupon(coupon)
	if err != nil {
		log.Println("err:", err)
		return responses.ClientReponse(http.StatusNotFound, "error while creating coupon", err, nil) // ? code ?
	}
	return responses.ClientReponse(http.StatusCreated, "coupon created successfully", nil, coupon)
}

func (c CouponUsecase) DeleteCoupon(couponID uint) responses.Response {
	err := c.CouponRepo.DeleteCoupon(uint(couponID))
	if err != nil {
		log.Println("err :", err)
		return responses.ClientReponse(http.StatusNotFound, "coupon not deleted ", err, nil) // code ?
	}
	return responses.ClientReponse(http.StatusOK, "coupon deleted successfully", nil, nil) // code ?
}

func (c CouponUsecase) ListAllCoupon(pageDet models.Pagination) responses.Response {
	pageDet.Size = 10
	list, err := c.CouponRepo.ListAllCoupon(pageDet)
	if err != nil {
		return responses.ClientReponse(http.StatusNotFound, "couldn't get coupons", err, nil) // code ?
	}
	if len(list) == 0 {
		return responses.ClientReponse(http.StatusNotFound, "coupons not found", nil, nil)
	}
	return responses.ClientReponse(http.StatusOK, "coupon returived successfully", nil, list) // code ?
}
