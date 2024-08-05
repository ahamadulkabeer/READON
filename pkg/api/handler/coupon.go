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

// CreateNewCoupon godoc
// @Summary Create a new coupon
// @Description Create a new coupon with the provided details
// @Tags Coupon
// @Accept json
// @Produce json
// @Param coupon body models.Coupon true "New Coupon"
// @Success 201 {object} responses.Response{data=models.ListCoupons} "Coupon created successfully"
// @Failure 400 {object} responses.Response "Bad Request"
// @Failure 500 {object} responses.Response "Internal Server Error"
// @Router /coupons [post]
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

// DeleteCoupon godoc
// @Summary Delete a coupon
// @Description Delete a coupon by its ID
// @Tags Coupon
// @Accept json
// @Produce json
// @Param id path int true "Coupon ID"
// @Success 200 {object} responses.Response{data=models.ListCoupons} "Coupon deleted successfully"
// @Failure 400 {object} responses.Response "Bad Request"
// @Failure 404 {object} responses.Response "Not Found"
// @Failure 500 {object} responses.Response "Internal Server Error"
// @Router /coupons/{id} [delete]
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

// ListAllCoupon godoc
// @Summary List all coupons
// @Description Retrieve a paginated list of all coupons .Admin Authentication required.
// @Tags Coupon
// @Accept json
// @Produce json
// @Param page query models.Pagination false "Pagination"
// @Success 200 {object} responses.Response{data=models.PaginatedListCoupons} "Coupons successfully retrieved"
// @Failure 404 {object} responses.Response "Coupons not found"
// @Failure 500 {object} responses.Response "Internal Server Error"
// @Router /coupons [get]
func (cr CouponHandler) ListAllCoupon(c *gin.Context) {
	var pageDet models.Pagination
	c.Bind(&pageDet)
	response := cr.CouponUseCase.ListAllCoupon(pageDet)
	c.JSON(response.StatusCode, response)
	fmt.Println("response", response)
}

func (cr CouponHandler) IssueCouponToUser(c *gin.Context) {
	userIDStr := c.Param("userId")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.ClientReponse(http.StatusBadRequest,
			"coludnt process request ", "invlaid parameters", err.Error()))
		return
	}
	couponIDStr := c.Param("couponId")
	couponID, err := strconv.Atoi(couponIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.ClientReponse(http.StatusBadRequest,
			"coludnt process request ", "invlaid parameters", err.Error()))
		return
	}
	fmt.Println("ids Couponid ", userID, couponID)

	response := cr.CouponUseCase.IssueCoupon(uint(userID), uint(couponID))
	c.JSON(response.StatusCode, response)
}

// ListCouponsbyUser godoc
// @Summary List user's coupons
// @Description Retrieve a list of coupons for a specific user. User authentication required.
// @Tags Coupon
// @Accept json
// @Produce json
// @Success 200 {object} responses.Response{data=[]models.ListUserCoupons} "Coupons successfully retrieved"
// @Failure 404 {object} responses.Response "Not found"
// @Failure 500 {object} responses.Response "Internal Server Error"
// @Router /users/coupons [get]
func (cr CouponHandler) ListCouponsbyUser(c *gin.Context) {
	userID := c.GetInt("userId")
	responses := cr.CouponUseCase.ListCouponsbyUser(uint(userID))
	c.JSON(responses.StatusCode, responses)
}
