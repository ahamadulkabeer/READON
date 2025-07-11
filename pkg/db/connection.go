package db

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	domain "readon/pkg/domain"
)

func ConnectDatabase(dsn string) (*gorm.DB, error) {

	db, dbErr := gorm.Open(postgres.Open(dsn), &gorm.Config{
		SkipDefaultTransaction: true,
	})

	err := db.AutoMigrate(&domain.User{})
	if err != nil {
		log.Fatalf("Failed to migrate database schema:User: %v", err)
	}
	db.AutoMigrate(&domain.Admin{})
	db.AutoMigrate(&domain.Book{})
	db.AutoMigrate(&domain.Category{})
	db.AutoMigrate(&domain.Otp{})
	db.AutoMigrate(&domain.Bookcover{})
	db.AutoMigrate(&domain.Address{})
	db.AutoMigrate(&domain.PaymentMethod{})
	db.AutoMigrate(&domain.Cart{})
	db.AutoMigrate(&domain.Order{})
	db.AutoMigrate(&domain.OrderItems{})
	err = db.AutoMigrate(&domain.Coupon{})
	if err != nil {
		log.Fatalf("Failed to migrate database schema:User: %v", err)
	}
	err = db.AutoMigrate(&domain.UserCoupon{})
	if err != nil {
		log.Fatalf("Failed to migrate database schema:User: %v", err)
	}
	err = db.AutoMigrate(&domain.WalletHistory{})
	if err != nil {
		log.Fatalf("Failed to migrate database schema:User: %v", err)
	}
	return db, dbErr
}
