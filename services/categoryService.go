package services

import (
	"februaryMVCProject/entities"
	"februaryMVCProject/repositories"
)

type CategoryService interface {
	AddCategory(category entities.Category)
	DeleteCategory(id int)
	UpdateCategory(updCategory entities.Category)
	GetAllCategories() []entities.Category
	GetCategory(id int) entities.Category
}

type categoryService struct {
	categoryRepository repositories.CategoryRepository
}

func NewCategoryService(categoryRepository repositories.CategoryRepository) CategoryService {
	return &categoryService{categoryRepository: categoryRepository}
}

func (categoryService *categoryService) AddCategory(category entities.Category) {
	categoryService.categoryRepository.AddCategory(category)
}

func (categoryService *categoryService) DeleteCategory(id int) {
	categoryService.categoryRepository.DeleteCategory(id)
}

func (categoryService *categoryService) UpdateCategory(updCategory entities.Category) {
	categoryService.categoryRepository.UpdateCategory(updCategory)
}

func (categoryService *categoryService) GetAllCategories() []entities.Category {
	return categoryService.categoryRepository.GetAllCategories()
}

func (categoryService *categoryService) GetCategory(id int) entities.Category {
	return categoryService.categoryRepository.GetCategory(id)
}
