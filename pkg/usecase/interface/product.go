package interfaces

import (
	"readon/pkg/api/responses"
	"readon/pkg/models"
)

type ProductUseCase interface {
	//ListProducts() responses.Response
	ListProductsForUser(*models.Pagination) responses.Response
	Addproduct(pdct models.Product) responses.Response
	EditProduct(pdct models.ProductUpdate) responses.Response
	AddBookCover(image []byte, book_id int) responses.Response
	GetProduct(bookId int) responses.Response
	DeleteProduct(int) responses.Response
	ListBookCovers(bookId int) responses.Response
}
