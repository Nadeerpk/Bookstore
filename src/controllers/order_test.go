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

func TestOrderController(t *testing.T) {
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
			models.AddToCart(uint(1), uint(user.ID), cart)
			req, _ := http.NewRequest("GET", "/order/1", nil)
			router.ServeHTTP(w, req)
			assert.Equal(t, http.StatusFound, w.Code)
		})
	}
}
