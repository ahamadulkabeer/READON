package repository

import (
	"fmt"
	"log"
	"readon/pkg/domain"
	"readon/pkg/models"
	interfaces "readon/pkg/repository/interface"

	"gorm.io/gorm"
)

type productDatabase struct {
	DB *gorm.DB
}

// Initilising repository
// this is a 

func NewProductRepository(DB *gorm.DB) interfaces.ProductRepository {
	return &productDatabase{
		DB,
	}
}

func (c *productDatabase) ListProducts() ([]models.ListingBook, error) {
	var list []models.ListingBook

	db := c.DB.Table("books").Select("id,title, author,rating").Find(&list) //.Limit(8)
	fmt.Println("list :", list)
	return list, db.Error
}

func (c *productDatabase) ListProductsForUser(pageDet models.Pagination, offset int) ([]models.ListingBook, int, error) {
	var list []models.ListingBook
	var count int64
	log.Println("pagedat in repo", pageDet)
	log.Println("offset :", offset)
	query := c.DB.Table("books").Select("id,title, author,rating").Where(" title ILIKE  ? ", fmt.Sprintf("%%%s%%", pageDet.Search))
	if pageDet.Filter != 0 {
		query = query.Where("category_id = ?", pageDet.Filter)
	}
	err := query.Offset(offset).Limit(pageDet.Size).Find(&list).Error
	if err != nil {
		return list, 0, err
	}
	log.Println("list :", list)
	query = c.DB.Table("books").Select("id,title, author,rating").Where(" title ILIKE  ? ", fmt.Sprintf("%%%s%%", pageDet.Search))
	if pageDet.Filter != 0 {
		query = query.Where("category_id = ?", pageDet.Filter)
	}
	err = query.Count(&count).Error
	if err != nil {
		return list, 0, err
	}
	return list, int(count), err
}

func (c *productDatabase) AddProduct(product domain.Book) error {

	err := c.DB.Save(&product).Error

	return err
}

func (c productDatabase) AddImage(image []byte) error {

	err := c.DB.Exec("INSERT INTO bookcovers(image,book_id) VALUES ( ?, ?)", image, 1).Error
	return err
}

func (c productDatabase) GetProduct(bookId int) (domain.Book, error) {
	var books domain.Book
	err := c.DB.Table("books").Where("id = ?", bookId).Find(&books).Error
	return books, err
}

func (c productDatabase) DeleteProduct(book domain.Book) error {
	err := c.DB.Delete(&book).Error

	return err
}
