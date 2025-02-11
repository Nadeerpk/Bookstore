package main

import (
	"bookstore/internal/delivery/http/handlers"
	"bookstore/internal/delivery/http/routes"
	"bookstore/internal/repository"
	"bookstore/internal/usecase"
	"bookstore/src/models"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	godotenv.Load()
	username := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	dbname := os.Getenv("DB_NAME")
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", username, password, host, port, dbname)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&models.User{}, &models.Book{}, &models.Category{}, &models.Order{}, &models.Cart{})

	userRepo := repository.NewUserRepository(db)
	bookRepo := repository.NewBookRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)
	cartRepo := repository.NewCartRepository(db)
	orderRepo := repository.NewOrderRepository(db)
	reviewRepo := repository.NewReviewRepository(db)

	userUseCase := usecase.NewUserUsecase(userRepo)
	bookUseCase := usecase.NewBookUsecase(bookRepo)
	categoryUseCase := usecase.NewCategoryUsecase(categoryRepo)
	cartUseCase := usecase.NewCartUsecase(cartRepo)
	orderUseCase := usecase.NewOrderUsecase(orderRepo)
	reviewUseCase := usecase.NewReviewUsecase(reviewRepo)

	userhandler := handlers.NewUserHandler(userUseCase)
	bookhandler := handlers.NewBookHandler(bookUseCase, userUseCase, categoryUseCase)
	categoryhandler := handlers.NewCategoryHandler(categoryUseCase)
	carthandler := handlers.NewCartHandler(cartUseCase)
	orderhandler := handlers.NewOrderHandler(orderUseCase, cartUseCase)
	reviewhandler := handlers.NewReviewHandler(reviewUseCase)

	r := routes.SetupRoutes(db, userhandler, bookhandler, carthandler, orderhandler, reviewhandler, categoryhandler)
	r.Run(":8081")

}
