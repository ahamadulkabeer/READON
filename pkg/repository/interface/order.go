package interfaces

import (
	"readon/pkg/domain"
)

type OrderRepository interface {
	CreateOrder(order domain.Order, userId int) error
	CancelOrder(orderId, userId int) error
	GetOrders(UserID int) ([]domain.Order, error)
	GetOrder(userId, orderId int) (domain.Order, error)
}
