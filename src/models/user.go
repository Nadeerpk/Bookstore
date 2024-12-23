package models

import (
	"bookstore/src/config"
	"fmt"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var db *gorm.DB

type User struct {
	gorm.Model
	ID       uint   `json:"id" form:"id" gorm:"primary_key"`
	Name     string `json:"username" form:"name" binding:"required" gorm:"not null"`
	Password string `json:"password" form:"password" binding:"required" gorm:"not null"`
	Email    string `json:"email" form:"email" binding:"required" gorm:"not null;unique"`
}

func init() {
	db = config.Getdb()
	db.AutoMigrate(&User{})
}

func (u *User) CreateUser() *User {
	db.Create(&u)
	return u
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
