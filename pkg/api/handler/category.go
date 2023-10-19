package handler

import (
	"fmt"
	"net/http"
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
	list, err := cr.CategoryUsecase.ListCategories()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err,
			"list":  list,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "ok !",
		"list":   list,
	})
}

func (cr CategoryHandler) AddCategory(c *gin.Context) {
	var newCategory models.Newcategory
	err := c.Bind(&newCategory)
	fmt.Println("newcategory  :", newCategory)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "error while binding user input",
		})
		return
	}

	addedCategory, err := cr.CategoryUsecase.AddCategory(newCategory.Name)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"newcategory": newCategory,
			"error":       err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":      "successfully added new category",
		"newcategory": addedCategory,
	})
}

func (cr CategoryHandler) UpdateCategory(c *gin.Context) {
	paramsId := c.Param("id")
	id, err := strconv.Atoi(paramsId)
	fmt.Println("idd:", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Cannot parse category_id",
		})
		return
	}
	var newcategory models.Newcategory
	c.Bind(&newcategory)
	fmt.Println("newcategory.Name:", newcategory.Name)
	addedCategory, err := cr.CategoryUsecase.UpdateCategory(id, newcategory.Name)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":       err.Error(),
			"status":      "update filed",
			"newcategory": addedCategory,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":      "successfully updated new category",
		"newcategory": addedCategory,
		"error":       err,
	})

}

func (cr CategoryHandler) DeleteCategory(c *gin.Context) {
	paramsId := c.Param("id")
	id, err := strconv.Atoi(paramsId)
	fmt.Println("idd:", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Cannot parse category_id",
			"err":   err,
		})
		return
	}
	err = cr.CategoryUsecase.DeleteCategory(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":       err.Error(),
			"status":      "delete filed",
			"category_id": id,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":      "successfully deleted the category",
		"category_id": id,
		"error":       err,
	})

}
