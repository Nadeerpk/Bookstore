package controllers

import (
	"bookstore/src/models"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AddBookController(c *gin.Context) {
	Book := &models.Book{}
	c.ShouldBind(Book)

	file, err := c.FormFile("image")
	if err != nil {
		fmt.Printf("Error getting file: %v\n", err)
	} else {
		src, err := file.Open()
		if err != nil {
			fmt.Printf("Error opening file: %v\n", err)
		} else {
			defer src.Close()

			imageBytes := make([]byte, file.Size)
			n, err := src.Read(imageBytes)
			if err != nil {
				fmt.Printf("Error reading file: %v\n", err)
			} else {
				fmt.Printf("Read %d bytes from image\n", n)
				Book.Image = imageBytes
			}
		}
	}

	Book.CreateBook()
	c.Redirect(http.StatusFound, "/books")
}

func UpdateBookController(c *gin.Context) {
	Book := &models.Book{}
	c.ShouldBind(Book)
	fmt.Println(Book)
	_ = Book.UpdateBook()
	c.HTML(http.StatusOK, "books.html", gin.H{})
}

func BooksController(c *gin.Context) {
	var books []models.Book
	if err := models.GetAllBooks(&books); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching books"})
		return
	}
	user := &models.User{}
	user_id := c.GetFloat64("user_id")
	models.GetUserByID(user, uint(user_id))
	c.HTML(http.StatusOK, "books.html", gin.H{
		"books":   books,
		"isadmin": user.Role == "admin",
	})
}

func DeleteBookController(c *gin.Context) {
	id := c.Param("id")
	var book models.Book
	if err := models.GetBookByID(id, &book); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching book details"})
		return
	}
	if err := book.DeleteBook(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting book"})
		return
	}
	c.Redirect(http.StatusFound, "/books")
}

func BookSearchController(c *gin.Context) {
	title := c.Query("title")
	author := c.Query("author")
	genre := c.Query("genre")
	isbn := c.Query("isbn")
	availability := c.Query("availability")
	yearFrom := c.Query("year_from")
	yearTo := c.Query("year_to")
	title_sort := c.Query("title_sort")
	author_sort := c.Query("author_sort")
	year_sort := c.Query("year_sort")
	var books []models.Book
	if err := models.SearchBooks(title, author, genre, isbn, availability, yearFrom, yearTo,
		title_sort, author_sort, year_sort, &books); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching books"})
		return
	}
	c.HTML(http.StatusOK, "books.html", gin.H{
		"books": books,
	})
}
