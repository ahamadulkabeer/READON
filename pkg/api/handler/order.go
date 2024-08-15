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
)

type OrderHAndler struct {
	OrderUseCase services.OrderUseCase
}

func NewOrderHandler(usecase services.OrderUseCase) *OrderHAndler {
	return &OrderHAndler{
		OrderUseCase: usecase,
	}
}

func (cr OrderHAndler) AddOrder(c *gin.Context) {
	PaymentMethoadIdstr := c.Query("paymentMethodId")
	addressidstr := c.Query("addressId")
	paymentMethoadId, err := strconv.Atoi(PaymentMethoadIdstr)
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.ClientReponse(http.StatusBadRequest,
			"Error while converting paymentMethoadId", err.Error(), nil))
		return
	}
	addressID, err := strconv.Atoi(addressidstr)
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.ClientReponse(http.StatusBadRequest,
			"Error while converting addressId", err.Error(), nil))
		return
	}

	couponsApplied := c.Query("coupons")
	var couponsSlice []string
	err = json.Unmarshal([]byte(couponsApplied), &couponsSlice)
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.ClientReponse(http.StatusBadRequest,
			"Invalid coupons JSON format", err.Error(), nil))
		return
	}

	fmt.Println("coupon applied  : ", couponsSlice)

	userID := c.GetInt("userId")

	response := cr.OrderUseCase.CreateOrder(userID, addressID, paymentMethoadId, couponsSlice)

	c.JSON(response.StatusCode, response)
}

// retrys payment when payment failed
func (cr OrderHAndler) RetryOrder(c *gin.Context) {
	orderidstr := c.Param("orderId")
	orderId, err := strconv.Atoi(orderidstr)
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.ClientReponse(http.StatusBadRequest,
			"error getting param : orderId", err.Error(), nil))
		return
	}
	userID := c.GetInt("userId")

	response := cr.OrderUseCase.RetryOrder(userID, orderId)
	c.JSON(response.StatusCode, response)
}

func (cr OrderHAndler) CancelOrder(c *gin.Context) {

	orderIdStr := c.Param("orderId")
	orderID, err := strconv.Atoi(orderIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.ClientReponse(http.StatusBadRequest,
			"error getting param : orderId", err.Error(), nil))
		return
	}
	userID := c.GetInt("userId")
	response := cr.OrderUseCase.CancelOrder(userID, orderID)
	c.JSON(response.StatusCode, response)
}

func (cr OrderHAndler) GetOrder(c *gin.Context) {
	orderIdStr := c.Param("orderId")
	orderID, err := strconv.Atoi(orderIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.ClientReponse(http.StatusBadRequest,
			"error getting param : orderId", err.Error(), nil))
		return
	}
	userID := c.GetInt("userId")
	response := cr.OrderUseCase.GetOrder(userID, orderID)
	c.JSON(response.StatusCode, response)
}

func (cr OrderHAndler) ListOrders(c *gin.Context) {

	var pageDetails models.Pagination

	err := c.Bind(&pageDetails)
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.ClientReponse(http.StatusBadRequest,
			"error while binding pagination data ", err.Error(), nil))
	}

	userID := c.GetInt("userId")

	response := cr.OrderUseCase.ListOrders(userID, pageDetails)
	c.JSON(response.StatusCode, response)
}

func (cr OrderHAndler) GetAllOrders(c *gin.Context) {
	var pageDetails models.Pagination

	err := c.Bind(&pageDetails)
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.ClientReponse(http.StatusBadRequest,
			"error while binding pagination data ", err.Error(), nil))
	}

	response := cr.OrderUseCase.GetAllOrders(pageDetails.Filter)
	c.JSON(response.StatusCode, response)
}

