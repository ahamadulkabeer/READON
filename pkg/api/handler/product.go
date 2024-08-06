package handler

import (
	"io"
	"net/http"
	"readon/pkg/api/helpers"
	"readon/pkg/api/responses"
	"readon/pkg/models"
	services "readon/pkg/usecase/interface"
	"strconv"

	gin "github.com/gin-gonic/gin"
)

type ProductHandler struct {
	productUseCase services.ProductUseCase
}

func NewProductHandler(usecase services.ProductUseCase) *ProductHandler {
	return &ProductHandler{
		productUseCase: usecase,
	}
}

// func (cr *ProductHandler) ListProducts(c *gin.Context) {
// 	response := cr.productUseCase.ListProducts()
// 	c.JSON(response.StatusCode, response)
// }

func (cr ProductHandler) ListProductsForUSer(c *gin.Context) { // listing , search , explore user side
	var pagedetails models.Pagination
	pagedetails.Page = 1

	err := c.BindQuery(&pagedetails)
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.ClientReponse(http.StatusBadRequest,
			"Error while parsing data", err.Error(), nil))
		return
	}

	response := cr.productUseCase.ListProductsForUser(&pagedetails)

	c.JSON(http.StatusOK, response)

}

func (cr *ProductHandler) Addproduct(c *gin.Context) {

	var product models.Product

	if err := c.Request.ParseMultipartForm(10 << 20); err != nil {
		c.JSON(http.StatusBadRequest, responses.ClientReponse(http.StatusBadRequest,
			"Unable to parse formData", err.Error(), nil))
		return
	}

	file, _, err := c.Request.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.ClientReponse(http.StatusBadRequest,
			"Error while parsing the image file", err.Error(), nil))
		return
	}

	defer file.Close()

	imagefile, err := io.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.ClientReponse(http.StatusBadRequest,
			"Error reading the image file", err.Error(), nil))
		return
	}

	imagefile, err = helpers.CropImage(imagefile)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.ClientReponse(http.StatusInternalServerError,
			"Error processing image file", err.Error(), nil))
		return
	}

	product.Image = imagefile
	product.Title = c.PostForm("title")
	product.Author = c.PostForm("author")
	product.About = c.PostForm("about")
	product.Price, err = strconv.ParseFloat(c.PostForm("price"), 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.ClientReponse(http.StatusBadRequest,
			"Error while parsing price ", err.Error(), nil))
		return
	}
	product.CategoryID, err = strconv.Atoi(c.PostForm("category"))
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.ClientReponse(http.StatusBadRequest,
			"Error while parsing category id", err.Error(), nil))
		return
	}

	// dont have to explicitly encode into base64 as gin will automatically do it.
	//imageDataBase64 := base64.StdEncoding.EncodeToString(product.Image)
	response := cr.productUseCase.Addproduct(product)

	c.JSON(response.StatusCode, response)
}

func (cr ProductHandler) EditProductDet(c *gin.Context) {
	var product models.ProductUpdate
	var err error
	product.ID, err = strconv.Atoi(c.Param("bookId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.ClientReponse(http.StatusBadRequest,
			"couldn't parse bookId", err.Error(), nil))
		return
	}
	err = c.Bind(&product)
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.ClientReponse(http.StatusBadRequest,
			"couldn't bind form data", err.Error(), nil))
		return
	}

	response := cr.productUseCase.EditProduct(product)
	c.JSON(response.StatusCode, response)

}

func (cr ProductHandler) AddBookCover(c *gin.Context) {

	bookId, err := strconv.Atoi(c.Param("bookId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.ClientReponse(http.StatusBadRequest,
			"Error while converting book id", err.Error(), nil))
		return
	}

	file, _, err := c.Request.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.ClientReponse(http.StatusBadRequest,
			"Error while parsing the image file", err.Error(), nil))
		return
	}

	defer file.Close()

	imagefile, err := io.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.ClientReponse(http.StatusBadRequest,
			"Error reading the image file", err.Error(), nil))
		return
	}

	imagefile, err = helpers.CropImage(imagefile)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.ClientReponse(http.StatusInternalServerError,
			"Error processing image file", err.Error(), nil))
		return
	}

	response := cr.productUseCase.AddBookCover(imagefile, bookId)

	c.JSON(response.StatusCode, response)
}

func (cr ProductHandler) GetProduct(c *gin.Context) {

	bookId, err := strconv.Atoi(c.Param("bookId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.ClientReponse(http.StatusBadRequest,
			"Error while converting book id", err.Error(), nil))
		return
	}

	response := cr.productUseCase.GetProduct(bookId)
	c.JSON(response.StatusCode, response)
}

func (cr ProductHandler) DeleteProduct(c *gin.Context) {
	bookId, err := strconv.Atoi(c.Param("bookId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.ClientReponse(http.StatusBadRequest,
			"Error while converting book id", err.Error(), nil))
		return
	}
	response := cr.productUseCase.DeleteProduct(bookId)
	c.JSON(response.StatusCode, response)
}

func (cr ProductHandler) ListBookCovers(c *gin.Context) {
	bookId, err := strconv.Atoi(c.Param("bookId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.ClientReponse(http.StatusBadRequest,
			"Error while converting book id", err.Error(), nil))
		return
	}
	response := cr.productUseCase.ListBookCovers(bookId)
	c.JSON(response.StatusCode, response)

}
