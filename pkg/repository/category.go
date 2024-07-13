package repository

import (
	domain "readon/pkg/domain"
	interfaces "readon/pkg/repository/interface"

	"gorm.io/gorm"
)

type categoryDatabase struct {
	DB *gorm.DB
}

func NewCategoryRepository(DB *gorm.DB) interfaces.CategoryRepository {
	return &categoryDatabase{
		DB: DB,
	}
}

func (c categoryDatabase) AddCategory(newcategory string) (domain.Category, error) {
	newCategory := &domain.Category{Name: newcategory}
	err := c.DB.Create(newCategory).Error
	if err != nil {
		return *newCategory, err
	}
	return *newCategory, nil

}

func (c categoryDatabase) CheckCategory(category string) (bool, error) {
	existingCategory := &domain.Category{Name: category}
	result := c.DB.Where("name = ?", category).First(&existingCategory)
	if result.Error != nil {
		return false, result.Error
	}
	if result.RowsAffected != 0 {
		return false, nil
	}
	return true, nil
}

func (c categoryDatabase) GetCategoryById(categoryID int) (bool, error) {
	var Category domain.Category
	result := c.DB.Where("id = ?", categoryID).First(&Category)
	if result.Error != nil {
		return false, result.Error
	}
	if result.RowsAffected != 0 {
		return false, nil
	}
	return true, nil
}

func (c categoryDatabase) UpdateCategory(IDToChange uint, newCategory string) error {

	err := c.DB.Model(domain.Category{}).Where("id = ?", IDToChange).Set("name", newCategory).Error
	if err != nil {
		return err
	}
	return nil
}

func (c categoryDatabase) DeleteCategory(categoryID int) error {
	var categoryToDelete domain.Category
	result := c.DB.Where("id = ?", categoryID).First(&categoryToDelete)
	if result.Error != nil {
		return result.Error
	}

	if err := c.DB.Delete(&categoryToDelete).Error; err != nil {
		return err
	}

	return nil
}

func (c categoryDatabase) ListCategories(limit int) ([]domain.Category, error) {
	var categories []domain.Category
	err := c.DB.Limit(limit).Find(&categories).Error
	if err != nil {
		return categories, err
	}
	return categories, nil
}
