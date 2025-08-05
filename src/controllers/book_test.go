package controllers_test

import (
	"bookstore/src/controllers"
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

func BookController(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.LoadHTMLGlob("../../templates/*")
	routes.SetupRoutes(router)
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	user := models.User{
		Name: "testuser", Email: "test@example.com",
		Password: string(hashedPassword), Role: "admin"}

	if err := user.CreateUser(); err != nil {
		t.Fatal("Failed to create test user:", err)
	}

	var createdUser models.User
	models.GetUserByEmail(user.Email, &createdUser)
	createdUser.Password = "password123"
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
			name: "Success - Authenticated User - Add a new Book",
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
			book := models.Book{
				Title:         "Test Book",
				Author:        "test author",
				Price:         122,
				CategoryID:    models.GetCategories()[0].ID,
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
			models.DeleteBookIsbn("99999999")
		})
	}
}
