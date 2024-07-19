package repository

import (
	"readon/pkg/domain"
	interfaces "readon/pkg/repository/interface"

	"gorm.io/gorm"
)

type CartDatabase struct {
	DB *gorm.DB
}

func NewCartRepository(DB *gorm.DB) interfaces.CartRepository {
	return &CartDatabase{
		DB: DB,
	}
}

func (c CartDatabase) AddItem(item domain.Cart) error {
	err := c.DB.Model(&domain.Cart{}).Create(&item).Error
	return err
}

func (c CartDatabase) UpdateQty(userID, bookID, nQty int) error {
	err := c.DB.Model(&domain.Cart{}).Where("user_id = ? AND book_id = ? ", userID, bookID).Update("Quantity", nQty).Error
	if err != nil {
		return err
	}
	return nil
}
func (c CartDatabase) DeleteItem(userID, bookID int) error {
	err := c.DB.Where("user_id = ? AND book_id = ?", userID, bookID).Delete(&domain.Cart{}).Error
	if err != nil {
		return err
	}
	return nil
}

func (c CartDatabase) GetItems(userID int) ([]domain.Cart, error) {
	var list []domain.Cart
	err := c.DB.Model(&domain.Cart{}).Where("user_id = ?", userID).Preload("Book").Preload("Book.Category").Find(&list).Error
	if err != nil {
		return list, err
	}
	return list, err
}

func (c CartDatabase) CheckForItem(userID, bookID int) (int, error) {
	var qty int64
	err := c.DB.Model(&domain.Cart{}).Select("quantity").Where("user_id = ? AND book_id = ?", userID, bookID).Find(&qty).Error
	if err != nil {
		return 0, err
	}
	return int(qty), nil
}

func (c CartDatabase) GetTotalCartPrice(userID int) (float64, error) {
	var totalPrice float64
	err := c.DB.Model(&domain.Cart{}).Where("user_id = ?", userID).Select("SUM(price) as total_price").Scan(&totalPrice).Error
	if err != nil {
		return 0.0, err
	}
	return totalPrice, nil
}

// i say this is cart repo . but never once i used cart id :)
