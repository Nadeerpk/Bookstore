package routes

import (
	"bookstore/src/controllers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	router.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", gin.H{})
	})
	router.POST("/register", controllers.RegisterController)
	router.POST("/login", controllers.LoginController)
}
