package usecase

import (
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

func (c ProductUseCase) ListProducts() ([]models.ListingBook, error) {
	listofbooks, err := c.productRepo.ListProducts()

	return listofbooks, err
}

func (c ProductUseCase) ListProductsForUser(pageDet *models.Pagination) ([]models.ListingBook, error) {
	pageDet.Size = 5

	offset := pageDet.Size * (pageDet.NewPage - 1)
	listofbooks, numOfResults, err := c.productRepo.ListProductsForUser(*pageDet, offset)
	pageDet.Lastpage = numOfResults / pageDet.Size
	if numOfResults%pageDet.Size != 0 {
		pageDet.Lastpage++
	}
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

func (c ProductUseCase) GetProduct(bookId int) (domain.Book, error) {
	return c.productRepo.GetProduct(bookId)
}

func (c ProductUseCase) DeleteProduct(bookID int) error {
	book, err := c.productRepo.GetProduct(bookID)
	if err != nil {
		return err
	}
	err = c.productRepo.DeleteProduct(book)
	return err
}
