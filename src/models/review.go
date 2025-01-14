package models

import (
	"bookstore/src/config"
)

type Review struct {
	ID      uint   `json:"id" form:"id" gorm:"primary_key"`
	BookID  uint   `json:"book_id" form:"book_id"`
	UserID  uint   `json:"user_id" form:"user_id"`
	User    User   `gorm:"foreignKey:UserID"`
	Comment string `json:"comment" form:"comment"`
}

func init() {
	db = config.Getdb()
	db.AutoMigrate(&Review{})
}

func AddReview(review *Review) {
	db.Create(&review)
}

func GetReviewsByBookID(bookID uint, reviews *[]Review) error {
	err := db.Preload("User").Where("book_id = ?", bookID).Find(&reviews).Error
	return err
}
