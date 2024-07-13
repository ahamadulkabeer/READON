package usecase

import (
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

	err := helpers.ValidateCategory(newCategory)
	if err != nil {
		return responses.ClientReponse(http.StatusBadRequest, "couldn't add new category", err, nil)
	}
	exist, err := c.CategoryRepo.CheckCategory(newCategory)
	if err != nil {
		statusCode, _ := errorhandler.HandleDatabaseError(err)
		return responses.ClientReponse(statusCode, "couldn't add new category", err, nil)
	}
	if exist {
		return responses.ClientReponse(http.StatusConflict, "couldn't add new category", "category '"+newCategory+"' already exist", nil)
	}
	category, err := c.CategoryRepo.AddCategory(newCategory)
	if err != nil {
		statusCode, _ := errorhandler.HandleDatabaseError(err)
		return responses.ClientReponse(statusCode, "couldn't add new category", err, nil)
	}
	return responses.ClientReponse(http.StatusOK, "category created successfully", nil, category)
}

func (c CategoryUseCase) UpdateCategory(IDToUpdate uint, newCategory string) responses.Response {

	err := helpers.ValidateCategory(newCategory)
	if err != nil {
		return responses.ClientReponse(http.StatusBadRequest, "couldn't update category", err, nil)
	}
	exist, err := c.CategoryRepo.CheckCategory(newCategory)
	if err != nil {
		statusCode, _ := errorhandler.HandleDatabaseError(err)
		return responses.ClientReponse(statusCode, "couldn't update category", err, nil)
	}
	if !exist {
		return responses.ClientReponse(http.StatusNotFound, "couldn't update category", "category doesn't exist", nil)
	}

	err = c.CategoryRepo.UpdateCategory(IDToUpdate, newCategory)
	if err != nil {
		statusCode, _ := errorhandler.HandleDatabaseError(err)
		return responses.ClientReponse(statusCode, "couldn't update category", err, nil)
	}
	return responses.ClientReponse(http.StatusOK, "category updated ", nil, nil)
}

func (c CategoryUseCase) DeleteCategory(categoryID int) responses.Response {
	exist, err := c.CategoryRepo.GetCategoryById(categoryID)
	if err != nil {
		statusCode, _ := errorhandler.HandleDatabaseError(err)
		return responses.ClientReponse(statusCode, "couldn't delete category", err, nil)
	}
	if !exist {
		return responses.ClientReponse(http.StatusNotFound, "couldn't delete category", "category doesn't exist", nil)
	}
	err = c.CategoryRepo.DeleteCategory(categoryID)
	if err != nil {
		statusCode, _ := errorhandler.HandleDatabaseError(err)
		return responses.ClientReponse(statusCode, "couldn't delete category", err, nil)
	}
	return responses.ClientReponse(http.StatusOK, "category deleted ", nil, nil)
}

func (c CategoryUseCase) ListCategories() responses.Response {
	list, err := c.CategoryRepo.ListCategories(100)
	if err != nil {
		statusCode, _ := errorhandler.HandleDatabaseError(err)
		return responses.ClientReponse(statusCode, "couldn't fetch list of categories", err, nil)
	}
	return responses.ClientReponse(http.StatusOK, "categories fetched ", nil, list)
}
