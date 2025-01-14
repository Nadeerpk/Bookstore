package controllers

import (
	"bookstore/src/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func ShowCartController(c *gin.Context) {
	user_id := c.GetFloat64("user_id")
	var cart models.Cart
	models.GetCart(&cart, uint(user_id))
	c.HTML(http.StatusOK, "cart.html", gin.H{"cart": cart})
}

func AddToCartController(c *gin.Context) {
	book_id_str := c.Param("book_id")
	user_id := c.GetFloat64("user_id")
	book_id, _ := strconv.ParseUint(book_id_str, 10, 64)
	var cart models.Cart
	err := models.AddToCart(uint(book_id), uint(user_id), cart)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Redirect(http.StatusFound, "/cart")
}

func RemoveFromCartController(c *gin.Context) {
	book_id_str := c.Param("book_id")
	user_id := c.GetFloat64("user_id")
	book_id, _ := strconv.ParseInt(book_id_str, 10, 64)
	models.DeleteFromCart(uint(book_id), uint(user_id))
	c.Redirect(http.StatusFound, "/cart")
}
