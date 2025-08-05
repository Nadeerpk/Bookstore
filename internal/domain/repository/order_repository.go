package repository

import (
	"bookstore/internal/domain/models"
	"bookstore/src/utils"

	"gorm.io/gorm"
)

type orderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) OrderRepository {
	return &orderRepository{db: db}
}

func (r *orderRepository) AddOrder(order *models.Order) error {
	r.db.Create(&order)
	r.db.Model(&order).Preload("User").Preload("Book").First(&order)

	subject := "Subject: Your Order Has Been Placed!\n"
	body := "Your order from the Bookstore for " + order.Book.Title + " has been placed successfully.\n"
	to := []string{order.User.Email}
	err := utils.Send_mail(to, subject, body)
	if err != nil {
		return err
	}
	return nil
}

func (r *orderRepository) GetOrdersByUserID(userID uint, orders *[]models.Order) error {
	err := r.db.Preload("Book").Where("user_id = ?", userID).Find(&orders).Error
	return err
}
