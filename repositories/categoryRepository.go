package repositories

import (
	"februaryMVCProject/entities"
	"gorm.io/gorm"
)

type CategoryRepository interface {
	AddCategory(category entities.Category)
	DeleteCategory(id int)
	UpdateCategory(updCategory entities.Category)
	GetAllCategories() []entities.Category
	GetCategory(id int) entities.Category
}

type categoryRepository struct {
	gormDB *gorm.DB
}

func NewCategoryRepository(gormDB *gorm.DB) CategoryRepository {
	return &categoryRepository{gormDB: gormDB}
}

func (categoryRepository *categoryRepository) AddCategory(category entities.Category) {
	categoryRepository.gormDB.Create(&category)
}

func (categoryRepository *categoryRepository) DeleteCategory(id int) {
	categoryRepository.gormDB.Delete(&entities.Category{}, id)
}

func (categoryRepository *categoryRepository) UpdateCategory(updCategory entities.Category) {
	categoryRepository.gormDB.Save(&updCategory)
}

func (categoryRepository *categoryRepository) GetAllCategories() []entities.Category {
	var categories []entities.Category
	categoryRepository.gormDB.Find(&categories)
	return categories
}

func (categoryRepository *categoryRepository) GetCategory(id int) entities.Category {
	var category entities.Category
	categoryRepository.gormDB.First(&category, id)
	return category
}
