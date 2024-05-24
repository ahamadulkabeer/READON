package handler

import (
	"fmt"
	"net/http"
	"readon/pkg/api/responses"
	"readon/pkg/models"
	interfaces "readon/pkg/usecase/interface"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CouponHandler struct {
	CouponUseCase interfaces.CouponUsecase
}

func NewCouponHandler(usecase interfaces.CouponUsecase) *CouponHandler {
	return &CouponHandler{
		CouponUseCase: usecase,
	}
}

func (cr CouponHandler) CreateNewCoupon(c *gin.Context) {
	var newCoupon models.Coupon
	err := c.Bind(&newCoupon)
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.ClientReponse(http.StatusBadRequest,
			"couldnt bind form data please check the input", err, nil))
		return
	}
	fmt.Println("coupon : ", newCoupon)

	response := cr.CouponUseCase.CreateNewCoupon(newCoupon)
	c.JSON(response.StatusCode, response)
	fmt.Println("response", response)
}

func (cr CouponHandler) DeleteCoupon(c *gin.Context) {
	couponIdStr := c.Param("id")
	couponID, err := strconv.Atoi(couponIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.ClientReponse(http.StatusBadRequest,
			"error while getting coupon id ", err.Error(), nil))
		return
	}
	fmt.Println("couponID :", couponID)
	response := cr.CouponUseCase.DeleteCoupon(uint(couponID))
	c.JSON(response.StatusCode, response)
	fmt.Println("response", response)
}

func (cr CouponHandler) ListAllCoupon(c *gin.Context) {
	var pageDet models.Pagination
	c.Bind(&pageDet)
	response := cr.CouponUseCase.ListAllCoupon(pageDet)
	c.JSON(response.StatusCode, response)
	fmt.Println("response", response)
}
