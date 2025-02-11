package models

type User struct {
	ID       uint     `json:"id" form:"id" gorm:"primary_key"`
	Name     string   `json:"username" form:"name" binding:"required" gorm:"not null"`
	Password string   `json:"password" form:"password" binding:"required" gorm:"not null"`
	Email    string   `json:"email" form:"email" binding:"required" gorm:"not null;unique"`
	Role     string   `json:"role" form:"role" binding:"required" gorm:"not null"`
	Reviews  []Review `json:"reviews" form:"reviews" gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	Orders   []Order  `json:"orders" form:"orders" gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}
