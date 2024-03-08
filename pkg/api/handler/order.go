package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"readon/pkg/models"
	services "readon/pkg/usecase/interface"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
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
	// responsee may be a string or aRazororderId
	fmt.Println("works here  hndlr")
	response, err := cr.OrderUseCase.CreateOrder(userid, addressid, paymentmethoadid)
	if err != nil {
		errResponse := models.ErrorResponse{
			Err:    err.Error(),
			Status: "could not place order ",
			Hint:   "please try again",
		}
		c.JSON(http.StatusInternalServerError, errResponse)
		return
	}
	fmt.Println("till end here check")
	if paymentmethoadid == 2 {
		/*c.HTML(http.StatusOK, "paymentpage.html", gin.H{
			"orderid":    response,
			"FinalPrice": 1,
			"Username":   "anyname",
		})*/
		c.JSON(200, "https://www.razorpay.com/payment-link/"+string(response))
		return
	}
	c.JSON(http.StatusOK, response)
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
// @Success 200 {object} models.OrdersListing "Order details"
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
	var orders models.OrdersListing
	copier.Copy(&orders, &order)
	c.JSON(http.StatusOK, orders)
}

// @Summary list a users orders
// @Description Retrieve all orders for a user.
// @Tags order
// @Produce json
// @Param userid query int true "User ID"
// @Success 200 {object} []models.OrdersListing "Order details"
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /user/listorder [get]
func (cr OrderHAndler) ListOrders(c *gin.Context) {
	useridstr := c.Query("userid")
	userid, err := strconv.Atoi(useridstr)
	if err != nil {
		errResponse := models.ErrorResponse{
			Err:    err.Error(),
			Status: "Error while converting data",
			Hint:   "please try again",
		}
		c.JSON(http.StatusBadRequest, errResponse)
		return
	}
	list, err := cr.OrderUseCase.ListOrders(userid)
	if err != nil {
		errResponse := models.ErrorResponse{
			Err:    err.Error(),
			Status: "couldn't retreive order",
			Hint:   "please try again",
		}
		c.JSON(http.StatusInternalServerError, errResponse)
		return
	}
	var orders []models.OrdersListing
	copier.Copy(&orders, &list)
	c.JSON(http.StatusOK, orders)
}

// @Summary Get all orders
// @Description Retrieve all orders based on a filter.
// @Tags order
// @Produce json
// @Param filter query int true "Filter parameter"
// @Success 200 {array} []models.OrdersListing "List of orders"
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /admin/allorders [get]
func (cr OrderHAndler) GetAllOrders(c *gin.Context) {
	filterStr := c.Query("filter")
	filter, err := strconv.Atoi(filterStr)
	if err != nil {
		errResponse := models.ErrorResponse{
			Err:    err.Error(),
			Status: "Error while converting data",
			Hint:   "please try again",
		}
		c.JSON(http.StatusBadRequest, errResponse)
		return
	}

	list, err := cr.OrderUseCase.GetAllOrders(filter)
	if err != nil {
		errResponse := models.ErrorResponse{
			Err:    err.Error(),
			Status: "couldn't retreive orders",
			Hint:   "please try again",
		}
		c.JSON(http.StatusInternalServerError, errResponse)
		return
	}
	var orderslist []models.OrdersListing
	copier.Copy(&orderslist, &list)
	c.JSON(http.StatusOK, orderslist)
}

// to handle weebhook from razorpay on paymentcaptured and payment failed
func (cr OrderHAndler) VerifyPayment(c *gin.Context) {

	var body map[string]interface{}

	if err := json.NewDecoder(c.Request.Body).Decode(&body); err != nil {
		fmt.Println("error while decoding JSON body:", err)
		c.JSON(http.StatusBadRequest, "error while decoding JSON body")
		return
	}

	razorOrderID, ok := body["payload"].(map[string]interface{})["payment"].(map[string]interface{})["entity"].(map[string]interface{})["order_id"].(string)
	if !ok || razorOrderID == "" {
		fmt.Println("order_id not found in payload")
		c.JSON(http.StatusBadRequest, "order_id not found in payload")
		return
	}
	err := cr.OrderUseCase.VerifyPayment(razorOrderID)
	if err != nil {
		errResponse := models.ErrorResponse{
			Err:    err.Error(),
			Status: "could not update payment details",
			Hint:   "please try again",
		}
		fmt.Println("err : ", errResponse)
		return
	}
	c.JSON(http.StatusOK, "payment verified")
}
