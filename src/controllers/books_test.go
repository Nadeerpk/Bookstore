package controllers_test

import (
	"bookstore/src/controllers"
	"bookstore/src/routes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGetBooks(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	routes.SetupRoutes(router)

	t.Run("Success - Get Books List", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/books", nil)
		token, _ := controllers.GenerateJWT(1)
		cookie := &http.Cookie{
			Name:  "jwt_token",
			Value: token,
			Path:  "/",
		}
		req.AddCookie(cookie)
		router.ServeHTTP(w, req)
		fmt.Println(w.Body.String())
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Header().Get("Content-Type"), "text/html")
		assert.Contains(t, w.Body.String(), "<title>Book List</title>")
		assert.Contains(t, w.Body.String(), "<th>Title</th>")
	})
}
