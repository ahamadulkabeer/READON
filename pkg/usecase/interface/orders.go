package interfaces

import "readon/pkg/domain"

type OrderUseCase interface {
	CreateOrder(userid, addressid, PaymentMethoadid int) error
	CancelOrder(userid, orderId int) error
	ListOrders(userid int) ([]domain.Order, error)
	GetOrder(userid, orderid int) (domain.Order, error)
}
