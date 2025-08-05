package handlers_test

import (
	"bookstore/internal/delivery/http/handlers"
	"bookstore/internal/delivery/http/routes"
	"bookstore/internal/domain/models"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

type MockBookUsecase struct {
	mock.Mock
}
type MockCategoryUsecase struct {
	mock.Mock
}

func (m *MockBookUsecase) GetAllBooks(books *[]models.Book) error {
	args := m.Called(books)
	return args.Error(0)
}
func (m *MockBookUsecase) GetBookByID(id uint, book *models.Book) error {
	args := m.Called(id, book)
	return args.Error(0)
}
func (m *MockBookUsecase) CreateBook(book *models.Book) *models.Book {
	args := m.Called(book)
	return args.Get(0).(*models.Book)
}
func (m *MockBookUsecase) UpdateBook(book *models.Book) error {
	args := m.Called(book)
	return args.Error(0)
}
func (m *MockBookUsecase) DeleteBook(book *models.Book) error {
	args := m.Called(book)
	return args.Error(0)
}
func (m *MockBookUsecase) DeleteBookIsbn(isbn string) {
	m.Called(isbn)
}
func (m *MockBookUsecase) SearchBooks(title, author, genre, isbn, availability, yearFrom, yearTo, titleSort, authorSort, yearSort string, books *[]models.Book) error {
	args := m.Called(title, author, genre, isbn, availability, yearFrom, yearTo, titleSort, authorSort, yearSort, books)
	return args.Error(0)
}
func (m *MockCategoryUsecase) GetCategories() []models.Category {
	args := m.Called()
	return args.Get(0).([]models.Category)
}
func (m *MockCategoryUsecase) AddCategory(category *models.Category) {
	m.Called(category)
}

var jwtKey = []byte("jwt-key")

func GenerateTestJWT(user_id uint) (string, error) {
	claims := &jwt.MapClaims{
		"user_id": user_id,
		"exp":     time.Now().Add(1 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}
func TestBookController(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	mockBookUsecase := new(MockBookUsecase)
	mockUserUsecase := new(MockUserUsecase)
	mockCategoryUsecase := new(MockCategoryUsecase)
	bookHandler := handlers.NewBookHandler(mockBookUsecase, mockUserUsecase, mockCategoryUsecase)
	router.LoadHTMLGlob("../../../../templates/*")
	protected := router.Group("/")
	protected.Use(routes.JwtHandler())
	protected.POST("/add-book", bookHandler.AddBook)
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	user := models.User{
		Name:     "testuser",
		Email:    "test@example.com",
		Password: string(hashedPassword),
		Role:     "admin",
		ID:       1, // Set an initial ID
	}

	mockUserUsecase.On("CreateUser", &user).Return(nil)
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
	mockUserUsecase.On("AuthenticateUser", mock.MatchedBy(func(u *models.User) bool {
		return u.Password == "password123"
	})).Return(true)
	token, _ := GenerateTestJWT(user.ID)
	mockUserUsecase.On("GenerateJWT", user.ID).Return(token, nil)
	mockCategoryUsecase.On("GetCategories").Return([]models.Category{{ID: 1, Name: "test"}})
	mockBookUsecase.On("CreateBook", mock.AnythingOfType("*models.Book")).Return(&models.Book{ID: 1})
	mockBookUsecase.On("DeleteBookIsbn", "99999999").Return()
	if err := mockUserUsecase.CreateUser(&user); err != nil {
		t.Fatal("Failed to create test user:", err)
	}

	var createdUser models.User
	mockUserUsecase.GetUserByEmail(user.Email, &createdUser)
	createdUser.Password = "password123"
	if createdUser.ID == 0 {
		t.Fatal("User was not properly created")
	}
	user.ID = createdUser.ID
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
			name: "Success - Authenticated User - Add a new Book",
			setupAuth: func(req *http.Request) {
				token, _ := mockUserUsecase.GenerateJWT(user.ID)
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
			book := models.Book{
				Title:         "Test Book",
				Author:        "test author",
				Price:         122,
				CategoryID:    mockCategoryUsecase.GetCategories()[0].ID,
				Isbn:          "99999999",
				PublishedDate: "01-01-2000",
				Availability:  true,
			}
			jsonValue, _ := json.Marshal(book)
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/add-book", bytes.NewBuffer(jsonValue))
			req.Header.Set("Content-Type", "application/json")
			test.setupAuth(req)
			router.ServeHTTP(w, req)
			assert.Equal(t, http.StatusFound, w.Code)
			assert.Equal(t, w.Header().Get("location"), "/books")
			mockBookUsecase.DeleteBookIsbn("99999999")
		})
	}
}
