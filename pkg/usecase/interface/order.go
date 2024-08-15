package interfaces

import (
	"readon/pkg/api/responses"
	"readon/pkg/models"
)

type OrderUseCase interface {
	CreateOrder(userID, addressID, paymentMethoadID int, coupons []string) responses.Response
	RetryOrder(userID, orderID int) responses.Response
	CancelOrder(userID, orderID int) responses.Response
	ListOrders(userID int, pagination models.Pagination) responses.Response
	GetOrder(userID, orderID int) responses.Response
	GetAllOrders(filter string) responses.Response
	VerifyPayment(paymentData models.PaymentVerificationData) responses.Response
	//DeletefailedRazorOrder(userID int) error
	GetInvoiveData(userID, orderID int) responses.Response
	GetChartData(pageDetails models.Pagination) responses.Response

	GetTopTenCategory(filter models.Pagination) responses.Response
	GetTopTenBooks(filter models.Pagination) responses.Response

	GetDataForPaymentpage(orderID int) models.PaymentPageData
}
