package controllers

import (
	"bookstore/src/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

var jwtKey = []byte("jwt-key")

func LoginController(c *gin.Context) {
	User := &models.User{}
	c.ShouldBind(User)
	authenticated := models.AuthenticateUser(User)
	if !authenticated {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication failed"})
		return
	}
	var user models.User
	models.GetUser(&user, User)
	token, err := GenerateJWT(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating token"})
		return
	}
	c.SetCookie("jwt_token", token, 3600, "/", "localhost", false, true)

	c.Redirect(http.StatusFound, "/books")
}

func GenerateJWT(user_id uint) (string, error) {
	claims := &jwt.MapClaims{
		"user_id": user_id,
		"exp":     time.Now().Add(1 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

func LogoutController(c *gin.Context) {
	c.SetCookie("jwt_token", "", -1, "/", "localhost", false, true)
	c.Redirect(http.StatusFound, "/login")
}
