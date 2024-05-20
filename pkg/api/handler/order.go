package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"readon/pkg/api/responses"
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
	PaymentMethoadIdstr := c.Query("paymentMethodId")
	addressidstr := c.Query("addressId")
	paymentMethoadId, err := strconv.Atoi(PaymentMethoadIdstr)
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.RespondWithError(http.StatusBadRequest,
			"Error while converting paymentMethoadId", err.Error()))
		return
	}
	addressID, err := strconv.Atoi(addressidstr)
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.RespondWithError(http.StatusBadRequest,
			"Error while converting addressId", err.Error()))
		return
	}
	userID := c.GetInt("userId")

	// response may be a string razor order id
	razorOrderId, err := cr.OrderUseCase.CreateOrder(userID, addressID, paymentMethoadId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.RespondWithError(http.StatusBadRequest,
			"could not place order ", err.Error()))
		return
	}
	if paymentMethoadId == 1 {
		c.JSON(http.StatusOK, responses.RespondWithSuccess(http.StatusOK, "order placed.", nil))
		return
	}
	c.JSON(http.StatusOK, responses.RespondWithSuccess(http.StatusOK, "order placed.", gin.H{
		"RazopayOrderId ": razorOrderId,
	}))
}

// retrying payment when payment failed
func (cr OrderHAndler) RetryOrder(c *gin.Context) {
	orderidstr := c.Param("orderId")
	orderId, err := strconv.Atoi(orderidstr)
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.RespondWithError(http.StatusBadRequest,
			"error getting param : orderId", err.Error()))
		return
	}
	userID := c.GetInt("userId")

	razorOrderId, err := cr.OrderUseCase.RetryOrder(userID, orderId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.RespondWithError(http.StatusInternalServerError,
			"error while retrying order ", err.Error()))
		return
	}
	c.JSON(http.StatusOK, responses.RespondWithSuccess(http.StatusOK, "order placed.", gin.H{
		"RazopayOrderId ": razorOrderId,
	}))
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

	orderIdStr := c.Param("orderId")
	orderID, err := strconv.Atoi(orderIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.RespondWithError(http.StatusBadRequest,
			"error getting param : orderId", err.Error()))
		return
	}
	userID := c.GetInt("userId")
	err = cr.OrderUseCase.CancelOrder(userID, orderID)
	if err != nil {

		c.JSON(http.StatusInternalServerError, responses.RespondWithError(http.StatusInternalServerError,
			"order not cancelled", err))
		return
	}
	c.JSON(http.StatusOK, responses.RespondWithSuccess(http.StatusOK,
		"oder cancelled : "+orderIdStr, nil))
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
	orderIdStr := c.Param("orderId")
	orderID, err := strconv.Atoi(orderIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.RespondWithError(http.StatusBadRequest,
			"error getting param : orderId", err.Error()))
		return
	}
	userID := c.GetInt("userId")
	order, err := cr.OrderUseCase.GetOrder(userID, orderID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.RespondWithError(http.StatusInternalServerError,
			"couldn't retreive order", err.Error()))
		return
	}

	c.JSON(http.StatusOK, responses.RespondWithSuccess(http.StatusOK, "order retrived ", order))
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

	var pageDetails models.Pagination

	err := c.Bind(&pageDetails)
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.RespondWithError(http.StatusBadRequest,
			"error while binding pagination data ", err.Error()))
	}

	userID := c.GetInt("userId")

	listOfOrders, err := cr.OrderUseCase.ListOrders(userID, pageDetails)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.RespondWithError(http.StatusInternalServerError,
			"couldn't retreive orders", err.Error()))
		return
	}

	c.JSON(http.StatusOK, responses.RespondWithSuccess(http.StatusOK, "order retrived ", listOfOrders))
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
	var pageDetails models.Pagination

	err := c.Bind(&pageDetails)
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.RespondWithError(http.StatusBadRequest,
			"error while binding pagination data ", err.Error()))
	}

	list, err := cr.OrderUseCase.GetAllOrders(pageDetails.Filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.RespondWithError(http.StatusInternalServerError,
			"couldn't retreive orders", err.Error()))
		return
	}
	var orderslist []models.OrdersListing
	copier.Copy(&orderslist, &list)
	c.JSON(http.StatusOK, responses.RespondWithSuccess(http.StatusOK, "order retrived ", list))
}

