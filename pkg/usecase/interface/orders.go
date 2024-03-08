package interfaces

import "readon/pkg/domain"

type OrderUseCase interface {
	CreateOrder(userid, addressid, PaymentMethoadid int) (string, error)
	ConfirmOrder(userid, addressid, PaymentMethoadid int, order domain.Order) error
	CancelOrder(userid, orderId int) error
	ListOrders(userid int) ([]domain.Order, error)
	GetOrder(userid, orderid int) (domain.Order, error)
	GetAllOrders(filter int) ([]domain.Order, error)
	VerifyPayment(paymentID string) error
}
