package handlers

import (
	"bookstore/internal/domain/models"
	"bookstore/internal/usecase"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	OrderUsecase usecase.OrderUseCase
	CartUsecase  usecase.CartUseCase
}

func NewOrderHandler(orderUsecase usecase.OrderUseCase, cartUsecase usecase.CartUseCase) *OrderHandler {
	return &OrderHandler{OrderUsecase: orderUsecase, CartUsecase: cartUsecase}
}
func (h *OrderHandler) AddOrder(c *gin.Context) {
	book_id_str := c.Param("book_id")
	user_id := c.GetFloat64("user_id")
	book_id, _ := strconv.ParseUint(book_id_str, 10, 64)
	order := &models.Order{}
	order.UserID = uint(user_id)
	order.BookID = uint(book_id)
	err := h.OrderUsecase.AddOrder(order)
	if err != nil {
		c.Redirect(http.StatusFound, "/cart")
		return
	}
	h.CartUsecase.DeleteFromCart(uint(book_id), uint(user_id))
	c.Redirect(http.StatusFound, "/cart")
}

func (h *OrderHandler) OrderHistory(c *gin.Context) {
	user_id := c.GetFloat64("user_id")
	var orders []models.Order
	h.OrderUsecase.GetOrdersByUserID(uint(user_id), &orders)
	c.HTML(200, "orders.html", gin.H{"orders": orders})
}
