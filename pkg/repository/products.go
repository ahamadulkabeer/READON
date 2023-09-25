package repository

import (
	"context"
	"fmt"
	"readon/pkg/models"
	interfaces "readon/pkg/repository/interface"

	"gorm.io/gorm"
)

type productDatabase struct {
	DB *gorm.DB
}

// Initilising repository

func NewProductRepository(DB *gorm.DB) interfaces.ProductRepository {
	return &productDatabase{
		DB,
	}
}

func (c *productDatabase) ListProducts(ctx context.Context) ([]models.ListingBook, error) {
	var list []models.ListingBook

	fmt.Println("got the products")
	db := c.DB.Table("books").Select("title, author,rating").Limit(8).Find(&list)
	fmt.Println("list :", list)
	return list, db.Error
}
