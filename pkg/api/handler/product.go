package handler

import (
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"readon/pkg/models"
	services "readon/pkg/usecase/interface"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	productUseCase services.ProductUseCase
}

func NewProductHandler(usecase services.ProductUseCase) *ProductHandler {
	return &ProductHandler{
		productUseCase: usecase,
	}
}

// when using search ()
func (cr *ProductHandler) ListProducts(c *gin.Context) {
	fmt.Println("listing products")

	list, err := cr.productUseCase.ListProducts(c.Request.Context())
	if err != nil {
		fmt.Println("err while getting product list :", err)
		c.JSON(http.StatusNotFound, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "request working",
		"list":   list,
	})
}

func (cr *ProductHandler) Addproduct(c *gin.Context) {
	var product models.Product
	/*rawBody, err := c.GetRawData()

	if err != nil {
		// Handle error
	}
	fmt.Println("Raw JSON Body:", string(rawBody))*/

	if err := c.Request.ParseMultipartForm(10 << 20); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unable to parse form"})
		return
	}
	file, _, err := c.Request.FormFile("image")

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error retrieving the file"})
		return
	}
	defer file.Close()

	imagefile, err := io.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error reading the file"})
		return
	}
	product.Image = imagefile

	product.Name = c.PostForm("name")
	product.Author = c.PostForm("author")
	product.About = c.PostForm("about")
	product.CategoryID, err = strconv.Atoi(c.PostForm("category"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error while converting category id"})
		return
	}

	// dont have to explicitly encode into base64 as gin will automatically do it.
	imageDataBase64 := base64.StdEncoding.EncodeToString(product.Image)

	producterr, imageerr := cr.productUseCase.Addproduct(product)
	if producterr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"err":    err,
			"status": "product not added",
		})
		return
	}
	if imageerr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"err":    err.Error(),
			"status": "product added  BUt image not added pleasse upload image seperatly",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"err":    err,
		"status": "product added with image",
		"img":    imageDataBase64,
	})
}

func (cr ProductHandler) GetProduct(c *gin.Context) {
	bookId := c.Param("id")
	id, err := strconv.Atoi(bookId)
	fmt.Println("id:", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Cannot parse category_id",
			"err":   err,
		})
		return
	}
	books, err := cr.productUseCase.GetProduct(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"err":    err,
			"status": "Book not found",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"err":    err,
		"status": "product added with image",
		"Book":   books,
	})
}
