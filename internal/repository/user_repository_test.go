package repository_test

import (
	"bookstore/internal/domain/models"
	"bookstore/internal/repository"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type UserRepositoryTestSuite struct {
	suite.Suite
	db       *gorm.DB
	repo     repository.UserRepository
	testUser *models.User
}

func (suite *UserRepositoryTestSuite) SetupTest() {
	// MySQL connection configuration
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		"testuser",       // username
		"password",       // password
		"localhost",      // host
		"3306",           // port
		"bookstore_test", // database name
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	assert.NoError(suite.T(), err)

	// Clean database before each test
	db.Exec("DROP TABLE IF EXISTS users")
	db.AutoMigrate(&models.User{})

	suite.db = db
	suite.repo = repository.NewUserRepository(db)

	// Create test user
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("testpass"), bcrypt.DefaultCost)
	suite.testUser = &models.User{
		Name:     "testuser",
		Password: string(hashedPassword),
		Email:    "test@example.com",
		Role:     "user",
	}
}

func (suite *UserRepositoryTestSuite) TearDownTest() {
	// Clean up after each test
	db, err := suite.db.DB()
	if err == nil {
		db.Close()
	}
}

func TestUserRepository(t *testing.T) {
	suite.Run(t, new(UserRepositoryTestSuite))
}

func (suite *UserRepositoryTestSuite) TestCreateUser() {
	t := suite.T()

	err := suite.repo.CreateUser(suite.testUser)
	assert.NoError(t, err)

	var savedUser models.User
	suite.db.First(&savedUser, suite.testUser.ID)
	assert.Equal(t, suite.testUser.Name, savedUser.Name)
}

func (suite *UserRepositoryTestSuite) TestAuthenticateUser() {
	t := suite.T()

	suite.repo.CreateUser(suite.testUser)

	// Test valid credentials
	loginUser := &models.User{
		Name:     "testuser",
		Password: "testpass",
	}
	assert.True(t, suite.repo.AuthenticateUser(loginUser))

	// Test invalid password
	loginUser.Password = "wrongpass"
	assert.False(t, suite.repo.AuthenticateUser(loginUser))
}

func (suite *UserRepositoryTestSuite) TestGetUser() {
	t := suite.T()

	suite.repo.CreateUser(suite.testUser)

	var fetchedUser models.User
	requestUser := &models.User{Name: suite.testUser.Name}
	suite.repo.GetUser(&fetchedUser, requestUser)

	assert.Equal(t, suite.testUser.Name, fetchedUser.Name)
	assert.Equal(t, suite.testUser.Email, fetchedUser.Email)
}

func (suite *UserRepositoryTestSuite) TestGetUserByID() {
	t := suite.T()

	suite.repo.CreateUser(suite.testUser)

	var fetchedUser models.User
	suite.repo.GetUserByID(&fetchedUser, suite.testUser.ID)

	assert.Equal(t, suite.testUser.Name, fetchedUser.Name)
	assert.Equal(t, suite.testUser.Email, fetchedUser.Email)
}

func (suite *UserRepositoryTestSuite) TestDeleteUserByName() {
	t := suite.T()

	suite.repo.CreateUser(suite.testUser)

	err := suite.repo.DeleteUserByName(suite.testUser.Name)
	assert.NoError(t, err)

	var count int64
	suite.db.Model(&models.User{}).Where("name = ?", suite.testUser.Name).Count(&count)
	assert.Equal(t, int64(0), count)
}

func (suite *UserRepositoryTestSuite) TestGetUserByEmail() {
	t := suite.T()

	suite.repo.CreateUser(suite.testUser)

	var fetchedUser models.User
	suite.repo.GetUserByEmail(suite.testUser.Email, &fetchedUser)

	assert.Equal(t, suite.testUser.Name, fetchedUser.Name)
	assert.Equal(t, suite.testUser.Email, fetchedUser.Email)
}
