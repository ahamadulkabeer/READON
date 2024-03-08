package interfaces

import (
	"readon/pkg/domain"
	"readon/pkg/models"
)

type ProductRepository interface {
	ListProducts() ([]models.ListingBook, error)
	ListProductsForUser(models.Pagination, int) ([]models.ListingBook, error)
	GetTotalNoOfproducts(pageDet models.Pagination) (int, error)
	AddProduct(product domain.Book) (int, error)
	EditProduct(product domain.Book) (domain.Book, error)
	AddImage([]byte, int) error
	GetProduct(int) (models.ListingBook, error)
	DeleteProduct(domain.Book) error

	ListBookCovers(bookId int) ([][]byte, error)
	GetPrice(bookId int) (float64, error)
}
