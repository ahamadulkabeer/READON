package usecase

import (
	"fmt"
	"readon/pkg/domain"
	"readon/pkg/models"
	interfaces "readon/pkg/repository/interface"
	services "readon/pkg/usecase/interface"

	"github.com/jinzhu/copier"
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
	numOfResults, err := c.productRepo.GetTotalNoOfproducts(*pageDet)
	if err != nil {
		return nil, err
	}
	// populating book details in a slice of models.listingbook object
	listofbooks, err := c.productRepo.ListProductsForUser(*pageDet, offset)
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
		Price:      pdct.Price,
	}

	bookId, addbookerr := c.productRepo.AddProduct(product)
	if addbookerr != nil {
		return
	}
	addimgerr = c.productRepo.AddImage(pdct.Image, bookId)
	return
}

func (c ProductUseCase) EditProduct(pdct models.ProductUpdate) (models.ProductUpdate, error) {

	/*oproduct, err := c.productRepo.GetProduct(pdct.Id)
	if err != nil {
		return pdct, err
	}
	fmt.Println("oprodect :", oproduct)*/

	var product = domain.Book{
		ID:         uint(pdct.Id),
		Title:      pdct.Name,
		Author:     pdct.Author,
		About:      pdct.About,
		CategoryID: pdct.CategoryID,
		Price:      pdct.Price,
	}
	fmt.Println("product category id :", product.CategoryID)
	product, err := c.productRepo.EditProduct(product)
	copier.Copy(&pdct, &product)
	return pdct, err

}

func (c ProductUseCase) GetProduct(bookId int) (models.ListingBook, error) {
	return c.productRepo.GetProduct(bookId)
}

func (c ProductUseCase) AddBookCover(image []byte, bookId int) error {
	_, err := c.productRepo.GetProduct(bookId)
	if err != nil {
		return err
	}
	return c.productRepo.AddImage(image, bookId)
}

func (c ProductUseCase) DeleteProduct(bookID int) error {
	listingbook, err := c.productRepo.GetProduct(bookID)
	if err != nil {
		return err
	}
	var book domain.Book
	copier.Copy(&book, listingbook)
	err = c.productRepo.DeleteProduct(book)
	return err
}

func (c ProductUseCase) ListBookCovers(bookId int) ([][]byte, error) {
	return c.productRepo.ListBookCovers(bookId)
}
