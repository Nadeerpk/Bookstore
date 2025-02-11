package routes

import (
	"bookstore/internal/delivery/http/handlers"
	"bookstore/src/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"gorm.io/gorm"
)

var jwtKey = []byte("jwt-key")

func SetupRoutes(db *gorm.DB, userHandler *handlers.UserHandler, bookHandler *handlers.BookHandler,
	cartHandler *handlers.CartHandler, orderHandler *handlers.OrderHandler,
	reviewHandler *handlers.ReviewHandler, categoryHandler *handlers.CategoryHandler) *gin.Engine {
	router := gin.Default()
	router.LoadHTMLGlob("templates/*")
	router.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", gin.H{})
	})
	router.GET("/register", func(c *gin.Context) {
		c.HTML(200, "register.html", gin.H{})
	})
	router.GET("/login", func(c *gin.Context) {
		c.HTML(200, "login.html", gin.H{})
	})
	router.POST("/register", userHandler.RegisterUser)
	router.POST("/login", userHandler.LoginUser)
	userGroup := router.Group("/", jwtHandler())
	userGroup.GET("/logout", userHandler.LogoutUser)
	userGroup.GET("/books", bookHandler.ListBooks)
	userGroup.GET("/book-search", bookHandler.BookSearch)

	userGroup.POST("/add-to-cart/:book_id", cartHandler.AddToCart)
	userGroup.GET("/cart", cartHandler.ShowCart)
	userGroup.GET("/categories", categoryHandler.GetCategories)

	userGroup.POST("add-review/:book_id", reviewHandler.AddReview)
	userGroup.GET("/reviews/:book_id", reviewHandler.GetReviews)
	userGroup.GET("/order/:book_id", orderHandler.AddOrder)
	userGroup.GET("/remove-from-cart/:book_id", cartHandler.RemoveFromCart)
	userGroup.GET("/order-history", orderHandler.OrderHistory)
	adminGroup := userGroup.Group("/", adminHandler())
	adminGroup.POST("/add-book", bookHandler.AddBook)
	adminGroup.GET("/add-book", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "add-book.html", gin.H{"categories": models.GetCategories()})
	})
	adminGroup.GET("/edit-book/:id", bookHandler.EditBook)
	adminGroup.POST("/edit-book/:id", bookHandler.UpdateBook)
	adminGroup.POST("/delete-book/:id", bookHandler.DeleteBook)
	adminGroup.POST("/add-category", categoryHandler.AddCategory)
	return router
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
