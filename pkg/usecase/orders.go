package usecase

import (
	"readon/pkg/domain"
	interfaces "readon/pkg/repository/interface"
	services "readon/pkg/usecase/interface"
)

type OrderUseCase struct {
	OrderRepo interfaces.OrderRepository
	CartRepo  interfaces.CartRepository
}

func NewOrderUseCase(orepo interfaces.OrderRepository, crepo interfaces.CartRepository) services.OrderUseCase {
	return &OrderUseCase{
		OrderRepo: orepo,
		CartRepo:  crepo,
	}
}

func (c OrderUseCase) CreateOrder(userid, addressid, PaymentMethoadid int) error {
	var order domain.Order
	cart, err := c.CartRepo.GetItems(userid)
	if err != nil {
		return err
	}
	for i := range cart {
		order.UserId = cart[i].UserId
		order.BookId = cart[i].BookId
		order.Quantity = cart[i].Quantity
		order.PaymentMethoadId = uint(PaymentMethoadid)
		order.AdressId = uint(addressid)
		//order.TotalPrice =

		err = c.OrderRepo.CreateOrder(order, userid)
	}
	if err != nil {
		return err
	}
	return nil

}
func (c OrderUseCase) CancelOrder(userid, orderId int) error {
	err := c.OrderRepo.CancelOrder(orderId, userid)
	if err != nil {
		return err
	}
	return nil
}

func (c OrderUseCase) ListOrders(userid int) ([]domain.Order, error) {
	list, err := c.OrderRepo.GetOrders(userid)
	if err != nil {
		return list, err
	}
	return list, err

}

func (c OrderUseCase) GetOrder(userid, orderid int) (domain.Order, error) {
	order, err := c.OrderRepo.GetOrder(userid, orderid)
	if err != nil {
		return order, err
	}
	return order, nil
}
