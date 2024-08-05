package usecase

import (
	"fmt"
	"log"
	"net/http"
	"readon/pkg/api/errorhandler"
	"readon/pkg/api/helpers"
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
		return responses.ClientReponse(http.StatusBadRequest, "error while binding data", err.Error(), nil)
	}
	coupon.ValidFrom, err = time.Parse("2006-01-02", newCoupon.ValidFrom)
	if err != nil {
		fmt.Println("err :", err)
		return responses.ClientReponse(http.StatusBadRequest, "error while binding 'validfrom'", err.Error(), nil)
	}
	coupon.ValidTill, err = time.Parse("2006-01-02", newCoupon.ValidTill)
	if err != nil {
		fmt.Println("err :", err)
		return responses.ClientReponse(http.StatusBadRequest, "error while binding 'validtill'", err.Error(), nil)
	}

	fmt.Println("copied :", coupon)

	// create new coupon
	coupon, err = c.CouponRepo.CreateNewCoupon(coupon)
	if err != nil {
		statusCode, _ := errorhandler.HandleDatabaseError(err)
		return responses.ClientReponse(statusCode, "error while creating coupon", err.Error(), nil)
	}

	// response
	var createdCoupon models.ListCoupons
	copier.Copy(&createdCoupon, &coupon)
	return responses.ClientReponse(http.StatusCreated, "coupon created successfully", nil, createdCoupon)
}

func (c CouponUsecase) DeleteCoupon(couponID uint) responses.Response {

	// check if the coupon exist
	coupon, err := c.CouponRepo.GetCouponByID(couponID)
	if err != nil {
		statusCode, _ := errorhandler.HandleDatabaseError(err)
		log.Println("err :", err)
		return responses.ClientReponse(statusCode, "coupon not deleted ", err.Error(), nil)
	}
	if coupon.ID == 0 {
		return responses.ClientReponse(http.StatusNotFound, fmt.Sprintf("coupon with  id : %d doesn't exist", couponID), "no record found", nil)
	}

	// delete coupon
	err = c.CouponRepo.DeleteCoupon(uint(couponID))
	if err != nil {
		statusCode, _ := errorhandler.HandleDatabaseError(err)
		log.Println("err :", err)
		return responses.ClientReponse(statusCode, "coupon not deleted ", err.Error(), nil)
	}
	//response
	var deletedCoupon models.ListCoupons
	copier.Copy(&deletedCoupon, &coupon)
	return responses.ClientReponse(http.StatusOK, "coupon deleted successfully", nil, deletedCoupon)
}

func (c CouponUsecase) ListAllCoupon(pageDet models.Pagination) responses.Response {

	// retirieves coupons
	pageDet.Size = 10
	list, err := c.CouponRepo.ListAllCoupon(pageDet)
	if err != nil {
		statusCode, _ := errorhandler.HandleDatabaseError(err)
		return responses.ClientReponse(statusCode, "couldn't retrive coupons", err.Error(), nil)
	}

	// if there is co coupon
	if len(list) == 0 {
		return responses.ClientReponse(http.StatusNotFound, "coupons not found", nil, nil)
	}

	// responses
	var coupons []models.ListCoupons
	copier.Copy(&coupons, &list)
	return responses.ClientReponse(http.StatusOK, "coupon returived successfully", nil, models.PaginatedListCoupons{coupons, pageDet})
}

func (c CouponUsecase) IssueCoupon(userID, couponID uint) responses.Response {
	coupon, err := c.CouponRepo.GetCouponByID(couponID)
	if err != nil {
		statusCode, _ := errorhandler.HandleDatabaseError(err)
		return responses.ClientReponse(statusCode, fmt.Sprint("coupon with id : ", couponID, " not doesn't exist"),
			err.Error(), nil)
	}
	couponCode := helpers.GenerateCouponCode(coupon.Prefix)
	fmt.Println("coupon code :", couponCode)
	userCoupon := domain.UserCoupon{
		UserID:     userID,
		CouponID:   couponID,
		CouponCode: couponCode,
	}
	err = c.CouponRepo.IssueCoupon(&userCoupon)
	if err != nil {
		statusCode, _ := errorhandler.HandleDatabaseError(err)
		return responses.ClientReponse(statusCode,
			"couldnt issue coupon to user", err.Error(), nil)
	}
	var userCoup models.ListUserCoupons
	copier.Copy(&userCoup, &userCoupon)
	return responses.ClientReponse(http.StatusCreated,
		fmt.Sprint("coupon issued to user : ", couponCode), nil, userCoup)
}

func (c CouponUsecase) ListCouponsbyUser(userID uint) responses.Response {

	// retrieves users coupons
	list, err := c.CouponRepo.ListCouponsbyUser(userID)
	statusCode, _ := errorhandler.HandleDatabaseError(err)
	if err != nil {
		return responses.ClientReponse(statusCode,
			fmt.Sprint("couldnt retieve coupons on user id : ", userID), err.Error(), nil)

	}

	// if no coupon found
	if len(list) == 0 {
		return responses.ClientReponse(http.StatusOK,
			"No coupons found for user", nil, nil)
	}

	//response
	var userCoupons []models.ListUserCoupons
	copier.Copy(&userCoupons, &list)
	return responses.ClientReponse(http.StatusOK,
		"coupon successfully retrived ", nil, userCoupons)
}
