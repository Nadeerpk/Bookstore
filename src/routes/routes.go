package routes

import (
	"bookstore/src/controllers"
	"bookstore/src/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

var jwtKey = []byte("jwt-key")

func SetupRoutes(router *gin.Engine) {
	router.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", gin.H{})
	})
	router.GET("/register", func(c *gin.Context) {
		c.HTML(200, "register.html", gin.H{})
	})
	router.GET("/login", func(c *gin.Context) {
		c.HTML(200, "login.html", gin.H{})
	})
	router.POST("/register", controllers.RegisterController)
	router.POST("/login", controllers.LoginController)
	userGroup := router.Group("/", jwtHandler())
	userGroup.GET("/logout", controllers.LogoutController)
	userGroup.GET("/books", controllers.BooksController)
	userGroup.GET("/book-search", controllers.BookSearchController)

	userGroup.POST("/add-to-cart/:book_id", controllers.AddToCartController)
	userGroup.GET("/cart", controllers.ShowCartController)
	userGroup.GET("/categories", controllers.GetCategories)

	userGroup.POST("add-review/:book_id", controllers.AddReview)
	userGroup.GET("/reviews/:book_id", controllers.GetReviews)
	userGroup.GET("/order/:book_id", controllers.AddOrderController)
	userGroup.GET("/remove-from-cart/:book_id", controllers.RemoveFromCartController)
	userGroup.GET("/order-history", controllers.OrderHistoryController)
	adminGroup := userGroup.Group("/", adminHandler())
	adminGroup.POST("/add-book", controllers.AddBookController)
	adminGroup.GET("/add-book", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "add-book.html", gin.H{"categories": models.GetCategories()})
	})
	adminGroup.GET("/edit-book/:id", controllers.EditBookController)
	adminGroup.POST("/edit-book/:id", controllers.EditBookController)
	adminGroup.POST("/delete-book/:id", controllers.DeleteBookController)
	adminGroup.POST("/add-category", controllers.AddCategory)
}

func jwtHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := c.Cookie("jwt_token")
		if err != nil {
			c.Redirect(http.StatusFound, "/login")
			c.Abort()
			return
		}
		claims := jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil || !token.Valid {
			c.Redirect(http.StatusFound, "/login")
			c.Abort()
			return
		}
		user_id := claims["user_id"].(float64)
		c.Set("user_id", user_id)
		c.Next()
	}
}

func adminHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		user_id := c.GetFloat64("user_id")
		var user models.User
		models.GetUserByID(&user, uint(user_id))
		if user.Role != "admin" {
			c.Redirect(http.StatusFound, "/login")
			c.Abort()
			return
		}
		c.Next()
	}
}
