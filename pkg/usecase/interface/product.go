package interfaces

import (
	domain "readon/pkg/domain"
	"readon/pkg/models"
)

type ProductUseCase interface {
	ListProducts() ([]models.ListingBook, error)
	ListProductsForUser(*models.Pagination) ([]models.ListingBook, error)
	Addproduct(pdct models.Product) (error, error)
	GetProduct(bookId int) (domain.Book, error)
	DeleteProduct(int) error
}
