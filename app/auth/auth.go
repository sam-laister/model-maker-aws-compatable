package auth

import (
	"errors"
	"log"

	"firebase.google.com/go/v4/auth"
	"gorm.io/gorm"

	model "github.com/Soup666/diss-api/model"
)

// AuthService provides authentication services
type AuthService struct {
	DB       *gorm.DB
	FireAuth *auth.Client
}

// Login authenticates a user with the provided credentials and returns a Firebase custom token
func (s *AuthService) Login(authToken string) (*model.User, error) {

	var user model.User
	err := s.DB.Where("firebase_uid = ?", authToken).First(&user).Error

	// Get the user from the database
	if err != nil {
		if err == gorm.ErrRecordNotFound {

			// First time - create user
			result := s.DB.Create(&model.User{FirebaseUid: authToken})
			if result.Error != nil {
				log.Printf("failed to insert user into database: %+v", result.Error)
				return nil, errors.New("internal server error")
			}
			err := s.DB.Where("firebase_uid = ?", authToken).First(&user).Error
			if err != nil {
				log.Printf("failed to get user by email from database: %v", err)
				return nil, errors.New("internal server error")
			}
			return &user, nil
		}
		log.Printf("failed to get user by email from database: %v", err)
		return nil, errors.New("internal server error")
	}

	return &user, nil
}

// Register creates a new user with the provided credentials and returns token
func (s *AuthService) Register(user model.User) (*model.User, error) {
	// Create a new user in the database
	result := s.DB.Create(&user)
	if result.Error != nil {
		log.Printf("failed to insert user into database: %+v", result.Error)
		return nil, errors.New("internal server error")
	}

	return &user, nil
}
