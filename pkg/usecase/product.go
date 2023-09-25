package usecase

import (
	"context"
	"readon/pkg/models"
	interfaces "readon/pkg/repository/interface"
	services "readon/pkg/usecase/interface"
)

type ProductUseCase struct {
	productRepo interfaces.ProductRepository
}

func NewProductUseCase(repo interfaces.ProductRepository) services.ProductUseCase {
	return &ProductUseCase{
		productRepo: repo,
	}
}

func (c ProductUseCase) ListProducts(ctx context.Context) ([]models.ListingBook, error) {
	listofbooks, err := c.productRepo.ListProducts(ctx)

	return listofbooks, err
}
