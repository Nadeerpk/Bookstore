package handlers_test

import (
	"bookstore/internal/delivery/http/handlers"
	"bookstore/internal/domain/models"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

func (m *MockUserUsecase) AuthenticateUser(user *models.User) bool {
	args := m.Called(user)
	return args.Bool(0)
}
func (m *MockUserUsecase) GetUserByEmail(email string, user *models.User) {
	m.Called(email, user)
}
func (m *MockUserUsecase) DeleteUserByName(email string) error {
	args := m.Called(email)
	return args.Error(0)
}

func (m *MockUserUsecase) GenerateJWT(user_id uint) (string, error) {
	args := m.Called(user_id)
	return args.String(0), args.Error(1)
}
func (m *MockUserUsecase) GetUser(user *models.User, RequestUser *models.User) {
	m.Called(user, RequestUser)
}
func (m *MockUserUsecase) GetUserByID(user *models.User, user_id uint) {
	m.Called(user, user_id)
}
func TestLoginController(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	mockUserUsecase := new(MockUserUsecase)
	userHandler := handlers.NewUserHandler(mockUserUsecase)

	router.POST("/login", userHandler.LoginUser)
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	user := models.User{
		Name:     "testuser",
		Email:    "test@example.com",
		Password: string(hashedPassword),
		Role:     "user",
	}

	// Setup mock expectations
	mockUserUsecase.On("CreateUser", &user).Return(nil)
	mockUserUsecase.On("GetUserByEmail", user.Email, mock.AnythingOfType("*models.User")).Run(func(args mock.Arguments) {
		// args[1] is the user pointer parameter
		userArg := args[1].(*models.User)
		*userArg = user
	}).Return()
	mockUserUsecase.On("GetUser", &models.User{}, mock.AnythingOfType("*models.User")).Return(nil)
	mockUserUsecase.On("DeleteUserByName", user.Name).Return(nil)
	mockUserUsecase.On("AuthenticateUser", mock.MatchedBy(func(u *models.User) bool {
		return u.Password == "password123"
	})).Return(true)
	token, _ := GenerateTestJWT(user.ID)
	mockUserUsecase.On("GenerateJWT", user.ID).Return(token, nil)
	if err := mockUserUsecase.CreateUser(&user); err != nil {
		t.Fatal("Failed to create test user:", err)
	}

	var createdUser models.User
	mockUserUsecase.GetUserByEmail(user.Email, &createdUser)
	createdUser.Password = "password123"
	t.Cleanup(func() {
		mockUserUsecase.DeleteUserByName(createdUser.Name)
	})

	jsonValue, _ := json.Marshal(createdUser)
	t.Run("login", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonValue))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusFound, w.Code)
		assert.Equal(t, w.Header().Get("location"), "/books")
	})
}
