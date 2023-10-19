package repository

import (
	"context"
	"fmt"
	"readon/pkg/domain"
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

	db := c.DB.Table("books").Select("title, author,rating").Find(&list) //.Limit(8)
	fmt.Println("list :", list)
	return list, db.Error
}

func (c *productDatabase) AddProduct(product domain.Book) error {

	err := c.DB.Save(&product).Error

	return err
}

func (c productDatabase) AddImage(image []byte) error {

	err := c.DB.Exec("INSERT INTO bookcovers(image,book_id) VALUES ( ?, ?)", image, 1).Error
	return err
}

func (c productDatabase) GetProducts(bookId int) ([]domain.Book, error) {
	var books []domain.Book
	err := c.DB.Table("books").Where("id = ?", bookId).Find(&books).Error
	return books, err
}
