package repository

import (
	"bookstore/internal/domain/models"

	"gorm.io/gorm"
)

type cartRepository struct {
	db *gorm.DB
}

func NewCartRepository(db *gorm.DB) CartRepository {
	return &cartRepository{db: db}
}
func (r *cartRepository) GetCart(cart *models.Cart, userID uint) error {
	err := r.db.Preload("Items.Book").Where("user_id = ?", userID).First(cart).Error
	return err
}

func (r *cartRepository) AddToCart(bookID, user_id uint, cart models.Cart) error {
	var book models.Book

	if err := r.db.Where("id = ?", bookID).First(&book).Error; err != nil {
		return err
	}

	if err := r.db.Where("user_id = ?", user_id).Find(&cart); err != nil {
		cart = models.Cart{UserID: user_id}
		r.db.Create(&cart)
		cartitem := models.CartItem{CartID: cart.ID, BookID: bookID}
		r.db.Create(&cartitem)
		return nil
	}
	CartItem := models.CartItem{CartID: cart.ID, BookID: bookID}
	err := r.db.Create(&CartItem)
	return err.Error
}
func (r *cartRepository) DeleteFromCart(bookID, user_id uint) error {
	var cart models.Cart
	if err := r.db.Where("user_id = ?", user_id).First(&cart).Error; err != nil {
		return err
	}
	err := r.db.Where("cart_id = ? AND book_id = ?", cart.ID, bookID).Delete(&models.CartItem{}).Error
	return err
}
