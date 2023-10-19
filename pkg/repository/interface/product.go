package interfaces

import (
	"context"
	"readon/pkg/domain"
	"readon/pkg/models"
)

type ProductRepository interface {
	ListProducts(ctx context.Context) ([]models.ListingBook, error)
	AddProduct(product domain.Book) error
	AddImage(image []byte) error
	GetProducts(int) ([]domain.Book, error)
}
