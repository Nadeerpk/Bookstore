package repository

import (
	"bookstore/internal/domain/models"
	"strings"

	"gorm.io/gorm"
)

type bookRepository struct {
	db *gorm.DB
}

func NewBookRepository(db *gorm.DB) BookRepository {
	return &bookRepository{db: db}
}
func (r *bookRepository) GetAllBooks(books *[]models.Book) error {
	err := r.db.Preload("Category").Find(&books)
	return err.Error
}

func (r *bookRepository) GetBookByID(id string, book *models.Book) error {
	err := r.db.Preload("Category").Where("id = ?", id).First(&book).Error
	return err
}

func (r *bookRepository) CreateBook(book *models.Book) *models.Book {
	r.db.Create(&book)
	return book
}

func (r *bookRepository) UpdateBook(book *models.Book) error {
	if err := r.db.First(&models.Book{}, book.ID).Error; err != nil {
		return err
	}

	return r.db.Model(&models.Book{ID: book.ID}).Updates(models.Book{
		Title:         book.Title,
		Author:        book.Author,
		Price:         book.Price,
		Isbn:          book.Isbn,
		PublishedDate: book.PublishedDate,
		Availability:  book.Availability,
		CategoryID:    book.CategoryID,
		Image:         book.Image,
	}).Error
}

func (r *bookRepository) DeleteBook(book *models.Book) error {
	err := r.db.Delete(&book)
	return err.Error
}
func (r *bookRepository) DeleteBookIsbn(isbn string) {
	r.db.Where("isbn = ?", isbn).Delete(&models.Book{})
}

func (r *bookRepository) SearchBooks(title, author, genre, isbn, availability, yearFrom, yearTo,
	title_sort, author_sort, year_sort string, books *[]models.Book) error {
	query := r.db.Preload("Category").Joins("Category")

	if title != "" {
		query = query.Where("LOWER(books.title) LIKE ?", "%"+strings.ToLower(title)+"%")
	}
	if author != "" {
		query = query.Where("LOWER(books.author) LIKE ?", "%"+strings.ToLower(author)+"%")
	}
	if genre != "" {
		query = query.Where("LOWER(name) LIKE ?", "%"+strings.ToLower(genre)+"%")
	}
	if isbn != "" {
		query = query.Where("LOWER(books.isbn) LIKE ?", "%"+strings.ToLower(isbn)+"%")
	}
	if availability != "" {
		query = query.Where("books.availability = ?", availability)
	}
	if yearFrom != "" {
		query = query.Where("SUBSTRING(books.published_date, 1, 4) >= ?", yearFrom)
	}
	if yearTo != "" {
		query = query.Where("SUBSTRING(books.published_date, 1, 4) <= ?", yearTo)
	}
	if title_sort != "" {
		query = query.Order("books.title " + title_sort)
	}
	if author_sort != "" {
		query = query.Order("books.author " + author_sort)
	}
	if year_sort != "" {
		query = query.Order("books.published_date " + year_sort)
	}

	err := query.Debug().Find(&books).Error
	return err
}
func (r *bookRepository) GetCategories() []models.Category {
	var categories []models.Category
	r.db.Find(&categories)
	return categories
}
