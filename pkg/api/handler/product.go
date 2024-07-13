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

// ListProducts godoc
// @Summary List products
// @Description Get a list of products
// @Tags product
// @Produce json
// @Success 200 {array} models.Product
// @Failure 500 {object} models.ErrorResponse
// @Router /user/books [get]
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

// ListProductsForUser godoc
// @Summary List, search, and explore products for a user
// @Description Get a list of products with pagination details for user side
// @Produce json
// @Tags product
// @Param page query int false "Page number for pagination (default: 1)"
// @Param filter query int false "filters the books by category , category id "
// @Param search query string false "Search keyword for products"
// @Success 200 {object} models.BooksListResponse
// @Failure 400 {object} models.ErrorResponse
// @Router /user/listbooks [get]
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

// AddProduct godoc
// @Summary Add a new product with an image
// @Description Add a new product with its details and an associated image
// @Tags product
// @Accept mpfd
// @Produce json
// @Param name formData string true "Product name"
// @Param author formData string true "Product author"
// @Param about formData string true "Product description"
// @Param category formData int true "Product category ID"
// @Param price formData float64 true "Price"
// @Param image formData file true "Product image"
// @Success 200 {string} string "Product added with image"
// @Failure 400 {object} models.ErrorResponse "Invalid request or form data"
// @Failure 500 {object} models.ErrorResponse "Error while adding product or image"
// @Router /admin/addproduct [post]
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

// @Summary Edit product details
// @Description Edit product details by providing JSON payload
// @Tags product
// @Security ApiKeyAuth
// @ID EditProductDet
// @Produce json
// @Param id formData int true "Product id"
// @Param name formData string false "Product name"
// @Param author formData string false "Product author"
// @Param about formData string  false "Product description"
// @Param price formData float64 true "Price"
// @Param category formData int true  "Product category ID"
// @Success 200 {object} models.ProductUpdate "OK"
// @Failure 500 {object} models.ErrorResponse "Internal Server Error"
// @Router /admin/editproduct [put]
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

// @Summary Add a book cover image
// @Description Add a book cover image for a specific book.
// @Tags product
// @Accept multipart/form-data
// @Produce json
// @Param id path int true "Book ID to associate with the cover image"
// @Param image formData file true "Book cover image file"
// @Success 200 {string} string
// @Failure 400 {object} models.ErrorResponse "Error while converting category id" or "Error while getting the image file"
// @Failure 500 {object} models.ErrorResponse "Error reading the file" or "Product added but image not added"
// @Router /admin/addcover/{id} [post]
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

// GetProduct godoc
// @Summary Get details of a specific product
// @Description Get details of a product using its ID
// @Tags product
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {object} models.ListingBook "Product details"
// @Failure 400 {object} models.ErrorResponse "Invalid product ID"
// @Failure 404 {object} models.ErrorResponse "Product not found"
// @Router /user/book/{id} [get]
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

// GetProduct godoc
// @Summary Get details of a specific product
// @Description Get details of a product using its ID
// @Tags product
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 string string
// @Failure 400 {object} models.ErrorResponse "Invalid product ID"
// @Failure 404 {object} models.ErrorResponse "Product not found"
// @Router /admin/deletebook/{id} [Delete]
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

// @Summary List book cover images
// @Description List book cover images for a specific book by its ID.
// @Tags product
// @Produce json
// @Param id path int true "Book ID to retrieve cover images for"
// @Success 200 {object} []byte  "Covers retrieved"
// @Failure 400 {object} models.ErrorResponse "Error while converting book ID"
// @Failure 500 {object} models.ErrorResponse "Error while getting cover images"
// @Router /admin/listbookcovers/{id} [get]
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
