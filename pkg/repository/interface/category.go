package interfaces

import "readon/pkg/domain"

type CategoryRepository interface {
	AddCategory(string) (domain.Category, error)
	CheckCategory(currrent string) error
	GetCategoryById(idtoch int) error
	UpdateCategory(int, string) (domain.Category, error)
	DeleteCategory(categoryID int) error
	ListCategories(limit int) ([]domain.Category, error)
}
