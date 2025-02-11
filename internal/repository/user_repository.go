package repository

import (
	"bookstore/internal/domain/models"
	"fmt"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}
func (r *userRepository) CreateUser(u *models.User) error {
	result := r.db.Create(&u)
	return result.Error
}
func (r *userRepository) AuthenticateUser(RequestUser *models.User) bool {
	var user models.User
	r.db.Find(&user, "name = ?", RequestUser.Name)
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(RequestUser.Password))
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}
func (r *userRepository) GetUser(user *models.User, RequestUser *models.User) {
	r.db.Find(user, "name = ?", RequestUser.Name)
}

func (r *userRepository) GetUserByID(user *models.User, user_id uint) {
	r.db.Find(&user, user_id)
}
func (r *userRepository) DeleteUserByName(name string) error {
	var user models.User
	if err := r.db.Where("name = ?", name).First(&user).Error; err != nil {
		return err
	}
	result := r.db.Where("name = ?", name).Delete(&models.User{})
	return result.Error
}
func (r *userRepository) GetUserByEmail(email string, user *models.User) {
	r.db.Where("email = ?", email).Find(user)
}
