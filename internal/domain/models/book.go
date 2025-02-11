package models

import (
	"encoding/base64"

	"html/template"
)

type Book struct {
	ID            uint       `json:"id" form:"id" gorm:"primary_key"`
	Title         string     `json:"title" form:"title" binding:"required" gorm:"not null"`
	Author        string     `json:"author" form:"author" binding:"required" gorm:"not null"`
	Price         float64    `json:"price" form:"price" binding:"required" gorm:"not null"`
	Category      Category   `gorm:"foreignKey:CategoryID"`
	Isbn          string     `json:"isbn" form:"isbn" binding:"required" gorm:"not null; unique"`
	PublishedDate string     `json:"published_date" form:"published_date" binding:"required" gorm:"not null"`
	Availability  bool       `json:"availability" form:"availability" binding:"required" gorm:"not null"`
	CategoryID    uint       `json:"category_id" form:"category_id"`
	CartItems     []CartItem `gorm:"foreignKey:BookID;constraint:OnDelete:CASCADE"`
	Reviews       []Review   `json:"reviews" form:"reviews" gorm:"foreignKey:BookID;constraint:OnDelete:CASCADE"`
	Orders        []Order    `json:"orders" form:"orders" gorm:"foreignKey:BookID;constraint:OnDelete:CASCADE"`
	Image         []byte     `json:"image" form:"image" gorm:"type:mediumblob"`
}

func (b *Book) GetImageBase64() template.URL {
	if len(b.Image) == 0 {
		return ""
	}
	return template.URL("data:image/jpeg;base64," + base64.StdEncoding.EncodeToString(b.Image))
}
