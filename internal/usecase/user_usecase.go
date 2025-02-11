package usecase

import (
	"bookstore/internal/domain/models"
	"bookstore/internal/repository"
	"time"

	"github.com/golang-jwt/jwt"
)

var jwtKey = []byte("jwt-key")

type UserUseCase interface {
	CreateUser(user *models.User) error
	AuthenticateUser(RequestUser *models.User) bool
	GetUser(user *models.User, RequestUser *models.User)
	GetUserByID(user *models.User, user_id uint)
	DeleteUserByName(name string) error
	GetUserByEmail(email string, user *models.User)
	GenerateJWT(user_id uint) (string, error)
}
type userUseCase struct {
	UserRepo repository.UserRepository
}

func NewUserUsecase(userRepo repository.UserRepository) UserUseCase {
	return &userUseCase{UserRepo: userRepo}
}
func (u *userUseCase) CreateUser(user *models.User) error {
	return u.UserRepo.CreateUser(user)
}
func (u *userUseCase) AuthenticateUser(RequestUser *models.User) bool {
	return u.UserRepo.AuthenticateUser(RequestUser)
}
func (u *userUseCase) GetUser(user *models.User, RequestUser *models.User) {
	u.UserRepo.GetUser(user, RequestUser)
}
func (u *userUseCase) GetUserByID(user *models.User, user_id uint) {
	u.UserRepo.GetUserByID(user, user_id)
}
func (u *userUseCase) DeleteUserByName(name string) error {
	return u.UserRepo.DeleteUserByName(name)
}
func (u *userUseCase) GetUserByEmail(email string, user *models.User) {
	u.UserRepo.GetUserByEmail(email, user)
}
func (u *userUseCase) GenerateJWT(user_id uint) (string, error) {
	claims := &jwt.MapClaims{
		"user_id": user_id,
		"exp":     time.Now().Add(1 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}
