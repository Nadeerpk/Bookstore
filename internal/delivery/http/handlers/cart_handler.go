package handlers

import (
	"bookstore/internal/domain/models"
	"bookstore/internal/usecase"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CartHandler struct {
	CartUsecase usecase.CartUseCase
}

func NewCartHandler(cartUsecase usecase.CartUseCase) *CartHandler {
	return &CartHandler{CartUsecase: cartUsecase}
}
func (h *CartHandler) ShowCart(c *gin.Context) {
	user_id := c.GetFloat64("user_id")
	var cart models.Cart
	h.CartUsecase.GetCart(&cart, uint(user_id))
	c.HTML(http.StatusOK, "cart.html", gin.H{"cart": cart})
}

func (h *CartHandler) AddToCart(c *gin.Context) {
	book_id_str := c.Param("book_id")
	user_id := c.GetFloat64("user_id")
	book_id, _ := strconv.ParseUint(book_id_str, 10, 64)
	var cart models.Cart
	err := h.CartUsecase.AddToCart(uint(book_id), uint(user_id), cart)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Redirect(http.StatusFound, "/cart")
}

func (h *CartHandler) RemoveFromCart(c *gin.Context) {
	book_id_str := c.Param("book_id")
	user_id := c.GetFloat64("user_id")
	book_id, _ := strconv.ParseInt(book_id_str, 10, 64)
	h.CartUsecase.DeleteFromCart(uint(book_id), uint(user_id))
	c.Redirect(http.StatusFound, "/cart")
}
