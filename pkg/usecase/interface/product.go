package interfaces

import (
	"context"
	"readon/pkg/models"
)

type ProductUseCase interface {
	ListProducts(ctx context.Context) ([]models.ListingBook, error)
}
