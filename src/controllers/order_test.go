package controllers_test

import (
	"bookstore/src/controllers"
	"bookstore/src/models"
	"bookstore/src/routes"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestOrderController(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.LoadHTMLGlob("../../templates/*")
	routes.SetupRoutes(router)
	user := &models.User{
		Name:     "testuser",
		Email:    "nadeer@qburst.com",
		Password: "password123",
		Role:     "user",
	}
	if err := user.CreateUser(); err != nil {
		t.Fatal("Failed to create test user:", err)
	}
	var createdUser models.User
	models.GetUserByEmail(user.Email, &createdUser)
	if createdUser.ID == 0 {
		t.Fatal("User was not properly created")
	}
	user.ID = createdUser.ID
	t.Cleanup(func() {
		models.DeleteUserByName(user.Name)
		models.DeleteBookIsbn("123456789")
	})
	book := &models.Book{Title: "Test Book 1", Author: "Author 1", Price: 29.99, Isbn: "123456789",
		PublishedDate: "01-01-2000", CategoryID: models.GetCategories()[0].ID}
	book.CreateBook()
	var createdBook models.Book
	models.GetBookByIsbn(book.Isbn, &createdBook)
	book_id_str := strconv.FormatUint(uint64(createdBook.ID), 10)
	tests := []struct {
		name       string
		setupAuth  func(*http.Request)
		wantStatus int
		wantBooks  bool
	}{
		{
			name: "Success - Authenticated User",
			setupAuth: func(req *http.Request) {
				token, _ := controllers.GenerateJWT(user.ID)
				req.AddCookie(&http.Cookie{
					Name:  "jwt_token",
					Value: token,
					Path:  "/",
				})
			},
			wantStatus: http.StatusOK,
			wantBooks:  true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			var cart models.Cart
			models.AddToCart(uint(createdBook.ID), uint(user.ID), cart)
			req, _ := http.NewRequest("GET", "/order/"+book_id_str, nil)
			test.setupAuth(req)
			router.ServeHTTP(w, req)
			assert.Equal(t, http.StatusFound, w.Code)
			assert.Equal(t, w.Header().Get("location"), "/cart")
		})
	}
}
