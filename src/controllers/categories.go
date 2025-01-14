package controllers

import (
	"bookstore/src/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetCategories(c *gin.Context) {
	categories := models.GetCategories()
	c.HTML(http.StatusOK, "categories.html", gin.H{"categories": categories})
}
func AddCategory(c *gin.Context) {
	category := &models.Category{}
	c.ShouldBind(category)
	models.AddCategory(category)
	c.Redirect(http.StatusFound, "/categories")
}
func UpdateCategory(c *gin.Context) {
	category := &models.Category{}
	c.ShouldBind(category)
	models.UpdateCategory(category)
	c.Redirect(http.StatusFound, "/categories")
}
func DeleteCategory(c *gin.Context) {
	category := &models.Category{}
	c.ShouldBind(category)
	models.DeleteCategory(category)
	c.Redirect(http.StatusFound, "/categories")
}
