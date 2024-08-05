package handler

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"readon/pkg/api/helpers"
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

func (cr *ProductHandler) ListProducts(c *gin.Context) {

	list, err := cr.productUseCase.ListProducts()
	if err != nil {
		errResponse := models.ErrorResponse{
			Err:    err.Error(),
			Status: "err while getting product list",
			Hint:   "please try again",
		}

		c.JSON(http.StatusInternalServerError, errResponse)
		return
	}
	for i := range list {
		log.Println(i, " : ", list[i])
	}
	c.JSON(http.StatusOK, list)
}

func (cr ProductHandler) ListProductsForUSer(c *gin.Context) { // listing , search , explore user side
	var pagedetails models.Pagination
	pagedetails.NewPage = 1

	err := c.BindQuery(&pagedetails)
	if err != nil {
		errResponse := models.ErrorResponse{
			Err:    err.Error(),
			Status: "Error while parsing data",
			Hint:   "please try again",
		}
		c.JSON(http.StatusBadRequest, errResponse)
		return
	}

	log.Println("query got :", pagedetails)

	list, err := cr.productUseCase.ListProductsForUser(&pagedetails)
	if err != nil {
		errResponse := models.ErrorResponse{
			Err:    err.Error(),
			Status: "err while getting product list",
			Hint:   "please try again",
		}

		c.JSON(http.StatusInternalServerError, errResponse)
		return
	}
	for i := range list {
		log.Println(i, " : ", list[i])
	}
	Response := models.BooksListResponse{
		Pagination: pagedetails,
		List:       list,
	}
	c.JSON(http.StatusOK, Response)

}

func (cr *ProductHandler) Addproduct(c *gin.Context) {
	var product models.Product
	/*rawBody, err := c.GetRawData()

	if err != nil {
		// Handle error
	}
	fmt.Println("Raw JSON Body:", string(rawBody))*/
	log.Println("request is runnig")
	if err := c.Request.ParseMultipartForm(10 << 20); err != nil {
		errResponse := models.ErrorResponse{
			Err:    err.Error(),
			Status: "Unable to parse form",
			Hint:   "please try again",
		}
		c.JSON(http.StatusBadRequest, errResponse)
		return
	}

	file, _, err := c.Request.FormFile("image")

	if err != nil {
		errResponse := models.ErrorResponse{
			Err:    err.Error(),
			Status: "Error while getting the image file",
			Hint:   "please try again",
		}
		c.JSON(http.StatusBadRequest, errResponse)
		return
	}

	defer file.Close()

	imagefile, err := io.ReadAll(file)
	if err != nil {
		errResponse := models.ErrorResponse{
			Err:    err.Error(),
			Status: "Error reading the file",
			Hint:   "please try again",
		}
		c.JSON(http.StatusInternalServerError, errResponse)
		return
	}

	imagefile, err = helpers.CropImage(imagefile)
	if err != nil {
		errResponse := models.ErrorResponse{
			Err:    err.Error(),
			Status: "Error while processing iimage file",
			Hint:   "please try again",
		}
		c.JSON(http.StatusInternalServerError, errResponse)
		return
	}

	product.Image = imagefile
	product.Title = c.PostForm("title")
	product.Author = c.PostForm("author")
	product.About = c.PostForm("about")
	product.Price, err = strconv.ParseFloat(c.PostForm("price"), 64)
	product.CategoryID, err = strconv.Atoi(c.PostForm("category"))
	if err != nil {
		errResponse := models.ErrorResponse{
			Err:    err.Error(),
			Status: "Error while converting category id",
			Hint:   "please try again",
		}
		c.JSON(http.StatusBadRequest, errResponse)
		return
	}

	// dont have to explicitly encode into base64 as gin will automatically do it.
	//imageDataBase64 := base64.StdEncoding.EncodeToString(product.Image)
	log.Println("neccessary data collected without error")
	producterr, imageerr := cr.productUseCase.Addproduct(product)
	if producterr != nil {
		errResponse := models.ErrorResponse{
			Err:    err.Error(),
			Status: "product not added",
			Hint:   "please try again",
		}
		c.JSON(http.StatusInternalServerError, errResponse)
		return
	}
	log.Println("book added succesfully")
	if imageerr != nil {
		errResponse := models.ErrorResponse{
			Err:    err.Error(),
			Status: "product added  BUt image not added pleasse upload image seperatly",
			Hint:   "please try again",
		}
		c.JSON(http.StatusInternalServerError, errResponse)
		return
	}
	response := "product added with image"
	c.JSON(http.StatusOK, response)
}

