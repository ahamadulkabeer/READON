package interfaces

import (
	"context"
	domain "readon/pkg/domain"
	"readon/pkg/models"
)

type ProductUseCase interface {
	ListProducts(ctx context.Context) ([]models.ListingBook, error)
	Addproduct(pdct models.Product) (error, error)
	GetProduct(bookId int) ([]domain.Book, error)
}
