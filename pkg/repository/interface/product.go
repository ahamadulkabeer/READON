package interfaces

import (
	"readon/pkg/domain"
	"readon/pkg/models"
)

type ProductRepository interface {
	ListProducts() ([]models.ListingBook, error)
	ListProductsForUser(models.Pagination, int) ([]models.ListingBook, int, error)
	AddProduct(product domain.Book) error
	AddImage(image []byte) error
	GetProduct(int) (domain.Book, error)
	DeleteProduct(domain.Book) error
}
