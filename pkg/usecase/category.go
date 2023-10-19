package usecase

import (
	"errors"
	"fmt"
	"readon/pkg/api/middleware"
	"readon/pkg/domain"
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

func (c CategoryUseCase) AddCategory(newcategory string) (domain.Category, error) {
	var newCategory domain.Category
	err := c.CategoryRepo.CheckCategory(newcategory)
	if err == nil {
		fmt.Println("error == nil catogory already exist getcategory by id")
		err := errors.New("category already exist!")
		fmt.Println(err)
		return newCategory, err
	}
	newCategory, err = c.CategoryRepo.AddCategory(newcategory)
	if err != nil {
		return newCategory, err
	}
	return newCategory, nil
}

func (c CategoryUseCase) UpdateCategory(idtoch int, newcategory string) (domain.Category, error) {

	var category = &domain.Category{Name: newcategory}
	err := middleware.ValidateCategory(newcategory)
	if err != nil {
		return *category, err
	}
	err = c.CategoryRepo.CheckCategory(newcategory)
	if err == nil {
		return *category, err
	}
	*category, err = c.CategoryRepo.GetCategoryById(idtoch)
	if err != nil {
		return *category, errors.New("category to update deos not exist")
	}
	*category, err = c.CategoryRepo.UpdateCategory(idtoch, newcategory)
	if err != nil {
		return *category, err
	}
	return *category, err
}

func (c CategoryUseCase) DeleteCategory(categoryID int) error {
	_, err := c.CategoryRepo.GetCategoryById(categoryID)
	if err != nil {
		return errors.New("category deos not exist")
	}
	err = c.CategoryRepo.DeleteCategory(categoryID)
	if err != nil {
		return err
	}
	return nil
}

func (c CategoryUseCase) ListCategories() ([]domain.Category, error) {
	list, err := c.CategoryRepo.ListCategories(8)
	return list, err
}
