package models

import (
	"bookstore/src/config"

	"gorm.io/gorm"
)

type Category struct {
	gorm.Model
	ID    uint   `json:"id" form:"id" gorm:"primary_key"`
	Name  string `json:"name" form:"name" binding:"required" gorm:"not null"`
	Books []Book `gorm:"foreignKey:CategoryID"`
}

func init() {
	db = config.Getdb()
	db.AutoMigrate(&Category{})
}
func GetCategories() []Category {
	var categories []Category
	db.Find(&categories)
	return categories
}
func AddCategory(category *Category) {
	db.Create(&category)
}

func DeleteCategory(category *Category) {
	db.Delete(&category)
}

func UpdateCategory(category *Category) {
	db.Save(&category)
}
