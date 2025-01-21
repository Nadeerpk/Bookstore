package models

import (
	"bookstore/src/config"
	"fmt"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var db *gorm.DB

type User struct {
	ID       uint     `json:"id" form:"id" gorm:"primary_key"`
	Name     string   `json:"username" form:"name" binding:"required" gorm:"not null"`
	Password string   `json:"password" form:"password" binding:"required" gorm:"not null"`
	Email    string   `json:"email" form:"email" binding:"required" gorm:"not null;unique"`
	Role     string   `json:"role" form:"role" binding:"required" gorm:"not null"`
	Reviews  []Review `json:"reviews" form:"reviews" gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	Orders   []Order  `json:"orders" form:"orders" gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}

func init() {
	db = config.Getdb()
	db.AutoMigrate(&User{})
}

func (u *User) CreateUser() error {
	result := db.Create(&u)
	return result.Error
}

func AuthenticateUser(RequestUser *User) bool {
	var user User
	db.Find(&user, "name = ?", RequestUser.Name)
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(RequestUser.Password))
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}
func GetUser(user *User, RequestUser *User) {
	db.Find(user, "name = ?", RequestUser.Name)
}

func GetUserByID(user *User, user_id uint) {
	db.Find(&user, user_id)
}
func DeleteUserByName(name string) error {
	var user User
	if err := db.Where("name = ?", name).First(&user).Error; err != nil {
		return err
	}
	result := db.Where("name = ?", name).Delete(&User{})
	return result.Error
}
func GetUserByEmail(email string, user *User) {
	db.Where("email = ?", email).Find(user)
}
