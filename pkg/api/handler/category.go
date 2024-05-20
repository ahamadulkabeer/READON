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

// ListCategories godoc
// @Summary Lists all the categories
// @Description Get a list of categories
// @Tags categories
// @Accept json
// @Produce json
// @Success 200 {array} domain.Category
// @Failure 500 {object} models.ErrorResponse
// @Router /admin/categorylist [get]
func (cr CategoryHandler) ListCategories(c *gin.Context) {
	list, err := cr.CategoryUsecase.ListCategories()
	if err != nil {
		errResponse := models.ErrorResponse{
			Err:    err.Error(),
			Status: "InternalServerError",
			Hint:   "please try again",
		}
		c.JSON(http.StatusInternalServerError, errResponse)
		return
	}
	c.JSON(http.StatusOK, list)
}

// AddCategory godoc
// @Summary Add a new category
// @Description Add a new category with the provided name , no duplicate allowed
// @Tags categories
// @Accept json
// @Produce json
// @Param request body models.Newcategory true "New category information"
// @Success 200 {string} string "Successfully added new category"
// @Failure 400 {object} models.ErrorResponse "BadRequest"
// @Failure 500 {object} models.ErrorResponse "InternalServerError"
// @Router /admin/addcategory [post]
func (cr CategoryHandler) AddCategory(c *gin.Context) {
	var newCategory models.Newcategory
	err := c.Bind(&newCategory)
	if err != nil {
		errResponse := models.ErrorResponse{
			Err:    err.Error(),
			Status: "error while binding user input",
			Hint:   "please try again",
		}
		c.JSON(http.StatusBadRequest, errResponse)
		return
	}

	addedCategory, err := cr.CategoryUsecase.AddCategory(newCategory.Name)
	if err != nil {
		errResponse := models.ErrorResponse{
			Err:    err.Error(),
			Status: "couldn't add  category",
			Hint:   "please try again",
		}
		c.JSON(http.StatusInternalServerError, errResponse)
		return
	}
	response := "successfully added new category : " + addedCategory
	c.JSON(http.StatusOK, response)
}

// UpdateCategory godoc
// @Summary Update a category
// @Description Update an existing category of provided ID with new category
// @Tags categories
// @Accept json
// @Produce json
// @Param id path int true "Category ID to update"
// @Param request body models.Newcategory true "New category information"
// @Success 200 {string} string "Successfully updated category"
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /admin/updatecategory/{id} [put]
func (cr CategoryHandler) UpdateCategory(c *gin.Context) {
	paramsId := c.Param("categoryId")
	id, err := strconv.Atoi(paramsId)
	fmt.Println("idd:", id)
	if err != nil {
		errResponse := models.ErrorResponse{
			Err:    err.Error(),
			Status: "Cannot parse category_id",
			Hint:   "please try again",
		}
		c.JSON(http.StatusBadRequest, errResponse)
		return
	}
	var newcategory models.Newcategory
	c.Bind(&newcategory)
	fmt.Println("newcategory.Name:", newcategory.Name)
	addedCategory, err := cr.CategoryUsecase.UpdateCategory(id, newcategory.Name)
	if err != nil {
		errResponse := models.ErrorResponse{
			Err:    err.Error(),
			Status: "update filed : " + addedCategory.Name,
			Hint:   "please try again",
		}
		c.JSON(http.StatusInternalServerError, errResponse)
		return
	}
	response := "successfully updated new category : " + addedCategory.Name
	c.JSON(http.StatusOK, response)
}

// DeleteCategory godoc
// @Summary Delete a category
// @Description Delete an existing category with the provided ID
// @Tags categories
// @Accept json
// @Produce json
// @Param id path int true "Category ID to delete"
// @Success 200 {string} string "Successfully deleted category"
// @Failure 400 {object} models.ErrorResponse "Cannot parse category_id"
// @Failure 500 {object} models.ErrorResponse "Delete failed" or "Please try again"
// @Router /admin/deletecategory/{id} [delete]
func (cr CategoryHandler) DeleteCategory(c *gin.Context) {
	paramsId := c.Param("categoryId")
	id, err := strconv.Atoi(paramsId)
	if err != nil {
		errResponse := models.ErrorResponse{
			Err:    err.Error(),
			Status: "Cannot parse category_id",
			Hint:   "please try again",
		}
		c.JSON(http.StatusBadRequest, errResponse)
		return
	}
	err = cr.CategoryUsecase.DeleteCategory(id)
	if err != nil {
		errResponse := models.ErrorResponse{
			Err:    err.Error(),
			Status: "Delete filed , id : " + paramsId,
			Hint:   "please try again",
		}
		c.JSON(http.StatusInternalServerError, errResponse)
		return
	}
	response := "successfully deleted the category  , id :" + paramsId
	c.JSON(http.StatusOK, response)

}
