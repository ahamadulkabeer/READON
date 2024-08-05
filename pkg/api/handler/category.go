package handler

import (
	"net/http"
	"readon/pkg/api/responses"
	"readon/pkg/models"
	services "readon/pkg/usecase/interface"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CategoryHandler struct {
	CategoryUsecase services.CategoryUsecase
}

func NewCategoryHandler(usecase services.CategoryUsecase) *CategoryHandler {
	return &CategoryHandler{
		CategoryUsecase: usecase,
	}
}

func (cr CategoryHandler) ListCategories(c *gin.Context) {
	response := cr.CategoryUsecase.ListCategories()

	c.JSON(http.StatusOK, response)
}

func (cr CategoryHandler) AddCategory(c *gin.Context) {
	var newCategory models.Newcategory
	err := c.Bind(&newCategory)
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.ClientReponse(http.StatusBadRequest,
			"error binding while category details", err, nil))
		return
	}

	response := cr.CategoryUsecase.AddCategory(newCategory.Name)

	c.JSON(http.StatusOK, response)
}

func (cr CategoryHandler) UpdateCategory(c *gin.Context) {
	paramsId := c.Param("categoryId")
	id, err := strconv.Atoi(paramsId)
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.ClientReponse(http.StatusBadRequest,
			"error binding params", err, nil))
		return
	}
	var newcategory models.Newcategory
	err = c.Bind(&newcategory)
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.ClientReponse(http.StatusBadRequest,
			"error binding while category details", err, nil))
		return
	}
	response := cr.CategoryUsecase.UpdateCategory(uint(id), newcategory.Name)

	c.JSON(http.StatusOK, response)
}

func (cr CategoryHandler) DeleteCategory(c *gin.Context) {
	paramsId := c.Param("categoryId")
	id, err := strconv.Atoi(paramsId)
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.ClientReponse(http.StatusBadRequest,
			"error binding params", err, nil))
		return
	}
	response := cr.CategoryUsecase.DeleteCategory(id)

	c.JSON(http.StatusOK, response)

}
