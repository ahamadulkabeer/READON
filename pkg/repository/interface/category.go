package interfaces

import "readon/pkg/domain"

type CategoryRepository interface {
	AddCategory(string) (domain.Category, error)
	CheckCategory(current string) (bool, error)
	GetCategoryById(categoryID int) (domain.Category, error)
	UpdateCategory(IDToChange uint, newCategory string) error
	DeleteCategory(categoryID int) error
	ListCategories(limit int) ([]domain.Category, error)
}
