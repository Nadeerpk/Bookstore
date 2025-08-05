package usecase

import (
	"bookstore/internal/domain/models"
	"bookstore/internal/domain/repository"
)

type ReviewUseCase interface {
	AddReview(review *models.Review)
	GetReviewsByBookID(bookID uint, reviews *[]models.Review) error
}
type reviewUseCase struct {
	ReviewRepo repository.ReviewRepository
}

func NewReviewUsecase(reviewRepo repository.ReviewRepository) ReviewUseCase {
	return &reviewUseCase{ReviewRepo: reviewRepo}
}
func (c *reviewUseCase) AddReview(review *models.Review) {
	c.ReviewRepo.AddReview(review)
}
func (c *reviewUseCase) GetReviewsByBookID(bookID uint, reviews *[]models.Review) error {
	return c.ReviewRepo.GetReviewsByBookID(bookID, reviews)
}
