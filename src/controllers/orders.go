package controllers

import (
	"bookstore/src/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func AddOrderController(c *gin.Context) {
	book_id_str := c.Param("book_id")
	user_id := c.GetFloat64("user_id")
	book_id, _ := strconv.ParseUint(book_id_str, 10, 64)
	order := &models.Order{}
	order.UserID = uint(user_id)
	order.BookID = uint(book_id)
	models.AddOrder(order)
	models.DeleteFromCart(uint(book_id), uint(user_id))
	c.Redirect(http.StatusFound, "/cart")
}
