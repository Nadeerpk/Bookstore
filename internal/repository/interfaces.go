package repository

import (
	"bookstore/internal/domain/models"
)

type UserRepository interface {
	CreateUser(u *models.User) error
	AuthenticateUser(RequestUser *models.User) bool
	GetUser(user *models.User, RequestUser *models.User)
	GetUserByID(user *models.User, user_id uint)
	DeleteUserByName(name string) error
	GetUserByEmail(email string, user *models.User)
}
type BookRepository interface {
	GetAllBooks(books *[]models.Book) error
	GetBookByID(id string, book *models.Book) error
	CreateBook(book *models.Book) *models.Book
	UpdateBook(book *models.Book) error
	DeleteBook(book *models.Book) error
	DeleteBookIsbn(isbn string)
	SearchBooks(title, author, genre, isbn, availability, yearFrom, yearTo, titleSort, authorSort, yearSort string, books *[]models.Book) error
	GetCategories() []models.Category
}
type CartRepository interface {
	GetCart(cart *models.Cart, userID uint) error
	AddToCart(bookID, user_id uint, cart models.Cart) error
	DeleteFromCart(bookID, user_id uint) error
}
type ReviewRepository interface {
	AddReview(review *models.Review)
	GetReviewsByBookID(bookID uint, reviews *[]models.Review) error
}
type CategoryRepository interface {
	AddCategory(category *models.Category)
	DeleteCategory(category *models.Category)
	UpdateCategory(category *models.Category)
	GetCategories() []models.Category
}
type OrderRepository interface {
	AddOrder(order *models.Order) error
	GetOrdersByUserID(userID uint, orders *[]models.Order) error
}
