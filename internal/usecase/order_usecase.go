package usecase

import (
	"bookstore/internal/domain/models"
	"bookstore/internal/domain/repository"
)

type OrderUseCase interface {
	AddOrder(order *models.Order) error
	GetOrdersByUserID(userID uint, orders *[]models.Order) error
}
type orderUseCase struct {
	OrderRepo repository.OrderRepository
}

func NewOrderUsecase(orderRepo repository.OrderRepository) OrderUseCase {
	return &orderUseCase{OrderRepo: orderRepo}
}
func (u *orderUseCase) AddOrder(order *models.Order) error {
	return u.OrderRepo.AddOrder(order)
}
func (u *orderUseCase) GetOrdersByUserID(userID uint, orders *[]models.Order) error {
	return u.OrderRepo.GetOrdersByUserID(userID, orders)
}
