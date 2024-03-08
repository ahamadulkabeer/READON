package repository

import (
	"errors"
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

	db := c.DB.Table("books").Select("books.id AS book_id, books.title, books.author, books.about, books.rating, books.premium,categories.id AS category_id,categories.name AS category_name").
		Joins("JOIN categories ON books.category_id = categories.id").
		//Joins("JOIN bookcovers ON books.id = bookcovers.book_id").
		Find(&list) //.Limit(8)
	fmt.Println("list :", list)
	return list, db.Error
}

func (c *productDatabase) ListProductsForUser(pageDet models.Pagination, offset int) ([]models.ListingBook, error) {
	var list []models.ListingBook

	log.Println("pagedat in repo", pageDet)
	log.Println("offset :", offset)
	//query := c.DB.Table("books").Select("id,title, author,rating").Where(" title ILIKE  ? ", fmt.Sprintf("%%%s%%", pageDet.Search))
	query := c.DB.Table("books").Select("books.id AS book_id, books.title, books.author, books.about, books.rating, books.premium, books.price, categories.id AS category_id,categories.name AS category_name, bookcovers.image ").
		Joins("JOIN categories ON books.category_id = categories.id").
		Joins("JOIN (SELECT DISTINCT ON (book_id) * FROM bookcovers ORDER BY book_id, id) AS bookcovers ON bookcovers.book_id  = books.id")

	if pageDet.Search != "" {
		query = query.Where(" books.title ILIKE  ? ", fmt.Sprintf("%%%s%%", pageDet.Search))
	}
	if pageDet.Filter != 0 {
		query = query.Where("books.category_id = ?", pageDet.Filter)
	}

	err := query.Offset(offset).Limit(pageDet.Size).Find(&list).Error
	if err != nil {
		return list, err
	}
	// for testing purpose
	for i := range list {
		list[i].Image = nil
	}

	//log.Println("list :", list)
	return list, err
}

func (c productDatabase) GetTotalNoOfproducts(pageDet models.Pagination) (int, error) {
	var count int64

	query := c.DB.Table("books").Select("id").Where(" title ILIKE  ? ", fmt.Sprintf("%%%s%%", pageDet.Search))
	if pageDet.Filter != 0 {
		query = query.Where("category_id = ?", pageDet.Filter)
	}
	err := query.Count(&count).Error
	if err != nil {
		return 0, err
	}
	return int(count), err
}

func (c *productDatabase) AddProduct(product domain.Book) (int, error) {

	err := c.DB.Save(&product).Error
	fmt.Println("product id :", product.ID)
	return int(product.ID), err
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

func (c productDatabase) GetProduct(bookId int) (models.ListingBook, error) {
	var book models.ListingBook
	var err error
	db := c.DB.Table("books").Select("books.id AS book_Id, books.title, books.author, books.about, books.rating, books.premium, books.price, categories, bookcovers.image").
		Joins("JOIN categories ON books.category_id = categories.id").
		Joins("JOIN bookcovers ON books.id = bookcovers.book_id").
		Where("books.id = ?", bookId).Find(&book)
	//log.Println("img : ", book.Image)
	if db.RowsAffected == 0 {
		err = errors.New("No record found")
	}
	return book, err
}

func (c productDatabase) DeleteProduct(book domain.Book) error {
	err := c.DB.Where("id = ?", book.ID).Delete(&book).Error

	return err
}

func (c productDatabase) ListBookCovers(bookId int) ([][]byte, error) {
	var coverList [][]byte
	err := c.DB.Table("bookcovers").Select("image").Where("book_id = ?", bookId).Find(&coverList).Error
	return coverList, err
}

func (c productDatabase) GetPrice(bookId int) (float64, error) {
	var price float64
	err := c.DB.Model(&domain.Book{}).Select("price").Where("id = ? ", bookId).First(&price).Error
	if err != nil {
		return 0.0, err
	}
	return price, err
}
