package services

import (
	"errors"
	"website-dummy/internal/config"
	"website-dummy/internal/models"
)

type UserService struct{}

// NewUserService creates a new user service
func NewUserService() *UserService {
	return &UserService{}
}

// CreateUser creates a new user with business logic validation
func (us *UserService) CreateUser(user *models.User) error {
	if user.Name == "" {
		return errors.New("name is required")
	}
	
	if user.Email == "" {
		return errors.New("email is required")
	}

	db := config.GetDB()
	if db == nil {
		return errors.New("database connection not available")
	}

	// Check if user with email already exists
	var existingUser models.User
	if err := db.Where("email = ?", user.Email).First(&existingUser).Error; err == nil {
		return errors.New("user with this email already exists")
	}

	return db.Create(user).Error
}

// GetUserByID retrieves a user by ID
func (us *UserService) GetUserByID(id uint) (*models.User, error) {
	db := config.GetDB()
	if db == nil {
		return nil, errors.New("database connection not available")
	}

	var user models.User
	if err := db.First(&user, id).Error; err != nil {
		return nil, errors.New("user not found")
	}

	return &user, nil
}

// GetAllUsers retrieves all users
func (us *UserService) GetAllUsers() ([]models.User, error) {
	db := config.GetDB()
	if db == nil {
		return nil, errors.New("database connection not available")
	}

	var users []models.User
	if err := db.Find(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
}
