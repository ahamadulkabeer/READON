package interfaces

import (
	"readon/pkg/domain"
	"readon/pkg/models"
)

type OrderUseCase interface {
	CreateOrder(userID, addressID, paymentMethoadID int, coupons []string) (string, error)
	RetryOrder(userID, orderID int) (string, error)
	CancelOrder(userID, orderID int) error
	ListOrders(userID int, pagination models.Pagination) ([]models.OrdersListing, error)
	GetOrder(userID, orderID int) (models.OrdersListing, error)
	GetAllOrders(filter int) ([]domain.Order, error)
	VerifyPayment(paymentData models.PaymentVerificationData) error
	//DeletefailedRazorOrder(userID int) error
	GetInvoiveData(userID, orderID int) (models.InvoiceData, error)
	GetChartData(pageDetails models.Pagination) (models.ChartData, error)

	GetTopTenCategory(filter models.Pagination) ([]models.TopTenCategory, error)
	GetTopTenBooks(filter models.Pagination) ([]models.ListingBook, error)
}
