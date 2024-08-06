package usecase

import (
	"net/http"
	"readon/pkg/api/errorhandler"
	"readon/pkg/api/responses"
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

// func (c ProductUseCase) ListProducts() responses.Response {
// 	listofbooks, err := c.productRepo.ListProducts()
// 	if err != nil {
// 		statusCode, _ := errorhandler.HandleDatabaseError(err)
// 		responses.ClientReponse(statusCode, "couldn't fetch list of books", err.Error(), nil)
// 	}
// 	return responses.ClientReponse(http.StatusOK, "list of books fetched", nil, listofbooks)
// }

func (c ProductUseCase) ListProductsForUser(pageDet *models.Pagination) responses.Response {

	// pagination details
	pageDet.Size = 6
	offset := pageDet.Size * (pageDet.Page - 1)
	var err error
	pageDet.NumberOfResults, err = c.productRepo.GetTotalNoOfproducts(*pageDet)
	if err != nil {
		statusCode, _ := errorhandler.HandleDatabaseError(err)
		return responses.ClientReponse(statusCode, "couldn't fetch list of books", err.Error(), nil)
	}

	// populating book details in a slice of models.listingbook object
	listOfBooks, err := c.productRepo.ListProductsForUser(*pageDet, offset)
	if err != nil {
		statusCode, _ := errorhandler.HandleDatabaseError(err)
		return responses.ClientReponse(statusCode, "couldn't fetch list of books", err.Error(), nil)
	}

	var ListOfBookList []models.ListBook
	copier.Copy(&ListOfBookList, &listOfBooks)

	// fetch book cover
	// for i := range ListOfBookList {
	// 	bookCover, err := c.productRepo.GetBookCover(ListOfBookList[i].ID)
	// 	if err != nil {
	// 		statusCode, _ := errorhandler.HandleDatabaseError(err)
	// 		return responses.ClientReponse(statusCode, "couldn't fetch list of books", err.Error(), nil)
	// 	}
	// 	ListOfBookList[i].Image = bookCover.Image
	// }

	//response
	return responses.ClientReponse(http.StatusOK, " list of books fetched", nil, models.PaginatedListBooks{
		Books:      ListOfBookList,
		Pagination: *pageDet,
	})
}

func (c ProductUseCase) Addproduct(pdct models.Product) responses.Response {
	var product = domain.Book{
		Title:      pdct.Title,
		Author:     pdct.Author,
		About:      pdct.About,
		CategoryID: pdct.CategoryID,
		Price:      pdct.Price,
	}

	// add product
	bookId, err := c.productRepo.AddProduct(&product)
	if err != nil {
		statusCode, _ := errorhandler.HandleDatabaseError(err)
		return responses.ClientReponse(statusCode, "couldn't add product ", err.Error(), nil)
	}

	// add image
	err = c.productRepo.AddImage(pdct.Image, bookId)
	if err != nil {
		statusCode, _ := errorhandler.HandleDatabaseError(err)
		return responses.ClientReponse(statusCode, "product added : couldn't add image ", err.Error(), nil)
	}

	// fetch the product
	product, err = c.productRepo.GetProduct(int(product.ID))
	if err != nil {
		_, _ = errorhandler.HandleDatabaseError(err)
	}

	// response
	var book models.ListBook
	copier.Copy(&book, &product)
	bookCover, err := c.productRepo.GetBookCover(int(product.ID))
	if err != nil {
		_, _ = errorhandler.HandleDatabaseError(err)
	}
	book.Image = bookCover.Image
	return responses.ClientReponse(http.StatusOK, "product added ", nil, book)
}

func (c ProductUseCase) EditProduct(pdct models.ProductUpdate) responses.Response {

	//initilise the book object
	var product = domain.Book{
		Title:      pdct.Name,
		Author:     pdct.Author,
		About:      pdct.About,
		CategoryID: pdct.CategoryID,
		Price:      pdct.Price,
	}
	product.ID = uint(pdct.ID)

	// edit product
	_, err := c.productRepo.EditProduct(product)
	if err != nil {
		statusCode, _ := errorhandler.HandleDatabaseError(err)
		return responses.ClientReponse(statusCode, "couldn't edit product ", err.Error(), nil)
	}

	// fetch the product
	product, err = c.productRepo.GetProduct(int(product.ID))
	if err != nil {
		_, _ = errorhandler.HandleDatabaseError(err)
	}

	// response
	var book models.ListBook
	copier.Copy(&book, &product)
	bookCover, err := c.productRepo.GetBookCover(int(product.ID))
	if err != nil {
		_, _ = errorhandler.HandleDatabaseError(err)
	}
	book.Image = bookCover.Image
	return responses.ClientReponse(http.StatusOK, "product edited", nil, book)
}

func (c ProductUseCase) GetProduct(bookId int) responses.Response {
	book, err := c.productRepo.GetProduct(bookId)
	if err != nil {
		statusCode, _ := errorhandler.HandleDatabaseError(err)
		return responses.ClientReponse(statusCode, "couldn't fetch product ", err.Error(), nil)
	}

	// response
	var bookListing models.ListBook
	copier.Copy(&bookListing, &book)
	bookCover, err := c.productRepo.GetBookCover(bookId)
	if err != nil {
		_, _ = errorhandler.HandleDatabaseError(err)
	}
	bookListing.Image = bookCover.Image
	return responses.ClientReponse(http.StatusOK, "book fetched", nil, bookListing)
}

func (c ProductUseCase) AddBookCover(image []byte, bookId int) responses.Response {

	// check count of existins book covers
	count, err := c.productRepo.GetNumberOfBookCovers(uint(bookId))
	if err != nil {
		statusCode, _ := errorhandler.HandleDatabaseError(err)
		return responses.ClientReponse(statusCode, "couldn't add book cover ", err.Error(), nil)
	}

	if count >= 3 {
		return responses.ClientReponse(http.StatusUnprocessableEntity, "couldn't add book cover ",
			"can't have more than 3 book covers", nil)
	}

	// add book cover
	err = c.productRepo.AddImage(image, bookId)
	if err != nil {
		statusCode, _ := errorhandler.HandleDatabaseError(err)
		return responses.ClientReponse(statusCode, "couldn't add book cover ", err.Error(), nil)
	}

	// get product
	book, err := c.productRepo.GetProduct(bookId)
	if err != nil {
		statusCode, _ := errorhandler.HandleDatabaseError(err)
		return responses.ClientReponse(statusCode, "couldn't add book cover ", err.Error(), nil)
	}

	// response
	var bookListing models.ListBook
	copier.Copy(&bookListing, &book)
	bookCover, err := c.productRepo.GetBookCover(bookId)
	if err != nil {
		_, _ = errorhandler.HandleDatabaseError(err)
	}
	bookListing.Image = bookCover.Image
	return responses.ClientReponse(http.StatusOK, "image added", nil, bookListing)
}

func (c ProductUseCase) DeleteProduct(bookID int) responses.Response {

	// get product
	book, err := c.productRepo.GetProduct(bookID)
	if err != nil {
		statusCode, _ := errorhandler.HandleDatabaseError(err)
		return responses.ClientReponse(statusCode, "couldn't delete product ", err.Error(), nil)
	}

	//response
	var bookListing models.ListBook
	copier.Copy(&bookListing, &book)
	bookCover, err := c.productRepo.GetBookCover(bookID)
	if err != nil {
		_, _ = errorhandler.HandleDatabaseError(err)
	}
	bookListing.Image = bookCover.Image

	// delete book
	err = c.productRepo.DeleteProduct(book)
	if err != nil {
		statusCode, _ := errorhandler.HandleDatabaseError(err)
		return responses.ClientReponse(statusCode, "couldn't delete product ", err.Error(), nil)
	}

	return responses.ClientReponse(http.StatusOK, "product deleted", nil, bookListing)
}

func (c ProductUseCase) ListBookCovers(bookId int) responses.Response {
	bookCover, err := c.productRepo.ListBookCovers(bookId)
	if err != nil {
		statusCode, _ := errorhandler.HandleDatabaseError(err)
		return responses.ClientReponse(statusCode, "couldn't fetch  book covers ", err.Error(), nil)
	}
	// response
	var bookCoverListing []models.ListBookCover
	copier.Copy(&bookCoverListing, &bookCover)
	return responses.ClientReponse(http.StatusOK, "book covers fetched ", nil, bookCoverListing)
}
