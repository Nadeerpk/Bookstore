package repository

import (
	"bookstore/internal/domain/models"

	"gorm.io/gorm"
)

type categoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &categoryRepository{db: db}
}
func (r *categoryRepository) GetCategories() []models.Category {
	var categories []models.Category
	r.db.Find(&categories)
	return categories
}
func (r *categoryRepository) AddCategory(category *models.Category) {
	r.db.Create(&category)
}

func (r *categoryRepository) DeleteCategory(category *models.Category) {
	r.db.Delete(&category)
}

func (r *categoryRepository) UpdateCategory(category *models.Category) {
	r.db.Save(&category)
}
