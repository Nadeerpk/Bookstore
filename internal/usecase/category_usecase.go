package usecase

import (
	"bookstore/internal/domain/models"
	"bookstore/internal/repository"
)

type CategoryUseCase interface {
	GetCategories() []models.Category
	AddCategory(category *models.Category)
}
type categoryUseCase struct {
	CategoryRepo repository.CategoryRepository
}

func NewCategoryUsecase(categoryRepo repository.CategoryRepository) CategoryUseCase {
	return &categoryUseCase{CategoryRepo: categoryRepo}
}
func (u *categoryUseCase) GetCategories() []models.Category {
	return u.CategoryRepo.GetCategories()
}
func (u *categoryUseCase) AddCategory(category *models.Category) {
	u.CategoryRepo.AddCategory(category)
}
