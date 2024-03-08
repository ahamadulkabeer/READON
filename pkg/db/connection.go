package db

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	config "readon/pkg/config"
	domain "readon/pkg/domain"
)

func ConnectDatabase(cfg config.Config) (*gorm.DB, error) {
	psqlInfo := fmt.Sprintf("host=%s user=%s dbname=%s port=%s password=%s", cfg.DBHost, cfg.DBUser, cfg.DBName, cfg.DBPort, cfg.DBPassword)
	db, dbErr := gorm.Open(postgres.Open(psqlInfo), &gorm.Config{
		SkipDefaultTransaction: true,
	})

	db.AutoMigrate(&domain.User{})
	db.AutoMigrate(&domain.Admin{})
	db.AutoMigrate(&domain.Book{})
	db.AutoMigrate(&domain.Category{})
	db.AutoMigrate(&domain.Otp{})
	db.AutoMigrate(&domain.Bookcover{})
	db.AutoMigrate(&domain.Address{})
	db.AutoMigrate(&domain.PaymentMethoad{})
	db.AutoMigrate(&domain.Cart{})
	db.AutoMigrate(&domain.Order{})

	return db, dbErr
}
