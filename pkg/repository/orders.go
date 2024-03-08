package repository

import (
	"fmt"
	domain "readon/pkg/domain"
	interfaces "readon/pkg/repository/interface"
	"time"

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
	fmt.Println("works here  repo")
	tx := c.DB.Begin()
	err := tx.Save(&order).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Where("user_id = ?", order.UserId).Delete(domain.Cart{}).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return err
}

func (c OrderDatabse) CancelOrder(orderId, userId int) error {
	err := c.DB.Where("user_id = ? AND id = ?", userId, orderId).Delete(&domain.Order{}).Error
	return err
}

func (c OrderDatabse) GetOrders(UserID int) ([]domain.Order, error) {
	var list []domain.Order
	err := c.DB.Model(&domain.Order{}).Where("user_id = ?", UserID).Find(&list).Error
	return list, err
}

func (c OrderDatabse) GetOrder(userId, orderId int) (domain.Order, error) {
	var order domain.Order
	err := c.DB.Model(&domain.Order{}).Where("user_id = ? AND id = ?", userId, orderId).First(&order).Error
	return order, err
}

func (c OrderDatabse) GetAllOrders(start, end time.Time) ([]domain.Order, error) {
	var list []domain.Order
	err := c.DB.Model(&domain.Order{}).Where("created_at >= ? AND created_at < ?", start, end).Find(&list).Error
	if err != nil {
		return nil, err
	}
	return list, err
}

func (c OrderDatabse) UpdatePaymentStatus(PaymentId string) error {
	err := c.DB.Model(&domain.Order{}).Where("payment_id = ?", PaymentId).Update("payment_status", "paid").Error
	if err != nil {
		return err
	}
	return nil
}
