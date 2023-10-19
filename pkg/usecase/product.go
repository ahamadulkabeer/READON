package usecase

import (
	"context"
	"readon/pkg/domain"
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

func (c ProductUseCase) Addproduct(pdct models.Product) (addbookerr, addimgerr error) {
	var product = domain.Book{
		Title:      pdct.Name,
		Author:     pdct.Author,
		About:      pdct.About,
		CategoryID: pdct.CategoryID,
	}

	addbookerr = c.productRepo.AddProduct(product)
	if addbookerr != nil {
		return
	}
	addimgerr = c.productRepo.AddImage(pdct.Image)
	return
}

func (c ProductUseCase) GetProduct(bookId int) ([]domain.Book, error) {
	return c.productRepo.GetProducts(bookId)
}
