package repository

import (
	"fmt"
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

	return *newCategory, err

}

func (c categoryDatabase) CheckCategory(category string) error {
	existingCategory := &domain.Category{Name: category}
	err := c.DB.Where("name = ?", category).First(&existingCategory).Error
	return err
}

func (c categoryDatabase) GetCategoryById(idtoch int) error {
	var Category domain.Category
	err := c.DB.Where("id = ?", idtoch).First(&Category).Error
	return err
}

func (c categoryDatabase) UpdateCategory(idtoch int, newctg string) (domain.Category, error) {
	existingCategory := domain.Category{}
	err := c.DB.First(&existingCategory, idtoch).Error
	if err != nil {
		return existingCategory, err
	}

	existingCategory.Name = newctg
	err = c.DB.Save(&existingCategory).Error

	return existingCategory, err
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
	fmt.Println("cat :", categories)
	return categories, err
}
