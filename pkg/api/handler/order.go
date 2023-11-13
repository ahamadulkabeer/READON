package handler

import (
	"net/http"
	"readon/pkg/models"
	services "readon/pkg/usecase/interface"
	"strconv"

	"github.com/gin-gonic/gin"
)

type OrderHAndler struct {
	OrderUseCase services.OrderUseCase
}

func NewOrderHandler(usecase services.OrderUseCase) *OrderHAndler {
	return &OrderHAndler{
		OrderUseCase: usecase,
	}
}

// @Summary Add an order
// @Description Place an order for a user with specified details.
// @Tags order
// @Produce json
// @Param userid query int true "User ID"
// @Param paymentmethoadid query int true "Payment Method ID"
// @Param adressid query int true "Address ID"
// @Success 200 {string} string "Order placed"
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /user/addorder [post]
func (cr OrderHAndler) AddOrder(c *gin.Context) {
	useridstr := c.Query("userid")
	PaymentMethoadIdstr := c.Query("paymentmethoadid")
	addressidstr := c.Query("adressid")
	userid, err := strconv.Atoi(useridstr)
	paymentmethoadid, err := strconv.Atoi(PaymentMethoadIdstr)
	addressid, err := strconv.Atoi(addressidstr)
	if err != nil {
		errResponse := models.ErrorResponse{
			Err:    err.Error(),
			Status: "Error while converting data",
			Hint:   "please try again",
		}
		c.JSON(http.StatusBadRequest, errResponse)
		return
	}
	err = cr.OrderUseCase.CreateOrder(userid, addressid, paymentmethoadid)
	if err != nil {
		errResponse := models.ErrorResponse{
			Err:    err.Error(),
			Status: "could not place order ",
			Hint:   "please try again",
		}
		c.JSON(http.StatusInternalServerError, errResponse)
		return
	}
	c.JSON(http.StatusOK, "oder placed")
}

// @Summary Cancel an order
// @Description Cancel an order for a user with specified details.
// @Tags order
// @Produce json
// @Param userid query int true "User ID"
// @Param orderid query int true "Order ID"
// @Success 200 {string} string "Order cancelled"
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /user/cancelorder [delete]
func (cr OrderHAndler) CancelOrder(c *gin.Context) {
	useridstr := c.Query("userid")
	orderidstr := c.Query("orderid")
	userid, err := strconv.Atoi(useridstr)
	orderid, err := strconv.Atoi(orderidstr)
	if err != nil {
		errResponse := models.ErrorResponse{
			Err:    err.Error(),
			Status: "Error while converting data",
			Hint:   "please try again",
		}
		c.JSON(http.StatusBadRequest, errResponse)
		return
	}
	err = cr.OrderUseCase.CancelOrder(userid, orderid)
	if err != nil {
		errResponse := models.ErrorResponse{
			Err:    err.Error(),
			Status: "order cancelled",
			Hint:   "please try again",
		}
		c.JSON(http.StatusInternalServerError, errResponse)
		return
	}
	c.JSON(http.StatusOK, "oder cancelled : "+orderidstr)
}

// @Summary Get order details
// @Description Retrieve details of a specific order for a user.
// @Tags order
// @Produce json
// @Param userid query int true "User ID"
// @Param orderid query int true "Order ID"
// @Success 200 {object} domain.Order "Order details"
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /user/getorder [get]
func (cr OrderHAndler) GetOrder(c *gin.Context) {
	useridstr := c.Query("userid")
	orderidstr := c.Query("orderid")
	userid, err := strconv.Atoi(useridstr)
	orderid, err := strconv.Atoi(orderidstr)
	if err != nil {
		errResponse := models.ErrorResponse{
			Err:    err.Error(),
			Status: "Error while converting data",
			Hint:   "please try again",
		}
		c.JSON(http.StatusBadRequest, errResponse)
		return
	}
	order, err := cr.OrderUseCase.GetOrder(userid, orderid)
	if err != nil {
		errResponse := models.ErrorResponse{
			Err:    err.Error(),
			Status: "couldn't retreive order",
			Hint:   "please try again",
		}
		c.JSON(http.StatusInternalServerError, errResponse)
		return
	}
	c.JSON(http.StatusOK, order)
}
