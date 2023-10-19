package handler

import (
	"io"
	"log"
	"net/http"
	"readon/pkg/models"
	services "readon/pkg/usecase/interface"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
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
	c.JSON(http.StatusOK, list)
}

// ListProductsForUser godoc
// @Summary List, search, and explore products for a user
// @Description Get a list of products with pagination details for user side
// @Produce json
// @Tags product
// @Param page query int false "Page number for pagination (default: 1)"
// @Param filter query int false "Sort order for products (e.g., name ASC, price DESC)"
// @Param search query string false "Search keyword for products"
// @Success 200 {object} models.BooksListResponse
// @Failure 400 {object} models.ErrorResponse
// @Router /user/listbooks [post]
func (cr ProductHandler) ListProductsForUSer(c *gin.Context) { // listing , search , explore user side
	var pagedetails models.Pagination
	pagedetails.NewPage = 1

	err := c.BindQuery(&pagedetails)
	if err != nil {
		errResponse := models.ErrorResponse{
			Err:    err.Error(),
			Status: "Error while converting category id",
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

	product.Image = imagefile
	product.Name = c.PostForm("name")
	product.Author = c.PostForm("author")
	product.About = c.PostForm("about")
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
	bookId := c.Param("id")

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
	var response models.ListingBook
	copier.Copy(&response, book)

	c.JSON(http.StatusOK, response)
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
	bookId := c.Param("id")

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
