package usecase

import (
	"bookstore/internal/domain/models"
	"bookstore/internal/domain/repository"
)

type CartUseCase interface {
	GetCart(cart *models.Cart, userID uint) error
	AddToCart(bookID, user_id uint, cart models.Cart) error
	DeleteFromCart(bookID, user_id uint) error
}
type cartUseCase struct {
	CartRepo repository.CartRepository
}

func NewCartUsecase(cartRepo repository.CartRepository) CartUseCase {
	return &cartUseCase{CartRepo: cartRepo}
}

func (u *cartUseCase) GetCart(cart *models.Cart, userID uint) error {
	return u.CartRepo.GetCart(cart, userID)
}
func (u *cartUseCase) AddToCart(bookID, user_id uint, cart models.Cart) error {
	return u.CartRepo.AddToCart(bookID, user_id, cart)
}
func (u *cartUseCase) DeleteFromCart(bookID, user_id uint) error {
	return u.CartRepo.DeleteFromCart(bookID, user_id)
}
