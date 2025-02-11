package models

type Review struct {
	ID      uint   `json:"id" form:"id" gorm:"primary_key"`
	BookID  uint   `json:"book_id" form:"book_id"`
	UserID  uint   `json:"user_id" form:"user_id"`
	User    User   `gorm:"foreignKey:UserID"`
	Comment string `json:"comment" form:"comment"`
}
