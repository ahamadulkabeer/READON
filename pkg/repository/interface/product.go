package interfaces

import (
	"context"
	"readon/pkg/models"
)

type ProductRepository interface {
	ListProducts(ctx context.Context) ([]models.ListingBook, error)
}
