package models

import (
	"bookstore/src/config"
	"bookstore/src/utils"

	"gorm.io/gorm"
)

type Order struct {
	gorm.Model
	ID     uint `json:"id" form:"id" gorm:"primary_key"`
	UserID uint `json:"user_id" form:"user_id"`
	User   User `gorm:"foreignKey:UserID"`
	BookID uint `json:"book_id" form:"book_id"`
	Book   Book `gorm:"foreignKey:BookID"`
}

func init() {
	db = config.Getdb()
	db.AutoMigrate(&Order{})
}

func AddOrder(order *Order) {
	db.Create(&order)
	db.Model(&order).Preload("User").Preload("Book").First(&order)

	subject := "Subject: Your Order Has Been Placed!\n"
	body := "Your order from the Bookstore for " + order.Book.Title + " has been placed successfully.\n"
	to := []string{order.User.Email}
	utils.Send_mail(to, subject, body)
}
