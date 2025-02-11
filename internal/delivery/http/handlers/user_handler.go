package handlers

import (
	"bookstore/internal/domain/models"
	"bookstore/internal/usecase"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	UserUsecase usecase.UserUseCase
}

func NewUserHandler(userUseCase usecase.UserUseCase) *UserHandler {
	return &UserHandler{UserUsecase: userUseCase}
}
func (h *UserHandler) RegisterUser(c *gin.Context) {
	User := &models.User{}
	if err := c.ShouldBind(User); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if User.Name == "" || User.Password == "" || User.Email == "" || User.Role == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "All fields are required"})
		return
	}

	if !strings.Contains(User.Email, "@") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email format"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(User.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error hashing password"})
		return
	}

	User.Password = string(hashedPassword)
	if err := h.UserUsecase.CreateUser(User); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating user"})
		return
	}

	c.HTML(http.StatusOK, "login.html", gin.H{})
}

func (h *UserHandler) LoginUser(c *gin.Context) {
	User := &models.User{}
	c.ShouldBind(User)
	authenticated := h.UserUsecase.AuthenticateUser(User)
	if !authenticated {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication failed"})
		return
	}
	var user models.User
	h.UserUsecase.GetUser(&user, User)
	token, err := h.UserUsecase.GenerateJWT(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating token"})
		return
	}
	c.SetCookie("jwt_token", token, 3600, "/", "localhost", false, true)

	c.Redirect(http.StatusFound, "/books")
}

func (h *UserHandler) LogoutUser(c *gin.Context) {
	c.SetCookie("jwt_token", "", -1, "/", "localhost", false, true)
	c.Redirect(http.StatusFound, "/login")
}
