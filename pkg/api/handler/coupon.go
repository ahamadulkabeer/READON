package handler

import (
	"fmt"
	"readon/pkg/models"
	interfaces "readon/pkg/usecase/interface"

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
	c.Bind(&newCoupon)
	fmt.Println("coupon : ", newCoupon)
}
