package interfaces

import (
	"readon/pkg/domain"
	"readon/pkg/models"
)

type ProductUseCase interface {
	ListProducts() ([]models.ListingBook, error)
	ListProductsForUser(*models.Pagination) ([]models.ListingBook, error)
	Addproduct(pdct models.Product) (error, error)
	EditProduct(pdct models.ProductUpdate) (models.ProductUpdate, error)
	AddBookCover(image []byte, book_id int) error
	GetProduct(bookId int) (domain.Book, error)
	DeleteProduct(int) error
	ListBookCovers(bookId int) ([][]byte, error)
}
