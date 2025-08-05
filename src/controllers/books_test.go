package controllers_test

import (
	"bookstore/src/controllers"
	"bookstore/src/models"
	"bookstore/src/routes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func createTestBooks() error {
	books := []models.Book{
		{Title: "Test Book 1", Author: "Author 1", Price: 29.99, Isbn: "123456789",
			PublishedDate: "01-01-2000", CategoryID: models.GetCategories()[0].ID},
		{Title: "Test Book 2", Author: "Author 2", Price: 19.99, Isbn: "1234567891",
			PublishedDate: "01-01-2000", CategoryID: models.GetCategories()[0].ID}}
	for _, book := range books {
		book.CreateBook()
	}
	return nil
}

func GetBooks(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.LoadHTMLGlob("../../templates/*")
	routes.SetupRoutes(router)

	user := &models.User{
		Name:     "testuser",
		Email:    "test@example.com",
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

	if err := createTestBooks(); err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		models.DeleteUserByName(user.Name)
	})

	tests := []struct {
		name       string
		setupAuth  func(*http.Request)
		wantStatus int
		wantBooks  bool
	}{
		{
			name: "Success - Authenticated User - List all books",
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
		{
			name:       "Failure - No Authentication - List all books",
			setupAuth:  func(req *http.Request) {},
			wantStatus: http.StatusFound,
			wantBooks:  false,
		},
		{
			name: "Failure - Invalid Token - List all books",
			setupAuth: func(req *http.Request) {
				req.AddCookie(&http.Cookie{
					Name:  "jwt_token",
					Value: "invalid_token",
					Path:  "/",
				})
			},
			wantStatus: http.StatusFound,
			wantBooks:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/books", nil)
			tt.setupAuth(req)

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)

			if tt.wantBooks {
				assert.Equal(t, tt.wantStatus, http.StatusOK)
				assert.Contains(t, w.Body.String(), "Test Book 1")
				assert.Contains(t, w.Body.String(), "Test Book 2")
				assert.NotContains(t, w.Body.String(), "Edit Book")
			}
			models.DeleteBookIsbn("123456789")
			models.DeleteBookIsbn("1234567891")
		})
	}
}
