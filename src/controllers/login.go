package controllers

import (
	"bookstore/src/models"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

var jwtKey = []byte("jwt-key")

func LoginController(c *gin.Context) {
	User := &models.User{}
	c.ShouldBind(User)
	fmt.Println(User)
	authenticated := models.AuthenticateUser(User)
	if !authenticated {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication failed"})
		return
	}
	token, err := GenerateJWT(User.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating token"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Successfully logged in", "jwt-token": token})

}
func GenerateJWT(username string) (string, error) {
	claims := &jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(1 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}
