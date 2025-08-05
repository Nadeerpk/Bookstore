package handlers_test

import (
	"bookstore/internal/delivery/http/handlers"
	"bookstore/internal/delivery/http/routes"
	"bookstore/internal/domain/models"
	"bookstore/src/controllers"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func createTestBooks(books []models.Book, mockBook *MockBookUsecase) error {
	for _, book := range books {
		mockBook.CreateBook(&book)
	}
	return nil
}

func TestGetBooks(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.LoadHTMLGlob("../../../../templates/*")

	mockBookUsecase := new(MockBookUsecase)
	mockUserUsecase := new(MockUserUsecase)
	mockCategoryUsecase := new(MockCategoryUsecase)
	bookHandler := handlers.NewBookHandler(mockBookUsecase, mockUserUsecase, mockCategoryUsecase)

	// Create a protected group with JWT middleware
	protected := router.Group("/")
	protected.Use(routes.JwtHandler())
	protected.GET("/books", bookHandler.ListBooks)

	user := &models.User{
		Name:     "testuser",
		Email:    "test@example.com",
		Password: "password123",
		Role:     "user",
	}
	mockUserUsecase.On("CreateUser", user).Return(nil)
	mockUserUsecase.On("GetUserByEmail", user.Email, mock.AnythingOfType("*models.User")).Run(func(args mock.Arguments) {
		userArg := args[1].(*models.User)
		// Copy all fields including ID
		*userArg = models.User{
			ID:       1, // Explicitly set ID
			Name:     user.Name,
			Email:    user.Email,
			Password: user.Password,
			Role:     user.Role,
		}
	}).Return()
	mockUserUsecase.On("DeleteUserByName", user.Name).Return(nil)
	mockUserUsecase.On("GetUserByID", mock.AnythingOfType("*models.User"), uint(1)).Return()
	mockBookUsecase.On("DeleteBookIsbn", "123456789").Return()
	mockBookUsecase.On("DeleteBookIsbn", "1234567891").Return()
	mockCategoryUsecase.On("GetCategories").Return([]models.Category{{ID: 1, Name: "Test Category"}})
	mockCategoryUsecase.On("GetCategories").Return([]models.Category{{ID: 1, Name: "Test Category"}})
	mockBookUsecase.On("CreateBook", mock.AnythingOfType("*models.Book")).Return(&models.Book{ID: 1})
	books := []models.Book{
		{Title: "Test Book 1", Author: "Author 1", Price: 29.99, Isbn: "123456789",
			PublishedDate: "01-01-2000", CategoryID: mockCategoryUsecase.GetCategories()[0].ID},
		{Title: "Test Book 2", Author: "Author 2", Price: 19.99, Isbn: "1234567891",
			PublishedDate: "01-01-2000", CategoryID: mockCategoryUsecase.GetCategories()[0].ID}}
	mockBookUsecase.On("GetAllBooks", mock.AnythingOfType("*[]models.Book")).Run(func(args mock.Arguments) {
		booksArg := args[0].(*[]models.Book)
		*booksArg = books
	}).Return(nil)
	if err := mockUserUsecase.CreateUser(user); err != nil {
		t.Fatal("Failed to create test user:", err)
	}

	var createdUser models.User
	mockUserUsecase.GetUserByEmail(user.Email, &createdUser)
	if createdUser.ID == 0 {
		t.Fatal("User was not properly created")
	}

	user.ID = createdUser.ID

	if err := createTestBooks(books, mockBookUsecase); err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		mockUserUsecase.DeleteUserByName(user.Name)
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
			mockBookUsecase.DeleteBookIsbn("123456789")
			mockBookUsecase.DeleteBookIsbn("1234567891")
		})
	}
}
