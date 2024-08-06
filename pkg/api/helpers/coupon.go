package helpers

import (
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"readon/pkg/api/errorhandler"
	"readon/pkg/domain"
	interfaces "readon/pkg/repository/interface"
	"time"
)

func GenerateCouponCode(prefix string) string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"

	// Initialize random number generator
	rand.Seed(time.Now().UnixNano())

	// Generate the code by selecting random characters from the charset
	code := make([]byte, 10)
	for i := range code {
		code[i] = charset[rand.Intn(len(charset))]
	}

	return prefix + string(code)
}

// preparing data for calculation
func CalculateCouponDiscount(coupenRepo interfaces.CouponRepository, couponCodes []string, cart *[]domain.Cart, order *domain.Order) (int, string, error) {

	statusCode, couponsAndDiscounts, message, err := mapCouponsAndDiscounts(couponCodes, coupenRepo, order)
	if err != nil {
		return statusCode, message, err //
	}
	coupons := []domain.Coupon{}
	for couponID := range couponsAndDiscounts {
		coupon, err := coupenRepo.GetCouponByID(couponID)
		if err != nil {
			statusCode, err := errorhandler.HandleDatabaseError(err)
			return statusCode, "db err : couldnt retrieve coupon ", err
		}
		coupons = append(coupons, coupon)

	}
	TotalDiscount := calculateDiscoundFromCart(cart, &coupons)
	order.DiscountedPrice = order.TotalPrice - TotalDiscount
	order.TotalDiscount = TotalDiscount
	return http.StatusOK, "", nil
}

// actual discount calculation
func calculateDiscoundFromCart(cart *[]domain.Cart, coupons *[]domain.Coupon) float64 {
	TotalDiscount := 0.0
	for _, orderItem := range *cart {
		for _, coupon := range *coupons {
			if coupon.ApplicableOn == "any" || coupon.ApplicableOn == "" {
				TotalDiscount += calculateDiscound(&coupon, orderItem.Price)
			} else if coupon.ApplicableOn == "category" && orderItem.Book.Category.Name == coupon.ApplicableCategory {
				TotalDiscount += calculateDiscound(&coupon, orderItem.Price)
			} else if coupon.ApplicableOn == "product" && orderItem.Book.Title == coupon.ApplicableProduct {
				TotalDiscount += calculateDiscound(&coupon, orderItem.Price)
			}
		}
	}
	return TotalDiscount
}

func calculateDiscound(coupon *domain.Coupon, price float64) float64 {
	discount := 0.0
	if coupon.DiscountType == "percentage" {
		discount += (price / 100) * float64(coupon.DiscountAmount)
	} else {
		discount += float64((*coupon).DiscountAmount)
	}
	return discount
}

// check for business logic breaches
func mapCouponsAndDiscounts(couponCodes []string, coupenRepo interfaces.CouponRepository, order *domain.Order) (int, map[uint]float64, string, error) {

	couponsAndDiscounts := make(map[uint]float64)

	for _, couponCode := range couponCodes {
		fmt.Println("coupn :", couponCode)
		found, userCoupon, err := coupenRepo.UserHasCoupon(uint(order.UserID), couponCode)
		fmt.Println("user Coupon :", userCoupon)
		if err != nil {
			statusCode, err := errorhandler.HandleDatabaseError(err)
			return statusCode, map[uint]float64{}, "db error while searching coupon code", err
		}
		if !found {
			fmt.Println("coupon :", couponCode, "not found")
			return http.StatusNotFound, map[uint]float64{}, "invalid coupon code ", nil
		}
		if _, f := couponsAndDiscounts[userCoupon.CouponID]; f {
			return http.StatusUnprocessableEntity, map[uint]float64{}, "cant have duplicate coupons and same kind ", errors.New("can only apply same type of coupon once per order")
		}

		couponsAndDiscounts[userCoupon.CouponID] = 0.0

	}
	return http.StatusOK, couponsAndDiscounts, "", nil
}
