package db

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	config "readon/config"
	domain "readon/pkg/domain"
)

func ConnectDatabase(cfg config.Config) (*gorm.DB, error) {
	psqlInfo := fmt.Sprintf("host=%s user=%s dbname=%s port=%s password=%s", cfg.DBHost, cfg.DBUser, cfg.DBName, cfg.DBPort, cfg.DBPassword)
	db, dbErr := gorm.Open(postgres.Open(psqlInfo), &gorm.Config{
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
