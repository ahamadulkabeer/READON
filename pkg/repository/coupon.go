package repository

import (
	"fmt"
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

func (c CouponDatabase) ListCoupons(userID int, pageDet models.Pagination) ([]domain.UserCoupon, error) {
	var list []domain.UserCoupon
	err := c.DB.Where("userid = ?", userID).Find(&list).Error
	if err != nil {
		fmt.Println("Db err :", err.Error())
		return list, err

	}
	return list, nil
}

func (c CouponDatabase) GetCouponByID(couponID uint) (domain.Coupon, error) {
	var coupon domain.Coupon
	err := c.DB.Model(&domain.Coupon{}).Where("id = ?", couponID).First(&coupon).Error
	if err != nil {
		return domain.Coupon{}, err
	}
	return coupon, nil
}

func (c CouponDatabase) IssueCoupon(userCoupon domain.UserCoupon) error {
	err := c.DB.Create(&userCoupon).Error
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (c CouponDatabase) ListCouponsbyUser(userID uint) ([]domain.UserCoupon, error) {
	var list []domain.UserCoupon
	err := c.DB.Where("user_id = ?", userID).Find(&list).Error
	if err != nil {
		log.Println(err)
		return []domain.UserCoupon{}, err
	}
	return list, nil
}

func (c CouponDatabase) UserHasCoupon(userID uint, couponCode string) (bool, domain.UserCoupon, error) {
	var userCoupon domain.UserCoupon
	err := c.DB.Where("user_id = ? AND coupon_code = ?", userID, couponCode). //Preload("Coupon").
											Find(&userCoupon).Error
	if err != nil {
		return false, domain.UserCoupon{}, err
	}
	return true, userCoupon, nil
}

func (c CouponDatabase) MarkCouponAsRedemed(couponCode string, orderID uint) error {
	err := c.DB.Model(&domain.UserCoupon{}).Where("coupon_code = ?", couponCode).
		Updates(map[string]interface{}{"redeemed": true, "redeemed_on": orderID}).Error
	if err != nil {
		fmt.Println("db err :", err)
		return err
	}
	return nil
}

func (c CouponDatabase) DeleteCouponUser(couponCode string) error {
	err := c.DB.Where("coupon_code = ?", couponCode).Delete(&domain.UserCoupon{}).Error

	if err != nil {
		fmt.Println("db err :", err)
		return err
	}
	return nil
}

func (c CouponDatabase) MarkCouponAsNotRedeemed(orderID uint) error {
	err := c.DB.Model(&domain.UserCoupon{}).Where("redeemed_on = ?", orderID).
		Updates(map[string]interface{}{"redeemed": false, "redeemed_on": nil}).Error

	if err != nil {
		log.Println("db err : ", err)
		return err
	}
	return nil
}
