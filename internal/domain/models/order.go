package models

import (
	"gorm.io/gorm"
)

type Order struct {
	gorm.Model
	ID     uint `json:"id" form:"id" gorm:"primary_key"`
	UserID uint `json:"user_id" form:"user_id"`
	User   User `gorm:"foreignKey:UserID"`
	BookID uint `json:"book_id" form:"book_id"`
	Book   Book `gorm:"foreignKey:BookID"`
}
