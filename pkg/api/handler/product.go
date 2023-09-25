package handler

import (
	"fmt"
	"net/http"
	services "readon/pkg/usecase/interface"

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
