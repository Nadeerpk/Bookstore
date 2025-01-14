package controllers

import (
	"bookstore/src/models"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func EditBookController(c *gin.Context) {
	id := c.Param("id")
	var book models.Book
	if err := models.GetBookByID(id, &book); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching book details"})
		return
	}
	if c.Request.Method == http.MethodPost {
		if err := c.Request.ParseForm(); err != nil {
			fmt.Printf("Error parsing form: %v\n", err)
		}
		c.ShouldBind(&book)
		if err := book.UpdateBook(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating book"})
			return
		}
		c.Redirect(http.StatusFound, "/books")
		return
	}

	c.HTML(http.StatusOK, "edit-book.html", gin.H{
		"ID":            book.ID,
		"Title":         book.Title,
		"Author":        book.Author,
		"Genre":         book.Category.Name,
		"Price":         book.Price,
		"PublishedDate": book.PublishedDate,
		"Isbn":          book.Isbn,
		"Availability":  book.Availability,
		"categories":    models.GetCategories(),
	})
}
