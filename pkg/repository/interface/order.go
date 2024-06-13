package interfaces

import (
	"readon/pkg/domain"
	"readon/pkg/models"
	"time"
)

type OrderRepository interface {
	CreateOrder(order *domain.Order, cart []domain.Cart) error

	GetOrder(userID, orderID int) (domain.Order, error)
	ListOrders(userID int, pageDetails models.Pagination) ([]domain.Order, error)
	GetAllOrders(start, end time.Time) ([]domain.Order, error)
	ListOrderItems(orderID int) ([]domain.OrderItems, error)

	FindTopTenCategories(filter models.Pagination) ([]models.TopTenCategory, error)
	FindTopTenBooks(filter models.Pagination) ([]domain.Book, error)

	GetTotalAmountOfSpan(start, end time.Time, interval string) ([]models.SalesResult, error)

	GetPrice(userID, orderID int) (float64, error)
	FetchRazorOrderID(orderID int) (string, error)
	CheckPymentStatus(orderID int) (string, error)

	UpdatePaymentStatus(paymentData models.PaymentVerificationData) error
	UpdateRazorOrderId(userID, orderID int, RazorOrderID string) error

	//GetFailedRazorOrderIds(userID int) ([]string, error)

	CancelOrder(orderID, userID int) error
	DeleteOrder(orderID, userID int) error
}
