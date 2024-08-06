package repository

import (
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
// this is a

func NewProductRepository(DB *gorm.DB) interfaces.ProductRepository {
	return &productDatabase{
		DB,
	}
}

// func (c *productDatabase) ListProducts() ([]models.ListBook, error) {
// 	var list []models.ListBook

// 	db := c.DB.Table("books").Select("books.id AS book_id, books.title, books.author, books.about, books.rating, books.premium,categories.id AS category_id,categories.name AS category_name").
// 		Joins("JOIN categories ON books.category_id = categories.id").
// 		//Joins("JOIN bookcovers ON books.id = bookcovers.book_id").
// 		Find(&list) //.Limit(8)
// 	return list, db.Error
// }

func (c *productDatabase) ListProductsForUser(pageDet models.Pagination, offset int) ([]domain.Book, error) {
	var list []domain.Book
	query := c.DB.Model(&domain.Book{}).Preload("Category")
	if pageDet.Search != "" {
		query = query.Where(" books.title ILIKE  ? ", fmt.Sprintf("%%%s%%", pageDet.Search))
	}
	// filter by category
	if pageDet.Filter != "" {
		query = query.Where("books.category_id = ?", pageDet.Filter)
	}

	err := query.Offset(offset).Limit(pageDet.Size).Find(&list).Error
	if err != nil {
		return list, err
	}
	return list, err
}

func (c productDatabase) GetTotalNoOfproducts(pageDet models.Pagination) (int, error) {
	var count int64

	query := c.DB.Table("books").Select("id").Where(" title ILIKE  ? ", fmt.Sprintf("%%%s%%", pageDet.Search))
	if pageDet.Filter != "" {
		query = query.Where("category_id = ?", pageDet.Filter)
	}
	err := query.Count(&count).Error
	if err != nil {
		return 0, err
	}
	return int(count), err
}

func (c *productDatabase) AddProduct(product *domain.Book) (int, error) {

	err := c.DB.Create(product).Error
	if err != nil {
		return 0, err
	}
	return int(product.ID), nil
}

func (c *productDatabase) EditProduct(product domain.Book) (domain.Book, error) {
	err := c.DB.Save(&product).Error
	fmt.Println("product id :", product.ID)
	return product, err
}

func (c productDatabase) AddImage(image []byte, book_Id int) error {

	err := c.DB.Exec("INSERT INTO bookcovers(image,book_id) VALUES ( ?, ?)", image, book_Id).Error
	return err
}

func (c productDatabase) GetProduct(BookID int) (domain.Book, error) {
	var book domain.Book

	err := c.DB.Model(&domain.Book{}).Where("id = ?", BookID).Preload("Category").First(&book).Error
	if err != nil {
		return domain.Book{}, err
	}
	return book, nil
}
func (c productDatabase) DeleteProduct(book domain.Book) error {
	err := c.DB.Where("id = ?", book.ID).Delete(&book).Error

	return err
}

func (c productDatabase) ListBookCovers(bookId int) ([]domain.Bookcover, error) {
	var coverList []domain.Bookcover
	err := c.DB.Model(&domain.Bookcover{}).Where("book_id = ?", bookId).Find(&coverList).Error
	return coverList, err
}

func (c productDatabase) GetBookCover(bookId int) (domain.Bookcover, error) {
	var bookCover domain.Bookcover
	err := c.DB.Model(&domain.Bookcover{}).Where("book_id = ?", bookId).First(&bookCover).Error
	if err != nil {
		return domain.Bookcover{}, err
	}
	return bookCover, nil
}

func (c productDatabase) GetNumberOfBookCovers(bookID uint) (int, error) {
	var count int64
	err := c.DB.Model(&domain.Bookcover{}).Where("book_id = ?", bookID).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return int(count), nil
}

func (c productDatabase) GetPrice(bookId int) (float64, error) {
	var price float64
	err := c.DB.Model(&domain.Book{}).Select("price").Where("id = ? ", bookId).First(&price).Error
	if err != nil {
		return 0.0, err
	}
	return price, err
}
