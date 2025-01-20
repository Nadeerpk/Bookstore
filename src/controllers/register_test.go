package controllers_test

import (
	"bookstore/src/controllers"
	"bookstore/src/models"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestRegister(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name       string
		user       models.User
		statusCode int
	}{
		{
			name: "Valid Registration",
			user: models.User{
				Name:     "testuser",
				Password: "password123",
				Email:    "test@example.com",
				Role:     "user",
			},
			statusCode: http.StatusOK,
		},
		{
			name: "Missing Required Fields",
			user: models.User{
				Name:     "",
				Password: "password123",
				Email:    "test@example.com",
				Role:     "user",
			},
			statusCode: http.StatusBadRequest,
		},
		{
			name: "Invalid Email Format",
			user: models.User{
				Name:     "testuser2",
				Password: "password123",
				Email:    "invalid-email",
				Role:     "user",
			},
			statusCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.New()
			router.LoadHTMLGlob("../../templates/*")
			if tt.statusCode == http.StatusOK {
				t.Cleanup(func() {
					models.DeleteUserByName(tt.user.Name)
				})
			}
			router.POST("/register", controllers.RegisterController)

			jsonValue, _ := json.Marshal(tt.user)
			req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(jsonValue))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)
			assert.Equal(t, tt.statusCode, w.Code)
		})
	}
}
