package usecase

import (
	"bookstore/internal/domain/models"
	"bookstore/internal/domain/repository"
)

type BookUseCase interface {
	GetAllBooks(books *[]models.Book) error
	GetBookByID(id uint, book *models.Book) error
	CreateBook(book *models.Book) *models.Book
	UpdateBook(book *models.Book) error
	DeleteBook(book *models.Book) error
	DeleteBookIsbn(isbn string)
	SearchBooks(title, author, genre, isbn, availability, yearFrom, yearTo, titleSort, authorSort, yearSort string, books *[]models.Book) error
}
type bookUseCase struct {
	BookRepo repository.BookRepository
}

func NewBookUsecase(bookRepo repository.BookRepository) BookUseCase {
	return &bookUseCase{BookRepo: bookRepo}
}
func (u *bookUseCase) GetBookByID(id uint, book *models.Book) error {
	return u.BookRepo.GetBookByID(id, book)
}
func (u *bookUseCase) UpdateBook(book *models.Book) error {
	return u.BookRepo.UpdateBook(book)
}

func (u *bookUseCase) DeleteBook(book *models.Book) error {
	return u.BookRepo.DeleteBook(book)
}
func (u *bookUseCase) DeleteBookIsbn(isbn string) {
	u.BookRepo.DeleteBookIsbn(isbn)
}
func (u *bookUseCase) CreateBook(book *models.Book) *models.Book {
	return u.BookRepo.CreateBook(book)
}
func (u *bookUseCase) GetAllBooks(books *[]models.Book) error {
	return u.BookRepo.GetAllBooks(books)
}
func (u *bookUseCase) SearchBooks(title, author, genre, isbn, availability, yearFrom, yearTo, titleSort, authorSort, yearSort string, books *[]models.Book) error {
	return u.BookRepo.SearchBooks(title, author, genre, isbn, availability, yearFrom, yearTo, titleSort, authorSort, yearSort, books)
}
