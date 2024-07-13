package interfaces

import (
	"readon/pkg/api/responses"
)

type CategoryUsecase interface {
	AddCategory(string) responses.Response
	UpdateCategory(IDToUpdate uint, newCategory string) responses.Response
	DeleteCategory(categoryID int) responses.Response
	ListCategories() responses.Response
}
