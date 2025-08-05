package repository

import (
	"bookstore/internal/domain/models"

	"gorm.io/gorm"
)

type reviewRepository struct {
	db *gorm.DB
}

func NewReviewRepository(db *gorm.DB) ReviewRepository {
	return &reviewRepository{db: db}
}
func (r *reviewRepository) AddReview(review *models.Review) {
	r.db.Create(&review)
}

func (r *reviewRepository) GetReviewsByBookID(bookID uint, reviews *[]models.Review) error {
	err := r.db.Preload("User").Where("book_id = ?", bookID).Find(&reviews).Error
	return err
}