// to handle webhook from razorpay on paymentcaptured and payment failed
func (cr OrderHAndler) VerifyPayment(c *gin.Context) {

	var body map[string]interface{}

	if err := json.NewDecoder(c.Request.Body).Decode(&body); err != nil {
		fmt.Println("error while decoding JSON body:", err)
		c.JSON(http.StatusBadRequest, responses.ClientReponse(http.StatusBadRequest, "error while decoding JSON body", err.Error(), nil))
		return
	}
	p, _ := body["payload"].(map[string]interface{})["payment"].(map[string]interface{})["entity"].(map[string]interface{})
	fmt.Println("p ::: ", p)
	var verificationData models.PaymentVerificationData
	var ok bool
	verificationData.RazorOrderId, ok = body["payload"].(map[string]interface{})["payment"].(map[string]interface{})["entity"].(map[string]interface{})["order_id"].(string)
	if !ok || verificationData.RazorOrderId == "" {
		fmt.Println("order_id not found in payload")
		c.JSON(http.StatusBadRequest, responses.ClientReponse(http.StatusBadRequest, "order_id not found in payload", "order_id not found in payload", nil))
		return
	}
	verificationData.PaymentStatus, ok = body["payload"].(map[string]interface{})["payment"].(map[string]interface{})["entity"].(map[string]interface{})["status"].(string)
	if !ok || verificationData.RazorOrderId == "" {
		fmt.Println("payment status not found in payload")
		c.JSON(http.StatusBadRequest, responses.ClientReponse(http.StatusBadRequest, "status not found in payload", "status not found in payload", nil))
		return
	}
	verificationData.RazorPaymentId, ok = body["payload"].(map[string]interface{})["payment"].(map[string]interface{})["entity"].(map[string]interface{})["id"].(string)
	if !ok || verificationData.RazorOrderId == "" {
		fmt.Println("payment_id not found in payload")
		c.JSON(http.StatusBadRequest, responses.ClientReponse(http.StatusBadRequest, "id not found in payload", "id not found in payload", nil))
		return
	}
	response := cr.OrderUseCase.VerifyPayment(verificationData)

	c.JSON(response.StatusCode, response)
}

func (cr OrderHAndler) DownloadInvoice(c *gin.Context) {
	orderIdStr := c.Param("orderId")
	orderID, err := strconv.Atoi(orderIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.ClientReponse(http.StatusBadRequest,
			"error getting param : orderId", err.Error(), nil))
		return
	}
	userID := c.GetInt("userId")
	if userID == 0 {
		userID = 1
	}
	response := cr.OrderUseCase.GetInvoiveData(userID, orderID)
	if response.Error != nil {
		c.JSON(response.StatusCode, response)
		return
	}
	invoice := response.Data.(models.InvoiceData)
	c.HTML(http.StatusOK, "invoice.html", invoice)
}

func (cr OrderHAndler) GetChart(c *gin.Context) {
	var pageDetails models.Pagination

	err := c.Bind(&pageDetails)
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.ClientReponse(http.StatusBadRequest,
			"error while binding pagination data ", err.Error(), nil))
	}
	response := cr.OrderUseCase.GetChartData(pageDetails)

	c.JSON(response.StatusCode, response)

}

func (cr OrderHAndler) GetTopTen(c *gin.Context) {
	var pageDetails models.Pagination

	err := c.Bind(&pageDetails)
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.ClientReponse(http.StatusBadRequest,
			"error while binding pagination data ", err.Error(), nil))
	}
	if pageDetails.Filter == "category" {
		response := cr.OrderUseCase.GetTopTenCategory(pageDetails)
		c.JSON(response.StatusCode, response)
		return
	}

	if pageDetails.Filter == "book" {
		response := cr.OrderUseCase.GetTopTenBooks(pageDetails)
		c.JSON(response.StatusCode, response)
		return
	}

}

func (cr OrderHAndler) MakePayment(c *gin.Context) {
	orderID, err := strconv.Atoi(c.Param("orderId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.ClientReponse(http.StatusBadRequest,
			"error parsing query data", err.Error(), nil))
		return
	}
	data := cr.OrderUseCase.GetDataForPaymentpage(orderID)
	c.HTML(200, "paymentpage.html", data)

}
