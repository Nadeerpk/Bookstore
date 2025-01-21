package controllers_test

import (
	"bookstore/src/models"
	"bookstore/src/routes"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestLoginController(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	routes.SetupRoutes(router)
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	user := models.User{
		Name: "testuser", Email: "test@example.com",
		Password: string(hashedPassword), Role: "user"}

	if err := user.CreateUser(); err != nil {
		t.Fatal("Failed to create test user:", err)
	}
	var createdUser models.User
	models.GetUserByEmail(user.Email, &createdUser)
	createdUser.Password = "password123"
	t.Cleanup(func() {
		models.DeleteUserByName(createdUser.Name)
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
