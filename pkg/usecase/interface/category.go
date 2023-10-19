package interfaces

import domain "readon/pkg/domain"

type CategoryUsecase interface {
	AddCategory(string) (domain.Category, error)
	UpdateCategory(int, string) (domain.Category, error)
	DeleteCategory(categoryID int) error
	ListCategories() ([]domain.Category, error)
}
