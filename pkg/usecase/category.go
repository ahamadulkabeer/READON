package usecase

import (
	"fmt"
	"net/http"
	"readon/pkg/api/errorhandler"
	"readon/pkg/api/helpers"
	"readon/pkg/api/responses"
	interfaces "readon/pkg/repository/interface"
	services "readon/pkg/usecase/interface"
)

type CategoryUseCase struct {
	CategoryRepo interfaces.CategoryRepository
}

func NewCategoryUseCase(repo interfaces.CategoryRepository) services.CategoryUsecase {
	return &CategoryUseCase{
		CategoryRepo: repo,
	}
}

func (c CategoryUseCase) AddCategory(newCategory string) responses.Response {

	// validate category data
	err := helpers.ValidateCategory(newCategory)
	if err != nil {
		return responses.ClientReponse(http.StatusBadRequest, "couldn't add new category", err.Error(), nil)
	}

	// check if the same ategory exist
	exist, err := c.CategoryRepo.CheckCategory(newCategory)
	if err != nil {
		statusCode, _ := errorhandler.HandleDatabaseError(err)
		return responses.ClientReponse(statusCode, "couldn't add new category", err.Error(), nil)
	}
	if exist {
		return responses.ClientReponse(http.StatusConflict, "couldn't add new category", "category '"+newCategory+"' already exist", nil)
	}

	//add new category
	category, err := c.CategoryRepo.AddCategory(newCategory)
	if err != nil {
		statusCode, _ := errorhandler.HandleDatabaseError(err)
		return responses.ClientReponse(statusCode, "couldn't add new category", err.Error(), nil)
	}

	// response
	return responses.ClientReponse(http.StatusOK, "category created successfully", nil, category)
}

func (c CategoryUseCase) UpdateCategory(IDToUpdate uint, newCategory string) responses.Response {

	// valaidates category data
	err := helpers.ValidateCategory(newCategory)
	if err != nil {
		return responses.ClientReponse(http.StatusBadRequest, "couldn't update category", err.Error(), nil)
	}
	fmt.Println("new :", newCategory)
	// checks if the category to update already exist
	exist, err := c.CategoryRepo.CheckCategory(newCategory)
	if err != nil {
		statusCode, _ := errorhandler.HandleDatabaseError(err)
		return responses.ClientReponse(statusCode, "couldn't update category", err.Error(), nil)
	}
	if exist {
		return responses.ClientReponse(http.StatusUnprocessableEntity, "couldn't update category", "category already exist", nil)
	}

	//updates category
	err = c.CategoryRepo.UpdateCategory(IDToUpdate, newCategory)
	if err != nil {
		statusCode, _ := errorhandler.HandleDatabaseError(err)
		return responses.ClientReponse(statusCode, "couldn't update category", err.Error(), nil)
	}

	//response
	category, err := c.CategoryRepo.GetCategoryById(int(IDToUpdate))
	if err != nil {
		_, _ = errorhandler.HandleDatabaseError(err)
	}
	return responses.ClientReponse(http.StatusOK, "category updated ", nil, category)
}

func (c CategoryUseCase) DeleteCategory(categoryID int) responses.Response {

	// fetch category to delete
	category, err := c.CategoryRepo.GetCategoryById(categoryID)
	if err != nil {
		statusCode, _ := errorhandler.HandleDatabaseError(err)
		return responses.ClientReponse(statusCode, "couldn't delete category", err.Error(), nil)
	}
	if category.ID == 0 {
		return responses.ClientReponse(http.StatusNotFound, "couldn't delete category", "category doesn't exist", nil)
	}

	// delete category
	err = c.CategoryRepo.DeleteCategory(categoryID)
	if err != nil {
		statusCode, _ := errorhandler.HandleDatabaseError(err)
		return responses.ClientReponse(statusCode, "couldn't delete category", err.Error(), nil)
	}
	//response
	return responses.ClientReponse(http.StatusOK, "category deleted ", nil, category)
}

func (c CategoryUseCase) ListCategories() responses.Response {
	// fetches all categories
	list, err := c.CategoryRepo.ListCategories(100)
	if err != nil {
		statusCode, _ := errorhandler.HandleDatabaseError(err)
		return responses.ClientReponse(statusCode, "couldn't fetch list of categories", err.Error(), nil)
	}
	//response
	return responses.ClientReponse(http.StatusOK, "categories fetched ", nil, list)
}
