package controllers

import (
	"bookstore/src/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func AddReview(c *gin.Context) {
	review := &models.Review{}
	c.ShouldBind(review)
	user_id := c.GetFloat64("user_id")
	book_id_str := c.Param("book_id")
	book_id, _ := strconv.ParseUint(book_id_str, 10, 64)
	review.UserID = uint(user_id)
	review.BookID = uint(book_id)
	models.AddReview(review)
	c.Redirect(http.StatusFound, "/reviews/"+book_id_str)
}
func GetReviews(c *gin.Context) {
	book_id := c.Param("book_id")
	book_id_int, _ := strconv.ParseUint(book_id, 10, 32)
	var reviews []models.Review
	models.GetReviewsByBookID(uint(book_id_int), &reviews)

	c.HTML(200, "reviews.html", gin.H{"reviews": reviews,
		"book_id": book_id})
}
