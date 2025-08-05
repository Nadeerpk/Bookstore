package handlers_test

import (
	// "bookstore/src/controllers"
	// "bookstore/src/models"
	// "bookstore/src/routes"
	"bookstore/internal/delivery/http/handlers"
	"bookstore/internal/delivery/http/routes"
	"bookstore/internal/domain/models"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockOrderUseCase struct {
	mock.Mock
}
type MockCartUsecase struct {
	mock.Mock
}

func (m *MockOrderUseCase) AddOrder(order *models.Order) error {
	args := m.Called(order)
	return args.Error(0)
}
func (m *MockOrderUseCase) GetOrdersByUserID(userID uint, orders *[]models.Order) error {
	args := m.Called(userID, orders)
	return args.Error(0)
}
func (m *MockCartUsecase) GetCart(cart *models.Cart, userID uint) error {
	args := m.Called(cart, userID)
	return args.Error(0)
}
func (m *MockCartUsecase) AddToCart(bookID, user_id uint, cart models.Cart) error {
	args := m.Called(bookID, user_id, cart)
	return args.Error(0)
}
func (m *MockCartUsecase) DeleteFromCart(bookID, user_id uint) error {
	args := m.Called(bookID, user_id)
	return args.Error(0)
}

func TestOrderController(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	mockOrderUseCase := new(MockOrderUseCase)
	mockBookUseCase := new(MockBookUsecase)
	mockCategoryUseCase := new(MockCategoryUsecase)
	mockCartUseCase := new(MockCartUsecase)
	mockUserUseCase := new(MockUserUsecase)
	orderhandler := handlers.NewOrderHandler(mockOrderUseCase, mockCartUseCase)
	router.LoadHTMLGlob("../../../../templates/*")
	// routes.SetupRoutes(router)
	protected := router.Group("/")
	protected.Use(routes.JwtHandler())
	protected.GET("/order/:book_id", orderhandler.AddOrder)
	user := &models.User{
		ID:       1,
		Name:     "testuser",
		Email:    "nadeer@qburst.com",
		Password: "password123",
		Role:     "user",
	}
	mockUserUseCase.On("CreateUser", user).Return(nil)
	token, _ := GenerateTestJWT(user.ID)
	mockUserUseCase.On("GenerateJWT", user.ID).Return(token, nil)
	mockUserUseCase.On("GetUserByEmail", user.Email, mock.AnythingOfType("*models.User")).Run(func(args mock.Arguments) {
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
	mockUserUseCase.On("DeleteUserByName", user.Name).Return(nil)
	mockUserUseCase.On("AuthenticateUser", mock.MatchedBy(func(u *models.User) bool {
		return u.Password == "password123"
	})).Return(true)
	mockCategoryUseCase.On("GetCategories").Return([]models.Category{{ID: 1, Name: "test"}})
	mockBookUseCase.On("CreateBook", mock.AnythingOfType("*models.Book")).Return(&models.Book{ID: 1})
	mockBookUseCase.On("DeleteBookIsbn", "123456789").Return()
	mockBookUseCase.On("GetBookByID", mock.AnythingOfType("uint"), mock.AnythingOfType("*models.Book")).Return(nil)
	mockCartUseCase.On("AddToCart", mock.AnythingOfType("uint"), mock.AnythingOfType("uint"), mock.AnythingOfType("models.Cart")).Return(nil)
	mockOrderUseCase.On("AddOrder", mock.AnythingOfType("*models.Order")).Return(nil)
	mockCartUseCase.On("DeleteFromCart", mock.AnythingOfType("uint"), mock.AnythingOfType("uint")).Return(nil)
	if err := mockUserUseCase.CreateUser(user); err != nil {
		t.Fatal("Failed to create test user:", err)
	}
	var createdUser models.User
	mockUserUseCase.GetUserByEmail(user.Email, &createdUser)
	if createdUser.ID == 0 {
		t.Fatal("User was not properly created")
	}
	user.ID = createdUser.ID
	t.Cleanup(func() {
		mockUserUseCase.DeleteUserByName(user.Name)
		mockBookUseCase.DeleteBookIsbn("123456789")
	})
	book := &models.Book{Title: "Test Book 1", Author: "Author 1", Price: 29.99, Isbn: "123456789",
		PublishedDate: "01-01-2000", CategoryID: mockCategoryUseCase.GetCategories()[0].ID}
	mockBookUseCase.CreateBook(book)
	var createdBook models.Book
	mockBookUseCase.GetBookByID(book.ID, &createdBook)
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
				token, _ := mockUserUseCase.GenerateJWT(user.ID)
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
			mockCartUseCase.AddToCart(uint(createdBook.ID), uint(user.ID), cart)
			req, _ := http.NewRequest("GET", "/order/"+book_id_str, nil)
			test.setupAuth(req)
			router.ServeHTTP(w, req)
			assert.Equal(t, http.StatusFound, w.Code)
			assert.Equal(t, w.Header().Get("location"), "/cart")
		})
	}
}
