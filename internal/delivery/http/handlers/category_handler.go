package handlers

import (
	"bookstore/internal/domain/models"
	"bookstore/internal/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CategoryHandler struct {
	CategoryUsecase usecase.CategoryUseCase
}

func NewCategoryHandler(categoryUseCase usecase.CategoryUseCase) *CategoryHandler {
	return &CategoryHandler{CategoryUsecase: categoryUseCase}
}
func (h *CategoryHandler) GetCategories(c *gin.Context) {
	categories := h.CategoryUsecase.GetCategories()
	c.HTML(http.StatusOK, "categories.html", gin.H{"categories": categories})
}
func (h *CategoryHandler) AddCategory(c *gin.Context) {
	category := &models.Category{}
	c.ShouldBind(category)
	h.CategoryUsecase.AddCategory(category)
	c.Redirect(http.StatusFound, "/categories")
}
