package models

import (
	"gorm.io/gorm"
)

type Category struct {
	gorm.Model
	ID    uint   `json:"id" form:"id" gorm:"primary_key"`
	Name  string `json:"name" form:"name" binding:"required" gorm:"not null"`
	Books []Book `gorm:"foreignKey:CategoryID"`
}
