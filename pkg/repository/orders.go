package repository

import (
	domain "readon/pkg/domain"
	interfaces "readon/pkg/repository/interface"

	"gorm.io/gorm"
)

type OrderDatabse struct {
	DB *gorm.DB
}

func NewOrdersRepository(db *gorm.DB) interfaces.OrderRepository {
	return &OrderDatabse{
		DB: db,
	}
}

func (c OrderDatabse) CreateOrder(order domain.Order, userId int) error {
	//  ned a trascation here indstead
	err := c.DB.Save(&order).Error
	return err
}

func (c OrderDatabse) CancelOrder(orderId, userId int) error {
	err := c.DB.Where("user_id = ? AND order_id = ?", userId, orderId).Delete(&domain.Order{}).Error
	return err
}

func (c OrderDatabse) GetOrders(UserID int) ([]domain.Order, error) {
	var list []domain.Order
	err := c.DB.Model(&domain.Order{}).Where("user_id = ?", UserID).Find(&list).Error
	return list, err
}

func (c OrderDatabse) GetOrder(userId, orderId int) (domain.Order, error) {
	var order domain.Order
	err := c.DB.Model(&domain.Order{}).Where("user_id = ? AND order_id = ?", userId, orderId).First(&order).Error
	return order, err
}
