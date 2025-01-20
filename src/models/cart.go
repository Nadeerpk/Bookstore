package models

import (
	"bookstore/src/config"
)

type Cart struct {
	ID     uint       `json:"id" form:"id" gorm:"primary_key"`
	UserID uint       `json:"user_id" form:"user_id"`
	Items  []CartItem `gorm:"foreignKey:CartID;constraint:OnDelete:CASCADE"`
}

type CartItem struct {
	CartID   uint `json:"cart_id" form:"cart_id" gorm:"not null"`
	BookID   uint `json:"book_id" form:"book_id" gorm:"not null"`
	Quantity uint `json:"quantity" form:"quantity" gorm:"default:1"`
	Book     Book `json:"book" gorm:"foreignKey:BookID"`
}

func init() {
	db = config.Getdb()
	// db.Migrator().DropTable(&CartItem{}, &Cart{})
	db.AutoMigrate(&CartItem{}, &Cart{})
}

func GetCart(cart *Cart, userID uint) error {
	err := db.Preload("Items.Book").Where("user_id = ?", userID).First(cart).Error
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
