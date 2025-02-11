package models

type Cart struct {
	ID     uint       `json:"id" form:"id" gorm:"primary_key"`
	UserID uint       `json:"user_id" form:"user_id"`
	Items  []CartItem `gorm:"foreignKey:CartID;constraint:OnDelete:CASCADE"`
}

type CartItem struct {
	CartID   uint `json:"cart_id" form:"cart_id" gorm:"not null"`
	BookID   uint `json:"book_id" form:"book_id" gorm:"not null"`
	Quantity uint `json:"quantity" form:"quantity" gorm:"default:1"`
	Book     Book `json:"book" gorm:"foreignKey:BookID"`
}
