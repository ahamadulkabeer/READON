package interfaces

import (
	"readon/pkg/domain"
	"readon/pkg/models"
)

type ProductRepository interface {
	//ListProducts() ([]models.ListBook, error)
	ListProductsForUser(models.Pagination, int) ([]domain.Book, error)
	GetTotalNoOfproducts(pageDet models.Pagination) (int, error)
	AddProduct(product *domain.Book) (int, error)
	EditProduct(product domain.Book) (domain.Book, error)
	AddImage([]byte, int) error
	//GetProduct(bookID int) (models.ListingBook, error)
	GetProduct(BookID int) (domain.Book, error)
	DeleteProduct(domain.Book) error

	GetNumberOfBookCovers(bookID uint) (int, error)
	GetBookCover(bookId int) (domain.Bookcover, error)
	ListBookCovers(bookId int) ([]domain.Bookcover, error)
	GetPrice(bookId int) (float64, error)
}
