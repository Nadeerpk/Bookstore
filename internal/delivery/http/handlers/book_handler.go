package handlers

import (
	"bookstore/internal/domain/models"
	"bookstore/internal/usecase"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type BookHandler struct {
	BookUsecase     usecase.BookUseCase
	UserUsecase     usecase.UserUseCase
	CategoryUsecase usecase.CategoryUseCase
}

func NewBookHandler(bookUseCase usecase.BookUseCase, userUseCase usecase.UserUseCase,
	categoryUseCase usecase.CategoryUseCase) *BookHandler {
	return &BookHandler{
		BookUsecase:     bookUseCase,
		UserUsecase:     userUseCase,
		CategoryUsecase: categoryUseCase,
	}
}

func (h *BookHandler) EditBook(c *gin.Context) {
	id := c.Param("id")
	val, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		log.Fatalf("Error converting string to uint: %v", err)
	}
	uint_id := uint(val)
	var book models.Book
	if err := h.BookUsecase.GetBookByID(uint_id, &book); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching book details"})
		return
	}
	if c.Request.Method == http.MethodPost {
		if err := c.Request.ParseForm(); err != nil {
			fmt.Printf("Error parsing form: %v\n", err)
		}
		c.ShouldBind(&book)
		if err := h.BookUsecase.UpdateBook(&book); err != nil {
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
		"categories":    h.CategoryUsecase.GetCategories(),
	})
}
func (h *BookHandler) AddBook(c *gin.Context) {
	Book := &models.Book{}
	c.ShouldBind(Book)

	AttachImage(c, Book)

	h.BookUsecase.CreateBook(Book)
	c.Redirect(http.StatusFound, "/books")
}

func (h *BookHandler) UpdateBook(c *gin.Context) {
	Book := &models.Book{}
	c.ShouldBind(Book)
	AttachImage(c, Book)
	fmt.Println(Book)
	_ = h.BookUsecase.UpdateBook(Book)
	c.HTML(http.StatusOK, "books.html", gin.H{})
}

func (h *BookHandler) ListBooks(c *gin.Context) {
	var books []models.Book
	if err := h.BookUsecase.GetAllBooks(&books); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching books"})
		return
	}
	user := &models.User{}
	user_id := c.GetFloat64("user_id")
	h.UserUsecase.GetUserByID(user, uint(user_id))
	c.HTML(http.StatusOK, "books.html", gin.H{
		"books":   books,
		"isadmin": user.Role == "admin",
	})
}

func (h *BookHandler) DeleteBook(c *gin.Context) {
	id := c.Param("id")
	val, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		log.Fatalf("Error converting string to uint: %v", err)
	}
	uint_id := uint(val)
	var book models.Book
	if err := h.BookUsecase.GetBookByID(uint_id, &book); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching book details"})
		return
	}
	if err := h.BookUsecase.DeleteBook(&book); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting book"})
		return
	}
	c.Redirect(http.StatusFound, "/books")
}

func (h *BookHandler) BookSearch(c *gin.Context) {
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
	if err := h.BookUsecase.SearchBooks(title, author, genre, isbn, availability, yearFrom, yearTo,
		title_sort, author_sort, year_sort, &books); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching books"})
		return
	}
	c.HTML(http.StatusOK, "books.html", gin.H{
		"books": books,
	})
}

func AttachImage(c *gin.Context, Book *models.Book) {
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
}
