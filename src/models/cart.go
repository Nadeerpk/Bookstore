package models

import (
	"bookstore/src/config"

	"gorm.io/gorm"
)

type Cart struct {
	gorm.Model
	ID     uint    `json:"id" form:"id" gorm:"primary_key"`
	UserID uint    `json:"user_id" form:"user_id"`
	User   User    `json:"user" form:"user" gorm:"foreignKey:UserID"`
	Books  []*Book `json:"books" form:"books" gorm:"many2many:cart_items;constraint:OnDelete:CASCADE"`
}
type CartItem struct {
	gorm.Model
	CartID   uint `json:"cart_id" form:"cart_id" gorm:"primary_key"`
	BookID   uint `json:"book_id" form:"book_id" gorm:"primary_key"`
	Quantity uint `json:"quantity" form:"quantity" gorm:"default:1"`
}

func init() {
	db = config.Getdb()
	db.AutoMigrate(&Cart{}, &CartItem{})
}

func GetCart(cart *Cart, userID uint) error {
	err := db.Preload("Books", func(db *gorm.DB) *gorm.DB {
		return db.Joins("JOIN cart_items ON cart_items.book_id = books.id AND cart_items.deleted_at IS NULL")
	}).Where("user_id = ?", userID).
		First(cart).Error
	return err
}

func AddToCart(bookID, user_id uint, cart Cart) error {
	var book Book

	if err := db.Where("id = ?", bookID).First(&book).Error; err != nil {
		return err
	}

	if err := db.Where("user_id = ?", user_id).First(&cart).Error; err != nil {
		cart = Cart{UserID: user_id}
		db.Create(&cart)
		cartitem := CartItem{CartID: cart.ID, BookID: bookID}
		db.Create(&cartitem)
		return nil
	}
	CartItem := CartItem{CartID: cart.ID, BookID: bookID}
	err := db.Create(&CartItem)
	return err.Error
}
func DeleteFromCart(bookID, user_id uint) error {
	var cart Cart
	if err := db.Where("user_id = ?", user_id).First(&cart).Error; err != nil {
		return err
	}
	err := db.Where("cart_id = ? AND book_id = ?", cart.ID, bookID).Delete(&CartItem{}).Error
	return err
}
