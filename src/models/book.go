package models

import (
	"bookstore/src/config"
	"encoding/base64"
	"strings"

	"html/template"

	"gorm.io/gorm"
)

type Book struct {
	gorm.Model
	ID            uint     `json:"id" form:"id" gorm:"primary_key"`
	Title         string   `json:"title" form:"title" binding:"required" gorm:"not null"`
	Author        string   `json:"author" form:"author" binding:"required" gorm:"not null"`
	Price         float64  `json:"price" form:"price" binding:"required" gorm:"not null"`
	Category      Category `gorm:"foreignKey:CategoryID"`
	Isbn          string   `json:"isbn" form:"isbn" binding:"required" gorm:"not null; unique"`
	PublishedDate string   `json:"published_date" form:"published_date" binding:"required" gorm:"not null"`
	Availability  bool     `json:"availability" form:"availability" binding:"required" gorm:"not null"`
	Carts         []*Cart  `json:"carts" form:"carts" gorm:"many2many:cart_items"`
	CategoryID    uint     `json:"category_id" form:"category_id"`
	Reviews       []Review `json:"reviews" form:"reviews" gorm:"foreignKey:BookID"`
	Orders        []Order  `json:"orders" form:"orders" gorm:"foreignKey:BookID"`
	Image         []byte   `json:"image" form:"image" gorm:"type:mediumblob"`
}

func (b *Book) GetImageBase64() template.URL {
	if len(b.Image) == 0 {
		return ""
	}
	return template.URL("data:image/jpeg;base64," + base64.StdEncoding.EncodeToString(b.Image))
}

func init() {
	db = config.Getdb()
	db.AutoMigrate(&Book{})
}

func GetAllBooks(books *[]Book) error {
	err := db.Preload("Category").Find(&books)
	return err.Error
}

func GetBookByID(id string, book *Book) error {
	err := db.Preload("Category").Where("id = ?", id).First(&book).Error
	return err
}

func (book *Book) CreateBook() *Book {
	db.Create(&book)
	return book
}

func (book *Book) UpdateBook() error {
	if err := db.First(&Book{}, book.ID).Error; err != nil {
		return err
	}

	return db.Model(&Book{ID: book.ID}).Updates(Book{
		Title:         book.Title,
		Author:        book.Author,
		Price:         book.Price,
		Isbn:          book.Isbn,
		PublishedDate: book.PublishedDate,
		Availability:  book.Availability,
		CategoryID:    book.CategoryID,
	}).Error
}

func (book *Book) DeleteBook() error {
	err := db.Delete(&book)
	return err.Error
}

func SearchBooks(title, author, genre, isbn, availability, yearFrom, yearTo,
	title_sort, author_sort, year_sort string, books *[]Book) error {
	query := db.Preload("Category").Joins("Category")

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