func (cr ProductHandler) EditProductDet(c *gin.Context) {
	var product models.ProductUpdate
	var err error
	product.ID, err = strconv.Atoi(c.Param("bookId"))
	if err != nil {
		errResponse := models.ErrorResponse{
			Err:    err.Error(),
			Status: "couldn't parse bookId",
			Hint:   "please try again",
		}
		c.JSON(http.StatusInternalServerError, errResponse)
		return
	}
	c.Bind(&product)
	fmt.Println("product :", product)
	product, err = cr.productUseCase.EditProduct(product)
	if err != nil {
		errResponse := models.ErrorResponse{
			Err:    err.Error(),
			Status: "couldn't update product",
			Hint:   "please try again",
		}
		c.JSON(http.StatusInternalServerError, errResponse)
		return
	}

	c.JSON(http.StatusOK, product)

}

func (cr ProductHandler) AddBookCover(c *gin.Context) {

	bookId, err := strconv.Atoi(c.Param("bookId"))
	if err != nil {
		errResponse := models.ErrorResponse{
			Err:    err.Error(),
			Status: "Error while converting category id",
			Hint:   "please try again",
		}
		c.JSON(http.StatusBadRequest, errResponse)
		return
	}

	file, _, err := c.Request.FormFile("image")

	if err != nil {
		errResponse := models.ErrorResponse{
			Err:    err.Error(),
			Status: "Error while getting the image file",
			Hint:   "please try again",
		}
		c.JSON(http.StatusBadRequest, errResponse)
		return
	}

	defer file.Close()
	imagefile, err := io.ReadAll(file)
	if err != nil {
		errResponse := models.ErrorResponse{
			Err:    err.Error(),
			Status: "Error reading the file",
			Hint:   "please try again",
		}
		c.JSON(http.StatusInternalServerError, errResponse)
		return
	}
	imagefile, err = helpers.CropImage(imagefile)
	if err != nil {
		errResponse := models.ErrorResponse{
			Err:    err.Error(),
			Status: "Error while processing iimage file",
			Hint:   "please try again",
		}
		c.JSON(http.StatusInternalServerError, errResponse)
		return
	}
	// dont have to explicitly encode into base64 as gin will automatically do it when passin g as json.
	//imageDataBase64 := base64.StdEncoding.EncodeToString(imagefile)
	//fmt.Println("image : ", imageDataBase64)
	imageerr := cr.productUseCase.AddBookCover(imagefile, bookId)
	if imageerr != nil {
		errResponse := models.ErrorResponse{
			Err:    err.Error(),
			Status: "product added  BUt image not added pleasse upload image seperatly",
			Hint:   "please try again",
		}
		c.JSON(http.StatusInternalServerError, errResponse)
		return
	}
	response := "BookCover added successfully"
	c.JSON(http.StatusOK, response)
}

func (cr ProductHandler) GetProduct(c *gin.Context) {
	bookId := c.Param("bookId")

	id, err := strconv.Atoi(bookId)
	if err != nil {
		errResponse := models.ErrorResponse{
			Err:    err.Error(),
			Status: "Error while converting category id",
			Hint:   "please try again",
		}
		c.JSON(http.StatusBadRequest, errResponse)
		return
	}

	book, err := cr.productUseCase.GetProduct(id)
	if err != nil {
		errResponse := models.ErrorResponse{
			Err:    err.Error(),
			Status: "Book not found",
			Hint:   "please try again",
		}
		c.JSON(http.StatusNotFound, errResponse)
		return
	}
	//imageDataBase64 := base64.StdEncoding.EncodeToString(book.Image)
	//log.Println("img : ", imageDataBase64)
	//var response models.ListingBook
	//copier.Copy(&response, book)

	c.JSON(http.StatusOK, book)
}

func (cr ProductHandler) DeleteProduct(c *gin.Context) {
	bookId := c.Param("bookId")

	id, err := strconv.Atoi(bookId)
	if err != nil {
		errResponse := models.ErrorResponse{
			Err:    err.Error(),
			Status: "Error while converting  bookId ",
			Hint:   "please try again",
		}
		c.JSON(http.StatusBadRequest, errResponse)
		return
	}

	err = cr.productUseCase.DeleteProduct(id)
	if err != nil {
		errResponse := models.ErrorResponse{
			Err:    err.Error(),
			Status: "Book not found",
			Hint:   "please try again",
		}
		c.JSON(http.StatusNotFound, errResponse)
		return
	}
	response := "book deleted succresfully || on id : " + bookId
	c.JSON(http.StatusOK, response)
}

func (cr ProductHandler) ListBookCovers(c *gin.Context) {
	paramId := c.Param("bookId")
	bookId, err := strconv.Atoi(paramId)
	if err != nil {
		errResponse := models.ErrorResponse{
			Err:    err.Error(),
			Status: "Error while converting bookId",
			Hint:   "please try again",
		}
		c.JSON(http.StatusBadRequest, errResponse)
		return
	}
	list, err := cr.productUseCase.ListBookCovers(bookId)
	if err != nil {
		errResponse := models.ErrorResponse{
			Err:    err.Error(),
			Status: "err while getting product list",
			Hint:   "please try again",
		}

		c.JSON(http.StatusInternalServerError, errResponse)
		return
	}
	fmt.Println("image : ", list[0])
	c.JSON(http.StatusOK, list)

}
