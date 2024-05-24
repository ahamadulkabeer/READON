package repository

import (
	"log"
	"readon/pkg/domain"
	"readon/pkg/models"
	interfaces "readon/pkg/repository/interface"

	"gorm.io/gorm"
)

type CouponDatabase struct {
	DB *gorm.DB
}

func NewCouponRepository(db *gorm.DB) interfaces.CouponRepository {
	return &CouponDatabase{
		DB: db,
	}
}

func (c CouponDatabase) CreateNewCoupon(coupon domain.Coupon) (domain.Coupon, error) {
	err := c.DB.Create(&coupon).Error
	if err != nil {
		log.Println(err)
		return domain.Coupon{}, err
	}
	return coupon, nil
}

func (c CouponDatabase) DeleteCoupon(couponID uint) error {
	err := c.DB.Delete(&domain.Coupon{}, "id = ?", couponID).Error
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (c CouponDatabase) ListAllCoupon(pageDet models.Pagination) ([]domain.Coupon, error) {
	var list []domain.Coupon
	err := c.DB.Find(&list).Limit(pageDet.Size).Error
	if err != nil {
		log.Println(err)
		return []domain.Coupon{}, err
	}
	return list, err
}