// to handle weebhook from razorpay on paymentcaptured and payment failed
func (cr OrderHAndler) VerifyPayment(c *gin.Context) {

	var body map[string]interface{}

	if err := json.NewDecoder(c.Request.Body).Decode(&body); err != nil {
		fmt.Println("error while decoding JSON body:", err)
		c.JSON(http.StatusBadRequest, "error while decoding JSON body")
		return
	}
	p, _ := body["payload"].(map[string]interface{})["payment"].(map[string]interface{})["entity"].(map[string]interface{})
	fmt.Println("p ::: ", p)
	var verificationData models.PaymentVerificationData
	var ok bool
	verificationData.RazorOrderId, ok = body["payload"].(map[string]interface{})["payment"].(map[string]interface{})["entity"].(map[string]interface{})["order_id"].(string)
	if !ok || verificationData.RazorOrderId == "" {
		fmt.Println("order_id not found in payload")
		c.JSON(http.StatusBadRequest, "order_id not found in payload")
		return
	}
	verificationData.PaymentStatus, ok = body["payload"].(map[string]interface{})["payment"].(map[string]interface{})["entity"].(map[string]interface{})["status"].(string)
	if !ok || verificationData.RazorOrderId == "" {
		fmt.Println("payment status not found in payload")
		c.JSON(http.StatusBadRequest, "status not found in payload")
		return
	}
	verificationData.RazorPaymentId, ok = body["payload"].(map[string]interface{})["payment"].(map[string]interface{})["entity"].(map[string]interface{})["id"].(string)
	if !ok || verificationData.RazorOrderId == "" {
		fmt.Println("payment_id not found in payload")
		c.JSON(http.StatusBadRequest, "id not found in payload")
		return
	}
	err := cr.OrderUseCase.VerifyPayment(verificationData)
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

func (cr OrderHAndler) DownloadInvoice(c *gin.Context) {
	orderIdStr := c.Param("orderId")
	orderID, err := strconv.Atoi(orderIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.RespondWithError(http.StatusBadRequest,
			"error getting param : orderId", err.Error()))
		return
	}
	userID := c.GetInt("userId")

	//for bypassing cookie :to delete before hosting
	if userID == 0 {
		userID = 1
	}

	invoice, err := cr.OrderUseCase.GetInvoiveData(userID, orderID)
	// invoice, err := cr.OrderUseCase.GetInvoiveData(1, 1)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.RespondWithError(http.StatusInternalServerError,
			"couldn't retreive data", err.Error()))
		return
	}
	c.HTML(http.StatusOK, "invoice.html", invoice)
}

func (cr OrderHAndler) GetChart(c *gin.Context) {
	var pageDetails models.Pagination

	err := c.Bind(&pageDetails)
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.RespondWithError(http.StatusBadRequest,
			"error while binding pagination data ", err.Error()))
	}
	data, err := cr.OrderUseCase.GetChartData(pageDetails)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.RespondWithError(http.StatusInternalServerError,
			"couldn't retreive data", err.Error()))
		return
	}
	c.JSON(http.StatusOK, responses.RespondWithSuccess(http.StatusOK, "chart generated .", data))

}

func (cr OrderHAndler) GetTopTen(c *gin.Context) {
	var pageDetails models.Pagination

	err := c.Bind(&pageDetails)
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.RespondWithError(http.StatusBadRequest,
			"error while binding pagination data ", err.Error()))
	}
	if pageDetails.Filter == 1 {
		data, err := cr.OrderUseCase.GetTopTenCategory(pageDetails)
		if err != nil {
			c.JSON(http.StatusNotFound, nil)
			return
		}
		c.JSON(http.StatusOK, responses.RespondWithSuccess(http.StatusOK, "categories found ", data))
		return
	}

	if pageDetails.Filter == 2 {
		data, err := cr.OrderUseCase.GetTopTenBooks(pageDetails)
		if err != nil {
			c.JSON(http.StatusNotFound, nil)
			return
		}
		c.JSON(http.StatusOK, responses.RespondWithSuccess(http.StatusOK, "top ten books  ", data))
		return
	}

}
