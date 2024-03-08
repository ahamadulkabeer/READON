package interfaces

import (
	"readon/pkg/domain"
	"time"
)

type OrderRepository interface {
	CreateOrder(order domain.Order, userId int) error
	CancelOrder(orderId, userId int) error
	GetOrders(UserID int) ([]domain.Order, error)
	GetOrder(userId, orderId int) (domain.Order, error)
	GetAllOrders(start, end time.Time) ([]domain.Order, error)
	UpdatePaymentStatus(PaymentId string) error
}
