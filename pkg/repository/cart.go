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

// just getting the whole data ecxept the cart id pk
// i dont dont knoe why i amnot even using the pk here !!
func (c CartDatabase) AddItem(item domain.Cart, userId int) error {
	err := c.DB.Model(&domain.Cart{}).Save(&item).Error
	return err
}

func (c CartDatabase) UpdateQty(userId, bookId, NQty int) error {
	err := c.DB.Model(&domain.Cart{}).Where("user_id = ? AND book_id = ? ", userId, bookId).Update("Quantity", NQty).Error
	return err
}
func (c CartDatabase) DeleteItem(userId, bookId int) error {
	err := c.DB.Where("user_id = ? AND book_id = ?", userId, bookId).Delete(&domain.Cart{}).Error
	return err
}
func (c CartDatabase) GetItems(userId int) ([]domain.Cart, error) {
	var list []domain.Cart
	err := c.DB.Model(&domain.Cart{}).Where("user_id = ?", userId).Find(&list).Error

	return list, err
}

// a little bit confusion here
// it returnsthe qty as zero as there is no results ...
// and it interepted as  there is no record exist !
func (c CartDatabase) CheckForItem(userId, bookId int) (int, error) {
	var qty int64
	result := c.DB.Model(&domain.Cart{}).Select("quantity").Where("user_id = ? AND book_id = ?", userId, bookId).Find(&qty)
	if result.Error != nil {
		return 0, result.Error
	}
	return int(qty), nil
}
