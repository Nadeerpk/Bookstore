package controllers

import (
	"bookstore/src/models"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func RegisterController(c *gin.Context) {
	User := &models.User{}
	c.ShouldBind(User)
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(User.Password), bcrypt.DefaultCost)
	User.Password = string(hashedPassword)
	fmt.Println(User)
	_ = User.CreateUser()
	c.HTML(http.StatusOK, "login.html", gin.H{})
}
